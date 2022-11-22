package iwd

import (
	"github.com/godbus/dbus/v5"
)

const (
	IWD_SERVICE                 = "net.connman.iwd"
	IWD_AGENT_INTERFACE         = "net.connman.iwd.Agent"
	IWD_AGENT_MANAGER_INTERFACE = "net.connman.iwd.AgentManager"
	IWD_DEVICE_INTERFACE        = "net.connman.iwd.Device"
	IWD_WSC_INTERFACE           = "net.connman.iwd.SimpleConfiguration"
	IWD_SIGNAL_AGENT_INTERFACE  = "net.connman.iwd.SignalLevelAgent"
	IWD_AP_INTERFACE            = "net.connman.iwd.AccessPoint"
	IWD_ADHOC_INTERFACE         = "net.connman.iwd.AdHoc"

	IWD_P2P_INTERFACE                 = "net.connman.iwd.p2p.Device"
	IWD_P2P_PEER_INTERFACE            = "net.connman.iwd.p2p.Peer"
	IWD_P2P_SERVICE_MANAGER_INTERFACE = "net.connman.iwd.p2p.ServiceManager"
	IWD_P2P_WFD_INTERFACE             = "net.connman.iwd.p2p.Display"
	IWD_STATION_DEBUG_INTERFACE       = "net.connman.iwd.StationDebug"
	IWD_DPP_INTERFACE                 = "net.connman.iwd.DeviceProvisioning"

	IWD_AGENT_MANAGER_PATH = "/net/connman/iwd"
	IWD_TOP_LEVEL_PATH     = "/"

	M_AGENT_PATH dbus.ObjectPath = "/site/zbyte/iwd/agent"
)

type Iwd struct {
	Conn          *dbus.Conn
	Stations      []*Station
	KnownNetworks []*KnownNetwork
	Adapters      []*Adapter
	Networks      []*Network
	// 临时保存密码位置
	passphrasse map[dbus.ObjectPath]string
}

func NewIwd() (*Iwd, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	obj := &Iwd{
		Conn:          conn,
		Stations:      []*Station{},
		KnownNetworks: []*KnownNetwork{},
		Adapters:      []*Adapter{},
		Networks:      []*Network{},
		passphrasse:   make(map[dbus.ObjectPath]string),
	}
	//get all remote info
	if err := obj.updateInfo(); err != nil {
		return nil, err
	}
	//set agent
	if err := obj.setupAgent(); err != nil {
		return nil, err
	}
	return obj, nil
}

// set agent
func (obj *Iwd) setupAgent() error {
	proxy := obj.Conn.Object(IWD_SERVICE, IWD_AGENT_MANAGER_PATH)
	err := obj.Conn.Export(obj, M_AGENT_PATH, IWD_AGENT_INTERFACE)
	if err != nil {
		return err
	}
	return proxy.Call(IWD_AGENT_MANAGER_INTERFACE+".RegisterAgent", 0, M_AGENT_PATH).Err
}

// get remote info
func (obj *Iwd) updateInfo() error {
	var objects map[dbus.ObjectPath]map[string]map[string]dbus.Variant

	objManager := obj.Conn.Object(IWD_SERVICE, "/")
	err := objManager.Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&objects)

	if err != nil {
		return err
	}
	// test print
	// for path, v := range objects {
	// 	println(path)
	// 	for interName, values := range v {
	// 		// switch interName {
	// 		// case IWD_STATION_INTERFACE:

	// 		// }
	// 		print(interName + "   ")
	// 		fmt.Println(values)
	// 	}
	// 	println("==========")
	// }
	// return

	for path, v := range objects {
		if _, has := v[IWD_NETWORK_INTERFACE]; has {
			//network
			network := &Network{
				iwd:  obj,
				Obj:  obj.Conn.Object(IWD_NETWORK_INTERFACE, path),
				Path: path,
				// Name: s["Name"].Value().(string),
				// Connected: s["Connected"].Value().(bool),
				// Type: s["Type"].Value().(string),
			}
			obj.Networks = append(obj.Networks, network)
		} else if _, has := v[IWD_AGENT_INTERFACE]; has {
			//iwd
		} else if s, has := v[IWD_WIPHY_INTERFACE]; has {
			//phy
			adapter := &Adapter{
				iwd:            obj,
				Obj:            obj.Conn.Object(IWD_WIPHY_INTERFACE, path),
				Path:           path,
				Powered:        s["Powered"].Value().(bool),
				Name:           s["Name"].Value().(string),
				SupportedModes: s["SupportedModes"].Value().([]string),
				// Model:          "",
				// Vendor:         "",
			}
			if model, has := s["Model"]; has {
				adapter.Model = model.Value().(string)
			}
			if vendor, has := s["Vendor"]; has {
				adapter.Vendor = vendor.Value().(string)
			}
			obj.Adapters = append(obj.Adapters, adapter)
		} else if _, has := v[IWD_STATION_INTERFACE]; has {
			//station
			station := &Station{
				iwd:  obj,
				Obj:  obj.Conn.Object(IWD_SERVICE, path),
				Path: path,
				// State:    s["State"].Value().(string),
				// Scanning: s["Scanning"].Value().(bool),
				// ConnectedNetworkPath: "",
				Name: v[IWD_DEVICE_INTERFACE]["Name"].Value().(string),
			}
			//network
			// if network, has := s["ConnectedNetwork"]; has {
			// 	station.ConnectedNetworkPath = network.Value().(dbus.ObjectPath)
			// }
			obj.Stations = append(obj.Stations, station)
		} else if s, has := v[IWD_KNOWN_NETWORK_INTERFACE]; has {
			//known_network
			knownNetwork := &KnownNetwork{
				iwd:         obj,
				Obj:         obj.Conn.Object(IWD_SERVICE, path),
				Path:        path,
				Name:        s["Name"].Value().(string),
				Type:        s["Type"].Value().(string),
				Hidden:      s["Hidden"].Value().(bool),
				AutoConnect: s["AutoConnect"].Value().(bool),
				// LastConnectedTime: "",
			}
			if time, has := s["LastConnectedTime"]; has {
				knownNetwork.LastConnectedTime = time.Value().(string)
			}
			obj.KnownNetworks = append(obj.KnownNetworks, knownNetwork)
		}
	}
	return nil
}

// when network connect. save passphrase for [RequestPassphrase]
func (obj *Iwd) UpdatePassphrase(path dbus.ObjectPath, passphrasse string) {
	obj.passphrasse[path] = passphrasse
}

// =====================for remote
func (obj *Iwd) Release() *dbus.Error {
	println("Release")
	return nil
}
func (obj *Iwd) RequestPassphrase(path dbus.ObjectPath) (string, *dbus.Error) {
	if val, ok := obj.passphrasse[path]; ok {
		return val, nil
	}
	return "", dbus.NewError("Passphrase not set", nil)
}
func (obj *Iwd) RequestPrivateKeyPassphrase(path dbus.ObjectPath) (string, *dbus.Error) {
	println("RequestPrivateKeyPassphrase:" + path)
	return "", nil
}
func (obj *Iwd) RequestUserNameAndPassword(path dbus.ObjectPath) (string, string, *dbus.Error) {
	println("RequestUserNameAndPassword:" + path)
	return "", "", nil
}
func (obj *Iwd) RequestUserPassword(path dbus.ObjectPath, user string) (string, *dbus.Error) {
	println("RequestUserPassword:" + string(path) + ":" + user)
	return "", nil
}
func (obj *Iwd) Cancel(reason string) *dbus.Error {
	println("Cancel:" + reason)
	return nil
}
