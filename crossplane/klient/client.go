// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package klient

import (
	"context"

	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/DSD-DBS/kubernetes-e2e-test-framework/crossplane/resources/connectiondetails"
	"github.com/DSD-DBS/kubernetes-e2e-test-framework/klient"
)

// NewClientFromClaimConnectionDetails creates a new kube client from a
// kubeconfig that is exposed in connection details secret of a claim
// resource.
func NewClientFromClaimConnectionDetails(ctx context.Context, kube klient.Client, claim resource.CompositeClaim, connectionDetailsKey string) (klient.Client, error) {
	connectionDetails, err := connectiondetails.FromClaim(ctx, kube, claim)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get cluster connection details")
	}
	return newClientFromConnectionDetails(connectionDetails, connectionDetailsKey)
}

// NewClientFromCompositeConnectionDetails creates a new kube client from a
// kubeconfig that is exposed in connection details secret of a composite
// resource.
func NewClientFromCompositeConnectionDetails(ctx context.Context, kube klient.Client, composite resource.Composite, connectionDetailsKey string) (klient.Client, error) {
	connectionDetails, err := connectiondetails.FromComposite(ctx, kube, composite)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get cluster connection details")
	}
	return newClientFromConnectionDetails(connectionDetails, connectionDetailsKey)
}

// NewClientFromComposedConnectionDetails creates a new kube client from a
// kubeconfig that is exposed in connection details secret of a composed
// resource.
func NewClientFromComposedConnectionDetails(ctx context.Context, kube klient.Client, claim client.Object, composedResourceName string, composedResourceGVK schema.GroupVersionKind, connectionDetailsKey string) (klient.Client, error) {
	connectionDetails, err := connectiondetails.FromComposedByClaim(ctx, kube, claim, composedResourceName, composedResourceGVK)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get cluster connection details")
	}
	return newClientFromConnectionDetails(connectionDetails, connectionDetailsKey)
}

func newClientFromConnectionDetails(cd map[string][]byte, key string) (klient.Client, error) {
	configBytes, exists := cd[key]
	if !exists {
		return nil, errors.Errorf("no connection details for key %q", key)
	}
	return klient.NewClientFromConfigBytes(configBytes)
}
