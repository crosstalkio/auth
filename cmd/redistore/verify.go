package main

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/crosstalkio/auth"
	"github.com/crosstalkio/log"
	"github.com/gbrlsnchs/jwt/v3"
)

func verify(logger log.Sugar, store auth.APIKeyStore, bytes []byte) error {
	payload := &jwt.Payload{}
	hdr, err := jwt.Verify(bytes, jwt.None(), &payload)
	if err != nil {
		logger.Errorf("Invalid JWT: %s", err.Error())
		return err
	}
	id := payload.Subject
	if id == "" {
		err = fmt.Errorf("Missing JWT subject")
		logger.Errorf(err.Error())
		return err
	}
	key, err := store.GetAPIKey(id)
	if err != nil {
		return err
	}
	if key == nil {
		logger.Errorf("No such API key: %s", id)
		return err
	}
	var v interface{}
	err = key.ParseToken(bytes, &v)
	if err != nil {
		logger.Warningf("Failed to verify JWT with %s: %s", hdr.Algorithm, id)
		return nil
	}
	logger.Infof("JWT verified with %s: %s", hdr.Algorithm, id)
	splits := strings.Split(string(bytes), ".")
	part1, _ := base64.RawStdEncoding.DecodeString(splits[0])
	part2, _ := base64.RawStdEncoding.DecodeString(splits[1])
	logger.Infof("Header:  %s", part1)
	logger.Infof("Payload: %s", part2)
	return nil
}
