package main

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version = ""
)

var (
	validate = validator.New()
)

type parameters struct {
	ServiceName        string `validate:"required"`
	MasterPassword     string `validate:"required_without=MasterPasswordFile"`
	MasterPasswordFile string `validate:"required_without=MasterPassword,omitempty,file"`
	Length             uint8  `validate:"omitempty,max=40"`
	Version            int    `validate:"omitempty,min=1"`
}

var (
	defaultPasswordLength  = uint8(32)
	currentPasswordVersion = 1
)

var (
	rootCmd = &cobra.Command{
		Use:   "hd-passgen",
		Short: "HD password generator",
		Long:  "HD password generator",
	}

	passGenCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate password",
		Long:  "Generate password",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := parameters{
				ServiceName:        viper.GetString("service"),
				MasterPassword:     viper.GetString("master-password"),
				MasterPasswordFile: viper.GetString("master-password-file"),
				Length:             viper.GetUint8("length"),
				Version:            viper.GetInt("version"),
			}

			password, err := generatePassword(params)
			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stdout, "%s\n", password)

			return nil
		},
	}
)

func main() {
	if Version == "" {
		Version = "unknown"
	}

	rootCmd.AddCommand(passGenCmd)

	passGenCmd.PersistentFlags().StringP("service", "s", "", "service name")
	if err := viper.BindPFlag("service", passGenCmd.PersistentFlags().Lookup("service")); err != nil {
		fmt.Println("Failed to add flag: ", err.Error())
		os.Exit(1)
	}

	passGenCmd.PersistentFlags().StringP("master-password", "p", "", "master password")
	if err := viper.BindPFlag("master-password", passGenCmd.PersistentFlags().Lookup("master-password")); err != nil {
		fmt.Println("Failed to add flag: ", err.Error())
		os.Exit(1)
	}

	passGenCmd.PersistentFlags().StringP("master-password-file", "f", "", "master password file")
	if err := viper.BindPFlag("master-password-file", passGenCmd.PersistentFlags().Lookup("master-password-file")); err != nil {
		fmt.Println("Failed to add flag: ", err.Error())
		os.Exit(1)
	}

	passGenCmd.PersistentFlags().Uint8P("length", "l", defaultPasswordLength, "password length (default 32)")
	if err := viper.BindPFlag("length", passGenCmd.PersistentFlags().Lookup("length")); err != nil {
		fmt.Println("Failed to add flag: ", err.Error())
		os.Exit(1)
	}

	passGenCmd.PersistentFlags().IntP("version", "v", currentPasswordVersion, "password version")
	if err := viper.BindPFlag("version", passGenCmd.PersistentFlags().Lookup("version")); err != nil {
		fmt.Println("Failed to add flag: ", err.Error())
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

const (
	minPasswordLength = 8
	maxPasswordLenght = 40
)

func generatePassword(params parameters) (string, error) {
	if params.Version == 0 {
		params.Version = currentPasswordVersion
	}

	if err := validate.Struct(params); err != nil {
		return "", err
	}

	if params.Version != 0 && params.Version > currentPasswordVersion {
		return "", fmt.Errorf("unsupported password version: %d", params.Version)
	}

	if params.Length < minPasswordLength || params.Length > maxPasswordLenght {
		return "", fmt.Errorf(
			"invalid password length: %d; shoud be %d <= length <= %d",
			params.Length, minPasswordLength, maxPasswordLenght,
		)
	}

	switch params.Version {
	case 1:
		return generatePasswordV1(params)
	default:
		return "", fmt.Errorf("password generation for v%d not implemented", params.Version)
	}
}
