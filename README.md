<!--
 ~ SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
 ~ SPDX-License-Identifier: Apache-2.0
 -->

# Kubernetes E2E Test Framework

A Go library that provides extensions to [`github.com/kubernetes-sigs/e2e-framework`](https://github.com/kubernetes-sigs/e2e-framework) to implement end-to-end tests for kubernetes APIs.

## Usage

Use it like so to apply crossplane claims to your cluster and make assertions on it:

```go
package coolfeature

import (
	"testing"
	"time"
	
	"sigs.k8s.io/e2e-framework/pkg/features"
	e2efeatures "github.com/DSD-DBS/kubernetes-e2e-test-framework/features"
	crossplanefeatures "github.com/DSD-DBS/kubernetes-e2e-test-framework/crossplane/features"
)

func FeatureTest(t *testing.T) {
	crossplaneClaim := LoadClaimFromYaml(`apiVersion: ...`)

	features.New("Cool Feature")
	.Setup(e2efeatures.ApplyObject(crossplaneClaim))
	.Assess("await claim", crossplanefeatures.WaitForClaimReady(crossplaneClaim, 5*time.Minutes))
	.Assess("check claim status", e2efeatures.Assess(func(ctx context.Context, t *testing.T, cfg *envconf.Config) error {
		kube := cfg.Client()
		// Use kube client to check that your claim produces the correct managed resources.

		return nil
	}))
	.Assess("delete claim", crossplanefeatures.DeleteClaim(crossplaneClaim, 5*time.Minutes))
}
```

# Contributing

See our [Contributing Guidelines](./CONTRIBUTING.md).

# Licensing

Each file contains a license reference to one of the [included licenses](./LICENSES).
