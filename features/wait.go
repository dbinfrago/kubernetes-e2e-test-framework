// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package features

import (
	"context"
	"testing"
	"time"

	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/dbinfrago/kubernetes-e2e-test-framework/klient"
)

// WaitConfig specifies how a waiting assessment should be executed.
type WaitConfig struct {
	waitForOptions []wait.Option
}

// Apply the given waitopts.
func (c *WaitConfig) Apply(opts []WaitOption) {
	for _, o := range opts {
		o(c)
	}
}

// WaitOption modifies a WaitConfig.
type WaitOption func(c *WaitConfig)

// WaitWithIntervall defines WaitOption that defines the poll interval in which
// a wait condition is checked.
func WaitWithIntervall(intervall time.Duration) WaitOption {
	return func(c *WaitConfig) {
		c.waitForOptions = append(c.waitForOptions, wait.WithInterval(intervall))
	}
}

func WaitWithTimeout(timeout time.Duration) WaitOption {
	return func(c *WaitConfig) {
		c.waitForOptions = append(c.waitForOptions, wait.WithTimeout(timeout))
	}
}

type WaitForFunc func(ctx context.Context, kube klient.Client) (done bool, err error)

// WaitForWithKube repeatedly waits until waitFunc succeeds or a timeout is
// reached.
func WaitFor(waitFunc WaitForFunc, timeout time.Duration, waitOpts ...WaitOption) features.Func {
	return func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		kube := cfg.Client()

		waitCfg := WaitConfig{
			waitForOptions: []wait.Option{wait.WithTimeout(timeout)},
		}
		waitCfg.Apply(waitOpts)

		err := wait.For(func(ctx context.Context) (done bool, err error) {
			return waitFunc(ctx, kube)
		}, waitCfg.waitForOptions...)
		if err != nil {
			t.Errorf("failed to wait for condition: %s\n", err.Error())
		}
		return ctx
	}
}

// WaitForWithKube repeatedly waits until waitFunc succeeds or a timeout is
// reached.
func WaitForWithClient(waitFunc WaitForFunc, client klient.Client, timeout time.Duration, waitOpts ...WaitOption) features.Func {
	return func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		waitCfg := WaitConfig{
			waitForOptions: []wait.Option{wait.WithTimeout(timeout)},
		}
		waitCfg.Apply(waitOpts)

		err := wait.For(func(ctx context.Context) (done bool, err error) {
			return waitFunc(ctx, client)
		}, waitCfg.waitForOptions...)
		if err != nil {
			t.Errorf("failed to wait for condition: %s\n", err.Error())
		}
		return ctx
	}
}
