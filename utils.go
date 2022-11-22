package iwd

import (
	"github.com/godbus/dbus/v5"
)

var listenerMap = make(map[dbus.ObjectPath](func(*dbus.Signal)))

// init listener
func initChangeListener(conn *dbus.Conn) {
	c := make(chan *dbus.Signal, 10)
	conn.Signal(c)
	go func() {
		for v := range c {
			if f, has := listenerMap[v.Path]; has {
				f(v)
			}
		}
	}()
}

// change listener
func setupChangeListener(
	conn *dbus.Conn,
	path dbus.ObjectPath,
	props *map[string]dbus.Variant,
) error {
	if err := conn.AddMatchSignal(
		dbus.WithMatchObjectPath(path),
		dbus.WithMatchInterface(DBUS_PROPERTIES_INTERFACE),
	); err != nil {
		return err
	}
	listenerMap[path] = func(s *dbus.Signal) {
		for k, v := range s.Body[1].(map[string]dbus.Variant) {
			(*props)[k] = v
		}
	}
	return nil
}

// 获取所有props
func getAllProps(conn *dbus.Conn, inter string, path dbus.ObjectPath) (map[string]dbus.Variant, error) {
	obj := conn.Object(IWD_SERVICE, path)
	var props map[string]dbus.Variant
	if err := obj.Call(DBUS_PROPERTIES_INTERFACE+".GetAll", 0, inter).Store(&props); err != nil {
		return nil, err
	}
	return props, nil
}
