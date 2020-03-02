package main

import (
	"encoding/json"
	"fmt"

	"github.com/crosstalkio/auth"
	"github.com/crosstalkio/log"
)

func sign(logger log.Sugar, store auth.APIKeyStore, id string, bytes []byte) error {
	key, err := store.GetAPIKey(id)
	if err != nil {
		return err
	}
	if key == nil {
		err = fmt.Errorf("No such API key: %s", id)
		logger.Errorf("%s", err.Error())
		return err
	}
	payload := make(map[string]interface{})
	err = json.Unmarshal(bytes, &payload)
	if err != nil {
		logger.Errorf("Failed to unmarshal JSON: %s", err.Error())
		return err
	}
	payload["sub"] = key.ID
	token, err := key.CreateToken(&payload)
	if err != nil {
		return err
	}
	logger.Infof("JWT signed with %s: %s", key.Algorithm, key.ID)
	fmt.Printf("%s\n", token)
	return nil
}
