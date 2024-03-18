
package models

import (
	"fmt"
)

type DnotatedUnit struct {
	PartOf    string `yaml:"part-of"`
	Component string `yaml:"component"`
	Version   string `yaml:"version"`
	ManagedBy string `yaml:"managed-by"`
	Hostname  string `yaml:"hostname"`
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
