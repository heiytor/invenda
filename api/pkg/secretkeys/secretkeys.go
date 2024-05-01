package secretkeys

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	KindPrivateKey = "PRIVATE KEY"
	KindPublicKey  = "PUBLIC KEY"
)

var (
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
)

// Load reads the secret keys under "secrets" directory and setup [PrivateKey] and [PublicKey].
func Load() error {
	// Up three directories and enters "secrets"
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Join(filepath.Join(filepath.Join(filepath.Dir(file), ".."), ".."), ".."), "secrets")

	var err error

	PrivateKey, err = parse[ed25519.PrivateKey](dir+"/api-private.pem", KindPrivateKey)
	if err != nil {
		return err
	}

	PublicKey, err = parse[ed25519.PublicKey](dir+"/api-public.pem", KindPublicKey)
	return err
}

func parse[T any](file, kind string) (T, error) {
	// tmp will be returned if an error occurs
	var tmp T

	pemBytes, err := os.ReadFile(file)
	if err != nil {
		return tmp, err
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return tmp, fmt.Errorf("pem failed to decode pem")
	}

	if block.Type != kind {
		return tmp, fmt.Errorf("pem is not of type %s", kind)
	}

	switch kind {
	case KindPrivateKey:
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return tmp, err
		}
		return key.(T), nil
	case KindPublicKey:
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return tmp, err
		}
		return key.(T), nil
	default:
		return tmp, fmt.Errorf("kind %s is not a valid entry", kind)
	}
}
