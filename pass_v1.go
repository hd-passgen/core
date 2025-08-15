package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"os"
	"strings"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

const (
	// The number of words should be 12, 15, 18, 21 or 24
	defaultMnemonic = "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
)

func generatePasswordV1(params parameters) (string, error) {
	if params.MasterPasswordFile != "" {
		content, err := os.ReadFile(params.MasterPasswordFile)
		if err != nil {
			return "", fmt.Errorf("failed to read master password file: %w", err)
		}
		params.MasterPassword = strings.TrimSpace(string(content))
	}

	masterKey, err := bip32.NewMasterKey(
		bip39.NewSeed(defaultMnemonic, params.MasterPassword),
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate master key: %w", err)
	}

	// define index child key
	hash := sha256.Sum256([]byte(params.ServiceName))
	// TODO: fix collision risk
	index := binary.BigEndian.Uint32(hash[:4])

	key, err := masterKey.NewChildKey(index)
	if err != nil {
		return "", fmt.Errorf("failed to generate child key: %w", err)
	}

	h := hmac.New(sha256.New, key.Key)
	_, _ = h.Write([]byte(params.ServiceName))
	raw := h.Sum(nil)

	encoded := base64.RawURLEncoding.EncodeToString(raw)
	lengthWithoutVersion := params.Length - 2
	result := encoded[:lengthWithoutVersion]

	return fmt.Sprintf("%s01", result), nil
}
