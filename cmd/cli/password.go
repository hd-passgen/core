package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hd-passgen/core/internal/domains/password"
	"github.com/hd-passgen/core/internal/objects"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagMasterPassword     = "master-password"
	flagMasterPasswordFile = "master-file"
	flagServiceName        = "service"
	flagPasswordLength     = "length"
	flagVersion            = "version"
)

var passwordGenerateCommand = &cobra.Command{
	Use:   "generate",
	Short: "Generate password.",
	Long:  "Generate password.",
	RunE: func(_ *cobra.Command, _ []string) (err error) {
		input := objects.PasswordParams{
			ServiceName:        viper.GetString(flagServiceName),
			MasterPassword:     viper.GetString(flagMasterPassword),
			MasterPasswordFile: viper.GetString(flagMasterPasswordFile),
			Length:             viper.GetUint8(flagPasswordLength),
			Version:            viper.GetInt(flagVersion),
		}

		if err := Validator.Struct(input); err != nil {
			return fmt.Errorf("RunE: %w", err)
		}

		pass, err := password.Generate(input)
		if err != nil {
			return fmt.Errorf("RunE: %w", err)
		}

		fmt.Fprintf(os.Stdout, "%s\n", pass)

		return nil
	},
}

func addPasswordCommands() (err error) {
	rootCommand.AddCommand(passwordGenerateCommand)

	passwordGenerateCommand.PersistentFlags().StringP(flagServiceName, "s", "", "Service name")
	if err = viper.BindPFlag(flagServiceName, passwordGenerateCommand.PersistentFlags().Lookup(flagServiceName)); err != nil {
		return fmt.Errorf("addPasswordCommands: %w", err)
	}

	passwordGenerateCommand.PersistentFlags().StringP(flagMasterPassword, "p", "", "Master password")
	if err = viper.BindPFlag(flagMasterPassword, passwordGenerateCommand.PersistentFlags().Lookup(flagMasterPassword)); err != nil {
		return fmt.Errorf("addPasswordCommands: %w", err)
	}

	passwordGenerateCommand.PersistentFlags().StringP(flagMasterPasswordFile, "f", "", "Path to file with master password")
	if err = viper.BindPFlag(flagMasterPasswordFile, passwordGenerateCommand.PersistentFlags().Lookup(flagMasterPasswordFile)); err != nil {
		return fmt.Errorf("addPasswordCommands: %w", err)
	}

	passwordGenerateCommand.PersistentFlags().Uint8P(flagPasswordLength, "l", 0, "Length of password (default 32)")
	if err = viper.BindPFlag(flagPasswordLength, passwordGenerateCommand.PersistentFlags().Lookup(flagPasswordLength)); err != nil {
		return fmt.Errorf("addPasswordCommands: %w", err)
	}

	passwordGenerateCommand.PersistentFlags().StringP(flagVersion, "v", "", "Version of password")
	if err = viper.BindPFlag(flagVersion, passwordGenerateCommand.PersistentFlags().Lookup(flagVersion)); err != nil {
		return fmt.Errorf("addPasswordCommands: %w", err)
	}

	return nil
}

func init() {
	if err := addPasswordCommands(); err != nil {
		log.Fatal(err)
	}
}
