// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package ingress

import (
	"context"

	"github.com/pkg/errors"
	networkingv1 "k8s.io/api/networking/v1"

	"github.com/dbinfrago/kubernetes-e2e-test-framework/klient"
)

// IsALBAvailable determines if the latest ingress provisioned load balancer is available based on desired number of LBs.
func IsALBAvailable(ctx context.Context, kube klient.Client, desiredNumberOfAlbs int, name, namespace string) (bool, *networkingv1.Ingress, error) {
	ingress := &networkingv1.Ingress{}
	if err := klient.Get(ctx, kube, name, namespace, ingress); err != nil {
		return false, ingress, errors.Wrap(err, "cannot get ingress")
	}
	return len(ingress.Status.LoadBalancer.Ingress) == desiredNumberOfAlbs && ingress.Status.LoadBalancer.Ingress[desiredNumberOfAlbs-1].Hostname != "", ingress, nil
}
