package main

import (
	"github.com/crosstalkio/auth"
	"github.com/crosstalkio/log"
)

func del(logger log.Sugar, store auth.KeyStore, id string) error {
	deleted, err := store.DelKey(id)
	if err != nil {
		return err
	}
	if deleted {
		logger.Infof("API Key deleted: %s", id)
	} else {
		logger.Infof("API Key not exist: %s", id)
	}
	return nil
}
