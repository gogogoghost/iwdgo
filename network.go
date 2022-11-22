package iwd

import (
	"time"

	"github.com/godbus/dbus/v5"
)

const (
	IWD_NETWORK_INTERFACE = "net.connman.iwd.Network"
)

type Network struct {
	*IwdSub
}

func NewNetwork(
	iwd *Iwd,
	path dbus.ObjectPath,
	props *map[string]dbus.Variant,
) (*Network, error) {
	for _, o := range iwd.Networks {
		if o.Obj.Path() == path {
			return o, nil
		}
	}
	if props == nil {
		// 自行读取所有prop
		p, err := getAllProps(iwd.Conn, IWD_NETWORK_INTERFACE, path)
		if err != nil {
			return nil, err
		}
		props = &p
	}
	setupChangeListener(iwd.Conn, path, props)
	o := &Network{
		IwdSub: &IwdSub{
			iwd:   iwd,
			Obj:   iwd.Conn.Object(IWD_SERVICE, path),
			Props: props,
		},
	}
	iwd.Networks = append(iwd.Networks, o)
	return o, nil
}

//======================props

func (obj *Network) Name() string {
	return obj.Get("Name").(string)
}

func (obj *Network) Type() string {
	return obj.Get("Type").(string)
}

func (obj *Network) Connected() bool {
	return obj.Get("Connected").(bool)
}

func (obj *Network) KnownNetwork() (*KnownNetwork, error) {
	path := obj.Get("KnownNetwork")
	if path == nil {
		return nil, nil
	}
	return NewKnownNetwork(obj.iwd, path.(dbus.ObjectPath), nil)
}

//======================methods

func (obj *Network) Connect(passphrase string) error {
	obj.iwd.UpdatePassphrase(obj.Obj.Path(), passphrase)
	return obj.Obj.Call(IWD_NETWORK_INTERFACE+".Connect", 0).Err
}

//======================ext

func (obj *Network) WaitForConnected() error {
	for {
		connected := obj.Connected()
		if connected {
			return nil
		} else {
			time.Sleep(time.Second * 1)
		}
	}
}
