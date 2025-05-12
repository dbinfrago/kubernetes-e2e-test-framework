// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package features

import (
	"context"
	"slices"
	"testing"
	"time"

	apimachinerywait "k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/DSD-DBS/kubernetes-e2e-test-framework/crossplane/internal/meta"
	"github.com/DSD-DBS/kubernetes-e2e-test-framework/klient"
)

// WaitForClaimReady is a feature that waits until the claim, composite and
// all composed resources have the conditions "Synced" and "Ready".
func WaitForClaimReady(claim client.Object, timeout time.Duration, waitOpts ...WaitOption) features.Func {
	return func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		kube := cfg.Client()

		waitCfg := WaitConfig{
			waitForOptions: []wait.Option{wait.WithTimeout(timeout)},
		}
		waitCfg.Apply(waitOpts)

		if err := wait.For(waitForConditionReadyAndSynced(kube, claim, waitCfg.ignoreComposedByCompositionResourceName), waitCfg.waitForOptions...); err != nil {
			t.Errorf("failed waiting for resources to become ready: %s\n", err.Error())

			// collect the resource tree and output the YAML of every unready
			// resource
			claim, composite, composed, err := collectResourceTree(ctx, kube, claim)
			if err != nil {
				t.Errorf("cannot collect unready resources: %s\n", err.Error())
			} else {
				t.Errorf("unready resources:\n%s\n", prettyPrintObjects(combineObjectsToSlice(claim, composite, composed), skipReadyAndSynced))
			}
		}
		return ctx
	}
}

func waitForConditionReadyAndSynced(kube klient.Client, sourceClaim client.Object, ignoreComposedByCompositionResourceName []string) apimachinerywait.ConditionWithContextFunc {
	return func(ctx context.Context) (bool, error) {
		claim, composite, composed, err := collectResourceTree(ctx, kube, sourceClaim)
		if err != nil {
			return false, err
		}
		// composite is nil if wait for MR, but is not nil anymore in isObjectSyncedAndReady due to interface cast
		if !isObjectSyncedAndReady(claim) || (composite != nil && !isObjectSyncedAndReady(composite)) {
			return false, nil
		}
		for _, o := range composed {
			if o == nil {
				continue
			}
			if slices.Contains(ignoreComposedByCompositionResourceName, meta.GetCompositionResourceName(o)) {
				continue // skip checks for resources that are explicitly ignored
			}
			if !isObjectSyncedAndReady(o) {
				return false, nil
			}
		}
		return true, nil
	}
}
