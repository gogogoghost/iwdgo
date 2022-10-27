package iwd

import (
	"github.com/godbus/dbus/v5"
)

const (
	IWD_NETWORK_INTERFACE = "net.connman.iwd.Network"
)

type Network struct {
	obj       dbus.BusObject
	Path      dbus.ObjectPath
	Name      string
	Connected bool
	Type      string
}
