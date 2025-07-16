package password

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"

	"github.com/hd-passgen/core/internal/constants"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Generate(masterPassword, serviceName string) (result string, err error) {
	seed := bip39.NewSeed(constants.DefaultMnemonic, masterPassword)

	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return "", fmt.Errorf("Generate: failed to generate master key: %w", err)
	}

	hash := sha256.Sum256([]byte(serviceName))

	// TODO: collision risk
	index := binary.BigEndian.Uint32(hash[:4])

	key, err := masterKey.NewChildKey(index)
	if err != nil {
		return "", fmt.Errorf("Generate: failed to generate child key: %w", err)
	}

	h := hmac.New(sha256.New, key.Key)
	h.Write([]byte(serviceName))
	result = base64.RawURLEncoding.EncodeToString(h.Sum(nil))[:32]

	return result, nil
}
