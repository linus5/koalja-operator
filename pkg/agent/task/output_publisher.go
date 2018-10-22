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
	"strings"
	"time"

	"github.com/dchest/uniuri"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"

	koalja "github.com/AljabrIO/koalja-operator/pkg/apis/koalja/v1alpha1"
	"github.com/AljabrIO/koalja-operator/pkg/constants"
	"github.com/AljabrIO/koalja-operator/pkg/event"
	ptask "github.com/AljabrIO/koalja-operator/pkg/task"
)

// outputPublisher is responsible for publishing events to one or more EventPublishers.
type outputPublisher struct {
	log                zerolog.Logger
	spec               *koalja.TaskSpec
	outputAddressesMap map[string][]string            // map[outputName]eventPublisherAddresses
	eventChannels      map[string][]chan *event.Event // map[outputName][]chan event
}

// newOutputPublisher initializes a new outputPublisher.
func newOutputPublisher(log zerolog.Logger, spec *koalja.TaskSpec, pod *corev1.Pod) (*outputPublisher, error) {
	outputAddressesMap := make(map[string][]string)
	eventChannels := make(map[string][]chan *event.Event)
	for _, tos := range spec.Outputs {
		annKey := constants.CreateOutputLinkAddressesAnnotationName(tos.Name)
		addressesStr := pod.GetAnnotations()[annKey]
		if addressesStr == "" {
			return nil, fmt.Errorf("No output addresses annotation found for input '%s'", tos.Name)
		}
		addresses := strings.Split(addressesStr, ",")
		outputAddressesMap[tos.Name] = addresses
		eventChans := make([]chan *event.Event, 0, len(addresses))
		for range addresses {
			eventChans = append(eventChans, make(chan *event.Event))
		}
		eventChannels[tos.Name] = eventChans
	}
	return &outputPublisher{
		log:                log,
		spec:               spec,
		outputAddressesMap: outputAddressesMap,
		eventChannels:      eventChannels,
	}, nil
}

// Run the output publisher until the given context is canceled.
func (op *outputPublisher) Run(ctx context.Context) error {
	g, lctx := errgroup.WithContext(ctx)
	for _, tos := range op.spec.Outputs {
		tos := tos // bring into scope
		addresses := op.outputAddressesMap[tos.Name]
		eventChans := op.eventChannels[tos.Name]
		for i, eventChan := range eventChans {
			eventChan := eventChan // bring into scope
			addr := addresses[i]
			g.Go(func() error {
				if err := op.runForOutput(lctx, tos, addr, eventChan); err != nil {
					return maskAny(err)
				}
				return nil
			})
		}
	}
	if err := g.Wait(); err != nil {
		return maskAny(err)
	}
	return nil
}

// PublishEvent pushes the given event onto the application output channel.
func (op *outputPublisher) PublishEvent(ctx context.Context, outputName string, evt event.Event) error {
	evt.ID = uniuri.New()
	evt.SourceTask = op.spec.Name
	evt.SourceTaskOutput = outputName

	eventChans, found := op.eventChannels[outputName]
	if !found {
		return fmt.Errorf("No channels found for output '%s'", outputName)
	}

	g, lctx := errgroup.WithContext(ctx)
	for _, eventChan := range eventChans {
		eventChan := eventChan // bring into scope
		g.Go(func() error {
			select {
			case eventChan <- &evt:
				// Done
				return nil
			case <-lctx.Done():
				// Context canceled
				return lctx.Err()
			}
		})
	}
	if err := g.Wait(); err != nil {
		return maskAny(err)
	}
	return nil
}

// OutputReady implements the notification endpoint for tasks with an "Auto" ready setting.
func (op *outputPublisher) OutputReady(ctx context.Context, req *ptask.OutputReadyRequest) (*ptask.OutputReadyResponse, error) {
	log := op.log.With().
		Str("output", req.GetOutputName()).
		Str("data", limitDataLength(req.GetEventData())).
		Logger()
	log.Debug().Msg("Received OutputReady request")

	evt := event.Event{
		Data: req.GetEventData(),
	}
	if err := op.PublishEvent(ctx, req.GetOutputName(), evt); err != nil {
		return nil, maskAny(err)
	}
	return &ptask.OutputReadyResponse{}, nil
}

// runForOutput keeps publishing events for the given output until the given context is canceled.
func (op *outputPublisher) runForOutput(ctx context.Context, tos koalja.TaskOutputSpec, addr string, eventChan chan *event.Event) error {
	defer close(eventChan)
	log := op.log.With().Str("address", addr).Str("output", tos.Name).Logger()

	runUntilError := func() error {
		epClient, err := CreateEventPublisherClient(addr)
		if err != nil {
			return maskAny(err)
		}
		defer epClient.CloseConnection()

		for {
			select {
			case evt := <-eventChan:
				// Publish event
				if _, err := epClient.Publish(ctx, &event.PublishRequest{Event: evt}); err != nil {
					return maskAny(err)
				}
				log.Debug().Str("id", evt.ID).Msg("published event")
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	for {
		if err := runUntilError(); ctx.Err() != nil {
			return ctx.Err()
		} else if err != nil {
			log.Error().Err(err).Msg("Error in publication loop")
		}
		select {
		case <-time.After(time.Second * 2):
			// Continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// limitDataLength returns the given data, trimmed to a reasonable length.
func limitDataLength(data string) string {
	maxLen := 128
	if len(data) > maxLen {
		return data[:maxLen] + "..."
	}
	return data
}
