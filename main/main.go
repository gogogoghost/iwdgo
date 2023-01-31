package main

import (
	"fmt"
	"time"

	iwd "github.com/gogogoghost/iwdgo"
)

func main() {
	// new instance
	instance, err := iwd.NewIwd()
	if err != nil {
		panic(err)
	}
	// get station
	station := instance.Stations[0]
	println(station.Name())
	// scan
	if err := station.Scan(); err != nil {
		panic(err)
	}
	fmt.Println("scanning...")
	// props auto update by signal, is not real time
	// wait for 1 second
	time.Sleep(time.Second * 1)
	station.WaitForScan()
	// get scan result
	list, err := station.GetOrderedNetworks()
	if err != nil {
		panic(err)
	}
	// find my wifi
	var ap *iwd.OrderedNetwork = nil
	for _, item := range list {
		name := item.Name()
		if name == "SZZY" {
			ap = item
			fmt.Println("found my wifi")
			break
		}
	}
	if ap == nil {
		panic("not found my wifi")
	}
	// if iwd has remember this wifi.
	kn, err := ap.KnownNetwork()
	if err != nil {
		panic(err)
	}
	if kn != nil {
		// forget it
		// avoid iwd using old password
		if err := kn.Forget(); err != nil {
			// ignore this error
			// sometime got a error
			// panic(err)
		}
		println("forget old password")
	}
	if err := ap.Connect("szzy12345678"); err != nil {
		// 此处就会验证密码
		panic(err)
	}
	fmt.Println("connecting...")
	if err := ap.WaitForConnected(); err != nil {
		panic(err)
	}
	fmt.Println("connected!")
}
