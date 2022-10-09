package iwd

import (
	"github.com/godbus/dbus/v5"
)

type Adapter struct {
	Path           dbus.ObjectPath
	Powered        bool
	Name           string
	SupportedModes []string
	Model          string
	Vendor         string
}
