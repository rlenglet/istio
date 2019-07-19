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
	k8s "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"istio.io/istio/pkg/kube/inject"
)

var (
	RootUserID   int64 = 0
	ProxyUserID  int64 = 1337
	ProxyGroupID int64 = 1337
	True               = true
	False              = false

	TestSpec = inject.SidecarInjectionSpec{
		PodRedirectAnnot:    nil,
		RewriteAppHTTPProbe: false,
		InitContainers: []k8s.Container{
			{
				Name:  "istio-init",
				Image: "docker.io/istio/proxy_init:1.2.2",
				Args: []string{
					"-p",
					"15001",
					"-u",
					"1337",
					"-m",
					"REDIRECT",
					"-i",
					"*",
					"-x",
					"",
					"-b",
					"5000",
					"-d",
					"15020",
				},
				Resources: k8s.ResourceRequirements{
					Limits: k8s.ResourceList{
						"cpu":    resource.MustParse("100m"),
						"memory": resource.MustParse("50Mi"),
					},
					Requests: k8s.ResourceList{
						"cpu":    resource.MustParse("10m"),
						"memory": resource.MustParse("10Mi"),
					},
				},
				SecurityContext: &k8s.SecurityContext{
					Capabilities: &k8s.Capabilities{
						Add: []k8s.Capability{"NET_ADMIN"},
					},
					RunAsUser:              &RootUserID,
					RunAsNonRoot:           &False,
				},
			},
		},
		Containers: []k8s.Container{
			{
				Name:  "istio-proxy",
				Image: "docker.io/istio/proxyv2:1.2.2",
				Args: []string{
					"proxy",
					"sidecar",
					"--domain",
					"default.svc.cluster.local",
					"--configPath",
					"/etc/istio/proxy",
					"--binaryPath",
					"/usr/local/bin/envoy",
					"--serviceCluster",
					"helloworld.default",
					"--drainDuration",
					"45s",
					"--parentShutdownDuration",
					"1m0s",
					"--discoveryAddress",
					"istio-pilot.istio-system:15010",
					"--zipkinAddress",
					"zipkin.istio-system:9411",
					"--dnsRefreshRate",
					"300s",
					"--connectTimeout",
					"10s",
					"--proxyAdminPort",
					"15000",
					"--concurrency",
					"2",
					"--controlPlaneAuthPolicy",
					"NONE",
					"--statusPort",
					"15020",
					"--applicationPorts",
					"5000",
				},
				Env: []k8s.EnvVar{
					{
						Name: "POD_NAME",
						ValueFrom: &k8s.EnvVarSource{
							FieldRef: &k8s.ObjectFieldSelector{
								FieldPath: "metadata.name",
							},
						},
					},
					{
						Name: "POD_NAMESPACE",
						ValueFrom: &k8s.EnvVarSource{
							FieldRef: &k8s.ObjectFieldSelector{
								FieldPath: "metadata.namespace",
							},
						},
					},
					{
						Name: "INSTANCE_IP",
						ValueFrom: &k8s.EnvVarSource{
							FieldRef: &k8s.ObjectFieldSelector{
								FieldPath: "status.podIP",
							},
						},
					},
					{
						Name: "ISTIO_META_POD_NAME",
						ValueFrom: &k8s.EnvVarSource{
							FieldRef: &k8s.ObjectFieldSelector{
								FieldPath: "metadata.name",
							},
						},
					},
					{
						Name: "ISTIO_META_CONFIG_NAMESPACE",
						ValueFrom: &k8s.EnvVarSource{
							FieldRef: &k8s.ObjectFieldSelector{
								FieldPath: "metadata.namespace",
							},
						},
					},
					{
						Name:  "ISTIO_META_INTERCEPTION_MODE",
						Value: "REDIRECT",
					},
					{
						Name:  "ISTIO_META_INCLUDE_INBOUND_PORTS",
						Value: "5000",
					},
					{
						Name:  "ISTIO_METAJSON_LABELS",
						Value: "{\"app\":\"helloworld\"}\n",
					},
				},
				VolumeMounts: nil, // TODO(rlenglet)
				Resources: k8s.ResourceRequirements{
					Limits: k8s.ResourceList{
						"cpu":    resource.MustParse("2000m"),
						"memory": resource.MustParse("1024Mi"),
					},
					Requests: k8s.ResourceList{
						"cpu":    resource.MustParse("100m"),
						"memory": resource.MustParse("128Mi"),
					},
				},
				SecurityContext: &k8s.SecurityContext{
					Capabilities: &k8s.Capabilities{
						Add: []k8s.Capability{"NET_ADMIN"},
					},
					RunAsUser:              &ProxyUserID,
					RunAsGroup:             &ProxyGroupID,
					ReadOnlyRootFilesystem: &True,
				},
			},
		},
		Volumes:          nil,
		DNSConfig:        nil,
		ImagePullSecrets: nil,
	}
)
