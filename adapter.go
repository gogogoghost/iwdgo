package iwd

import "github.com/godbus/dbus/v5"

const (
	IWD_WIPHY_INTERFACE = "net.connman.iwd.Adapter"
)

type Adapter struct {
	*IwdSub
}

func NewAdapter(
	iwd *Iwd,
	path dbus.ObjectPath,
	props *map[string]dbus.Variant,
) (*Adapter, error) {
	for _, o := range iwd.Adapters {
		if o.Obj.Path() == path {
			return o, nil
		}
	}
	if props == nil {
		// 自行读取所有prop
		p, err := getAllProps(iwd.Conn, IWD_WIPHY_INTERFACE, path)
		if err != nil {
			return nil, err
		}
		props = &p
	}
	setupChangeListener(iwd.Conn, path, props)
	o := &Adapter{
		IwdSub: &IwdSub{
			iwd:   iwd,
			Obj:   iwd.Conn.Object(IWD_SERVICE, path),
			Props: props,
		},
	}
	iwd.Adapters = append(iwd.Adapters, o)
	return o, nil
}

//======================props

func (obj *Adapter) Powered() bool {
	return obj.Get("Powered").(bool)
}

func (obj *Adapter) Name() string {
	return obj.Get("Name").(string)
}

func (obj *Adapter) Model() string {
	m := obj.Get("Model")
	if m == nil {
		return ""
	}
	return m.(string)
}

func (obj *Adapter) Vendor() string {
	m := obj.Get("Vendor")
	if m == nil {
		return ""
	}
	return m.(string)
}

func (obj *Adapter) SupportedModes() []string {
	return obj.Get("SupportedModes").([]string)
}

//======================methods
