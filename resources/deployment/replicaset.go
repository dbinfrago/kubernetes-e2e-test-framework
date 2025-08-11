// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package deployment

import (
	"context"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dbinfrago/kubernetes-e2e-test-framework/klient"
)

func GetActiveReplicaSetForDeployment(ctx context.Context, kube klient.Client, name, namespace string) (*appsv1.ReplicaSet, bool, error) {
	kubeClient := kube.Resources().GetControllerRuntimeClient()
	deploy := &appsv1.Deployment{}
	if err := klient.Get(ctx, kube, name, namespace, deploy); err != nil {
		return nil, false, errors.Wrap(err, "cannot get deployment")
	}

	replicaSetList := &appsv1.ReplicaSetList{}
	listOptions := []client.ListOption{
		client.InNamespace(deploy.Namespace),
	}
	if err := kubeClient.List(ctx, replicaSetList, listOptions...); err != nil {
		return nil, false, err
	}

	for _, replicaSet := range replicaSetList.Items {
		if replicaSet.Status.AvailableReplicas < 1 {
			continue
		}

		for _, owner := range replicaSet.OwnerReferences {
			if owner.Kind == "Deployment" && owner.Name == deploy.Name {
				return &replicaSet, true, nil
			}
		}
	}

	return nil, false, nil
}
