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

package kubeinject

import (
	"testing"

	k8s "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"istio.io/istio/pkg/kube/inject"
)

func TestValidateSidecarInjectionSpec(t *testing.T) {
	var tests = []struct {
		in        *inject.SidecarInjectionSpec
		wantError bool
	}{
		{
			&inject.SidecarInjectionSpec{
				InitContainers: []k8s.Container{
					{Name: ContainerNameIstioInit},
				},
				Containers: []k8s.Container{
					{Name: ContainerNameIstioProxy},
				},
				DNSConfig: &k8s.PodDNSConfig{
					Searches: []string{"example.com"},
				},
			},
			false,
		},
		{
			&inject.SidecarInjectionSpec{
				InitContainers: nil, // No init container.
				Containers: []k8s.Container{
					{Name: ContainerNameIstioProxy},
				},
			},
			true,
		},
		{
			&inject.SidecarInjectionSpec{
				InitContainers: []k8s.Container{
					{Name: ContainerNameIstioInit},
					{Name: "more-than-one-init-container"},
				},
				Containers: []k8s.Container{
					{Name: ContainerNameIstioProxy},
				},
			},
			true,
		},
		{
			&inject.SidecarInjectionSpec{
				InitContainers: []k8s.Container{
					{Name: "invalid-name"},
				},
				Containers: []k8s.Container{
					{Name: ContainerNameIstioProxy},
				},
			},
			true,
		},
		{
			&inject.SidecarInjectionSpec{
				InitContainers: []k8s.Container{
					{
						Name:  ContainerNameIstioInit,
						Ports: []k8s.ContainerPort{{}}, // Ports in init container.
					},
				},
				Containers: []k8s.Container{
					{Name: ContainerNameIstioProxy},
				},
			},
			true,
		},
		{
			&inject.SidecarInjectionSpec{
				InitContainers: []k8s.Container{ // Invalid container.
					{
						Name:    ContainerNameIstioInit,
						EnvFrom: []k8s.EnvFromSource{{}}, // EnvFrom is not supported.
					},
				},
				Containers: []k8s.Container{
					{Name: ContainerNameIstioProxy},
				},
			},
			true,
		},
		{
			&inject.SidecarInjectionSpec{
				InitContainers: []k8s.Container{
					{Name: ContainerNameIstioInit},
				},
				Containers: nil, // No container.
			},
			true,
		},
		{
			&inject.SidecarInjectionSpec{
				InitContainers: []k8s.Container{
					{Name: ContainerNameIstioInit},
				},
				Containers: []k8s.Container{
					{Name: ContainerNameIstioProxy},
					{Name: "more-than-one-container"},
				},
			},
			true,
		},
		{
			&inject.SidecarInjectionSpec{
				InitContainers: []k8s.Container{
					{Name: ContainerNameIstioInit},
				},
				Containers: []k8s.Container{
					{Name: "invalid-name"},
				},
			},
			true,
		},
		{
			&inject.SidecarInjectionSpec{
				InitContainers: []k8s.Container{
					{Name: ContainerNameIstioInit},
				},
				Containers: []k8s.Container{
					{
						Name:  ContainerNameIstioProxy,
						Ports: []k8s.ContainerPort{{}}, // Ports in container are OK.
					},
				},
			},
			false, // No error.
		},
		{
			&inject.SidecarInjectionSpec{
				InitContainers: []k8s.Container{
					{Name: ContainerNameIstioInit},
				},
				Containers: []k8s.Container{ // Invalid container.
					{
						Name:    ContainerNameIstioProxy,
						EnvFrom: []k8s.EnvFromSource{{}}, // EnvFrom is not supported.
					},
				},
			},
			true,
		},
		{
			&inject.SidecarInjectionSpec{
				InitContainers: []k8s.Container{
					{Name: ContainerNameIstioInit},
				},
				Containers: []k8s.Container{
					{Name: ContainerNameIstioProxy},
				},
				DNSConfig: &k8s.PodDNSConfig{ // Invalid DNS config.
					Nameservers: []string{"mydns.local"},
				},
			},
			true,
		},
	}

	for _, test := range tests {
		err := ValidateSidecarInjectionSpec(test.in)
		if test.wantError {
			if err == nil {
				t.Errorf("ValidateSidecarInjectionSpec(%#v) didn't fail", test.in)
			}
		} else {
			if err != nil {
				t.Errorf("ValidateSidecarInjectionSpec(%#v) failed: %v", test.in, err)
			}
		}
	}
}

func TestValidateContainer(t *testing.T) {
	var tests = []struct {
		in        *k8s.Container
		wantError bool
	}{
		{
			&k8s.Container{},
			false,
		},
		{
			&k8s.Container{
				Env: []k8s.EnvVar{
					{
						Name:  "TEST",
						Value: "abcd",
					},
				},
				Resources: k8s.ResourceRequirements{
					Limits: k8s.ResourceList{
						k8s.ResourceCPU:    resource.MustParse("2000m"),
						k8s.ResourceMemory: resource.MustParse("1024Mi"),
					},
					Requests: k8s.ResourceList{
						k8s.ResourceCPU:    resource.MustParse("100m"),
						k8s.ResourceMemory: resource.MustParse("128Mi"),
					},
				},
			},
			false,
		},
		{
			&k8s.Container{
				EnvFrom: []k8s.EnvFromSource{{}}, // EnvFrom is not supported.
			},
			true,
		},
		{
			&k8s.Container{
				Env: []k8s.EnvVar{
					{
						Name: "TEST",
						ValueFrom: &k8s.EnvVarSource{ // No FieldRef.
							ResourceFieldRef: &k8s.ResourceFieldSelector{},
						},
					},
				},
			},
			true,
		},
		{
			&k8s.Container{
				Env: []k8s.EnvVar{
					{
						Name:  "TEST",
						Value: "abcd",
					},
				},
				Resources: k8s.ResourceRequirements{
					Limits: k8s.ResourceList{
						k8s.ResourceStorage: resource.MustParse("1024Mi"),
					},
				},
			},
			true,
		},
		{
			&k8s.Container{
				Env: []k8s.EnvVar{
					{
						Name:  "TEST",
						Value: "abcd",
					},
				},
				Resources: k8s.ResourceRequirements{
					Requests: k8s.ResourceList{
						k8s.ResourceStorage: resource.MustParse("1024Mi"),
					},
				},
			},
			true,
		},
	}

	for _, test := range tests {
		err := ValidateContainer(test.in)
		if test.wantError {
			if err == nil {
				t.Errorf("ValidateContainer(%s) didn't fail", test.in)
			}
		} else {
			if err != nil {
				t.Errorf("ValidateContainer(%s) failed: %v", test.in, err)
			}
		}
	}
}

func TestValidateEnv(t *testing.T) {
	var tests = []struct {
		in        []k8s.EnvVar
		wantError bool
	}{
		{
			nil,
			false,
		},
		{
			[]k8s.EnvVar{},
			false,
		},
		{
			[]k8s.EnvVar{
				{
					Name:  "TEST",
					Value: "abcd",
				},
			},
			false,
		},
		{
			[]k8s.EnvVar{
				{
					Name: "TEST",
					// Empty Value.
				},
			},
			false,
		},
		{
			[]k8s.EnvVar{
				{
					Name: "TEST",
					ValueFrom: &k8s.EnvVarSource{
						FieldRef: &k8s.ObjectFieldSelector{
							FieldPath: "metadata.name",
						},
					},
				},
			},
			false,
		},
		{
			[]k8s.EnvVar{
				{
					Name:  "TEST",
					Value: "abcd", // Both Value and ValueFrom.
					ValueFrom: &k8s.EnvVarSource{
						FieldRef: &k8s.ObjectFieldSelector{
							FieldPath: "metadata.name",
						},
					},
				},
			},
			true,
		},
		{
			[]k8s.EnvVar{
				{
					Name: "TEST",
					ValueFrom: &k8s.EnvVarSource{ // No FieldRef.
						ResourceFieldRef: &k8s.ResourceFieldSelector{},
					},
				},
			},
			true,
		},
		{
			[]k8s.EnvVar{
				{
					Name: "TEST",
					ValueFrom: &k8s.EnvVarSource{
						FieldRef: &k8s.ObjectFieldSelector{
							FieldPath: "invalid.path", // Invalid path.
						},
					},
				},
			},
			true,
		},
	}

	for _, test := range tests {
		err := ValidateEnv(test.in)
		if test.wantError {
			if err == nil {
				t.Errorf("ValidateEnv(%s) didn't fail", test.in)
			}
		} else {
			if err != nil {
				t.Errorf("ValidateEnv(%s) failed: %v", test.in, err)
			}
		}
	}
}

func TestValidateResourceList(t *testing.T) {
	var tests = []struct {
		in        k8s.ResourceList
		wantError bool
	}{
		{
			nil,
			false,
		},
		{
			k8s.ResourceList{},
			false,
		},
		{
			k8s.ResourceList{
				k8s.ResourceCPU: resource.MustParse("2000m"),
			},
			false,
		},
		{
			k8s.ResourceList{
				k8s.ResourceMemory: resource.MustParse("1024Mi"),
			},
			false,
		},
		{
			k8s.ResourceList{
				k8s.ResourceCPU:    resource.MustParse("2000m"),
				k8s.ResourceMemory: resource.MustParse("1024Mi"),
			},
			false,
		},
		{
			k8s.ResourceList{
				k8s.ResourceCPU:     resource.MustParse("2000m"),
				k8s.ResourceStorage: resource.MustParse("1024Mi"),
			},
			true,
		},
		{
			k8s.ResourceList{
				k8s.ResourceStorage: resource.MustParse("1024Mi"),
			},
			true,
		},
	}

	for _, test := range tests {
		err := ValidateResourceList(test.in)
		if test.wantError {
			if err == nil {
				t.Errorf("ValidateResourceList(%#v) didn't fail", test.in)
			}
		} else {
			if err != nil {
				t.Errorf("ValidateResourceList(%#v) failed: %v", test.in, err)
			}
		}
	}
}

func TestValidatePodDNSConfig(t *testing.T) {
	var tests = []struct {
		in        *k8s.PodDNSConfig
		wantError bool
	}{
		{
			nil,
			false,
		},
		{
			&k8s.PodDNSConfig{},
			false,
		},
		{
			&k8s.PodDNSConfig{
				Searches: []string{"example.com"},
			},
			false,
		},
		{
			&k8s.PodDNSConfig{
				Nameservers: []string{"mydns.local"},
			},
			true,
		},
		{
			&k8s.PodDNSConfig{
				Options: []k8s.PodDNSConfigOption{{}},
			},
			true,
		},
	}

	for _, test := range tests {
		err := ValidatePodDNSConfig(test.in)
		if test.wantError {
			if err == nil {
				t.Errorf("ValidatePodDNSConfig(%s) didn't fail", test.in)
			}
		} else {
			if err != nil {
				t.Errorf("ValidatePodDNSConfig(%s) failed: %v", test.in, err)
			}
		}
	}
}
