// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package klient

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/e2e-framework/klient"
)

type Client = klient.Client

// NewClientFromConfigBytes creates a new kube client using the provided config.
func NewClientFromConfig(cfg *rest.Config) (Client, error) {
	if cfg.WrapTransport == nil {
		cfg.WrapTransport = newRetryTransportWrapper()
	}
	return klient.New(cfg)
}

// NewClientFromConfigBytes creates a new kube client using the provided config
// file.
func NewClientFromConfigBytes(configBytes []byte) (Client, error) {
	config, err := clientcmd.NewClientConfigFromBytes(configBytes)
	if err != nil {
		return nil, err
	}
	restConfig, err := config.ClientConfig()
	if err != nil {
		return nil, err
	}
	SetConfigParameter(restConfig)
	return NewClientFromConfig(restConfig)
}
