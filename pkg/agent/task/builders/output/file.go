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

package output

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/dchest/uniuri"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/AljabrIO/koalja-operator/pkg/agent/task"
	"github.com/AljabrIO/koalja-operator/pkg/annotatedvalue"
	koalja "github.com/AljabrIO/koalja-operator/pkg/apis/koalja/v1alpha1"
	"github.com/AljabrIO/koalja-operator/pkg/fs"
	"github.com/AljabrIO/koalja-operator/pkg/util"
)

type fileOutputBuilder struct{}

const (
	maxPublishFailures = 10 // Try 10 publish attempts before giving up
)

func init() {
	task.RegisterExecutorOutputBuilder(koalja.ProtocolFile, fileOutputBuilder{})
}

// Prepare output of a task for an output of File protocol.
func (b fileOutputBuilder) Build(ctx context.Context, cfg task.ExecutorOutputBuilderConfig, deps task.ExecutorOutputBuilderDependencies, target *task.ExecutorOutputBuilderTarget) error {
	// Prepare volume for output
	var nodeName string
	if target.NodeName != nil {
		nodeName = *target.NodeName
	}
	resp, err := deps.FileSystem.CreateVolumeForWrite(ctx, &fs.CreateVolumeForWriteRequest{
		EstimatedCapacity: 0,
		NodeName:          nodeName,
		Owner:             &cfg.OwnerRef,
		Namespace:         cfg.Pipeline.GetNamespace(),
	})
	if err != nil {
		return maskAny(err)
	}

	// Mount PVC or HostPath, depending on result
	volName := util.FixupKubernetesName(fmt.Sprintf("output-%s", cfg.OutputSpec.Name))
	if resp.GetVolumeName() != "" {
		// Get created PersistentVolume
		var pv corev1.PersistentVolume
		pvKey := client.ObjectKey{
			Name: resp.GetVolumeName(),
		}
		if err := deps.Client.Get(ctx, pvKey, &pv); err != nil {
			deps.Log.Warn().Err(err).Msg("Failed to get PersistentVolume")
			return maskAny(err)
		}

		// Add PV to resources for deletion list (if needed)
		if resp.DeleteAfterUse {
			target.Resources = append(target.Resources, &pv)
		}

		// Create PVC
		pvcName := util.FixupKubernetesName(fmt.Sprintf("%s-%s-%s-%s", cfg.Pipeline.GetName(), cfg.TaskSpec.Name, cfg.OutputSpec.Name, uniuri.NewLen(6)))
		storageClassName := pv.Spec.StorageClassName
		pvc := corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:            pvcName,
				Namespace:       cfg.Pipeline.GetNamespace(),
				OwnerReferences: []metav1.OwnerReference{cfg.OwnerRef},
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: pv.Spec.AccessModes,
				VolumeName:  resp.GetVolumeName(),
				Resources: corev1.ResourceRequirements{
					Requests: pv.Spec.Capacity,
				},
				StorageClassName: &storageClassName,
			},
		}
		if err := deps.Client.Create(ctx, &pvc); err != nil {
			return maskAny(err)
		}
		target.Resources = append(target.Resources, &pvc)

		// Add volume for the pod
		vol := corev1.Volume{
			Name: volName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: pvcName,
					ReadOnly:  false,
				},
			},
		}
		target.Pod.Spec.Volumes = append(target.Pod.Spec.Volumes, vol)
	} else if resp.GetVolumeClaimName() != "" {
		// Add PVC to resources for deletion list (if needed)
		if resp.DeleteAfterUse {
			// Get created PersistentVolumeClaim
			var pvc corev1.PersistentVolumeClaim
			pvcKey := client.ObjectKey{
				Name:      resp.GetVolumeName(),
				Namespace: cfg.Pipeline.GetNamespace(),
			}
			if err := deps.Client.Get(ctx, pvcKey, &pvc); err != nil {
				deps.Log.Warn().Err(err).Msg("Failed to get PersistentVolumeClaim")
				return maskAny(err)
			}
			target.Resources = append(target.Resources, &pvc)
		}

		// Add volume for the pod, unless such a volume already exists
		if vol, found := util.GetVolumeWithForPVC(&target.Pod.Spec, resp.GetVolumeClaimName()); !found {
			vol := corev1.Volume{
				Name: volName,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: resp.GetVolumeClaimName(),
						ReadOnly:  false,
					},
				},
			}
			target.Pod.Spec.Volumes = append(target.Pod.Spec.Volumes, vol)
		} else {
			volName = vol.Name
			vol.PersistentVolumeClaim.ReadOnly = false
		}
	} else if resp.GetVolumePath() != "" {
		// Mount VolumePath as HostPath volume
		dirType := corev1.HostPathDirectoryOrCreate
		// Add volume for the pod
		vol := corev1.Volume{
			Name: volName,
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: resp.GetVolumePath(),
					Type: &dirType,
				},
			},
		}
		target.Pod.Spec.Volumes = append(target.Pod.Spec.Volumes, vol)
		// Ensure pod is schedule on node
		if nodeName := resp.GetNodeName(); nodeName != "" {
			if target.Pod.Spec.NodeName == "" {
				target.Pod.Spec.NodeName = nodeName
			} else if target.Pod.Spec.NodeName != nodeName {
				// Found conflict
				deps.Log.Error().
					Str("pod-nodeName", target.Pod.Spec.NodeName).
					Str("pod-nodeNameRequest", nodeName).
					Msg("Conflicting pod node spec")
				return maskAny(fmt.Errorf("Conflicting Node assignment"))
			}
		}
	} else {
		// No valid respond
		return maskAny(fmt.Errorf("FileSystem service return invalid response"))
	}

	// Map volume in container fs namespace
	mountPath := filepath.Join("/koalja", "outputs", cfg.OutputSpec.Name)
	target.Container.VolumeMounts = append(target.Container.VolumeMounts, corev1.VolumeMount{
		Name:      volName,
		ReadOnly:  false,
		MountPath: mountPath,
		SubPath:   resp.GetSubPath(),
	})

	// Create template data
	localPath := "output"
	target.TemplateData = map[string]interface{}{
		"volumeName":      resp.GetVolumeName(),
		"volumeClaimName": resp.GetVolumeClaimName(),
		"volumePath":      resp.GetVolumePath(),
		"subPath":         resp.GetSubPath(),
		"mountPath":       mountPath,
		"nodeName":        resp.GetNodeName(),
		"path":            filepath.Join(mountPath, localPath),
		"base":            filepath.Base(filepath.Join(mountPath, localPath)),
		"dir":             filepath.Dir(filepath.Join(mountPath, localPath)),
	}

	// Prepare output processor
	target.OutputProcessor = &fileOutputProcessor{
		ReadyOnSucceeded: cfg.OutputSpec.Ready.IsSucceeded(),
		OutputName:       cfg.OutputSpec.Name,
		VolumeName:       resp.GetVolumeName(),
		VolumeClaimName:  resp.GetVolumeClaimName(),
		VolumePath:       resp.GetVolumePath(),
		SubPath:          resp.GetSubPath(),
		MountPath:        mountPath,
		NodeName:         resp.GetNodeName(),
		LocalPath:        localPath,
	}

	return nil
}

type fileOutputProcessor struct {
	ReadyOnSucceeded bool
	OutputName       string
	VolumeName       string
	VolumeClaimName  string
	VolumePath       string
	SubPath          string
	MountPath        string
	NodeName         string
	LocalPath        string
}

// Gets the name of the output this processor is intended for
func (p *fileOutputProcessor) GetOutputName() string {
	return p.OutputName
}

// Process a single file output
func (p *fileOutputProcessor) Process(ctx context.Context, cfg task.ExecutorOutputProcessorConfig, deps task.ExecutorOutputProcessorDependencies) error {
	if !p.ReadyOnSucceeded {
		// Nothing todo here
		return nil
	}
	log := deps.Log
	log.Debug().Msg("creating URI for output")
	resp, err := deps.FileSystem.CreateFileURI(ctx, &fs.CreateFileURIRequest{
		Scheme:          string(annotatedvalue.SchemeFile),
		VolumeName:      p.VolumeName,
		VolumeClaimName: p.VolumeClaimName,
		VolumePath:      p.VolumePath,
		SubPath:         p.SubPath,
		NodeName:        p.NodeName,
		LocalPath:       p.LocalPath,
		IsDir:           false, // TODO
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create file URI")
		return maskAny(err)
	}
	// Publish annotated value
	av := annotatedvalue.AnnotatedValue{Data: resp.GetURI()}
	delay := time.Millisecond * 100
	recentFailures := 0
	for {
		accepted, publishedAv, err := deps.OutputPublisher.Publish(ctx, p.OutputName, av, cfg.Snapshot)
		if err != nil {
			recentFailures++
			if recentFailures > maxPublishFailures {
				log.Error().Err(err).Msg("Publish failed too many times")
				return maskAny(err)
			}
			log.Debug().Err(err).Msg("Publish attempt failed")
		} else if accepted {
			// Output was accepted
			log.Debug().
				Str("annotatedvalue-id", publishedAv.GetID()).
				Msg("published annotated value")
			return nil
		} else {
			// Publish call succeeded, but output was not (yet) accepted
			recentFailures = 0
		}
		// Publication was not accepted, or call failed, try again soon.
		select {
		case <-time.After(delay):
			// Try again
			delay = util.Backoff(delay, 1.5, time.Minute)
		case <-ctx.Done():
			// Context canceled
			log.Debug().Err(ctx.Err()).Msg("Context canceled during Process")
			return ctx.Err()
		}
	}
}

// CreateFileURI creates a URI for the given file/dir
func (p *fileOutputProcessor) CreateFileURI(ctx context.Context, localPath string, isDir bool, cfg task.ExecutorOutputProcessorConfig, deps task.ExecutorOutputProcessorDependencies) (string, error) {
	log := deps.Log.With().
		Str("localPath", localPath).
		Bool("isDir", isDir).
		Logger()
	log.Debug().Msg("creating URI for file")
	resp, err := deps.FileSystem.CreateFileURI(ctx, &fs.CreateFileURIRequest{
		Scheme:          string(annotatedvalue.SchemeFile),
		VolumeName:      p.VolumeName,
		VolumeClaimName: p.VolumeClaimName,
		VolumePath:      p.VolumePath,
		SubPath:         p.SubPath,
		NodeName:        p.NodeName,
		LocalPath:       localPath,
		IsDir:           isDir,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create file URI")
		return "", maskAny(err)
	}
	return resp.GetURI(), nil
}
