package main

import (
	"fmt"
	"io"
	"log"

	"github.com/google/go-tpm/legacy/tpm2"
	"github.com/spf13/pflag"
)

var (
	tmp2Sock = pflag.String("sock", "/tmp/tpmstate/swtpm-sock", "swtpm socket")
)

const (
	tpm0   = "/dev/tpm0"
	tpmrm0 = "/dev/tpmrm0"
)

func main() {
	pflag.Parse()

	rwc, err := tpm2.OpenTPM(tpmrm0)
	if err != nil {
		log.Fatal(err)
	}
	// defer func() { _ = rwc.Close() }()

	fmt.Println("opened")

	// readRandom(rwc)
	secureStorage(rwc)
}

func readRandom(rwc io.ReadWriteCloser) {
	randomBytes, err := tpm2.GetRandom(rwc, 16)
	check(err)
	fmt.Printf("%x\n", randomBytes)
}

func secureStorage(rwc io.ReadWriteCloser) {
	primaryHandle, _, err := tpm2.CreatePrimary(
		rwc,
		tpm2.HandleOwner,
		tpm2.PCRSelection{},
		"",
		"",
		defaultKeyParams(),
	)
	check(err)
	defer func() { _ = tpm2.FlushContext(rwc, primaryHandle) }()

	secret := []byte("my super secret data")
	// pcrSelection := tpm2.PCRSelection{
	// 	Hash: tpm2.AlgSHA256,
	// 	PCRs: []int{0, 1, 2}, // Bind to PCRs 0,1,2
	// }

	priv, pub, err := tpm2.Seal(rwc, primaryHandle, "", "", []byte{}, secret)
	check(err)

	fmt.Println(priv, pub)
	// Seal the data
	// priv, pub, err := tpm2.Seal(rwc, primaryHandle, "", "", pcrSelection, secret)
	// check(err)

	// Later, unseal the data (will only work if PCRs haven't changed)
	// unsealed, err := tpm2.Unseal(rwc, primaryHandle, "", "", priv, pub)
	// check(err)

	// fmt.Printf("Unsealed data: %s\n", string(unsealed))
}

func defaultKeyParams() tpm2.Public {
	return tpm2.Public{
		Type:       tpm2.AlgRSA,
		NameAlg:    tpm2.AlgSHA256,
		Attributes: tpm2.FlagSign | tpm2.FlagSensitiveDataOrigin | tpm2.FlagUserWithAuth,
		RSAParameters: &tpm2.RSAParams{
			Sign: &tpm2.SigScheme{
				Alg:  tpm2.AlgRSASSA,
				Hash: tpm2.AlgSHA256,
			},
			KeyBits: 2048,
		},
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
