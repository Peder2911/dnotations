
package models

import (
	"fmt"
)

type DnotatedUnit struct {
	PartOf    string `yaml:"part-of"    json:"part-of"`
	Component string `yaml:"component"  json:"component"`
	Version   string `yaml:"version"    json:"version"`
	ManagedBy string `yaml:"managed-by" json:"managed-by"`
	Hostname  string `yaml:"hostname"   json:"hostname"`
}

func (u DnotatedUnit) String() string {
	return fmt.Sprintf(`
	Part of:    %s
	Component:  %s
	Version:    %s
	Managed by: %s
	Hostname:   %s
	---
	`, u.PartOf,u.Component,u.Version,u.ManagedBy,u.Hostname)
}
