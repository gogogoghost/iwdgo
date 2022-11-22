package iwd

import (
	"time"

	"github.com/godbus/dbus/v5"
)

const (
	IWD_NETWORK_INTERFACE = "net.connman.iwd.Network"
)

type Network struct {
	iwd  *Iwd
	Obj  dbus.BusObject
	Path dbus.ObjectPath
	// Name string
	// Connected bool
	// Type string
}

func (obj *Network) GetName() (string, error) {
	v, err := obj.Obj.GetProperty(IWD_NETWORK_INTERFACE + ".Name")
	if err != nil {
		return "", err
	}
	return v.Value().(string), nil
}

func (obj *Network) GetType() (string, error) {
	v, err := obj.Obj.GetProperty(IWD_NETWORK_INTERFACE + ".Type")
	if err != nil {
		return "", err
	}
	return v.Value().(string), nil
}

func (obj *Network) GetConnected() (bool, error) {
	v, err := obj.Obj.GetProperty(IWD_NETWORK_INTERFACE + ".Connected")
	if err != nil {
		return false, err
	}
	return v.Value().(bool), nil
}

func (obj *Network) WaitForConnected() error {
	for {
		connected, err := obj.GetConnected()
		if err != nil {
			return err
		}
		if connected {
			return nil
		} else {
			time.Sleep(time.Second * 1)
		}
	}
}

func (obj *Network) Connect(passphrase string) error {
	obj.iwd.UpdatePassphrase(obj.Path, passphrase)
	return obj.Obj.Call(IWD_NETWORK_INTERFACE+".Connect", 0).Err
}
