package iwd

import (
	"github.com/godbus/dbus/v5"
)

type KnownNetwork struct {
	Path              dbus.ObjectPath
	Name              string
	Type              string
	Hidden            bool
	AutoConnect       bool
	LastConnectedTime string
}
