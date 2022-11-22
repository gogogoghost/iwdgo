package iwd

import (
	"time"

	"github.com/godbus/dbus/v5"
)

const (
	IWD_STATION_INTERFACE = "net.connman.iwd.Station"
	IWD_DEVICE_INTERFACE  = "net.connman.iwd.Device"
)

type Station struct {
	*IwdSub
}

func NewStation(
	iwd *Iwd,
	path dbus.ObjectPath,
	props *map[string]dbus.Variant,
) (*Station, error) {
	for _, o := range iwd.Stations {
		if o.Obj.Path() == path {
			return o, nil
		}
	}
	if props == nil {
		// 自行读取所有prop
		p, err := getAllProps(iwd.Conn, IWD_STATION_INTERFACE, path)
		if err != nil {
			return nil, err
		}
		props = &p
		// 再读取一下它的device
		p2, err := getAllProps(iwd.Conn, IWD_DEVICE_INTERFACE, path)
		if err != nil {
			for k, v := range p2 {
				(*props)[k] = v
			}
		}
	}
	setupChangeListener(iwd.Conn, path, props)
	o := &Station{
		IwdSub: &IwdSub{
			iwd:   iwd,
			Obj:   iwd.Conn.Object(IWD_SERVICE, path),
			Props: props,
		},
	}
	iwd.Stations = append(iwd.Stations, o)
	return o, nil
}

type OrderedNetwork struct {
	Network
	SignalStrength int16
}

//======================props for station

func (obj *Station) Scanning() bool {
	return obj.Get("Scanning").(bool)
}

func (obj *Station) State() string {
	return obj.Get("State").(string)
}

func (obj *Station) ConnectedNetwork() (*Network, error) {
	path := obj.Get("ConnectedNetwork")
	if path == nil {
		return nil, nil
	}
	return NewNetwork(obj.iwd, path.(dbus.ObjectPath), nil)
}

//======================props for device

func (obj *Station) Name() string {
	return obj.Get("Name").(string)
}

func (obj *Station) Address() string {
	return obj.Get("Address").(string)
}

func (obj *Station) Power() bool {
	return obj.Get("Power").(bool)
}

func (obj *Station) Mode() string {
	return obj.Get("Mode").(string)
}

func (obj *Station) Adapter() (*Adapter, error) {
	path := obj.Get("Adapter").(dbus.ObjectPath)
	return NewAdapter(obj.iwd, path, nil)
}

//======================methods

func (obj *Station) Scan() error {
	return obj.Obj.Call(IWD_STATION_INTERFACE+".Scan", 0).Err
}

func (obj *Station) GetOrderedNetworks() ([]*OrderedNetwork, error) {
	res := []*OrderedNetwork{}
	list := [][]dbus.Variant{}
	err := obj.Obj.Call(IWD_STATION_INTERFACE+".GetOrderedNetworks", 0).Store(&list)
	if err != nil {
		return nil, err
	}
	for _, item := range list {
		path := item[0].Value().(dbus.ObjectPath)
		network, err := NewNetwork(
			obj.iwd,
			path,
			nil)
		if err != nil {
			continue
		}
		res = append(res, &OrderedNetwork{
			Network:        *network,
			SignalStrength: item[1].Value().(int16),
		})
	}
	return res, nil
}

func (obj *Station) Disconnect() error {
	return obj.Obj.Call(IWD_STATION_INTERFACE+".Disconnect", 0).Err
}

//======================ext

func (obj *Station) WaitForScan() {
	for {
		isScanning := obj.Scanning()
		if isScanning {
			time.Sleep(time.Second * 1)
		} else {
			return
		}
	}
}
