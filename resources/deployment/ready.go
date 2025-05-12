// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package deployment

import (
	"context"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/DSD-DBS/kubernetes-e2e-test-framework/klient"
)

// IsDeploymentAvailable determines if the available condition of a deployment
// is fulfilled.
func IsDeploymentAvailable(ctx context.Context, kube klient.Client, name, namespace string) (bool, error) {
	deploy := appsv1.Deployment{}
	if err := klient.Get(ctx, kube, name, namespace, &deploy); err != nil {
		return false, errors.Wrap(err, "cannot get deployment")
	}
	for _, c := range deploy.Status.Conditions {
		if c.Type == appsv1.DeploymentAvailable {
			return c.Status == corev1.ConditionTrue, nil
		}
	}
	return false, nil
}

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

func GetPodsForDeployment(ctx context.Context, kube klient.Client, name, namespace string) ([]corev1.Pod, error) {
	kubeClient := kube.Resources().GetControllerRuntimeClient()
	deploy := &appsv1.Deployment{}
	if err := klient.Get(ctx, kube, name, namespace, deploy); err != nil {
		return nil, errors.Wrap(err, "cannot get deployment")
	}
	labelSelector := labels.Set(deploy.Spec.Selector.MatchLabels)
	podList := &corev1.PodList{}
	listOptions := []client.ListOption{
		client.InNamespace(deploy.Namespace),
		client.MatchingLabelsSelector{
			Selector: labels.SelectorFromSet(labelSelector),
		},
	}
	if err := kubeClient.List(ctx, podList, listOptions...); err != nil {
		return nil, err
	}
	return podList.Items, nil
}
