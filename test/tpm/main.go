//go:build !windows && !darwin

package main

import (
	"log"

	"github.com/hd-passgen/core/pkg/tpm"
)

func main() {
	tpm, err := tpm.New()
	if err != nil {
		log.Fatal(err)
	}
	tpm.List()
}
