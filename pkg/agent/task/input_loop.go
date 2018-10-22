//
// Copyright © 2018 Aljabr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package task

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/eapache/go-resiliency/breaker"
	"github.com/golang/protobuf/ptypes"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"

	"github.com/AljabrIO/koalja-operator/pkg/agent/link"
	koalja "github.com/AljabrIO/koalja-operator/pkg/apis/koalja/v1alpha1"
	"github.com/AljabrIO/koalja-operator/pkg/constants"
	"github.com/AljabrIO/koalja-operator/pkg/event"
	"github.com/dchest/uniuri"
)

// inputLoop subscribes to all inputs of a task and build
// snapshots of incoming events, according to the policy on each input.
type inputLoop struct {
	log             zerolog.Logger
	spec            *koalja.TaskSpec
	inputAddressMap map[string]string // map[inputName]eventSourceAddress
	clientID        string
	snapshot        InputSnapshot
	mutex           sync.Mutex
	executionCount  int32
	execQueue       chan (*InputSnapshot)
	executor        Executor
}

// newInputLoop initializes a new input loop.
func newInputLoop(log zerolog.Logger, spec *koalja.TaskSpec, pod *corev1.Pod, executor Executor) (*inputLoop, error) {
	inputAddressMap := make(map[string]string)
	for _, tis := range spec.Inputs {
		annKey := constants.CreateInputLinkAddressAnnotationName(tis.Name)
		address := pod.GetAnnotations()[annKey]
		if address == "" {
			return nil, fmt.Errorf("No input address annotation found for input '%s'", tis.Name)
		}
		inputAddressMap[tis.Name] = address
	}
	return &inputLoop{
		log:             log,
		spec:            spec,
		inputAddressMap: inputAddressMap,
		clientID:        uniuri.New(),
		execQueue:       make(chan *InputSnapshot),
		executor:        executor,
	}, nil
}

// Run the input loop until the given context is canceled.
func (il *inputLoop) Run(ctx context.Context) error {
	defer close(il.execQueue)
	g, lctx := errgroup.WithContext(ctx)
	if len(il.spec.Inputs) > 0 {
		// Watch inputs
		for _, tis := range il.spec.Inputs {
			tis := tis // Bring in scope
			g.Go(func() error {
				return il.watchInput(lctx, tis)
			})
		}
	} else {
		// No inputs, run executor all the time
		g.Go(func() error {
			return il.runExecWithoutInputs(lctx)
		})
	}
	g.Go(func() error {
		return il.processExecQueue(lctx)
	})
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

// processExecQueue pulls snapshots from the exec queue and executes them.
func (il *inputLoop) processExecQueue(ctx context.Context) error {
	for {
		select {
		case snapshot, ok := <-il.execQueue:
			if !ok {
				return nil
			}
			if err := il.execOnSnapshot(ctx, snapshot); ctx.Err() != nil {
				return ctx.Err()
			} else if err != nil {
				il.log.Error().Err(err).Msg("Failed to execute task")
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// runExecWithoutInputs pulls snapshots from the exec queue and executes them.
func (il *inputLoop) runExecWithoutInputs(ctx context.Context) error {
	b := breaker.New(5, 1, time.Second*10)
	for {
		snapshot := &InputSnapshot{}
		if err := b.Run(func() error {
			if err := il.execOnSnapshot(ctx, snapshot); ctx.Err() != nil {
				return ctx.Err()
			} else if err != nil {
				il.log.Error().Err(err).Msg("Failed to execute task")
				return maskAny(err)
			}
			return nil
		}); ctx.Err() != nil {
			return ctx.Err()
		} else if err == breaker.ErrBreakerOpen {
			// Circuit break open
			select {
			case <-time.After(time.Second * 5):
				// Retry
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

// execOnSnapshot executes the task executor for the given snapshot.
func (il *inputLoop) execOnSnapshot(ctx context.Context, snapshot *InputSnapshot) error {
	if err := il.executor.Execute(ctx, snapshot); ctx.Err() != nil {
		return ctx.Err()
	} else if err != nil {
		il.log.Debug().Err(err).Msg("executor.Execute failed")
		return maskAny(err)
	} else {
		// Acknowledge all event in the snapshot
		il.log.Debug().Msg("acknowledging all events in snapshot")
		if err := snapshot.AckAll(ctx); err != nil {
			il.log.Error().Err(err).Msg("Failed to acknowledge events")
		}
		return nil
	}
}

// processEvent the event coming from the given input.
func (il *inputLoop) processEvent(ctx context.Context, e *event.Event, tis koalja.TaskInputSpec, ack func(context.Context, *event.Event) error) error {
	snapshotPolicy := tis.SnapshotPolicy
	il.mutex.Lock()
	defer il.mutex.Unlock()
	if snapshotPolicy.IsAll() {
		// Wait until snapshot does not have an event for given input
		for {
			if !il.snapshot.HasEvent(tis.Name) {
				break
			}
			// Wait a bit
			il.mutex.Unlock()
			select {
			case <-time.After(time.Millisecond * 50):
				// Retry
				il.mutex.Lock()
			case <-ctx.Done():
				// Context canceled
				il.mutex.Lock()
				return ctx.Err()
			}
		}
	}

	// Set the event in the snapshot
	if err := il.snapshot.Set(ctx, tis.Name, e, ack); err != nil {
		return err
	}

	// See if we should execute the task now
	tuple := il.snapshot.CreateTuple(len(il.spec.Inputs))
	if tuple == nil || (il.executionCount > 0 && snapshotPolicy.IsLatest()) {
		// Not all inputs have received an event yet or this event has a "Latest" policy and we've executed once or more
		return nil
	}

	// Clone the snapshot
	clonedSnapshot := il.snapshot.Clone()
	il.executionCount++

	// Prepare snapshot for next execution
	for _, inp := range il.spec.Inputs {
		if inp.SnapshotPolicy.IsAll() {
			// Delete event
			il.snapshot.Delete(inp.Name)
		} else {
			// Remove need to acknowledge event
			il.snapshot.RemoveAck(inp.Name)
		}
	}

	// Push snapshot into execution queue
	il.mutex.Unlock()
	il.execQueue <- clonedSnapshot
	il.mutex.Lock()

	return nil
}

// watchInput subscribes to the given input and gathers events until the given context is canceled.
func (il *inputLoop) watchInput(ctx context.Context, tis koalja.TaskInputSpec) error {
	// Create client
	address := il.inputAddressMap[tis.Name]

	// Prepare loop
	subscribeAndReadEventLoop := func(ctx context.Context, c link.EventSourceClient) error {
		defer c.CloseConnection()
		resp, err := c.Subscribe(ctx, &event.SubscribeRequest{
			ClientID: il.clientID,
		})
		if err != nil {
			return err
		}
		subscr := *resp.GetSubscription()

		ack := func(ctx context.Context, e *event.Event) error {
			if _, err := c.AckEvent(ctx, &event.AckEventRequest{
				Subscription: &subscr,
				EventID:      e.GetID(),
			}); err != nil {
				il.log.Error().Err(err).Msg("Failed to ack event")
				return err
			}
			return nil
		}

		for {
			resp, err := c.NextEvent(ctx, &event.NextEventRequest{
				Subscription: &subscr,
				WaitTimeout:  ptypes.DurationProto(time.Second * 30),
			})
			if ctx.Err() != nil {
				return ctx.Err()
			} else if err != nil {
				// handle err
				il.log.Error().Err(err).Msg("Failed to fetch next event")
			} else {
				// Process event (if any)
				if e := resp.GetEvent(); e != nil {
					if err := il.processEvent(ctx, e, tis, ack); err != nil {
						il.log.Error().Err(err).Msg("Failed to process event")
					}
				}
			}
		}
	}

	// Keep creating connection, subscribe and loop
	for {
		c, err := link.CreateEventSourceClient(address)
		if err == nil {
			if err := subscribeAndReadEventLoop(ctx, c); ctx.Err() != nil {
				return ctx.Err()
			} else if err != nil {
				il.log.Error().Err(err).Msg("Failure in subscribe & read event loop")
			}
		} else if ctx.Err() != nil {
			return ctx.Err()
		} else {
			il.log.Error().Err(err).Msg("Failed to create event source client")
		}
		// Wait a bit
		select {
		case <-time.After(time.Second * 5):
			// Retry
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
