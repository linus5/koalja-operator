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

package input

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/dchest/uniuri"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/AljabrIO/koalja-operator/pkg/agent/task"
	"github.com/AljabrIO/koalja-operator/pkg/annotatedvalue"
	"github.com/AljabrIO/koalja-operator/pkg/fs"
	"github.com/AljabrIO/koalja-operator/pkg/util"
)

type fileInputBuilder struct{}

func init() {
	task.RegisterExecutorInputBuilder(annotatedvalue.SchemeFile, fileInputBuilder{})
}

// Prepare input of a task for an input that uses a koalja-file scheme.
func (b fileInputBuilder) Build(ctx context.Context, cfg task.ExecutorInputBuilderConfig, deps task.ExecutorInputBuilderDependencies, target *task.ExecutorInputBuilderTarget) error {
	uri := cfg.AnnotatedValue.GetData()
	deps.Log.Debug().
		Int("sequence-index", cfg.AnnotatedValueIndex).
		Str("uri", uri).
		Str("input", cfg.InputSpec.Name).
		Str("task", cfg.TaskSpec.Name).
		Msg("Preparing file input value")

	// Prepare readonly volume for URI
	resp, err := deps.FileSystem.CreateVolumeForRead(ctx, &fs.CreateVolumeForReadRequest{
		URI:   uri,
		Owner: &cfg.OwnerRef,
	})
	if err != nil {
		return maskAny(err)
	}
	// TODO handle case where node is different
	if nodeName := resp.GetNodeName(); nodeName != "" {
		target.NodeName = &nodeName
	}

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
	pvcName := util.FixupKubernetesName(fmt.Sprintf("%s-%s-%s-%d-%s", cfg.Pipeline.GetName(), cfg.TaskSpec.Name, cfg.InputSpec.Name, cfg.AnnotatedValueIndex, uniuri.NewLen(6)))
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
	volName := util.FixupKubernetesName(fmt.Sprintf("input-%s-%d", cfg.InputSpec.Name, cfg.AnnotatedValueIndex))
	vol := corev1.Volume{
		Name: volName,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: pvcName,
				ReadOnly:  true,
			},
		},
	}
	target.Pod.Spec.Volumes = append(target.Pod.Spec.Volumes, vol)

	// Map volume in container fs namespace
	mountPath := filepath.Join("/koalja", "inputs", cfg.InputSpec.Name, strconv.Itoa(cfg.AnnotatedValueIndex))
	target.Container.VolumeMounts = append(target.Container.VolumeMounts, corev1.VolumeMount{
		Name:      volName,
		ReadOnly:  true,
		MountPath: mountPath,
	})

	// Create template data
	target.TemplateData = append(target.TemplateData, map[string]interface{}{
		"volumeName": resp.GetVolumeName(),
		"mountPath":  mountPath,
		"nodeName":   resp.GetNodeName(),
		"path":       filepath.Join(mountPath, resp.GetLocalPath()),
	})

	return nil
}