package query

import (
	"github.com/go-pg/pg/orm"

	"github.com/Syncano/orion/app/models"
	"github.com/Syncano/orion/pkg/storage"
)

// InstanceIndicatorManager represents Instance Indicator manager.
type InstanceIndicatorManager struct {
	*Manager
}

// NewInstanceIndicatorManager creates and returns new Instance Indicator manager.
func NewInstanceIndicatorManager(c storage.DBContext) *InstanceIndicatorManager {
	return &InstanceIndicatorManager{Manager: NewTenantManager(c)}
}

// ByInstanceAndType filters object filtered by instance and type.
func (mgr *InstanceIndicatorManager) ByInstanceAndType(o *models.InstanceIndicator) *orm.Query {
	return mgr.Query(o).Where("instance_id = ?", o.InstanceID).
		Where("type = ?", o.Type)
}