// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package features

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"sigs.k8s.io/e2e-framework/klient"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	pod "github.com/dbinfrago/kubernetes-e2e-test-framework/resources/pod"
)

// AssessExecInPod executes the given command in the specified container and
// checks if it executes successfully.
func AssessExecInPod(namespace, pod, container string, command []string) features.Func {
	return Assess(func(ctx context.Context, t *testing.T, cfg *envconf.Config) error {
		return assessExecInPod(ctx, t, cfg.Client(), namespace, pod, container, command)
	})
}

// AssessExecInPodWithClient executes the given command in the specified
// container using the provided kube client and checks if it executes
// successfully.
func AssessExecInPodWithClient(kube klient.Client, namespace, pod, container string, command []string) features.Func {
	return Assess(func(ctx context.Context, t *testing.T, cfg *envconf.Config) error {
		return assessExecInPod(ctx, t, kube, namespace, pod, container, command)
	})
}

func assessExecInPod(ctx context.Context, t *testing.T, kube klient.Client, namespace, podName, container string, command []string) error {
	stdout, stderr, err := pod.ExecWithClient(ctx, kube, namespace, podName, container, command)
	if err != nil {
		t.Errorf(
			"command did not execute successfully: %s\n\nBEGIN STDOUT\n%s\nEND STDOUT\n\nBEGIN STDERR\n%s\nEND STDERR\n",
			err.Error(),
			stdout.String(),
			stderr.String(),
		)
		return errors.Wrap(err, "command did not execute successfully")
	}
	return nil
}
