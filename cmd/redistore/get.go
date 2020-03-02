package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/crosstalkio/auth"
	"github.com/crosstalkio/log"
)

func get(logger log.Sugar, store auth.APIKeyStore, id string) error {
	key, err := store.GetAPIKey(id)
	if err != nil {
		return err
	}
	if key == nil {
		err = fmt.Errorf("No such API key: %s", id)
		logger.Errorf("%s", err.Error())
		return err
	}
	logger.Infof("Key ID: %s", id)
	if key.Secret != nil {
		if utf8.Valid(key.Secret) {
			logger.Infof("Secret: %s", key.Secret)
		} else {
			logger.Infof("Secret: %s", base64.StdEncoding.EncodeToString(key.Secret))
		}
	}
	if key.ECDSAKey != nil {
		logger.Infof("ECDSA:")
		dumpECDSAKey(logger, key.ECDSAKey)
	}
	return nil
}

func dumpECDSAKey(logger log.Sugar, key *ecdsa.PrivateKey) {
	prv, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		logger.Errorf("Failed to marshal ECDSA privatey key: %s\n", err.Error())
		os.Exit(1)
	}
	pub, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		logger.Errorf("Failed to marshal ECDSA public key: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("%s%s", pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: prv}), pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pub}))
}
