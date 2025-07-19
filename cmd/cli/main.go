package main

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

var (
	Version   = ""
	Validator = validator.New()
)

var rootCommand = &cobra.Command{
	Use:   "hd-passgen",
	Short: "HD password generator.",
	Long:  "HD password generator.",
}

func main() {
	if Version == "" {
		Version = "unknown"
	}

	if err := rootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}
