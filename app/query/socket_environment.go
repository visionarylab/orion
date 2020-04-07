package query

import (
	"fmt"

	"github.com/go-pg/pg/v9/orm"

	"github.com/Syncano/orion/app/models"
	"github.com/Syncano/orion/pkg/cache"
	"github.com/Syncano/orion/pkg/storage"
)

// SocketEnvironmentManager represents Socket Environment manager.
type SocketEnvironmentManager struct {
	*LiveManager
}

// NewSocketEnvironmentManager creates and returns new Socket Environment manager.
func NewSocketEnvironmentManager(c storage.DBContext) *SocketEnvironmentManager {
	return &SocketEnvironmentManager{LiveManager: NewLiveTenantManager(c)}
}

// OneByID outputs object filtered by ID.
func (m *SocketEnvironmentManager) OneByID(o *models.SocketEnvironment) error {
	return RequireOne(
		cache.SimpleModelCache(m.DB(), o, fmt.Sprintf("i=%d", o.ID), func() (interface{}, error) {
			return o, m.Query(o).Where("id = ?", o.ID).Select()
		}),
	)
}

// WithAccessQ outputs objects that entity has access to.
func (m *SocketEnvironmentManager) WithAccessQ(o interface{}) *orm.Query {
	return m.Query(o)
}
