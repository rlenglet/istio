// Copyright 2019 Istio Authors
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

package cri

const (
	// PodAnnotationKeyIstioManaged is the annotation set by a CRI proxy on any container it
	// manages directly.
	PodAnnotationKeyIstioManaged = "__istio__managed__"
)

// AddIstioManagedAnnotation adds into the given annotations an annotation which identifies the
// container as managed by a CRI proxy directly.
func AddIstioManagedAnnotation(annotations map[string]string) map[string]string {
	if annotations == nil {
		annotations = map[string]string{PodAnnotationKeyIstioManaged: ""}
	} else {
		annotations[PodAnnotationKeyIstioManaged] = ""
	}
	return annotations
}

// ContainsIstioManagedAnnotation return whether the given annotations include an annotation which
// identifies the container as managed by a CRI proxy directly.
func ContainsIstioManagedAnnotation(annotations map[string]string) bool {
	_, found := annotations[PodAnnotationKeyIstioManaged]
	return found
}
