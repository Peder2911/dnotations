
package io


import (
	"github.com/peder2911/dnotations/pkg/models"
	"github.com/godbus/dbus/v5"
	"context"
	"fmt"
	"sync"
)


type UnitFile struct {
	Path   string
	Status string
}

func (u UnitFile) String() string{
	return fmt.Sprintf("Path: %s\tStatus %s", u.Path, u.Status)
}

type DnotatedSystemd struct {
	conn *dbus.Conn
}

func NewDnotatedSystemd() (*DnotatedSystemd,error) {
	conn,err := dbus.ConnectSystemBus()
	if err != nil {
		return nil,err
	}
	return &DnotatedSystemd{conn},nil
}

func (s *DnotatedSystemd) object() dbus.BusObject {
	return s.conn.Object("org.freedesktop.systemd1","/org/freedesktop/systemd1")
}

func (s *DnotatedSystemd) mancall(ctx context.Context,method string, retvalues *[][]interface{}) error {
	obj := s.object()
	call := obj.CallWithContext(ctx, fmt.Sprintf("org.freedesktop.systemd1.Manager.%s",method),0)
	err := call.Store(retvalues)
	return err
}

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

func (s *DnotatedSystemd) Close() {
	s.conn.Close()
}

func (s *DnotatedSystemd) ListUnits(ctx context.Context) (*[]models.DnotatedUnit, error) {
	files,err := s.listUnitFiles(ctx)
	if err != nil {
		return nil,err
	}
	var units []models.DnotatedUnit
	c := make(chan *models.DnotatedUnit, len(*files))
	var wg sync.WaitGroup
	for _,f := range *files {
		wg.Add(1)
		go func(path string, out chan *models.DnotatedUnit) {
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
