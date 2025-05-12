// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package pod

import (
	"bytes"
	"context"

	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/e2e-framework/pkg/envconf"

	"github.com/DSD-DBS/kubernetes-e2e-test-framework/defaults"
	"github.com/DSD-DBS/kubernetes-e2e-test-framework/klient"
)

// Exec executes the given command in the specified container and returns
// the recorded data for stdout and stdin. An error is returned if the command
// execution fails.
func Exec(ctx context.Context, cfg *envconf.Config, namespace, pod, container string, command []string) (stdout, stderr *bytes.Buffer, err error) {
	return execInPod(ctx, cfg.Client(), namespace, pod, container, command)
}

// ExecInPodWithConfig executes the given command in the specified container and
// returns the recorded data for stdout and stdin. An error is returned if the
// command execution fails.
func ExecWithClient(ctx context.Context, kube klient.Client, namespace, pod, container string, command []string) (stdout, stderr *bytes.Buffer, err error) {
	return execInPod(ctx, kube, namespace, pod, container, command)
}

func execInPod(ctx context.Context, kube klient.Client, namespace, pod, container string, command []string) (stdout, stderr *bytes.Buffer, err error) {
	stdout = &bytes.Buffer{}
	stderr = &bytes.Buffer{}
	err = retry.OnError(
		defaults.DefaultBackoff,
		func(error) bool { return true },
		func() error {
			return kube.Resources().ExecInPod(ctx, namespace, pod, container, command, stdout, stderr)
		},
	)
	return stdout, stderr, err
}
