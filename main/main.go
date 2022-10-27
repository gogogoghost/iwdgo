package main

import (
	"fmt"

	iwd "github.com/gogogoghost/iwdgo"
)

func main() {
	iwd, err := iwd.NewIwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(iwd.Stations)
	err = iwd.Stations[0].Scan()
	if err != nil {
		panic(err)
	}
}
