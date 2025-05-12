// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package meta

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const AnnotationKeyCompositionResourceName = "crossplane.io/composition-resource-name"

// GetCompositionResourceName annotation of o.
func GetCompositionResourceName(o metav1.Object) string {
	return o.GetAnnotations()[AnnotationKeyCompositionResourceName]
}
