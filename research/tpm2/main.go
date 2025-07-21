package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	keyfile "github.com/foxboron/go-tpm-keyfiles"
	tpm2legacy "github.com/google/go-tpm/legacy/tpm2"
	"github.com/google/go-tpm/tpm"
	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpm2/transport"
	"github.com/google/go-tpm/tpm2/transport/linuxtpm"
	"github.com/google/go-tpm/tpm2/transport/simulator"
	"github.com/spf13/pflag"
)

const (
	tpm0Device   = "/dev/tpm0"
	tpmrm0Device = "/dev/tpmrm0"
)

var (
	tpmPath = pflag.String("tpm", tpm0Device, "")
)

func getDeviceTPM() transport.TPMCloser {
	rwc, err := linuxtpm.Open(*tpmPath)
	if err != nil {
		log.Fatal(err)
	}
	return rwc
}

func getSimulatorTPM() transport.TPMCloser {
	tpm, err := simulator.OpenSimulator()
	if err != nil {
		log.Fatal(err)
	}
	return tpm
}

func main() {
	pflag.Parse()
	// rwc, err := linuxtpm.Open(*tpmPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// rwc, err := tpm2legacy.OpenTPM(*tpmPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// listTPMKeys(rwc)
}

func loadKeyfile() {
	tpm := getSimulatorTPM()

	key, err := keyfile.NewLoadableKey(tpm, tpm2.TPMAlgECC, 256, []byte{}, keyfile.WithDescription("TPM Key"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", key)
	fmt.Println(string(key.Secret.Buffer))
}

func sealedData(tpm transport.TPMCloser) {
	msg := []byte("message")
	k, _ := keyfile.NewSealedData(tpm, msg, []byte(nil))
	data, _ := keyfile.UnsealData(tpm, k, []byte(nil))
	if bytes.Equal(data, msg) {
		fmt.Println("same message")
	}
}

func getRandomData() {

}

// func createAndStoreKey(rwc io.ReadWriteCloser) (tpmutil.Handle, []byte, error) {
// }

func listTPMKeys(rwc io.ReadWriteCloser) {
	_ = tpm2legacy.ClockInfo{} // for storing import of tpm2legacy module

	handles, err := tpm.GetKeys(rwc)
	if err != nil {
		log.Fatal("failed to get keys:", err)
	}

	fmt.Printf("%d keys loaded in TPM\n", len(handles))
	for i, handle := range handles {
		fmt.Printf(" (%d) key handle %d\n", i+1, handle)
	}
}

type TpmDevice struct {
	rwc io.ReadWriteCloser
}

func findPubKey(label string, tag string, hash []byte) error {
	return nil
}
