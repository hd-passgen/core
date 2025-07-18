package password

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"

	"github.com/hd-passgen/core/internal/constants"
	"github.com/samber/lo"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

var (
	ErrInvalidLength = fmt.Errorf("invalid length")
)

const (
	defaultLength = 32
	minLength     = 8
	maxLength     = 40
)

func Generate(masterPassword, serviceName string, length uint8) (result string, err error) {
	if lo.IsEmpty(length) {
		length = defaultLength
	}

	if length < minLength || length > maxLength {
		return "", fmt.Errorf("Generate: %w", ErrInvalidLength)
	}

	seed := bip39.NewSeed(constants.DefaultMnemonic, masterPassword)

	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return "", fmt.Errorf("Generate: failed to generate master key: %w", err)
	}

	hash := sha256.Sum256([]byte(serviceName))

	// TODO: fix collision risk
	index := binary.BigEndian.Uint32(hash[:4])

	key, err := masterKey.NewChildKey(index)
	if err != nil {
		return "", fmt.Errorf("Generate: failed to generate child key: %w", err)
	}

	h := hmac.New(sha256.New, key.Key)
	h.Write([]byte(serviceName))
	raw := h.Sum(nil)

	encoded := base64.RawURLEncoding.EncodeToString(raw)
	return encoded[:length], nil
}
