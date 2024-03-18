
package cli

import (
	"github.com/peder2911/dnotations/pkg/models"
)

type DnotatedUnitsListing struct {
	Units []models.DnotatedUnit `yaml:"units" json:"units"`
}
