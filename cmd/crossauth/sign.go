package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/crosstalkio/auth"
	"github.com/crosstalkio/log"
)

func sign(logger log.Sugar, store auth.APIKeyStore, id string, bytes []byte, ttl int64) error {
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
	payload["iss"] = key.ID
	if ttl > 0 {
		payload["exp"] = time.Now().Add(time.Second * time.Duration(ttl)).Unix()
	}
	token, err := key.CreateToken(&payload)
	if err != nil {
		return err
	}
	logger.Infof("JWT signed with %s: %s", key.Algorithm, key.ID)
	fmt.Printf("%s\n", token)
	return nil
}
