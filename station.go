package iwd

import (
	"time"

	"github.com/godbus/dbus/v5"
)

const (
	IWD_STATION_INTERFACE = "net.connman.iwd.Station"
)

type Station struct {
	iwd  *Iwd
	Obj  dbus.BusObject
	Path dbus.ObjectPath
	// State                string
	// Scanning             bool
	// ConnectedNetworkPath dbus.ObjectPath
	Name string
}

type OrderedNetwork struct {
	Network
	SignalStrength int16
}

func (obj *Station) Scan() error {
	return obj.Obj.Call(IWD_STATION_INTERFACE+".Scan", 0).Err
}

func (obj *Station) WaitForScan() error {
	for {
		isScanning, err := obj.IsScanning()
		if err != nil {
			return err
		}
		if isScanning {
			time.Sleep(time.Second * 1)
		} else {
			return nil
		}
	}
}

func (obj *Station) IsScanning() (bool, error) {
	v, err := obj.Obj.GetProperty(IWD_STATION_INTERFACE + ".Scanning")
	if err != nil {
		return false, err
	}
	return v.Value().(bool), nil
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
		res = append(res, &OrderedNetwork{
			Network: Network{
				iwd:  obj.iwd,
				Obj:  obj.iwd.Conn.Object(IWD_SERVICE, path),
				Path: path,
			},
			SignalStrength: item[1].Value().(int16),
		})
	}
	return res, nil
}
