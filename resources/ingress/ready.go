// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package ingress

import (
	"context"

	networkingv1 "k8s.io/api/networking/v1"

	"github.com/dbinfrago/kubernetes-e2e-test-framework/v2/klient"
)

// IsALBAvailable determines if the latest ingress provisioned load balancer is available based on desired number of LBs.
// It assumes that failure to get a ingress is not an error, since the ingress
// might not have been created yet, which just means that it is not available
// yet.
// This means that any other errors related to kubernetes API access are also 
// ignored.
func IsALBAvailable(ctx context.Context, kube klient.Client, desiredNumberOfAlbs int, name, namespace string) (bool, *networkingv1.Ingress, error) {
	ingress := &networkingv1.Ingress{}
	if err := klient.Get(ctx, kube, name, namespace, ingress); err != nil {
		return false, ingress, nil
	}
	return len(ingress.Status.LoadBalancer.Ingress) == desiredNumberOfAlbs && ingress.Status.LoadBalancer.Ingress[desiredNumberOfAlbs-1].Hostname != "", ingress, nil
}
