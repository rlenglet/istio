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

import (
	"reflect"
	"testing"

	api "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

func TestAppendPortMappings(t *testing.T) {
	var tests = []struct {
		mappings    []*api.PortMapping
		newMappings []*api.PortMapping
		want        []*api.PortMapping
	}{
		{
			nil,
			nil,
			nil,
		},
		{
			nil,
			[]*api.PortMapping{
				{Protocol: api.Protocol_TCP, ContainerPort: 80},
			},
			[]*api.PortMapping{
				{Protocol: api.Protocol_TCP, ContainerPort: 80},
			},
		},
		{
			[]*api.PortMapping{
				{Protocol: api.Protocol_TCP, ContainerPort: 80},
			},
			nil,
			[]*api.PortMapping{
				{Protocol: api.Protocol_TCP, ContainerPort: 80},
			},
		},
		{
			[]*api.PortMapping{
				{Protocol: api.Protocol_TCP, ContainerPort: 80},
			},
			[]*api.PortMapping{
				{Protocol: api.Protocol_UDP, ContainerPort: 80},
			},
			[]*api.PortMapping{
				{Protocol: api.Protocol_TCP, ContainerPort: 80},
				{Protocol: api.Protocol_UDP, ContainerPort: 80},
			},
		},
		{
			[]*api.PortMapping{
				{Protocol: api.Protocol_TCP, ContainerPort: 80},
			},
			[]*api.PortMapping{
				{Protocol: api.Protocol_TCP, ContainerPort: 443},
			},
			[]*api.PortMapping{
				{Protocol: api.Protocol_TCP, ContainerPort: 80},
				{Protocol: api.Protocol_TCP, ContainerPort: 443},
			},
		},
		{
			[]*api.PortMapping{
				{Protocol: api.Protocol_TCP, ContainerPort: 80},
				{Protocol: api.Protocol_TCP, ContainerPort: 443},
			},
			[]*api.PortMapping{
				{Protocol: api.Protocol_TCP, ContainerPort: 22},
				{Protocol: api.Protocol_TCP, ContainerPort: 443}, // Duplicate.
			},
			[]*api.PortMapping{
				{Protocol: api.Protocol_TCP, ContainerPort: 80},
				{Protocol: api.Protocol_TCP, ContainerPort: 443},
				{Protocol: api.Protocol_TCP, ContainerPort: 22},
			},
		},
		{
			[]*api.PortMapping{
				{Protocol: api.Protocol_TCP, ContainerPort: 80},
				{Protocol: api.Protocol_TCP, ContainerPort: 443},
			},
			[]*api.PortMapping{
				{Protocol: api.Protocol_TCP, ContainerPort: 22},
				{Protocol: api.Protocol_TCP, ContainerPort: 22}, // Duplicate.
				{Protocol: api.Protocol_TCP, ContainerPort: 22}, // Duplicate.
			},
			[]*api.PortMapping{
				{Protocol: api.Protocol_TCP, ContainerPort: 80},
				{Protocol: api.Protocol_TCP, ContainerPort: 443},
				{Protocol: api.Protocol_TCP, ContainerPort: 22},
			},
		},
	}

	for _, test := range tests {
		got := AppendPortMappings(test.mappings, test.newMappings...)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("AppendPortMappings(%s, %s...) = %s, want %s", test.mappings,
				test.newMappings, got, test.want)
		}
	}
}
