// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package klient

import (
	"context"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"

	e2edefaults "github.com/dbinfrago/kubernetes-e2e-test-framework/defaults"
)

// Get is a shorthand to retrieve an object using a [sigs.k8s.io/e2e-framework/klient.Client].
func Get(ctx context.Context, kube Client, name, namespace string, target client.Object) error {
	nn := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}
	return retry.OnError(
		e2edefaults.DefaultBackoff,
		func(error) bool { return true },
		func() error {
			return kube.Resources().GetControllerRuntimeClient().Get(ctx, nn, target)
		},
	)
}
