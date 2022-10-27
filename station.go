package iwd

import (
	"github.com/godbus/dbus/v5"
)

const (
	IWD_STATION_INTERFACE = "net.connman.iwd.Station"
	IWD_STATION_SCAN      = "net.connman.iwd.Station.Scan"
)

type Station struct {
	obj                  dbus.BusObject
	Path                 dbus.ObjectPath
	State                string
	Scanning             bool
	ConnectedNetworkPath dbus.ObjectPath
	Name                 string
}

func (obj *Station) Scan() error {
	println(obj.Path)
	return obj.obj.Call("net.connman.iwd.Station.Scan", 0).Err
}
