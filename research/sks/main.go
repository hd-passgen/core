//go:build !windows
// +build !windows

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/facebookincubator/flog"
	"github.com/facebookincubator/sks"
	"github.com/spf13/pflag"
)

const usage = `
usage: test`

func main() {
	pflag.Parse()
	pflag.Usage = func() {
		fmt.Println(usage)
	}

	flogCfg := &flog.Config{
		Verbosity:     "10",
		Vmodule:       "",
		TraceLocation: "",
	}
	if err := flogCfg.Set(); err != nil {
		log.Fatal("failed to set flog config:", err.Error())
	}

	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(0)
	}
	switch os.Args[1] {
	case "vendor":
		getVendor()
	case "store":
		storeKey()
	case "load":
		getKey()
	case "remove":
		removeKey()
	default:
		fmt.Println("unknown command:", os.Args[1])
	}

	// key := sks.FromLabelTag("test-key:test-tags")

	// signer, err := sks.NewKey(key.Label(), key.Tag(), false, true, []byte(nil))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// digest := make([]byte, 32)
	// rand.Read(digest)

	// signature, err := signer.Sign(nil, digest, nil)
	// if err != nil {
	// 	log.Fatal("failed to sign:", err.Error())
	// }
	// fmt.Println(string(signature))

	// signer.Remove()

	// key, err := sks.LoadKey("test-key", "test-tags", []byte(nil))
	// if err != nil {
	// 	log.Fatal("failed to load key:", err)
	// }
	// fmt.Println(key.Label(), key.Tag())
}

const (
	keyLabel = "key-label"
	keyTag   = "key-tag"
)

var (
	keyHash = []byte{0x01}
)

func getVendor() {
	vendor, err := sks.GetSecureHardwareVendorData()
	if err != nil {
		log.Fatal(err)
	}
	if vendor != nil {
		fmt.Println(
			"name:", vendor.VendorName,
			"info:", vendor.VendorInfo,
			"version:", vendor.Version,
			"is tpm 2.0 complient device?:", vendor.IsTPM20CompliantDevice,
		)
	}
}

func storeKey() {
	_, err := sks.NewKey(keyLabel, keyTag, true, true, keyHash)
	if err != nil {
		log.Fatal("failed to create key:", err)
	}
}

func getKey() {
	key, err := sks.LoadKey(keyLabel, keyTag, keyHash)
	if err != nil {
		log.Fatal("failed to load key:", err.Error())
	}
	fmt.Println(key.Label(), key.Tag())
}

func removeKey() {
	key := sks.FromLabelTag(keyLabel + ":" + keyTag)
	if err := key.Remove(); err != nil {
		log.Fatal("failed to remove key:", err.Error())
	}
}
