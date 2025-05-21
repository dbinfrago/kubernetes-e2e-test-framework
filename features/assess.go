// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package features

import (
	"context"
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/dsd-dbs/kubernetes-e2e-test-framework/klient"
)

// AssessKubeFunc is a simplified assess function.
type AssessKubeFunc func(ctx context.Context, t *testing.T, cfg *envconf.Config, kube klient.Client) error

// Assess is a shorthand assess function that invokes a delegate with a
// preconfigured kube client that performs a test step and returns an error in
// case it fails.
func AssessKube(assessFunc AssessKubeFunc) features.Func {
	return func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		kube := cfg.Client()
		if err := assessFunc(ctx, t, cfg, kube); err != nil {
			t.Errorf("assess failed: %s\n", err.Error())
		}
		return ctx
	}
}

// AssessFunc is a simplified assess function.
type AssessFunc func(ctx context.Context, t *testing.T, cfg *envconf.Config) error

// Assess is a shorthand assess function that simplifies error handling.
func Assess(assessFunc AssessFunc) features.Func {
	return func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		if err := assessFunc(ctx, t, cfg); err != nil {
			t.Errorf("assess failed: %s\n", err.Error())
		}
		return ctx
	}
}
