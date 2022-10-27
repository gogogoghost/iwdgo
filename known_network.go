package iwd

import (
	"github.com/godbus/dbus/v5"
)

const (
	IWD_KNOWN_NETWORK_INTERFACE = "net.connman.iwd.KnownNetwork"
)

type KnownNetwork struct {
	obj               dbus.BusObject
	Path              dbus.ObjectPath
	Name              string
	Type              string
	Hidden            bool
	AutoConnect       bool
	LastConnectedTime string
}
