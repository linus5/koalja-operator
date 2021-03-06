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

package v1alpha1

import (
	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
)

// TaskSpec holds the specification of a single task
type TaskSpec struct {
	// Name of the task
	Name string `json:"name" protobuf:"bytes,1,req,name=name"`
	// Type of task
	Type TaskType `json:"type,omitempty" protobuf:"bytes,2,opt,name=type"`
	// Inputs of the task
	Inputs []TaskInputSpec `json:"inputs,omitempty" protobuf:"bytes,3,rep,name=inputs"`
	// Outputs of the task
	Outputs []TaskOutputSpec `json:"outputs,omitempty" protobuf:"bytes,4,rep,name=outputs"`
	// Executor holds the spec of the execution part of the task
	Executor *v1.Container `json:"executor,omitempty" protobuf:"bytes,5,opt,name=executor"`
	// Policy determining how inputs are aggregated into snapshots
	SnapshotPolicy SnapshotPolicy `json:"snapshotPolicy,omitempty" protobuf:"bytes,6,opt,name=snapshotPolicy"`
	// Policy determining when the executor will be launched
	LaunchPolicy LaunchPolicy `json:"launchPolicy,omitempty" protobuf:"bytes,7,opt,name=launchPolicy"`
	// Service is set when the task executor must act as a network service.
	Service *ServiceSpec `json:"service,omitempty" protobuf:"bytes,8,opt,name=service"`
}

// TaskType identifies a well know type of task.
type TaskType string

// LaunchPolicy determines when the task agent will launch an executor.
type LaunchPolicy string

const (
	// LaunchPolicyAuto indicates that the executor will be launched
	// every time a new snapshot is available.
	// This is the default policy for tasks that have 1 or more inputs.
	LaunchPolicyAuto LaunchPolicy = "Auto"
	// LaunchPolicyCustom indicates that the executor will be launched
	// as soon as the task agent starts and that it will use the SnapshotService
	// provided by the task agent to query for the next available snapshot.
	// This is the default policy for tasks that have no inputs.
	LaunchPolicyCustom LaunchPolicy = "Custom"
	// LaunchPolicyRestart indicates that the executor will be launched
	// every time a new snapshot is available.
	// If a previous executor was still running, that will be stopped
	// first.
	LaunchPolicyRestart LaunchPolicy = "Restart"
)

// Validate that the given policy is a valid value.
func (lp LaunchPolicy) Validate() error {
	switch lp {
	case LaunchPolicyAuto, LaunchPolicyCustom, LaunchPolicyRestart, "":
		return nil
	default:
		return errors.Wrapf(ErrValidation, "Invalid LaunchPolicy '%s'", string(lp))
	}
}

// SnapshotPolicy determines how to sample annotated values into a tuple
// that is the input for task execution.
type SnapshotPolicy string

const (
	// SnapshotPolicySwapNew4Old policy replaces old values in the snapshot with new
	// values as soon as they come in.
	SnapshotPolicySwapNew4Old SnapshotPolicy = "SwapNew4Old"
	// SnapshotPolicyAllNew policy flushes all values when the task has been executed
	// on a snapshot and starts filling all inputs from scratch.
	SnapshotPolicyAllNew SnapshotPolicy = "AllNew"
	// SnapshotPolicySlidingWindow policy flushes a specifies amount of values out
	// of the snapshot when the task has been executed.
	SnapshotPolicySlidingWindow SnapshotPolicy = "SlidingWindow"
)

// IsSwapNew4Old returns true when the given policy is "SwapNew4Old"
func (sp SnapshotPolicy) IsSwapNew4Old() bool { return sp == SnapshotPolicySwapNew4Old }

// IsAllNew returns true when the given policy is "AllNew"
func (sp SnapshotPolicy) IsAllNew() bool { return sp == SnapshotPolicyAllNew || sp == "" }

// IsSlidingWindow returns true when the given policy is "SlidingWindow"
func (sp SnapshotPolicy) IsSlidingWindow() bool { return sp == SnapshotPolicySlidingWindow }

// Validate that the given readiness is a valid value.
func (sp SnapshotPolicy) Validate() error {
	switch sp {
	case SnapshotPolicySwapNew4Old, SnapshotPolicyAllNew, SnapshotPolicySlidingWindow, "":
		return nil
	default:
		return errors.Wrapf(ErrValidation, "Invalid SnapshotPolicy '%s'", string(sp))
	}
}

// InputByName returns the input of the task that has the given name.
// Returns false if not found.
func (ts TaskSpec) InputByName(name string) (TaskInputSpec, bool) {
	for _, x := range ts.Inputs {
		if x.Name == name {
			return x, true
		}
	}
	return TaskInputSpec{}, false
}

// SnapshotInputs returns a list of inputs that we use in the building of snapshots (leave out ones with MergeInto)
func (ts TaskSpec) SnapshotInputs() []TaskInputSpec {
	snapshotInputs := make([]TaskInputSpec, 0, len(ts.Inputs))
	for _, i := range ts.Inputs {
		if !i.HasMergeInto() {
			snapshotInputs = append(snapshotInputs, i)
		}
	}
	return snapshotInputs
}

// OutputByName returns the output of the task that has the given name.
// Returns false if not found.
func (ts TaskSpec) OutputByName(name string) (TaskOutputSpec, bool) {
	for _, x := range ts.Outputs {
		if x.Name == name {
			return x, true
		}
	}
	return TaskOutputSpec{}, false
}

// HasLaunchPolicyAuto returns true when the LaunchPolicy of the task
// has been set to Auto or no LaunchPolicy has been set and the
// task has one or more inputs.
func (ts TaskSpec) HasLaunchPolicyAuto() bool {
	if ts.LaunchPolicy == "" {
		return len(ts.Inputs) > 0
	}
	return ts.LaunchPolicy == LaunchPolicyAuto
}

// HasLaunchPolicyCustom returns true when the LaunchPolicy of the task
// has been set to Custom or no LaunchPolicy has been set and the
// task has no inputs.
func (ts TaskSpec) HasLaunchPolicyCustom() bool {
	if ts.LaunchPolicy == "" {
		return len(ts.Inputs) == 0
	}
	return ts.LaunchPolicy == LaunchPolicyCustom
}

// HasLaunchPolicyRestart returns true when the LaunchPolicy of the task
// has been set to Restart.
func (ts TaskSpec) HasLaunchPolicyRestart() bool {
	return ts.LaunchPolicy == LaunchPolicyRestart
}

// Validate the task in the context of the given pipeline spec.
// Return an error when an issue is found, nil when all ok.
func (ts TaskSpec) Validate(ps PipelineSpec) error {
	if err := ValidateName(ts.Name); err != nil {
		return maskAny(err)
	}
	if ts.Executor == nil && ts.Type == "" {
		return errors.Wrapf(ErrValidation, "Executor or Type expected in task '%s'", ts.Name)
	}
	if len(ts.Inputs) == 0 && len(ts.Outputs) == 0 {
		return errors.Wrapf(ErrValidation, "Task '%s' must have at least 1 input or 1 output", ts.Name)
	}
	if err := ts.SnapshotPolicy.Validate(); err != nil {
		return maskAny(err)
	}
	if err := ts.LaunchPolicy.Validate(); err != nil {
		return maskAny(err)
	}
	names := make(map[string]struct{})
	for _, x := range ts.Inputs {
		if _, found := names[x.Name]; found {
			return errors.Wrapf(ErrValidation, "Duplicate input name '%s' in task '%s'", x.Name, ts.Name)
		}
		names[x.Name] = struct{}{}
		if err := x.Validate(ps, ts); err != nil {
			return maskAny(err)
		}
	}
	names = make(map[string]struct{})
	for _, x := range ts.Outputs {
		if _, found := names[x.Name]; found {
			return errors.Wrapf(ErrValidation, "Duplicate output name '%s' in task '%s'", x.Name, ts.Name)
		}
		names[x.Name] = struct{}{}
		if err := x.Validate(ps); err != nil {
			return maskAny(err)
		}
	}
	if ts.Executor != nil {
		if ts.Executor.Image == "" {
			return errors.Wrapf(ErrValidation, "Executor of task '%s' must have an image", ts.Name)
		}
	}
	if ts.Service != nil {
		if err := ts.Service.Validate(); err != nil {
			return maskAny(err)
		}
	}
	return nil
}
