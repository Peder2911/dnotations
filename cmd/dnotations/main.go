/*
	dnotations

This is a (wip) tool for people who would like to organize their Systemd units with metadata and annotations.
This might be helpful to keep track of things like what services provide what functionality on a box, behind
which hostnames, and so on.

You do this by annotating your SystemD unit files with a comment header with some YAML data, like so:

```{example.service}
# annotations:
#   part-of: my-project
#   component: proxy
#   managed-by: someuser
#   version: 1.0.0
#   hostname: myservice.example.com
[Unit]
Description=A proxy server.
...
```
*/
package main

import (
	//"log"
	"log"
	"sync"

	"github.com/go-yaml/yaml"
	"github.com/godbus/dbus/v5"

	//"github.com/go-yaml/yaml"
	"context"
	"fmt"
	"os"
	"strings"
)

type UnitStatus struct {
	Name        string          // The primary unit name as string
	Description string          // The human readable description string
	LoadState   string          // The load state (i.e. whether the unit file has been loaded successfully)
	ActiveState string          // The active state (i.e. whether the unit is currently started or not)
	SubState    string          // The sub state (a more fine-grained version of the active state that is specific to the unit type, which the active state is not)
	Followed    string          // A unit that is being followed in its state by this unit, if there is any, otherwise the empty string.
	Path        dbus.ObjectPath // The unit object path
	JobId       uint32          // If there is a job queued for the job unit the numeric job id, 0 otherwise
	JobType     string          // The job type as string
	JobPath     dbus.ObjectPath // The job object path
}

func (u UnitStatus) String() string {
	return u.Name
}

type UnitFile struct {
	Path   string
	Status string
}

func (u UnitFile) String() string{
	return fmt.Sprintf("Path: %s\tStatus %s", u.Path, u.Status)
}

type DnotatedSystemd struct {
	Conn *dbus.Conn
}

func (s *DnotatedSystemd) object() dbus.BusObject {
	return s.Conn.Object("org.freedesktop.systemd1","/org/freedesktop/systemd1")
}

func (s *DnotatedSystemd) mancall(ctx context.Context,method string, retvalues *[][]interface{}) error {
	obj := s.object()
	call := obj.CallWithContext(ctx, fmt.Sprintf("org.freedesktop.systemd1.Manager.%s",method),0)
	err := call.Store(retvalues)
	return err
}

// func (s *DnotatedSystemd) ListUnits(ctx context.Context)(*[]UnitStatus, error){
// 	var result[][]interface{}
// 	err := s.mancall(ctx, "ListUnits",&result)
// 	if err != nil {
// 		return nil, err
// 	}
// 	units := make([]UnitStatus, len(result))
// 	for i,r := range result {
// 		units[i] = UnitStatus{
// 			r[0].(string),
// 			r[1].(string),
// 			r[2].(string),
// 			r[3].(string),
// 			r[4].(string),
// 			r[5].(string),
// 			r[6].(dbus.ObjectPath),
// 			r[7].(uint32),
// 			r[8].(string),
// 			r[9].(dbus.ObjectPath),
// 		}
// 	}
// 	return &units,err
// }

func (s *DnotatedSystemd) listUnitFiles(ctx context.Context)(*[]UnitFile, error){
	var result[][]interface{}
	err := s.mancall(ctx, "ListUnitFiles",&result)
	if err != nil {
		return nil, err
	}
	units := make([]UnitFile, len(result))
	for i,r := range result {
		units[i] = UnitFile{
			r[0].(string),
			r[1].(string),
		}
	}
	return &units,err
}

type DnotatedUnitHeader struct {
	Annotations DnotatedUnit `yaml:"annotations"`
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

func LoadDnotatedUnit(path string) (*DnotatedUnit, error){
	data,err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	header,err := ParseUnitHeader(data)
	if err != nil {
		return nil,err
	}

	log.Printf("Successfully loaded %s\n", path)
	return &header.Annotations,nil
}

func (s *DnotatedSystemd) ListUnits(ctx context.Context) (*[]DnotatedUnit, error) {
	files,err := s.listUnitFiles(ctx)
	if err != nil {
		return nil,err
	}
	var units []DnotatedUnit
	c := make(chan *DnotatedUnit, len(*files))
	var wg sync.WaitGroup
	for _,f := range *files {
		wg.Add(1)
		go func(path string, out chan *DnotatedUnit) {
			defer wg.Done()
			unit, err := LoadDnotatedUnit(path)
			if err != nil {
				out <- nil
				return
			}
			out <- unit
		}(f.Path, c)
	}
	wg.Wait()
	for i:=0;i<len(*files);i++{
		unit := <- c
		if unit != nil {
			units = append(units,*unit)
		}
	}
	return &units,nil
}

func main() {
	ctx := context.Background()
	conn,err := dbus.ConnectSystemBus()
	if err != nil {panic(err)}
	defer conn.Close()
	s := DnotatedSystemd{conn}

	units,err := s.ListUnits(ctx)
	log.Printf("Loaded %v units\n",len(*units))
	for _,u := range *units {
		fmt.Println(u)
	}
	if err != nil {
		panic(err)
	}
}
