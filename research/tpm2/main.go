package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/google/go-tpm/legacy/tpm2"
)

const device = "/dev/tpm0"

const (
	parentPassword = ""
	ownerPassword  = ""
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("not enougth args")
		os.Exit(0)
	}

	rwc, err := tpm2.OpenTPM(device)
	if err != nil {
		log.Fatal(err)
	}
	defer rwc.Close()

	switch os.Args[1] {
	case "store":
		if len(os.Args) < 3 {
			fmt.Println("not enougth args")
			os.Exit(0)
		}
		if err := storePassword(rwc, []byte(os.Args[2])); err != nil {
			log.Fatal("Failed to store password:", err.Error())
		}
	case "load":
		password, err := loadPassword(rwc)
		if err != nil {
			log.Fatal("Failed to load password:", err.Error())
		}
		fmt.Println(string(password))
	default:
		fmt.Println("unknown command")
		os.Exit(0)
	}
}

const (
	keySize = 1024
	// keySize = 2048
)

var tpmPublic = tpm2.Public{
	Type: tpm2.AlgRSA,
	// NameAlg: tpm2.AlgSHA256,
	NameAlg: tpm2.AlgSHA1,
	Attributes: tpm2.FlagFixedTPM | tpm2.FlagFixedParent | tpm2.FlagSensitiveDataOrigin |
		tpm2.FlagUserWithAuth | tpm2.FlagRestricted | tpm2.FlagDecrypt,

	RSAParameters: &tpm2.RSAParams{KeyBits: keySize},
}

func storePassword(rwc io.ReadWriteCloser, password []byte) error {
	// tpmPublic := tpm2.Public{
	// 	Type:    tpm2.AlgRSA,
	// 	NameAlg: tpm2.AlgSHA256,
	// 	Attributes: tpm2.FlagFixedTPM | tpm2.FlagFixedParent | tpm2.FlagSensitiveDataOrigin |
	// 		tpm2.FlagUserWithAuth | tpm2.FlagRestricted | tpm2.FlagDecrypt,
	// 	RSAParameters: &tpm2.RSAParams{KeyBits: 2048},
	// }

	primaryHandle, _, err := tpm2.CreatePrimary(
		rwc,
		tpm2.HandleOwner,
		tpm2.PCRSelection{},
		parentPassword,
		ownerPassword,
		tpmPublic,
	)
	if err != nil {
		return err
	}
	defer func() {
		if err := tpm2.FlushContext(rwc, primaryHandle); err != nil {
			log.Fatal("Failed to flush context:", err.Error())
		}
	}()

	// pcrSelection := tpm2.PCRSelection{
	// 	Hash: tpm2.AlgSHA256,
	// 	PCRs: []int{0, 1, 2, 3, 4, 5, 6, 7},
	// }

	// if err := tpm2.PolicyPCR(rwc, primaryHandle, []byte{}, pcrSelection); err != nil {
	// 	return fmt.Errorf("failed to set up policy pcr: %w", err)
	// }

	sealPrivate, sealPublic, err := tpm2.Seal(
		rwc,
		primaryHandle,
		parentPassword,
		ownerPassword,
		[]byte{},
		password,
	)
	if err != nil {
		return fmt.Errorf("failed to seal sensitive data: %w", err)
	}

	os.WriteFile("seal.priv", sealPrivate, 0644)
	os.WriteFile("seal.pub", sealPublic, 0644)

	return nil
}

func loadPassword(rwc io.ReadWriteCloser) ([]byte, error) {
	// tpmPublic := tpm2.Public{
	// 	Type:    tpm2.AlgRSA,
	// 	NameAlg: tpm2.AlgSHA256,
	// 	Attributes: tpm2.FlagFixedTPM | tpm2.FlagFixedParent | tpm2.FlagSensitiveDataOrigin |
	// 		tpm2.FlagUserWithAuth | tpm2.FlagRestricted | tpm2.FlagDecrypt,
	// 	RSAParameters: &tpm2.RSAParams{KeyBits: 2048},
	// }

	primaryHandle, _, err := tpm2.CreatePrimary(
		rwc,
		tpm2.HandleOwner,
		tpm2.PCRSelection{},
		parentPassword,
		ownerPassword,
		tpmPublic,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tpm2.FlushContext(rwc, primaryHandle); err != nil {
			log.Fatal("Failed to flush context:", err.Error())
		}
	}()

	// sealPublic, err := os.ReadFile("seal.pub")
	// if err != nil {
	// 	return nil, err
	// }

	unsealed, err := tpm2.Unseal(rwc, primaryHandle, parentPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to unseal: %w", err)
	}

	return unsealed, nil
}
