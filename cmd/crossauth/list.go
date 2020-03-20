package main

import (
	"github.com/crosstalkio/auth"
	"github.com/crosstalkio/log"
)

func list(logger log.Sugar, store auth.APIKeyStore) error {
	ids, err := store.ListAPIKeyIDs()
	if err != nil {
		return err
	}
	keys := make([]*auth.APIKey, len(ids))
	for i, id := range ids {
		keys[i], err = store.GetAPIKey(id)
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
