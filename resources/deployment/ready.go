// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package deployment

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/dbinfrago/kubernetes-e2e-test-framework/v2/klient"
)

// IsDeploymentAvailable determines if the available condition of a deployment
// is fulfilled.
// It assumes that failure to get a deployment is not an error, since the
// deployment might not have been created yet, which just means that it is not
// available yet.
// This means that any other errors related to kubernetes API access are also 
// ignored.
func IsDeploymentAvailable(ctx context.Context, kube klient.Client, name, namespace string) (bool, error) {
	deploy := appsv1.Deployment{}
	if err := klient.Get(ctx, kube, name, namespace, &deploy); err != nil {
		return false, nil
	}
	for _, c := range deploy.Status.Conditions {
		if c.Type == appsv1.DeploymentAvailable {
			return c.Status == corev1.ConditionTrue, nil
		}
	}
	return false, nil
}
