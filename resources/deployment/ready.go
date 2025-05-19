// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package deployment

import (
	"context"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/dsd-dbs/kubernetes-e2e-test-framework/klient"
)

// IsDeploymentAvailable determines if the available condition of a deployment
// is fulfilled.
func IsDeploymentAvailable(ctx context.Context, kube klient.Client, name, namespace string) (bool, error) {
	deploy := appsv1.Deployment{}
	if err := klient.Get(ctx, kube, name, namespace, &deploy); err != nil {
		return false, errors.Wrap(err, "cannot get deployment")
	}
	for _, c := range deploy.Status.Conditions {
		if c.Type == appsv1.DeploymentAvailable {
			return c.Status == corev1.ConditionTrue, nil
		}
	}
	return false, nil
}
