package api

import uuid "github.com/satori/go.uuid"

// NewUUID returns a new v1 uuid that is compatible with MySQL uuid v1.
func NewUUID() string {
	return uuid.Must(uuid.NewV1()).String()
}
