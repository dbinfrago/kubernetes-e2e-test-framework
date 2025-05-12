// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package schema

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EnsureObjectGVK for o and apply the first registered group version kind in
// the provided scheme if the GVK of o is empty.
func EnsureObjectGVK(scheme *runtime.Scheme, o client.Object) error {
	ok := o.GetObjectKind()
	if !ok.GroupVersionKind().Empty() {
		return nil
	}
	kinds, _, err := scheme.ObjectKinds(o)
	if err != nil {
		return errors.Wrap(err, "cannot resolve object kinds")
	}
	ok.SetGroupVersionKind(kinds[0])
	return nil
}
