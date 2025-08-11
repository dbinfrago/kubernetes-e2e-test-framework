// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package pod

import (
	"context"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"

	"github.com/dbinfrago/kubernetes-e2e-test-framework/klient"
)

func GetPod(ctx context.Context, kube klient.Client, name, namespace string) (*corev1.Pod, error) {
	pod := &corev1.Pod{}
	if err := klient.Get(ctx, kube, name, namespace, pod); err != nil {
		return nil, errors.Wrap(err, "cannot get pod")
	}
	return pod, nil
}
