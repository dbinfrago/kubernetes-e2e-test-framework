// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package pvc

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dbinfrago/kubernetes-e2e-test-framework/defaults"
	"github.com/dbinfrago/kubernetes-e2e-test-framework/klient"
)

// IsPersistentVolumeClaimStatus checks if PVC has specified status
// It assumes that failure to get a pvc is not an error, since the pvc might not
// have been created yet, which just means that it is not available yet.
// This means that any other errors related to kubernetes API access are also 
// ignored.
func IsPersistentVolumeClaimStatus(ctx context.Context, kube klient.Client, pvcStatus corev1.PersistentVolumeClaimPhase, name, namespace string) (bool, error) {
	pvc := &corev1.PersistentVolumeClaim{}
	if err := klient.Get(ctx, kube, name, namespace, pvc); err != nil {
		return false, nil
	}
	return pvc.Status.Phase == pvcStatus, nil
}

// IsPersistentVolumeVolumeModificationSuccessful checks PVC-related events to see if VolumeModification via
// Volumemodifier was successful
// It assumes that failure to get a pvc is not an error, since the pvc might not
// have been created yet, which just means that it is not available yet.
// This means that any other errors related to kubernetes API access are also 
// ignored.
//
// Cf: https://github.com/torredil/volume-modifier-for-k8s/blob/5eb7d23f72d688ae0b7d9db8019d3371f4e93289/pkg/controller/controller.go#L288
// https://aws.amazon.com/de/blogs/storage/simplifying-amazon-ebs-volume-migration-and-modification-using-the-ebs-csi-driver/
func IsPersistentVolumeVolumeModificationSuccessful(ctx context.Context, kube klient.Client, name, namespace string) (bool, error) {
	kubeClient := kube.Resources().GetControllerRuntimeClient()
	events := &corev1.EventList{}
	err := retry.OnError(defaults.DefaultBackoff, func(error) bool { return true }, func() error {
		return kubeClient.List(ctx, events, client.InNamespace(namespace), client.MatchingFields{"involvedObject.name": name})
	})
	if err != nil {
		return false, nil
	}
	for _, item := range events.Items {
		if item.Reason == "VolumeModificationSuccessful" {
			return true, nil
		}
	}
	return false, nil
}
