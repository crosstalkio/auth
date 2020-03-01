package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/crosstalkio/auth"
	"github.com/crosstalkio/log"
	"github.com/go-redis/redis"
)

func main() {
	logger := log.NewSugar(log.NewLogger(log.Color(log.GoLogger(log.Debug, os.Stderr, "", log.LstdFlags))))
	addr := flag.String("a", "127.0.0.1:6379", "")
	pass := flag.String("P", "", "")
	prefix := flag.String("p", "crosstalk/apikey/", "")
	flag.Usage = func() {
		fmt.Printf("Usage: %s <ID> [secret]\n", filepath.Base(os.Args[0]))
	}
	flag.Parse()
	n := len(os.Args)
	if n < 2 {
		flag.Usage()
		return
	}
	id := os.Args[1]
	secret := ""
	if n > 2 {
		secret = os.Args[2]
	}
	client := redis.NewClient(&redis.Options{
		Addr:     *addr,
		Password: *pass,
	})
	_, err := client.Ping().Result()
	if err != nil {
		logger.Errorf("Failed to ping redis: %s\n", err.Error())
		os.Exit(1)
	}
	store := auth.NewAPIKeyStore(logger, auth.NewGoRedisBlobStore(client, *prefix))
	if secret != "" {
		err = store.PutAPIKey(auth.NewAPIKey(logger, id, []byte(secret)))
		if err != nil {
			os.Exit(1)
		}
	} else {
		key, err := store.GetAPIKey(id)
		if err != nil {
			os.Exit(1)
		}
		if key == nil {
			fmt.Printf("No such API key: %s\n", id)
		} else {
			fmt.Printf("Key ID: %s\n", id)
			fmt.Printf("Secret: %s\n", key.Secret)
		}
	}
}
