// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"context"

	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dsd-dbs/kubernetes-e2e-test-framework/crossplane/internal/meta"
	"github.com/dsd-dbs/kubernetes-e2e-test-framework/internal/json"
	"github.com/dsd-dbs/kubernetes-e2e-test-framework/internal/schema"
	"github.com/dsd-dbs/kubernetes-e2e-test-framework/klient"
)

// GetCompositeFromClaim loads the referenced composite from a cluster.
func GetCompositeFromClaim(ctx context.Context, kube klient.Client, claim resource.CompositeClaim, target client.Object) error {
	ref := claim.GetResourceReference()
	if ref == nil {
		return errors.New("empty composite reference")
	}
	if ok := target.GetObjectKind(); ok.GroupVersionKind().Empty() {
		ok.SetGroupVersionKind(ref.GroupVersionKind())
	}
	return klient.Get(ctx, kube, ref.Name, "", target)
}

// GetCompositeFromClaim loads the referenced composite from a cluster.
func GetComposedFromClaim(ctx context.Context, kube klient.Client, claim client.Object, resourceName string, composed client.Object) error {
	return getComposed(ctx, kube, claim.GetName(), claim.GetNamespace(), resourceName, composed)
}

func GetComposedFromComposite(ctx context.Context, kube klient.Client, composite client.Object, resourceName string, composed client.Object) error {
	return getComposed(ctx, kube, meta.GetClaimName(composite), meta.GetClaimNamespace(composite), resourceName, composed)
}

func getComposed(ctx context.Context, kube klient.Client, claimName, claimNamespace, resourceName string, composed client.Object) error {
	if err := schema.EnsureObjectGVK(kube.Resources().GetScheme(), composed); err != nil {
		return err
	}
	ul := unstructured.UnstructuredList{}
	ul.SetGroupVersionKind(composed.GetObjectKind().GroupVersionKind())
	matchLabels := client.MatchingLabels{
		meta.LabelKeyClaimName:      claimName,
		meta.LabelKeyClaimNamespace: claimNamespace,
	}
	if err := kube.Resources().GetControllerRuntimeClient().List(ctx, &ul, matchLabels, client.MatchingFields{}); err != nil {
		return errors.Wrap(err, "cannot list objects")
	}
	// The resource is identified by the composition resource name annotation.
	// Pick up the first object that matches the given resourceName.
	foundComposed := unstructured.Unstructured{}
	found := false
	for _, u := range ul.Items {
		if meta.GetCompositionResourceName(&u) == resourceName {
			foundComposed = u
			found = true
			break
		}
	}
	if !found {
		return errors.Errorf("no object with resource name %q", resourceName)
	}
	if u, ok := composed.(runtime.Unstructured); ok {
		u.SetUnstructuredContent(foundComposed.Object)
		return nil
	}
	return json.Convert(foundComposed, composed)
}
