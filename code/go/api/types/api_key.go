package types

import (
	"time"

	"github.com/google/uuid"
)

// ApiKey struct definition.
type ApiKey struct {
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	KeyHash   string    `json:"key_hash" db:"key_hash"`
	Id        uuid.UUID `json:"id" db:"id"`
	SystemKey bool      `json:"system_key" db:"system_key"`
}
