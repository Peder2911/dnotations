/* dnotations

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
	"github.com/go-yaml/yaml"
	"github.com/godbus/dbus/v5"
	"log"

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

type UnitMetadata struct {
	Annotations map[string]string `yaml:"annotations"`
}

func (u UnitMetadata) String() string {
	kv := make([]string, len(u.Annotations))
	var i int
	for k,v := range u.Annotations{
		kv[i] = fmt.Sprintf("%s:\t%s",k,v)
		i++
	}
	return "Metadata:"+strings.Join(kv,",")+"\n"
}

type Unit struct {
	Name string
	Metadata UnitMetadata
}

type Units []Unit

func (u UnitStatus) String() string {
	return u.Name
}

func ListUnits(ctx context.Context, conn *dbus.Conn)(*[]UnitStatus, error){
	systemd := conn.Object("org.freedesktop.systemd1","/org/freedesktop/systemd1")
	var result [][]interface{}
	err := systemd.CallWithContext(ctx,"org.freedesktop.systemd1.Manager.ListUnits", 0).Store(&result)
	if err != nil {
		return nil, err
	}
	result_interface  := make([]interface{}, len(result))
	for i := range result {
		result_interface[i] = result[i]
	}

	unit_status := make([]UnitStatus, len(result))
	unit_status_interface := make([]interface{}, len(unit_status))
	for i := range unit_status{ 
		unit_status_interface[i] = &unit_status[i]
	}
	err = dbus.Store(result_interface, unit_status_interface...)
	if err != nil {
		return nil, err
	}
	return &unit_status, nil
}

func ParseUnitHeader(data []byte) (UnitMetadata, error) {
	var yaml_header_lines []string
	lines := strings.Split(string(data),"\n")

	var in_header bool = true
	var i int = 0
	var line string
	for in_header {
		line = lines[i]
		if line[0] == '#' {
			yaml_header_lines = append(yaml_header_lines, strings.TrimLeft(line, "#") )
		} else {
			in_header = false
		}
		i++
	}
	var unit_metadata UnitMetadata
	yaml.Unmarshal([]byte(strings.Join(yaml_header_lines, "\n")), &unit_metadata)
	return unit_metadata, nil

}

func LoadUnit(name string) (*Unit,error) {
	data,err := os.ReadFile(fmt.Sprintf("/etc/systemd/system/%s",name))
	if err != nil {
		return nil,err
	}

	unit_header,err := ParseUnitHeader(data)
	if err != nil {
		log.Println(fmt.Sprintf("Failed to parse header for unit %s: %s",name,err))
	}

	var unit Unit
	unit.Name = name
	unit.Metadata = unit_header
	return &unit,nil
}

func LoadUnits() (Units,error) {
	entries, err := os.ReadDir("/etc/systemd/system")
	if err != nil {
		return nil, err
	}
	var units Units
	for _,e := range entries {
		if e.Type().IsRegular() {
			unit,err := LoadUnit(e.Name())
			if err != nil {
				return nil, err
			}
			units = append(units, *unit)
		}
	} 
	return units, nil
}

func main() {
	units,err := LoadUnits()
	if err != nil {
		panic(err)
	}
	for _,unit := range units {
		println(unit.Name)
		println("\tannotations:")
		for k,v := range unit.Metadata.Annotations {
			println(fmt.Sprintf("\t\t%s:\t%s", k,v))
		}
	}
}
