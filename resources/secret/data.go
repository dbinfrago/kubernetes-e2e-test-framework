// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	corev1 "k8s.io/api/core/v1"

	"github.com/dbinfrago/kubernetes-e2e-test-framework/klient"
)

// GetSecretData from a kubernetes secret using the provided client.
func GetSecretData(ctx context.Context, kube klient.Client, name, namespace string) (map[string][]byte, error) {
	secret := &corev1.Secret{}
	if err := klient.Get(ctx, kube, name, namespace, secret); err != nil {
		return nil, err
	}
	return secret.Data, nil
}

// GetSecretDataKey from a kubernetes secret using the provided client.
func GetSecretDataKey(ctx context.Context, kube klient.Client, name, namespace, key string) ([]byte, bool, error) {
	data, err := GetSecretData(ctx, kube, name, namespace)
	if err != nil {
		return nil, false, err
	}
	keyData, exists := data[key]
	return keyData, exists, nil
}
