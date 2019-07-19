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
	"errors"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func TestWithField(t *testing.T) {
	newField := zap.String("z", "zz")

	var tests = []struct {
		logFields []zap.Field
		want      []zap.Field
	}{
		{
			nil,
			[]zap.Field{newField},
		},
		{
			[]zap.Field{},
			[]zap.Field{newField},
		},
		{
			[]zap.Field{zap.String("a", "aa")},
			[]zap.Field{zap.String("a", "aa"), newField},
		},
	}

	for _, test := range tests {
		got := WithField(newField, test.logFields...)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("WithField(%#v, %#v...) = %#v, want %#v", newField, test.logFields, got,
				test.want)
		}
	}
}

func TestWithError(t *testing.T) {
	err := errors.New("test error")

	var tests = []struct {
		logFields []zap.Field
		want      []zap.Field
	}{
		{
			nil,
			[]zap.Field{zap.Error(err)},
		},
		{
			[]zap.Field{},
			[]zap.Field{zap.Error(err)},
		},
		{
			[]zap.Field{zap.String("a", "aa")},
			[]zap.Field{zap.String("a", "aa"), zap.Error(err)},
		},
	}

	for _, test := range tests {
		got := WithError(err, test.logFields...)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("WithError(%v, %#v...) = %#v, want %#v", err, test.logFields, got, test.want)
		}
	}
}
