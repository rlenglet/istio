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

package main

import (
	"context"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"istio.io/pkg/collateral"
	"istio.io/pkg/env"
	"istio.io/pkg/log"
	"istio.io/pkg/version"

	"istio.io/istio/pilot/pkg/cri"
	"istio.io/istio/pilot/pkg/cri/kubeinject"
	"istio.io/istio/pkg/cmd"
)

var (
	loggingOptions = log.DefaultOptions()

	// hostIPVar contains the host's IP address.
	// It should be set in this pod's spec as:
	//- name: HOST_IP
	//  valueFrom:
	//    fieldRef:
	//      fieldPath: status.hostIP
	hostIPVar = env.RegisterStringVar("HOST_IP", "", "")

	rootCmd = &cobra.Command{
		Use:          "pilot-node-agent",
		Short:        "Istio Pilot node agent.",
		Long:         "Istio Pilot node agent runs on each k8s worker node to manage Istio sidecar containers.",
		SilenceUsage: true,
	}

	criProxyCmd = &cobra.Command{
		Use:   "cri-proxy",
		Short: "CRI proxy.",
		RunE: func(c *cobra.Command, args []string) error {
			cmd.PrintFlags(c.Flags())
			if err := log.Configure(loggingOptions); err != nil {
				return err
			}
			log.Infof("Version %s", version.Info.String())

			// TODO(rlenglet): Make this configurable.
			socket := "/var/run/istio-pilot-node-agent.sock"

			// TODO(rlenglet): Make this configurable or auto-detect it.
			criSocket := "/run/containerd/containerd.sock"

			// TODO(rlenglet): Make this configurable.
			excludedNamespaces := []string{"kube-system"}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			runtimeServiceClient, imageServiceClient, err := cri.NewClients(ctx, criSocket, criSocket)
			if err != nil {
				return err
			}

			injectionConfigurator := kubeinject.NewConfigurator(&kubeinject.Config{
				HostIP: hostIPVar.Get(),
				// TODO(rlenglet): Define a seccomp-profile-root flag.
				// Cf. https://github.com/kubernetes/kubernetes/blob/c7f9dd0bafe2f3328148a85eed9c18322b9f308e/cmd/kubelet/app/options/options.go#L397
				SeccompProfileRoot: "",
			})

			injector := cri.NewInjectorProxy(runtimeServiceClient, imageServiceClient,
				cri.DefaultLogger, injectionConfigurator, excludedNamespaces)

			s := cri.NewServer(socket, injector, injector)
			if err := s.Start(); err != nil {
				return err
			}

			cmd.WaitSignal(make(chan struct{}))
			return nil
		},
	}
)

func init() {
	// Attach the Istio logging options to the command.
	loggingOptions.AttachCobraFlags(rootCmd)

	cmd.AddFlags(rootCmd)

	rootCmd.AddCommand(criProxyCmd)
	rootCmd.AddCommand(version.CobraCommand())

	rootCmd.AddCommand(collateral.CobraCommand(rootCmd, &doc.GenManHeader{
		Title:   "Istio Pilot Node Agent",
		Section: "pilot-node-agent CLI",
		Manual:  "Istio Pilot Node Agent",
	}))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Errora(err)
		os.Exit(-1)
	}
}
