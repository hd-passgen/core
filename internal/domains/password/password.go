package password

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"os"
	"strings"

	"github.com/hd-passgen/core/internal/constants"
	"github.com/hd-passgen/core/internal/objects"
	"github.com/samber/lo"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

var (
	currentVersion = 1
)

const (
	defaultLength = 32
	minLength     = 8
	maxLength     = 40
)

var (
	ErrInvalidLength         = fmt.Errorf("invalid length")
	ErrVersionIsNotSupported = fmt.Errorf("version is not supported")
)

func validateGenerate(params objects.PasswordParams) (err error) {
	if lo.IsNotEmpty(params.Version) && params.Version > currentVersion {
		return fmt.Errorf("validateGenerate: %w", ErrVersionIsNotSupported)
	}

	if params.Length < minLength || params.Length > maxLength {
		return fmt.Errorf("validateGenerate: %w", ErrInvalidLength)
	}

	return nil
}

func setDefaultValues(params *objects.PasswordParams) {
	if lo.IsEmpty(params.Length) {
		params.Length = defaultLength
	}

	if lo.IsEmpty(params.Version) {
		params.Version = currentVersion
	}
}

func Generate(params objects.PasswordParams) (result string, err error) {
	setDefaultValues(&params)

	if err = validateGenerate(params); err != nil {
		return "", fmt.Errorf("Generate: %w", err)
	}

	switch params.Version {
	case 1:
		return generateV1(params)
	default:
		return "", fmt.Errorf("Generate: %w", ErrVersionIsNotSupported)
	}
}

func generateV1(params objects.PasswordParams) (result string, err error) {
	if lo.IsNotEmpty(params.MasterPasswordFile) {
		content, err := os.ReadFile(params.MasterPasswordFile)
		if err != nil {
			return "", fmt.Errorf("Generate: failed to read master password file: %w", err)
		}

		params.MasterPassword = strings.TrimSpace(string(content))
	}

	seed := bip39.NewSeed(constants.DefaultMnemonic, params.MasterPassword)

	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return "", fmt.Errorf("Generate: failed to generate master key: %w", err)
	}

	// define index child key
	hash := sha256.Sum256([]byte(params.ServiceName))
	// TODO: fix collision risk
	index := binary.BigEndian.Uint32(hash[:4])

	key, err := masterKey.NewChildKey(index)
	if err != nil {
		return "", fmt.Errorf("Generate: failed to generate child key: %w", err)
	}

	h := hmac.New(sha256.New, key.Key)
	h.Write([]byte(params.ServiceName))
	raw := h.Sum(nil)

	encoded := base64.RawURLEncoding.EncodeToString(raw)

	lengthWithoutVersion := params.Length - 2
	result = encoded[:lengthWithoutVersion]
	return fmt.Sprintf("%s01", result), nil
}
