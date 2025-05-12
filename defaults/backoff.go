// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package defaults

import (
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

var DefaultBackoff = wait.Backoff{
	Steps:    4,
	Duration: 50 * time.Millisecond,
	Factor:   3.0,
	Jitter:   0.1,
}
