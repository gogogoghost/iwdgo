package iwd

import "github.com/godbus/dbus/v5"

const (
	IWD_KNOWN_NETWORK_INTERFACE = "net.connman.iwd.KnownNetwork"
)

type KnownNetwork struct {
	*IwdSub
}

func NewKnownNetwork(
	iwd *Iwd,
	path dbus.ObjectPath,
	props *map[string]dbus.Variant,
) (*KnownNetwork, error) {
	for _, o := range iwd.KnownNetworks {
		if o.Obj.Path() == path {
			return o, nil
		}
	}
	if props == nil {
		// 自行读取所有prop
		p, err := getAllProps(iwd.Conn, IWD_KNOWN_NETWORK_INTERFACE, path)
		if err != nil {
			return nil, err
		}
		props = &p
	}
	setupChangeListener(iwd.Conn, path, props)
	o := &KnownNetwork{
		IwdSub: &IwdSub{
			iwd:   iwd,
			Obj:   iwd.Conn.Object(IWD_SERVICE, path),
			Props: props,
		},
	}
	iwd.KnownNetworks = append(iwd.KnownNetworks, o)
	return o, nil
}

//======================methods

func (obj *KnownNetwork) Forget() error {
	return obj.Obj.Call(IWD_KNOWN_NETWORK_INTERFACE+".Forget", 0).Err
}
