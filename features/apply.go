// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package features

import (
	"context"
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/dsd-dbs/kubernetes-e2e-test-framework/defaults"
	"github.com/dsd-dbs/kubernetes-e2e-test-framework/internal/schema"
	"github.com/dsd-dbs/kubernetes-e2e-test-framework/klient"
)

// ApplyObject returns a [sigs.k8s.io/e2e-framework/pkg/features.Func] that
// applies o on the cluster using server-side apply.
//
// It automatically sets the namespace of o to the preconfigured test namespace
// if o does not already have a namespace set.
//
// mod is an optional function that can be given to modify the
// object before applying it.
func ApplyObject(o client.Object, mod func(o client.Object)) features.Func {
	return Assess(func(ctx context.Context, t *testing.T, cfg *envconf.Config) error {
		return applyObject(ctx, t, cfg, cfg.Client().Resources().GetControllerRuntimeClient(), o, mod)
	})
}

// ApplyObjectWithClient returns a [sigs.k8s.io/e2e-framework/pkg/features.Func] that
// applies o on the cluster using server-side apply and the provided client
//
// It automatically sets the namespace of o to the preconfigured test namespace
// if o does not already have a namespace set.
func ApplyObjectWithClient(o client.Object, kube klient.Client, mod func(o client.Object)) features.Func {
	return Assess(func(ctx context.Context, t *testing.T, cfg *envconf.Config) error {
		return applyObject(ctx, t, cfg, kube.Resources().GetControllerRuntimeClient(), o, mod)
	})
}

func applyObject(ctx context.Context, t *testing.T, cfg *envconf.Config, kube client.Client, o client.Object, mod func(o client.Object)) error {
	// remove any managed fields in request for SSA
	o.SetManagedFields(nil)
	o.SetResourceVersion("")
	o.SetGeneration(0)

	// Set the namespace to the default test namespace if not already set
	if o.GetNamespace() == "" {
		isObjectNamespaced, err := kube.IsObjectNamespaced(o)
		if err != nil {
			return errors.Wrap(err, "cannot determine object scope")
		}
		if isObjectNamespaced {
			o.SetNamespace(cfg.Namespace())
		}
	}
	// Ensure the object contains a gvk
	if err := schema.EnsureObjectGVK(kube.Scheme(), o); err != nil {
		return errors.Wrap(err, "cannot set GVK from Scheme")
	}

	if mod != nil {
		mod(o)
	}

	return errors.Wrap(applyObjectSSA(ctx, kube, o, fieldOwnerFromT(t)), "cannot apply object")
}

func fieldOwnerFromT(t *testing.T) string {
	return fmt.Sprintf("test/%s", t.Name())
}

func applyObjectSSA(ctx context.Context, kube client.Client, o client.Object, fieldOwner string) error {
	return retry.OnError(defaults.DefaultBackoff, func(error) bool { return true }, func() error {
		return kube.Patch(ctx, o, client.Apply, client.FieldOwner(fieldOwner), client.ForceOwnership)
	})
}
