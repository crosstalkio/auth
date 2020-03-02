package main

import (
	"github.com/crosstalkio/auth"
	"github.com/crosstalkio/log"
)

func del(logger log.Sugar, store auth.APIKeyStore, id string) error {
	err := store.DelAPIKey(id)
	if err != nil {
		return err
	}
	logger.Infof("API Key deleted: %s", id)
	return nil
}
