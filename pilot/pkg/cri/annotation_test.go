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
)

func TestAddIstioManagedAnnotation(t *testing.T) {
	var tests = []struct {
		annotations map[string]string
		want        map[string]string
	}{
		{
			nil,
			map[string]string{PodAnnotationKeyIstioManaged: ""},
		},
		{
			map[string]string{},
			map[string]string{PodAnnotationKeyIstioManaged: ""},
		},
		{
			map[string]string{
				"abc": "def",
			},
			map[string]string{
				"abc":                        "def",
				PodAnnotationKeyIstioManaged: "",
			},
		},
		{
			map[string]string{
				"abc":                        "def",
				PodAnnotationKeyIstioManaged: "",
			},
			map[string]string{
				"abc":                        "def",
				PodAnnotationKeyIstioManaged: "",
			},
		},
		{
			map[string]string{
				"abc":                        "def",
				PodAnnotationKeyIstioManaged: "something",
			},
			map[string]string{
				"abc":                        "def",
				PodAnnotationKeyIstioManaged: "",
			},
		},
	}

	for _, test := range tests {
		got := AddIstioManagedAnnotation(test.annotations)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("AddIstioManagedAnnotation(%#v) = %#v, want %#v", test.annotations, got,
				test.want)
		}
	}
}

func TestContainsIstioManagedAnnotation(t *testing.T) {
	var tests = []struct {
		annotations map[string]string
		want        bool
	}{
		{
			nil,
			false,
		},
		{
			map[string]string{},
			false,
		},
		{
			map[string]string{
				"abc": "def",
			},
			false,
		},
		{
			map[string]string{
				"abc":                        "def",
				PodAnnotationKeyIstioManaged: "",
			},
			true,
		},
		{
			map[string]string{
				"abc":                        "def",
				PodAnnotationKeyIstioManaged: "something",
			},
			true,
		},
	}

	for _, test := range tests {
		got := ContainsIstioManagedAnnotation(test.annotations)
		if got != test.want {
			t.Errorf("ContainsIstioManagedAnnotation(%#v) = %t, want %t", test.annotations, got,
				test.want)
		}
	}
}
