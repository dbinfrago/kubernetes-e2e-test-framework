// SPDX-FileCopyrightText: Copyright DB InfraGO AG and contributors
// SPDX-License-Identifier: Apache-2.0

package json

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// Convert is a shortcut to convert an object into another by marshalling
// from into JSON and unmarshalling the result into to.
func Convert(from, to any) error {
	raw, err := json.Marshal(from)
	if err != nil {
		return errors.Wrap(err, "cannot marshal")
	}
	return errors.Wrap(json.Unmarshal(raw, to), "cannot unmarshal")
}
