// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package uuidutil

import (
	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/types"
)

const uuidNamespace = "metal.ironcore.dev"

func UUIDFromObjectKey(key types.NamespacedName) string {
	namespacedUUID := uuid.NewMD5(uuid.UUID{}, []byte(uuidNamespace))
	uuidFromKey := uuid.NewMD5(namespacedUUID, []byte(key.String()))
	return uuidFromKey.String()
}
