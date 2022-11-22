package main

import (
	"fmt"
	"os"

	iwd "github.com/gogogoghost/iwdgo"
)

func main() {
	pid := os.Getpid()
	fmt.Printf("进程 PID: %d \n", pid)
	instance, err := iwd.NewIwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(instance.Conn.Names())
	station := instance.Stations[0]
	err = station.Scan()
	if err != nil {
		panic(err)
	}
	fmt.Println("scanning...")
	err = station.WaitForScan()
	if err != nil {
		panic(err)
	}
	list, err := station.GetOrderedNetworks()
	if err != nil {
		panic(err)
	}
	var ap *iwd.OrderedNetwork = nil
	for _, item := range list {
		name, err := item.GetName()
		if err != nil {
			panic(err)
		}
		if name == "SZZY" {
			ap = item
			fmt.Println("found my wifi")
			break
		}
	}
	if ap == nil {
		panic("not found my wifi")
	}
	// time.Sleep(time.Hour)
	if err := ap.Connect("szzy123456"); err != nil {
		panic(err)
	}
	fmt.Println("connecting...")
	if err := ap.WaitForConnected(); err != nil {
		panic(err)
	}
	fmt.Println("connected!")
}
