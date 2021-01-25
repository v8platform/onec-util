package rac_addin

import (
	"errors"
	"fmt"
	"github.com/v8platform/rac"
	"reflect"
	"strconv"
	"strings"
)

type ServersCommandType string

func (c ServersCommandType) Check(params map[string]string) error {

	var err error

	switch c {

	case ServersRemoveCommand:

		if val, ok := params["--cluster"]; !ok || len(val) == 0 {
			err = errors.New("cluster must be identified")
		}

		if val, ok := params["--server"]; !ok || len(val) == 0 {
			err = errors.New("server must be identified")
		}

	case ServersInsertCommand, ServersListCommand:

		if val, ok := params["--cluster"]; !ok || len(val) == 0 {
			err = errors.New("cluster must be identified")
		}

	}

	return err
}

func (c ServersCommandType) Command() string {
	return string(c)
}

const (
	baseServersCommand ServersCommandType = "server"
	ServersListCommand                    = baseServersCommand + " list"
	//ClustersInfoCommand                      = baseServersCommand + " info"
	ServersInsertCommand = baseServersCommand + " insert"
	ServersRemoveCommand = baseServersCommand + " remove"
	//ClustersUpdateCommand                    = baseServersCommand + " update"
)

type ServersInsert struct {
	AgentHost                            string //agent-host                                : app
	AgentPort                            int    //agent-port                                : 1540
	PortRange                            string //port-range                                : 1560:1591
	Name                                 string //name                                      : "Центральный сервер"
	Using                                string //using                                     : main
	DedicateManagers                     string //dedicate-managers                         : none
	InfobasesLimit                       int64  //infobases-limit                           : 8
	MemoryLimit                          int64  //memory-limit                              : 0
	ConnectionsLimit                     int64  //connections-limit                         : 128
	SafeWorkingProcessesMemoryLimit      int64  //safe-working-processes-memory-limit       : 0
	SafeCallMemoryLimit                  int64  //safe-call-memory-limit                    : 0
	ClusterPort                          int    //cluster-port                              : 1541
	CriticalTotalMemory                  int64  //critical-total-memory                     : 0
	TemporaryAllowedTotalMemory          int64  //temporary-allowed-total-memory            : 0
	TemporaryAllowedTotalMemoryTimeLimit int64  //temporary-allowed-total-memory-time-limit : 300

}

type ServersList struct{}

func (_ ServersList) Command() rac.DoCommand {
	return ServersListCommand
}

func (i ServersList) Values() map[string]string {

	return map[string]string{}

}

func (i ServersList) Parse(res *rac.RawRespond) error {

	var list []rac.ServerInfo

	if !res.Status {
		return res.Error
	}

	err := rac.Unmarshal(res.Raw, &list)
	res.Error = err
	res.ParsedRespond = list

	return err

}

func (_ ServersInsert) Command() rac.DoCommand {
	return ServersInsertCommand
}

func (i ServersInsert) Values() map[string]string {

	val := map[string]string{}

	rv := reflect.ValueOf(&i)
	ri := reflect.Indirect(rv)

	rt := reflect.TypeOf(i)

	for i := 0; i < ri.NumField(); i++ {

		fieldName := rac.NameMapping(rt.Field(i).Name)

		tag := rt.Field(i).Tag.Get(rac.TagNamespace)

		tags := strings.Split(tag, ",")

		if len(tags) > 0 && len(tags[0]) > 0 {
			fieldName = tags[0]
		}

		if tags[0] == "-" {
			continue
		}

		paramName := "--" + fieldName

		value := ri.Field(i).Interface()

		switch v := value.(type) {
		case rac.YasNoBoolType:
			if v == rac.NullBool {
				continue
			}

			val[paramName] = v.String()

		case int:
			if v == 0 {
				continue
			}
			val[paramName] = fmt.Sprintf("%d", v)

		case bool:
			if !v {
				continue
			}
			val[paramName] = strconv.FormatBool(v)
		case string:
			if len(v) == 0 {
				continue
			}
			val[paramName] = v
		}

	}

	return val

}

func (i ServersInsert) Parse(res *rac.RawRespond) error {

	if !res.Status {
		return res.Error
	}

	return nil
}

type ServersRemove struct {
	UUID string
	rac.Auth
}

func (_ ServersRemove) Command() rac.DoCommand {
	return ServersRemoveCommand
}

func (i ServersRemove) Values() map[string]string {

	return map[string]string{
		"--server":       i.UUID,
		"--cluster-user": i.User,
		"--cluster-pwd":  i.Pwd,
	}

}

func (i ServersRemove) Parse(res *rac.RawRespond) error {

	if !res.Status {
		return res.Error
	}

	return nil
}

type ServersRespond struct {
	raw  *rac.RawRespond
	List []rac.ServerInfo
	Info rac.ServerInfo
}

func Servers(m *rac.Manager, what interface{}, opts ...interface{}) (respond ServersRespond, err error) {

	val, ok := what.(rac.Valued)

	if !ok {
		return respond, rac.ErrUnsupportedWhat
	}

	respond.raw, err = m.Do(val, opts...)

	if err != nil {
		return respond, err
	}

	switch v := respond.raw.Parsed().(type) {

	case rac.ServerInfo:
		respond.Info = v
		respond.List = append(respond.List, v)
	case []rac.ServerInfo:
		respond.List = v
		if len(v) == 1 {
			respond.Info = v[0]
		}

	}

	return respond, nil

}
