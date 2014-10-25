package main

import "fmt"
import "github.com/fuzxxl/nfc/2.0/nfc"
//import "github.com/fuzxxl/freefare/0.3/freefare"

func main() {

	nfc_device, nfc_error := nfc.Open("pn532_uart:/dev/ttyUSB0:115200")
	fmt.Println(nfc_device,nfc_error)
}