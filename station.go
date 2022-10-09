package iwd

import (
	"github.com/godbus/dbus/v5"
)

type Station struct {
	Path                 dbus.ObjectPath
	State                string
	Scanning             bool
	ConnectedNetworkPath dbus.ObjectPath
	Name                 string
}
