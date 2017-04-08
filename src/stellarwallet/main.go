/*

	MyStellarWaller is for stellar network wallet written in golang.

	It features:
		1. create account
		2. record history asset info and support web page draw graphics
		3. send payment

*/

package main

import (
	//"fmt"
	"time"
	"webui"
)

func main() {

	go webui.RunServer()

	//monitor.RunGrabData()
	for {
		time.Sleep(1 * time.Second)
	}
}
