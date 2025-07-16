package main

import (
	"fmt"
	"log"

	"github.com/facebookincubator/sks"
)

func main() {
	vendor, err := sks.GetSecureHardwareVendorData()
	if err != nil {
		log.Fatal(err)
	}
	if vendor != nil {
		fmt.Println(vendor.VendorName, vendor.VendorInfo, vendor.Version)
	}
}
