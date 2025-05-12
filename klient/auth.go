// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package klient

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/e2e-framework/klient"
)

// GetKubeUsername returns the actual kube user name in the same way
// "kubectl auth whoami" does it.
func GetKubeUsername(ctx context.Context, kube klient.Client) (string, error) {
	// We have to use unstructured here because the API types are not available
	// in this version of client-go (or are lacking the respective properties).
	ssar := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "authentication.k8s.io/v1",
			"kind":       "SelfSubjectReview",
		},
	}
	if err := kube.Resources().GetControllerRuntimeClient().Create(ctx, ssar); err != nil {
		return "", err
	}
	username, _, _ := unstructured.NestedString(ssar.Object, "status", "userInfo", "username")
	return username, nil
}
