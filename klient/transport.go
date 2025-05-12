// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package klient

import (
	"errors"
	"net/http"
	"time"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/transport"
)

func SetConfigParameter(cfg *rest.Config) {
	cfg.WrapTransport = newRetryTransportWrapper()
	cfg.DisableCompression = true //https://docs.aws.amazon.com/eks/latest/best-practices/scale-control-plane.html
	cfg.QPS = -1
	cfg.RateLimiter = nil
	cfg.Timeout = 5 * time.Minute
}

// retryTransport wraps an existing RoundTripper and retries on transient errors.
type retryTransport struct {
	RoundTripper http.RoundTripper
	MaxRetries   int
	Backoff      time.Duration
}

// newRetryTransportWrapper creates a transport.WrapperFunc that performs a
// retry if the http client connection is lost
// (which happens frequently on mdp environments using VDS).
func newRetryTransportWrapper() transport.WrapperFunc {
	return func(rt http.RoundTripper) http.RoundTripper {
		return &retryTransport{
			RoundTripper: rt,
			MaxRetries:   3,
			Backoff:      3 * time.Second,
		}
	}
}

// RoundTrip retries the request on transient errors.
func (rt *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	for range rt.MaxRetries {
		resp, err = rt.RoundTripper.RoundTrip(req)
		if err == nil || !isTransientError(err) {
			return resp, err
		}
		time.Sleep(rt.Backoff) // Backoff before retrying
	}
	return resp, err
}

// isTransientError checks if the error is transient (e.g., connection lost).
func isTransientError(err error) bool {
	return err != nil && (errors.Is(err, http.ErrServerClosed) || err.Error() == "http2: client connection lost" || err.Error() == "context deadline exceeded")
}
