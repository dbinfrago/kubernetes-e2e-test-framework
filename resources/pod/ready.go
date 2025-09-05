// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package pod

import (
	"context"

	corev1 "k8s.io/api/core/v1"

	"github.com/dbinfrago/kubernetes-e2e-test-framework/klient"
)

// IsPodAvailable determines if the available condition of a pod is fulfilled.
// It assumes that failure to get a pod is not an error, since the pod might not
// have been created yet, which just means that it is not available yet.
// This means that any other errors related to kubernetes API access are also 
// ignored.
func IsPodAvailable(ctx context.Context, kube klient.Client, name, namespace string) (bool, error) {
	pod := corev1.Pod{}
	if err := klient.Get(ctx, kube, name, namespace, &pod); err != nil {
		return false, nil
	}
	for _, c := range pod.Status.Conditions {
		if c.Type == corev1.ContainersReady {
			return c.Status == corev1.ConditionTrue, nil
		}
	}
	return false, nil
}
