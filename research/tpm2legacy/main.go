package main

import (
	"fmt"
	"log"

	"github.com/google/go-tpm/legacy/tpm2"
	"github.com/google/go-tpm/tpmutil"
)

const device = "/dev/tpm0"

func main() {
	// getRandomBytes()
	// generateRSA()
	all()
}

func getRandomBytes() {
	rwc, err := tpm2.OpenTPM(device)
	if err != nil {
		log.Fatal("Failed to open TPM device:", err.Error())
	}
	defer rwc.Close()
	randomBytes, err := tpm2.GetRandom(rwc, 32)
	if err != nil {
		log.Fatal("Failed to get random bytes:", err.Error())
	}
	fmt.Printf("%x\n", randomBytes)
}

func generateRSA() {
	rwc, err := tpm2.OpenTPM(device)
	if err != nil {
		log.Fatal("Failed to open TPM device:", err.Error())
	}
	defer rwc.Close()

	// Create a primary key template
	primaryTemplate := tpm2.Public{
		Type:       tpm2.AlgRSA,
		NameAlg:    tpm2.AlgSHA256,
		Attributes: tpm2.FlagFixedTPM | tpm2.FlagFixedParent | tpm2.FlagSensitiveDataOrigin | tpm2.FlagUserWithAuth | tpm2.FlagRestricted | tpm2.FlagDecrypt,
		RSAParameters: &tpm2.RSAParams{
			Symmetric: &tpm2.SymScheme{
				Alg:     tpm2.AlgAES,
				KeyBits: 128,
				Mode:    tpm2.AlgCFB,
			},
			KeyBits:     2048,
			ExponentRaw: 0,
			ModulusRaw:  make([]byte, 256),
		},
	}

	parentPassword := ""
	ownerPassword := ""

	primaryHandle, _, err := tpm2.CreatePrimary(
		rwc,
		tpm2.HandleOwner,
		tpm2.PCRSelection{},
		parentPassword,
		ownerPassword,
		primaryTemplate,
	)
	if err != nil {
		log.Fatal("Failed to create primary key:", err.Error())
	}
	defer tpm2.FlushContext(rwc, primaryHandle)

	fmt.Printf("Primary key created with handle: 0x%x\n", primaryHandle)

	keyTemplate := tpm2.Public{
		Type:       tpm2.AlgRSA,
		NameAlg:    tpm2.AlgSHA256,
		Attributes: tpm2.FlagFixedTPM | tpm2.FlagFixedParent | tpm2.FlagUserWithAuth | tpm2.FlagSign,
		RSAParameters: &tpm2.RSAParams{
			Sign: &tpm2.SigScheme{
				Alg:  tpm2.AlgRSASSA,
				Hash: tpm2.AlgSHA256,
			},
			KeyBits: 2048,
		},
	}

	// create key

	privateBlob, publicBlob, _, _, _, err := tpm2.CreateKey(
		rwc,
		primaryHandle,
		tpm2.PCRSelection{},
		parentPassword,
		ownerPassword,
		keyTemplate,
	)
	if err != nil {
		log.Fatal("Failed to create key:", err.Error())
	}

	keyHandle, _, err := tpm2.Load(rwc, primaryHandle, "", publicBlob, privateBlob)
	if err != nil {
		log.Fatal("Failed to load key:", err.Error())
	}
	defer tpm2.FlushContext(rwc, keyHandle)

	fmt.Printf("key loaded with handle: 0x%x\n", keyHandle)

	// making key persistent

	// Continuing from previous example...

	// Find an available persistent handle
	// persistentHandle := tpm2.Handle(0x81010000)

	// // Evict any existing key at this handle
	// _ = tpm2.EvictControl(rwc, "", tpm2.HandleOwner, persistentHandle, persistentHandle)

	// // Make the key persistent
	// err = tpm2.EvictControl(rwc, "", tpm2.HandleOwner, keyHandle, persistentHandle)
	// if err != nil {
	// 	log.Fatalf("Failed to make key persistent: %v", err)
	// }

	// fmt.Printf("Key made persistent at handle: 0x%x\n", persistentHandle)
}

func all() {
	// Open TPM
	tpm, err := tpm2.OpenTPM("/dev/tpm0")
	if err != nil {
		log.Fatalf("Opening TPM: %v", err)
	}
	defer tpm.Close()

	parentPasword := ""
	ownerPassword := ""

	// Create primary key
	primaryHandle, _, err := tpm2.CreatePrimary(tpm, tpm2.HandleOwner, tpm2.PCRSelection{}, parentPasword, ownerPassword, tpm2.Public{
		Type:       tpm2.AlgRSA,
		NameAlg:    tpm2.AlgSHA256,
		Attributes: tpm2.FlagFixedTPM | tpm2.FlagFixedParent | tpm2.FlagSensitiveDataOrigin | tpm2.FlagUserWithAuth | tpm2.FlagRestricted | tpm2.FlagDecrypt,
		RSAParameters: &tpm2.RSAParams{
			KeyBits: 2048,
		},
	})
	if err != nil {
		log.Fatalf("Creating primary key: %v", err)
	}
	defer tpm2.FlushContext(tpm, primaryHandle)

	// Create signing key
	privateBlob, publicBlob, _, _, _, err := tpm2.CreateKey(tpm, primaryHandle, tpm2.PCRSelection{}, parentPasword, ownerPassword, tpm2.Public{
		Type:       tpm2.AlgRSA,
		NameAlg:    tpm2.AlgSHA256,
		Attributes: tpm2.FlagFixedTPM | tpm2.FlagFixedParent | tpm2.FlagUserWithAuth | tpm2.FlagSign,
		RSAParameters: &tpm2.RSAParams{
			Sign: &tpm2.SigScheme{
				Alg:  tpm2.AlgRSASSA,
				Hash: tpm2.AlgSHA256,
			},
			KeyBits: 2048,
		},
	})
	if err != nil {
		log.Fatalf("Creating key: %v", err)
	}

	// Load key
	keyHandle, _, err := tpm2.Load(tpm, primaryHandle, "", publicBlob, privateBlob)
	if err != nil {
		log.Fatalf("Loading key: %v", err)
	}
	defer tpm2.FlushContext(tpm, keyHandle)

	// Make persistent
	persistentHandle := tpmutil.Handle(0x81010000)
	_ = tpm2.EvictControl(tpm, "", tpm2.HandleOwner, persistentHandle, persistentHandle)
	if err := tpm2.EvictControl(tpm, "", tpm2.HandleOwner, keyHandle, persistentHandle); err != nil {
		log.Fatalf("Making key persistent: %v", err)
	}

	// Sign data
	data := []byte("important data")
	sig, err := tpm2.Sign(tpm, persistentHandle, "", data, &tpm2.Ticket{}, &tpm2.SigScheme{
		Alg:  tpm2.AlgRSASSA,
		Hash: tpm2.AlgSHA256,
	})
	if err != nil {
		log.Fatalf("Signing: %v", err)
	}

	fmt.Printf("Signature: %x\n", sig.RSA.Signature)
}
