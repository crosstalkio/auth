package main

import (
	"github.com/crosstalkio/auth"
	"github.com/crosstalkio/log"
)

func add(logger log.Sugar, store auth.APIKeyStore, id, algo, secret string) error {
	key := auth.NewAPIKey(logger, id, nil)
	err := key.SetAlgorithm(auth.Algorithm(algo))
	if err != nil {
		return err
	}
	if secret != "" {
		key.Secret = []byte(secret)
	}
	err = store.PutAPIKey(key)
	if err != nil {
		return err
	}
	logger.Infof("API Key added with %s: %s", key.Algorithm, key.ID)
	return nil
}
