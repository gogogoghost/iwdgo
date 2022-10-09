package iwd

import (
	"github.com/godbus/dbus/v5"
)

type Network struct {
	Path      dbus.ObjectPath
	Name      string
	Connected bool
	Type      string
}
