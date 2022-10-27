package iwd

import (
	"github.com/godbus/dbus/v5"
)

const (
	IWD_WIPHY_INTERFACE = "net.connman.iwd.Adapter"
)

type Adapter struct {
	obj            dbus.BusObject
	Path           dbus.ObjectPath
	Powered        bool
	Name           string
	SupportedModes []string
	Model          string
	Vendor         string
}
