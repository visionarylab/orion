package models

import (
	"fmt"

	"github.com/Syncano/pkg-go/v2/database/fields"
)

// APIKey represents API Key model.
type APIKey struct {
	State
	tableName struct{} `pg:"apikeys_apikey"` // nolint

	IsLive bool `pg:"_is_live"`

	ID          int
	Key         string
	InstanceID  int
	Instance    *Instance
	Options     fields.Hstore
	CreatedAt   fields.Time
	Description string
}

func (m *APIKey) String() string {
	return fmt.Sprintf("APIKey<ID=%d Instance=%d>", m.ID, m.InstanceID)
}

// VerboseName returns verbose name for model.
func (m *APIKey) VerboseName() string {
	return "API Key"
}
