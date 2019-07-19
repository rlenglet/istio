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

// InMemoryLogger records the calls to LogCall.
type InMemoryLogger struct {
	Calls []Call
}

var _ Logger = &InMemoryLogger{}

type Call struct {
	Request  interface{}
	Response interface{}
	Err      error
}

func (l *InMemoryLogger) LogCall(request, response interface{}, err error) {
	l.Calls = append(l.Calls, Call{
		Request:  request,
		Response: response,
		Err:      err,
	})
}

// NopLogger is a logger that does nothing.
type NopLogger struct{}

var _ Logger = &NopLogger{}

func (l *NopLogger) LogCall(request, response interface{}, err error) {
}
