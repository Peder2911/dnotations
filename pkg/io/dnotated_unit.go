/* Load dnotated units.

*/
package io

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/peder2911/dnotations/pkg/models"
	"os"
	"strings"
)


func LoadDnotatedUnit(path string) (*models.DnotatedUnit, error){
	data,err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	header,err := ParseUnitHeader(data)
	if err != nil {
		return nil,err
	}

	return &header.Annotations,nil
}


type DnotatedUnitHeader struct {
	Annotations models.DnotatedUnit `yaml:"annotations"`
}


func ParseUnitHeader(data []byte) (*DnotatedUnitHeader, error) {
	var yaml_header_lines []string
	var header DnotatedUnitHeader
	lines := strings.Split(string(data),"\n")

	var in_header bool = true
	var i int = 0
	var line string
	if len(lines) == 0 {
		return nil,fmt.Errorf("Unit file had length 0")
	}
	for in_header {
		line = lines[i]
		if len(line) > 0 && line[0] == '#' {
			yaml_header_lines = append(yaml_header_lines, strings.TrimLeft(line, "#") )
		} else {
			in_header = false
		}

		i++
		if i > len(lines) {
			break
		}
	}
	if len(yaml_header_lines) == 0 {
		return nil, fmt.Errorf("Found no header in file.\n")
	}
	err := yaml.Unmarshal([]byte(strings.Join(yaml_header_lines, "\n")), &header)
	if err != nil {
		return nil, err
	}
	return &header, nil
}
