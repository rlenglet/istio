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

import api "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"

// AppendPortMappings appends each of the given port mappings into the slice only if it doesn't
// already contain a matching port mapping.
func AppendPortMappings(mappings []*api.PortMapping, newMappings ...*api.PortMapping) []*api.PortMapping {
loop:
	for _, newMapping := range newMappings {
		for _, mapping := range mappings {
			if mapping.Protocol == newMapping.Protocol &&
				mapping.ContainerPort == newMapping.ContainerPort {
				// Existing port mapping matches. Skip adding it.
				continue loop
			}
		}
		mappings = append(mappings, newMapping)
	}
	return mappings
}
