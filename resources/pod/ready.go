// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package pod

import (
	"context"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"

	"github.com/DSD-DBS/kubernetes-e2e-test-framework/klient"
)

// IsPodAvailable determines if the available condition of a pod is fulfilled.
func IsPodAvailable(ctx context.Context, kube klient.Client, name, namespace string) (bool, error) {
	pod := corev1.Pod{}
	if err := klient.Get(ctx, kube, name, namespace, &pod); err != nil {
		return false, errors.Wrap(err, "cannot get pod")
	}
	for _, c := range pod.Status.Conditions {
		if c.Type == corev1.ContainersReady {
			return c.Status == corev1.ConditionTrue, nil
		}
	}
	return false, nil
}
