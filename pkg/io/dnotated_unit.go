// Load and parse dnotated units.
package io

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/peder2911/dnotations/pkg/models"
	"os"
	"strings"
)

// Load a dnotated unit from a path by parsing the comment-
// header in the file, if present.
func LoadDnotatedUnit(path string) (*models.DnotatedUnit, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	header, err := parseUnitHeader(data)
	if err != nil {
		return nil, err
	}
	return &header.Annotations, nil
}

// Parse out all header documents (separated by ---).
func headerDocuments(data string) ([][]string){
	var in_header bool
	var documents [][]string = make([][]string,1)
	var doc int
	for _,line := range strings.Split(data, "\n") {
		in_header = len(line) > 0 && line[0] == '#'
		if ! in_header {
			break
		}
		if line == "#---" {
			documents = append(documents,make([]string,0))
			continue
		}
		doc = len(documents)-1
		documents[doc] = append(documents[doc],strings.TrimLeft(line, "#"))
	}
	return documents
}

// Parse out a unit header if it exists in the comment-header of the provided data.
func parseUnitHeader(data []byte) (*dnotatedUnitHeader, error) {
	docs := headerDocuments(string(data))
	var header dnotatedUnitHeader
	for _,doc := range docs {
		if len(doc) == 0 {
			continue
		}
		err := yaml.Unmarshal([]byte(strings.Join(doc,"\n")), &header)
		if err == nil {
			return &header, nil
		}
	}
	return nil, fmt.Errorf("No unit header found in file.")
}

// Expected YAML data in the comment header.
type dnotatedUnitHeader struct {
	Annotations models.DnotatedUnit `yaml:"annotations"`
}
