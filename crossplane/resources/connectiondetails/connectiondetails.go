// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package connectiondetails

import (
	"context"

	"github.com/crossplane/crossplane-runtime/pkg/resource"
	xpclaim "github.com/crossplane/crossplane-runtime/pkg/resource/unstructured/claim"
	xpcomposed "github.com/crossplane/crossplane-runtime/pkg/resource/unstructured/composed"
	xpcomposite "github.com/crossplane/crossplane-runtime/pkg/resource/unstructured/composite"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dsd-dbs/kubernetes-e2e-test-framework/crossplane/resources"
	"github.com/dsd-dbs/kubernetes-e2e-test-framework/defaults"
	"github.com/dsd-dbs/kubernetes-e2e-test-framework/klient"
	"github.com/dsd-dbs/kubernetes-e2e-test-framework/resources/secret"
)

type ConnectionDetails map[string][]byte

// FromClaimObject fetches the connection details exported as secret by a
// crossplane claim. It reloads the claim object from the API server based on
// the kind and name of claimObj. It returns nil if no secret object is found or
// the claim does not contain a reference.
//
// Use FromClaim if you already have an existing claim object with a reference
// to the connection secret.
func FromClaimObject(ctx context.Context, kube klient.Client, claimObj client.Object) (ConnectionDetails, error) {
	claim := xpclaim.Unstructured{}
	claim.SetGroupVersionKind(claimObj.GetObjectKind().GroupVersionKind())
	claim.SetName(claimObj.GetName())
	if claimObj.GetNamespace() != "" {
		claim.SetNamespace(claimObj.GetNamespace())
	}
	if err := retry.OnError(
		defaults.DefaultBackoff,
		func(error) bool { return true },
		func() error {
			return kube.Resources().GetControllerRuntimeClient().Get(ctx, client.ObjectKeyFromObject(&claim), &claim)
		},
	); err != nil {
		return nil, errors.Wrap(err, "cannot get claim")
	}
	return FromClaim(ctx, kube, &claim)
}

// FromClaim fetches the connection details exported as secret by a crossplane
// claim. It returns nil if no secret object is found or the claim does not
// contain a reference.
func FromClaim(ctx context.Context, kube klient.Client, claim resource.CompositeClaim) (ConnectionDetails, error) {
	ref := claim.GetWriteConnectionSecretToReference()
	if ref == nil {
		return nil, nil
	}
	return getConnectionDetails(ctx, kube, ref.Name, claim.GetNamespace())
}

// FromCompositeObject fetches the connection details exported as secret by a
// crossplane composite. It reloads the composite object from the API server
// based on the kind and name of compositeObj. It returns nil if no secret
// object is found or the composite does not contain a reference.
//
// Use FromComposite if you already have an existing composite object with a
// reference to the connection secret.
func FromCompositeObject(ctx context.Context, kube klient.Client, compositeObj client.Object) (ConnectionDetails, error) {
	composite := xpcomposite.Unstructured{}
	composite.SetGroupVersionKind(compositeObj.GetObjectKind().GroupVersionKind())
	composite.SetName(compositeObj.GetName())
	if compositeObj.GetNamespace() != "" {
		composite.SetNamespace(compositeObj.GetNamespace())
	}
	if err := retry.OnError(
		defaults.DefaultBackoff,
		func(error) bool { return true },
		func() error {
			return kube.Resources().GetControllerRuntimeClient().Get(ctx, client.ObjectKeyFromObject(&composite), &composite)
		},
	); err != nil {
		return nil, errors.Wrap(err, "cannot get claim")
	}
	return FromComposite(ctx, kube, &composite)
}

// FromComposite fetches the connection details exported as secret by a crossplane
// composite. It returns nil if no secret object is found or the composite does
// not contain a reference.
func FromComposite(ctx context.Context, kube klient.Client, composite resource.Composite) (ConnectionDetails, error) {
	ref := composite.GetWriteConnectionSecretToReference()
	if ref == nil {
		return nil, nil
	}
	return getConnectionDetails(ctx, kube, ref.Name, ref.Namespace)
}

// FromComposedObject fetches the connection details exported as secret by a
// crossplane composed resource. It returns nil if no secret object is found or
// the composed resource does not contain a reference.
//
// Use FromComposed if you already have an existing composed object with a
// reference to the connection secret.
func FromComposedObject(ctx context.Context, kube klient.Client, composedObj client.Object) (ConnectionDetails, error) {
	composed := xpcomposed.Unstructured{}
	composed.SetGroupVersionKind(composedObj.GetObjectKind().GroupVersionKind())
	composed.SetName(composedObj.GetName())
	if composedObj.GetNamespace() != "" {
		composed.SetNamespace(composedObj.GetNamespace())
	}
	if err := retry.OnError(
		defaults.DefaultBackoff,
		func(error) bool { return true },
		func() error {
			return kube.Resources().GetControllerRuntimeClient().Get(ctx, client.ObjectKeyFromObject(&composed), &composed)
		},
	); err != nil {
		return nil, errors.Wrap(err, "cannot get claim")
	}
	return FromComposed(ctx, kube, &composed)
}

// FromComposed fetches the connection details exported as secret by a crossplane
// composed resource. It returns nil if no secret object is found or the
// composed resource does not contain a reference.
func FromComposed(ctx context.Context, kube klient.Client, mg resource.Composed) (ConnectionDetails, error) {
	ref := mg.GetWriteConnectionSecretToReference()
	if ref == nil {
		return nil, nil
	}
	return getConnectionDetails(ctx, kube, ref.Name, ref.Namespace)
}

func getConnectionDetails(ctx context.Context, kube klient.Client, secretName, secretNamespace string) (ConnectionDetails, error) {
	data, err := secret.GetSecretData(ctx, kube, secretName, secretNamespace)
	return data, resource.IgnoreNotFound(err)
}

// FromComposedByClaim fetches the connection details exported as secret by
// a crossplane composed resource that is referenced by a claim and its
// composite. It returns nil if no secret object is found or the composed
// resource does not contain a reference.
func FromComposedByClaim(ctx context.Context, kube klient.Client, claim client.Object, resourceName string, resourceGVK schema.GroupVersionKind) (ConnectionDetails, error) {
	composed := xpcomposed.Unstructured{}
	composed.SetGroupVersionKind(resourceGVK)
	if err := retry.OnError(
		defaults.DefaultBackoff,
		func(error) bool { return true },
		func() error {
			return resources.GetComposedFromClaim(ctx, kube, claim, resourceName, &composed)
		},
	); err != nil {
		return nil, errors.Wrap(err, "cannot get claim")
	}
	return FromComposed(ctx, kube, &composed)
}

// FromComposedByClaim fetches the connection details exported as secret by
// a crossplane composed resource that is referenced by a composite resource.
// It returns nil if no secret object is found or the composed resource does
// not contain a reference.
func FromComposedByComposite(ctx context.Context, kube klient.Client, composite client.Object, resourceName string, resourceGVK schema.GroupVersionKind) (ConnectionDetails, error) {
	composed := xpcomposed.Unstructured{}
	composed.SetGroupVersionKind(resourceGVK)
	if err := retry.OnError(
		defaults.DefaultBackoff,
		func(error) bool { return true },
		func() error {
			return resources.GetComposedFromComposite(ctx, kube, composite, resourceName, &composed)
		},
	); err != nil {
		return nil, errors.Wrap(err, "cannot get claim")
	}
	return FromComposed(ctx, kube, &composed)
}
