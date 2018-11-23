package serializers

import (
	"strconv"
	"time"

	"github.com/jackc/pgx/pgtype"
	geom "github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkbhex"

	"github.com/Syncano/orion/app/models"
	"github.com/Syncano/orion/pkg/settings"
	"github.com/Syncano/orion/pkg/util"
)

// DataObjectSerializer ...
type DataObjectSerializer struct {
	Class *models.Class
}

// Response ...
func (s DataObjectSerializer) Response(i interface{}) interface{} {
	o := i.(*models.DataObject)
	base := map[string]interface{}{
		"id":         o.ID,
		"created_at": &o.CreatedAt,
		"updated_at": &o.UpdatedAt,
		"revision":   o.Revision,
	}

	processDataObjectFields(s.Class, o, base)
	return base
}

func processDataObjectFields(class *models.Class, o *models.DataObject, m map[string]interface{}) {
	// Serialize hstore fields.
	var data map[string]pgtype.Text
	var val interface{}
	if !o.Data.IsNull() {
		data = o.Data.Map
	}

	for _, field := range class.ComputedSchema() {
		val = nil
		if v, ok := data[field.Mapping]; ok && v.Status != pgtype.Null {
			val = dataObjectFieldResponse(field, v.String)
		}
		m[field.FName] = val
	}
}

// dataObjectFieldResponse returns representation structure of value (for JSON serialization).
// nolint: gocyclo
func dataObjectFieldResponse(f *models.DataObjectField, val string) interface{} {
	// If field is string or text type - return as is.
	if f.FType == models.FieldStringType || f.FType == models.FieldTextType {
		return val
	} else if val == "" {
		// If non-string, non-text field is empty - return nil.
		return nil
	}

	switch f.FType {
	case models.FieldIntegerType:
		if v, err := strconv.Atoi(val); err == nil {
			return v
		}

	case models.FieldFloatType:
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			return v
		}

	case models.FieldBooleanType:
		return util.IsTrue(val)

	case models.FieldDatetimeType:
		if v, err := f.FromString(val); err == nil {
			t := v.(time.Time)

			return struct {
				Type  string      `json:"type"`
				Value models.Time `json:"value"`
			}{
				Type:  f.FType,
				Value: models.NewTime(&t),
			}
		}

	case models.FieldFileType:
		return struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		}{
			Type:  f.FType,
			Value: settings.API.MediaPrefix + val,
		}

	case models.FieldReferenceType:
		if v, err := f.FromString(val); err == nil {
			return struct {
				Type   string `json:"type"`
				Target string `json:"target"`
				Value  int    `json:"value"`
			}{
				Type:   f.FType,
				Target: f.Target,
				Value:  v.(int),
			}
		}

	case models.FieldRelationType:
		if v, err := f.FromString(val); err == nil {
			return struct {
				Type   string `json:"type"`
				Target string `json:"target"`
				Value  []int  `json:"value"`
			}{
				Type:   f.FType,
				Target: f.Target,
				Value:  v.([]int),
			}
		}

	case models.FieldObjectType:
		fallthrough
	case models.FieldArrayType:
		if v, err := f.FromString(val); err == nil {
			return v
		}

	case models.FieldGeopointType:
		if g, err := ewkbhex.Decode(val); err == nil {
			p := g.(*geom.Point)
			return struct {
				Type      string  `json:"type"`
				Longitude float64 `json:"longitude"`
				Latitude  float64 `json:"latitude"`
			}{
				Type:      f.FType,
				Longitude: p.X(),
				Latitude:  p.Y(),
			}
		}
	}
	return val
}