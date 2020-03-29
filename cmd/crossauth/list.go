package main

import (
	"github.com/crosstalkio/auth"
	"github.com/crosstalkio/log"
)

func list(logger log.Sugar, store auth.KeyStore) error {
	ids, err := store.ListKeyIDs()
	if err != nil {
		return err
	}
	keys := make([]*auth.Key, len(ids))
	for i, id := range ids {
		keys[i], err = store.GetKey(id)
		if err != nil {
			return err
		}
	}
	logger.Infof("Queried %d keys", len(keys))
	for _, key := range keys {
		dumpKey(logger, key)
	}
	return nil
}
