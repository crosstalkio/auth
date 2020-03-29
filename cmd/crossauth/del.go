package main

import (
	"github.com/crosstalkio/auth"
	"github.com/crosstalkio/log"
)

func del(logger log.Sugar, store auth.KeyStore, id string) error {
	err := store.DelKey(id)
	if err != nil {
		return err
	}
	logger.Infof("API Key deleted: %s", id)
	return nil
}
