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

	"github.com/dbinfrago/kubernetes-e2e-test-framework/klient"
)

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
