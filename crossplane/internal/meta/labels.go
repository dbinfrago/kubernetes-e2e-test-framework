// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package meta

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	LabelKeyClaimName      = "crossplane.io/claim-name"
	LabelKeyClaimNamespace = "crossplane.io/claim-namespace"
)

// GetClaimName stored in the label of o or an empty string of it
// is not defined.
func GetClaimName(o metav1.Object) string {
	return o.GetLabels()[LabelKeyClaimName]
}

// GetClaimNamespace stored in the label of o or an empty string of it
// is not defined.
func GetClaimNamespace(o metav1.Object) string {
	return o.GetLabels()[LabelKeyClaimNamespace]
}
