/*
Copyright 2018 Aljabr Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package stub

import (
	"context"
	"sync"

	"github.com/AljabrIO/koalja-operator/pkg/constants"

	"github.com/golang/protobuf/ptypes"

	"github.com/AljabrIO/koalja-operator/pkg/agent/pipeline"
	"github.com/AljabrIO/koalja-operator/pkg/event"
	"github.com/AljabrIO/koalja-operator/pkg/event/registry"
	"github.com/AljabrIO/koalja-operator/pkg/event/tree"
	"github.com/AljabrIO/koalja-operator/pkg/util/retry"
	"github.com/rs/zerolog"
)

// outputStore is an in-memory implementation of an event queue.
type outputStore struct {
	log         zerolog.Logger
	registry    registry.EventRegistryClient
	events      []*tree.EventTree
	eventsMutex sync.Mutex
}

// newOutputStore creates a new output store
func newOutputStore(log zerolog.Logger, r registry.EventRegistryClient) *outputStore {
	return &outputStore{
		log:      log,
		registry: r,
	}
}

// Publish an event
func (s *outputStore) Publish(ctx context.Context, req *event.PublishRequest) (*event.PublishResponse, error) {
	s.log.Debug().Interface("event", req.Event).Msg("Publish request")
	// Try to record event in registry
	e := *req.GetEvent()
	e.Link = "" // the end
	if err := retry.Do(ctx, func(ctx context.Context) error {
		s.log.Debug().Msg("RecordEvent attempt start")
		if _, err := s.registry.RecordEvent(ctx, &e); err != nil {
			s.log.Debug().Err(err).Msg("RecordEvent attempt failed")
			return maskAny(err)
		}
		return nil
	}, retry.Timeout(constants.TimeoutRecordEvent)); err != nil {
		s.log.Error().Err(err).Msg("Failed to record event")
		return nil, maskAny(err)
	}

	// Build event tree
	evtTree, err := tree.Build(ctx, e, s.registry)
	if err != nil {
		return nil, maskAny(err)
	}

	// Now put event in in-memory list
	s.eventsMutex.Lock()
	defer s.eventsMutex.Unlock()
	s.events = append(s.events, evtTree)
	return &event.PublishResponse{}, nil
}

// GetOutputEvents returns all events (resulting from task outputs that
// are not connected to inputs of other tasks) that match the given filter.
func (s *outputStore) GetOutputEvents(ctx context.Context, req *pipeline.OutputEventsRequest) (*pipeline.OutputEvents, error) {
	s.log.Debug().Interface("req", req).Msg("GetOutputEvents request")
	s.eventsMutex.Lock()
	defer s.eventsMutex.Unlock()

	resp := &pipeline.OutputEvents{}
	for _, tree := range s.events {
		if isMatch(tree, req) {
			evt := tree.Event // Create clone
			resp.Events = append(resp.Events, &evt)
		}
	}

	return resp, nil
}

// isMatch returns true when the given event tree matches the given request.
func isMatch(e *tree.EventTree, req *pipeline.OutputEventsRequest) bool {
	createdAt, _ := ptypes.Timestamp(e.Event.GetCreatedAt())
	if tsPB := req.GetCreatedAfter(); tsPB != nil {
		ts, _ := ptypes.Timestamp(tsPB)
		if createdAt.Before(ts) {
			return false
		}
	}
	if tsPB := req.GetCreatedBefore(); tsPB != nil {
		ts, _ := ptypes.Timestamp(tsPB)
		if createdAt.After(ts) {
			return false
		}
	}
	if taskNames := req.GetTaskNames(); len(taskNames) > 0 {
		found := false
		for _, name := range taskNames {
			if e.Event.GetSourceTask() == name {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	if ids := req.GetEventIDs(); len(ids) > 0 {
		found := false
		for _, id := range ids {
			if e.ContainsID(id) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
