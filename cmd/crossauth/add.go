package main

import (
	"github.com/crosstalkio/auth"
	"github.com/crosstalkio/log"
)

func add(logger log.Sugar, store auth.APIKeyStore, id, algo, secret string) error {
	key := auth.NewAPIKey(id, nil)
	err := key.SetAlgorithm(auth.Algorithm(algo))
	if err != nil {
		logger.Errorf("Failed to set algorithm: %s", err.Error())
		return err
	}
	if secret != "" {
		key.Secret = []byte(secret)
	}
	err = store.PutAPIKey(key)
	if err != nil {
		return err
	}
	dumpKey(logger, key)
	return nil
}
