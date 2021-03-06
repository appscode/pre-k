/*
Copyright The Kmodules Authors.

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

package v1

import (
	"fmt"

	"github.com/appscode/go/types"
	"github.com/imdario/mergo"
	jsoniter "github.com/json-iterator/go"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var json = jsoniter.ConfigFastest

func AddFinalizer(m metav1.ObjectMeta, finalizer string) metav1.ObjectMeta {
	for _, name := range m.Finalizers {
		if name == finalizer {
			return m
		}
	}
	m.Finalizers = append(m.Finalizers, finalizer)
	return m
}

func HasFinalizer(m metav1.ObjectMeta, finalizer string) bool {
	for _, name := range m.Finalizers {
		if name == finalizer {
			return true
		}
	}
	return false
}

func RemoveFinalizer(m metav1.ObjectMeta, finalizer string) metav1.ObjectMeta {
	// https://github.com/golang/go/wiki/SliceTricks#filtering-without-allocating
	r := m.Finalizers[:0]
	for _, name := range m.Finalizers {
		if name != finalizer {
			r = append(r, name)
		}
	}
	m.Finalizers = r
	return m
}

func EnsureContainerDeleted(containers []core.Container, name string) []core.Container {
	for i, c := range containers {
		if c.Name == name {
			return append(containers[:i], containers[i+1:]...)
		}
	}
	return containers
}

func UpsertContainer(containers []core.Container, upsert core.Container) []core.Container {
	for i, container := range containers {
		if container.Name == upsert.Name {
			err := mergo.MergeWithOverwrite(&container, upsert)
			if err != nil {
				panic(err)
			}
			containers[i] = container
			return containers
		}
	}
	return append(containers, upsert)
}

func UpsertContainers(containers []core.Container, addons []core.Container) []core.Container {
	var out = containers
	for _, c := range addons {
		out = UpsertContainer(out, c)
	}
	return out
}

func UpsertVolume(volumes []core.Volume, nv ...core.Volume) []core.Volume {
	upsert := func(v core.Volume) {
		for i, vol := range volumes {
			if vol.Name == v.Name {
				volumes[i] = v
				return
			}
		}
		volumes = append(volumes, v)
	}

	for _, volume := range nv {
		upsert(volume)
	}
	return volumes

}

func UpsertVolumeClaim(volumeClaims []core.PersistentVolumeClaim, upsert core.PersistentVolumeClaim) []core.PersistentVolumeClaim {
	for i, vc := range volumeClaims {
		if vc.Name == upsert.Name {
			volumeClaims[i] = upsert
			return volumeClaims
		}
	}
	return append(volumeClaims, upsert)
}

func EnsureVolumeDeleted(volumes []core.Volume, name string) []core.Volume {
	for i, v := range volumes {
		if v.Name == name {
			return append(volumes[:i], volumes[i+1:]...)
		}
	}
	return volumes
}

func UpsertVolumeMount(mounts []core.VolumeMount, nv ...core.VolumeMount) []core.VolumeMount {
	upsert := func(m core.VolumeMount) {
		for i, vol := range mounts {
			if vol.Name == m.Name {
				mounts[i] = m
				return
			}
		}
		mounts = append(mounts, m)
	}

	for _, mount := range nv {
		upsert(mount)
	}
	return mounts
}

func EnsureVolumeMountDeleted(mounts []core.VolumeMount, name string) []core.VolumeMount {
	for i, v := range mounts {
		if v.Name == name {
			return append(mounts[:i], mounts[i+1:]...)
		}
	}
	return mounts
}

func UpsertVolumeMountByPath(mounts []core.VolumeMount, nv core.VolumeMount) []core.VolumeMount {
	for i, vol := range mounts {
		if vol.MountPath == nv.MountPath {
			mounts[i] = nv
			return mounts
		}
	}
	return append(mounts, nv)
}

func EnsureVolumeMountDeletedByPath(mounts []core.VolumeMount, mountPath string) []core.VolumeMount {
	for i, v := range mounts {
		if v.MountPath == mountPath {
			return append(mounts[:i], mounts[i+1:]...)
		}
	}
	return mounts
}

func UpsertEnvVars(vars []core.EnvVar, nv ...core.EnvVar) []core.EnvVar {
	upsert := func(env core.EnvVar) {
		for i, v := range vars {
			if v.Name == env.Name {
				vars[i] = env
				return
			}
		}
		vars = append(vars, env)
	}

	for _, env := range nv {
		upsert(env)
	}
	return vars
}

func EnsureEnvVarDeleted(vars []core.EnvVar, name string) []core.EnvVar {
	for i, v := range vars {
		if v.Name == name {
			return append(vars[:i], vars[i+1:]...)
		}
	}
	return vars
}

func UpsertMap(maps, upsert map[string]string) map[string]string {
	if maps == nil {
		maps = make(map[string]string)
	}
	for k, v := range upsert {
		maps[k] = v
	}
	return maps
}

func MergeLocalObjectReferences(l1, l2 []core.LocalObjectReference) []core.LocalObjectReference {
	result := make([]core.LocalObjectReference, 0, len(l1)+len(l2))
	m := make(map[string]core.LocalObjectReference)
	for _, ref := range l1 {
		m[ref.Name] = ref
		result = append(result, ref)
	}
	for _, ref := range l2 {
		if _, found := m[ref.Name]; !found {
			result = append(result, ref)
		}
	}
	return result
}

func EnsureOwnerReference(meta metav1.Object, owner *core.ObjectReference) {
	if owner == nil ||
		owner.APIVersion == "" ||
		owner.Kind == "" ||
		owner.Name == "" ||
		owner.UID == "" {
		return
	}
	if meta.GetNamespace() != owner.Namespace {
		panic(fmt.Errorf("owner %s %s must be from the same namespace as object %s", owner.Kind, owner.Name, meta.GetName()))
	}

	ownerRefs := meta.GetOwnerReferences()

	fi := -1
	for i, ref := range ownerRefs {
		if ref.Kind == owner.Kind && ref.Name == owner.Name {
			fi = i
			break
		}
	}
	if fi == -1 {
		ownerRefs = append(ownerRefs, metav1.OwnerReference{})
		fi = len(ownerRefs) - 1
	}
	ownerRefs[fi].APIVersion = owner.APIVersion
	ownerRefs[fi].Kind = owner.Kind
	ownerRefs[fi].Name = owner.Name
	ownerRefs[fi].UID = owner.UID
	if ownerRefs[fi].BlockOwnerDeletion == nil {
		ownerRefs[fi].BlockOwnerDeletion = types.FalseP()
	}

	meta.SetOwnerReferences(ownerRefs)
}

func RemoveOwnerReference(meta metav1.Object, owner *core.ObjectReference) {
	ownerRefs := meta.GetOwnerReferences()
	for i, ref := range ownerRefs {
		if ref.Kind == owner.Kind && ref.Name == owner.Name {
			ownerRefs = append(ownerRefs[:i], ownerRefs[i+1:]...)
			break
		}
	}
	meta.SetOwnerReferences(ownerRefs)
}

func IsOwnedByRef(o runtime.Object, owner *core.ObjectReference) bool {
	obj, err := meta.Accessor(o)
	if err != nil {
		return false
	}

	return o.GetObjectKind().GroupVersionKind() == owner.GroupVersionKind() &&
		obj.GetName() == owner.Name &&
		(string(owner.UID) == "" || obj.GetUID() == owner.UID)
}

func IsOwnedBy(o1 runtime.Object, o2 runtime.Object) bool {
	obj, err := meta.Accessor(o1)
	if err != nil {
		return false
	}
	owner, err := meta.Accessor(o2)
	if err != nil {
		return false
	}
	return o1.GetObjectKind().GroupVersionKind() == o2.GetObjectKind().GroupVersionKind() &&
		obj.GetName() == owner.GetName() &&
		(string(owner.GetUID()) == "" || obj.GetUID() == owner.GetUID())
}
