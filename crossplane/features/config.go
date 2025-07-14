// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package features

import (
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/e2e-framework/klient/wait"
)

// WaitConfig specifies how a waiting assessment should be executed.
type WaitConfig struct {
	waitForOptions                          []wait.Option
	ignoreComposedByCompositionResourceName []string
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

// WaitIgnoreComposedByCompositionResourceName causes the composed resources
// with the given name to be considered as ready during a WaitFor operation.
func WaitIgnoreComposedByCompositionResourceName(names ...string) WaitOption {
	return func(c *WaitConfig) {
		c.ignoreComposedByCompositionResourceName = append(c.ignoreComposedByCompositionResourceName, names...)
	}
}

// groupKindsWithoutConditions that don't have synced or ready conditions.
// Use GK instead of GVK because it should apply to all schema versions.
//
// Extend this list using RegisterKindWithoutCondition
var groupKindsWithoutConditions = map[schema.GroupKind]interface{}{
	{Group: "aws.crossplane.io", Kind: "ProviderConfig"}:              nil,
	{Group: "gitlab.crossplane.io", Kind: "ProviderConfig"}:           nil,
	{Group: "grafana.crossplane.io", Kind: "ProviderConfig"}:          nil,
	{Group: "argocd.crossplane.io", Kind: "ProviderConfig"}:           nil,
	{Group: "helm.crossplane.io", Kind: "ProviderConfig"}:             nil,
	{Group: "kubernetes.crossplane.io", Kind: "ProviderConfig"}:       nil,
	{Group: "aws.upbound.io", Kind: "ProviderConfig"}:                 nil,
	{Group: "rbac.authorization.k8s.io", Kind: "ClusterRole"}:         nil,
	{Group: "apiextensions.crossplane.io", Kind: "EnvironmentConfig"}: nil,
	{Group: "apiextensions.crossplane.io", Kind: "Usage"}:             nil,
}

func RegisterKindWithoutCondition(kind schema.GroupKind) {
	groupKindsWithoutConditions[kind] = nil
}

func RegisterKindsWithoutCondition(kinds []schema.GroupKind) {
	for _, kind := range kinds {
		RegisterKindWithoutCondition(kind)
	}
}
