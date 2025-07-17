package tpm

import (
	"fmt"
	"io"

	"github.com/google/go-tpm/legacy/tpm2"
	"github.com/google/go-tpm/tpmutil"
)

var tpmPaths = []string{
	"/dev/tpm0",
	"/dev/tpmrm0",
}

func openTPM() (io.ReadWriteCloser, error) {
	for _, path := range tpmPaths {
		rwc, err := tpm2.OpenTPM(path)
		if err == nil {
			return rwc, nil
		}
	}
	return nil, fmt.Errorf("no TPM 2.0 device found")
}

type TPM struct {
	rwc io.ReadWriteCloser
}

func New() (*TPM, error) {
	// if emulated {
	// 	tpmCloser, err := simulator.OpenSimulator()
	// 	if err != nil {
	// 		return TPM{}, err
	// 	}
	// }

	rwc, err := openTPM()
	if err != nil {
		return nil, err
	}

	return &TPM{
		rwc: rwc,
	}, nil
}

func (tpm *TPM) List() {
	// handles := []tpmutil.Handle{}

	// TPM 2.0 persistent handle range (spec defined)
	startHandle := tpmutil.Handle(0x81000000)
	endHandle := tpmutil.Handle(0x810FFFFF)

	for handle := startHandle; handle <= endHandle; handle++ {
		public, name, _, err := tpm2.ReadPublic(tpm.rwc, handle)
		if err == nil {
			// name, err := public.Name()
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// fmt.Println(name)

			fmt.Println(string(name))
			fmt.Println(public.Type.String())
			// handles = append(handles, handle)
			fmt.Printf("0x%x\n", handle)

			fmt.Println()
		}
	}
}
