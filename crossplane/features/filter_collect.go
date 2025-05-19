// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package features

import (
	"bytes"
	"context"
	"fmt"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	xpresource "github.com/crossplane/crossplane-runtime/pkg/resource"
	xpclaim "github.com/crossplane/crossplane-runtime/pkg/resource/unstructured/claim"
	xpcomposed "github.com/crossplane/crossplane-runtime/pkg/resource/unstructured/composed"
	xpcomposite "github.com/crossplane/crossplane-runtime/pkg/resource/unstructured/composite"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	"github.com/dsd-dbs/kubernetes-e2e-test-framework/klient"
)

func hasObjectStatusConditions(o client.Object) bool {
	_, onBlackList := groupKindsWithoutConditions[o.GetObjectKind().GroupVersionKind().GroupKind()]
	return !onBlackList
}

type objectWithConditions interface {
	client.Object
	xpresource.Conditioned
}

func isObjectSyncedAndReady(o objectWithConditions) bool {
	if o == nil {
		return false
	}
	if !hasObjectStatusConditions(o) {
		return true
	}
	if o.GetCondition(xpv1.TypeSynced).Status != corev1.ConditionTrue {
		return false
	}
	if o.GetCondition(xpv1.TypeReady).Status != corev1.ConditionTrue {
		return false
	}
	return true
}

// collectResourceTree for the given claim
func collectResourceTree(ctx context.Context, kube klient.Client, claim client.Object) (*xpclaim.Unstructured, *xpcomposite.Unstructured, []*xpcomposed.Unstructured, error) {
	var claimOnCluster *xpclaim.Unstructured
	var compositeOnCluster *xpcomposite.Unstructured
	var composedOnCluster []*xpcomposed.Unstructured //nolint:prealloc

	// Retrieve the actual claim on the cluster
	claimOnCluster = xpclaim.New(xpclaim.WithGroupVersionKind(claim.GetObjectKind().GroupVersionKind()))
	if err := klient.Get(ctx, kube, claim.GetName(), claim.GetNamespace(), claimOnCluster); err != nil {
		return claimOnCluster, compositeOnCluster, composedOnCluster, errors.Wrap(xpresource.IgnoreNotFound(err), "cannot get claim")
	}

	// Retrieve the actual composite from the reference specified in the claim
	compositeRef := claimOnCluster.GetResourceReference()
	if compositeRef == nil {
		return claimOnCluster, compositeOnCluster, composedOnCluster, nil
	}
	compositeOnCluster = xpcomposite.New(xpcomposite.WithGroupVersionKind(compositeRef.GroupVersionKind()))
	if err := klient.Get(ctx, kube, compositeRef.Name, "", compositeOnCluster); err != nil {
		return claimOnCluster, compositeOnCluster, composedOnCluster, errors.Wrap(xpresource.IgnoreNotFound(err), "cannot get composite")
	}

	// Retrieve the state of the actual composed resources by specified in
	// the composite
	for _, ref := range compositeOnCluster.GetResourceReferences() {
		childAtCluster := xpcomposed.New(xpcomposed.FromReference(ref))
		if err := klient.Get(ctx, kube, ref.Name, ref.Namespace, childAtCluster); err != nil {
			if kerrors.IsNotFound(err) {
				continue
			}
			return claimOnCluster, compositeOnCluster, composedOnCluster, errors.Wrapf(err, "cannot get object %q with name %q", ref.GroupVersionKind().String(), ref.Name)
		}
		composedOnCluster = append(composedOnCluster, childAtCluster)
	}
	return claimOnCluster, compositeOnCluster, composedOnCluster, nil
}

func combineObjectsToSlice(claim *xpclaim.Unstructured, composite *xpcomposite.Unstructured, composed []*xpcomposed.Unstructured) []client.Object {
	s := make([]client.Object, len(composed)+2)
	s[0] = claim
	s[1] = composite
	for i, c := range composed {
		s[i+2] = c
	}
	return s
}

type objectFilter func(o client.Object) bool

func prettyPrintObjects(objects []client.Object, filter objectFilter) string {
	buf := &bytes.Buffer{}
	for _, o := range objects {
		if o == nil || (filter != nil && filter(o)) {
			continue
		}
		raw, err := yaml.Marshal(o)
		if err != nil {
			fmt.Fprintf(buf, "---\nerror: marshalling object %s/%s: %s\n", o.GetObjectKind().GroupVersionKind().String(), o.GetName(), err.Error())
		} else {
			fmt.Fprintf(buf, "---\n%s\n", string(raw))
		}
	}
	return buf.String()
}

func skipReadyAndSynced(o client.Object) bool {
	cond, ok := o.(objectWithConditions)
	if !ok {
		return false
	}
	return isObjectSyncedAndReady(cond)
}
