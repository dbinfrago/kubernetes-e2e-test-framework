// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package features

import (
	"context"
	"testing"
	"time"

	xpclaim "github.com/crossplane/crossplane-runtime/pkg/resource/unstructured/claim"
	"github.com/crossplane/function-sdk-go/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachinerywait "k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/DSD-DBS/kubernetes-e2e-test-framework/klient"
)

// DeleteClaim deletes the given claim object and waits until the object and
// all subresources have been deleted.
// It uses foreground cascading delete policy to delete the claim and the
// underlying resources.
// It does not cancel if the passed timeout duration is zero.
func DeleteClaim(claim client.Object, timeout time.Duration, waitOpts ...WaitOption) features.Func {
	return func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		kube := cfg.Client()
		kubeClient := kube.Resources().GetControllerRuntimeClient()

		// Set the namespace to the default test namespace if not already set
		if claim.GetNamespace() == "" {
			isObjectNamespaced, err := kubeClient.IsObjectNamespaced(claim)
			if err != nil {
				t.Fatal(errors.Wrap(err, "cannot determine object scope").Error())
			}
			if isObjectNamespaced {
				claim.SetNamespace(cfg.Namespace())
			}
		}

		deleteCtx, cancel := contextWithOptionalTimeout(ctx, timeout)
		defer cancel()

		// CNP claims use cascading delete so the claim object will be the last
		// one deleted after all resources have been deleted
		if err := kubeClient.Delete(deleteCtx, claim, deleteForeground()); err != nil {
			t.Errorf("failed to delete resource: %s\n", err.Error())

			claim, composite, composed, err := collectResourceTree(ctx, kube, claim)
			if err != nil {
				t.Errorf("cannot collect undeleted resources: %s\n", err.Error())
			} else {
				t.Errorf("undeleted resources:\n%s\n", prettyPrintObjects(combineObjectsToSlice(claim, composite, composed), nil))
			}
		}

		waitCfg := WaitConfig{
			waitForOptions: []wait.Option{wait.WithTimeout(timeout)},
		}
		waitCfg.Apply(waitOpts)

		if err := wait.For(isClaimDeleted(cfg.Client(), claim), waitCfg.waitForOptions...); err != nil {
			t.Errorf("failed waiting for resources to become deleted: %s\n", err.Error())
		}
		return ctx
	}
}

func DeleteClaims(claims []client.Object, timeout time.Duration, waitOpts ...WaitOption) features.Func {
	return func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		kube := cfg.Client()

		return deleteClaims(ctx, t, kube, claims, timeout, waitOpts...)
	}
}

func DeleteClaimsWithClient(claims []client.Object, kube klient.Client, timeout time.Duration, waitOpts ...WaitOption) features.Func {
	return func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		return deleteClaims(ctx, t, kube, claims, timeout, waitOpts...)
	}
}

func deleteClaims(ctx context.Context, t *testing.T, kube klient.Client, claims []client.Object, timeout time.Duration, waitOpts ...WaitOption) context.Context {
	kubeclient := kube.Resources().GetControllerRuntimeClient()
	deleteCtx, cancel := contextWithOptionalTimeout(ctx, timeout)
	defer cancel()

	for _, claim := range claims {
		go func() {
			// Crossplane claims use cascading delete so the claim object will be the last
			// one deleted after all resources have been deleted
			if err := kubeclient.Delete(deleteCtx, claim, deleteForeground()); err != nil {
				t.Errorf("failed to delete resource: %s\n", err.Error())

				claim, composite, composed, errs := collectResourceTree(ctx, kube, claim)
				if errs != nil {
					t.Errorf("cannot collect undeleted resources: %s\n", errs.Error())
				} else {
					t.Errorf("undeleted resources:\n%s\n", prettyPrintObjects(combineObjectsToSlice(claim, composite, composed), nil))
				}
			}
		}()
	}

	waitCfg := WaitConfig{
		waitForOptions: []wait.Option{wait.WithTimeout(timeout)},
	}
	waitCfg.Apply(waitOpts)

	if err := wait.For(areClaimsDeleted(kube, claims), waitCfg.waitForOptions...); err != nil {
		t.Errorf("failed waiting for resources to become deleted: %s\n", err.Error())
	}
	return ctx
}

func deleteForeground() client.DeleteOption {
	return &client.DeleteOptions{
		PropagationPolicy: ptr.To(metav1.DeletePropagationForeground),
	}
}

func contextWithOptionalTimeout(parentCtx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout == 0 {
		return parentCtx, func() {}
	}
	return context.WithTimeout(parentCtx, timeout)
}

func isClaimDeleted(kube klient.Client, sourceClaim client.Object) apimachinerywait.ConditionWithContextFunc {
	return func(ctx context.Context) (bool, error) {
		claimOnCluster := xpclaim.New(xpclaim.WithGroupVersionKind(sourceClaim.GetObjectKind().GroupVersionKind()))
		if err := klient.Get(ctx, kube, sourceClaim.GetName(), sourceClaim.GetNamespace(), claimOnCluster); err == nil {
			return false, nil
		}
		return true, nil
	}
}

func areClaimsDeleted(kube klient.Client, sourceClaims []client.Object) apimachinerywait.ConditionWithContextFunc {
	return func(ctx context.Context) (bool, error) {
		for _, sourceClaim := range sourceClaims {
			claimOnCluster := xpclaim.New(xpclaim.WithGroupVersionKind(sourceClaim.GetObjectKind().GroupVersionKind()))
			if err := klient.Get(ctx, kube, sourceClaim.GetName(), sourceClaim.GetNamespace(), claimOnCluster); err == nil {
				return false, nil
			}
		}
		return true, nil
	}
}
