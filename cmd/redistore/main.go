package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/crosstalkio/auth"
	"github.com/crosstalkio/log"
	"github.com/go-redis/redis"
)

var errInvalid = fmt.Errorf("Invalid argument")
var basename string

func addUsage() {
	fmt.Printf("Usage: %s add <ID> <algorithm> [secret]\n", basename)
}

func getUsage() {
	fmt.Printf("Usage: %s get <ID>\n", basename)
}

func delUsage() {
	fmt.Printf("Usage: %s del <ID>\n", basename)
}

func signUsage() {
	fmt.Printf("Usage: %s sign <ID> <json string|file>\n", basename)
}

func verifyUsage() {
	fmt.Printf("Usage: %s verify <JWT string|file>\n", basename)
}

func main() {
	logger := log.NewSugar(log.NewLogger(log.Color(log.GoLogger(log.Debug, os.Stderr, "", log.LstdFlags))))
	storeUrl := flag.String("store", "redis://127.0.0.1:6379/crosstalk/apikey/", "")
	flag.Usage = func() {
		fmt.Printf("Usage: %s add|get|del|sign|verify <args...>\n", basename)
	}
	flag.Parse()
	basename = filepath.Base(os.Args[0])
	u, err := url.Parse(*storeUrl)
	if err != nil {
		logger.Errorf("Invalid store URL: %s", storeUrl)
		os.Exit(1)
	}
	if u.Scheme != "redis" {
		err = fmt.Errorf("Unkown store: %s", u.Scheme)
		logger.Errorf("%s", err.Error())
		os.Exit(1)
	}
	addr := fmt.Sprintf("%s:%s", u.Hostname(), u.Port())
	pass, _ := u.User.Password()
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
	})
	_, err = client.Ping().Result()
	if err != nil {
		logger.Errorf("Failed to ping redis: %s\n", err.Error())
		os.Exit(1)
	}
	store := auth.NewAPIKeyStore(logger, auth.NewGoRedisBlobStore(client, u.Path))
	err = handle(logger, store)
	if err != nil {
		os.Exit(1)
	}
	return
}

func handle(logger log.Sugar, store auth.APIKeyStore) error {
	cmd := flag.Arg(0)
	switch cmd {
	case "add":
		id := flag.Arg(1)
		algo := flag.Arg(2)
		secret := flag.Arg(3)
		if id == "" || algo == "" {
			addUsage()
			return errInvalid
		}
		return add(logger, store, id, algo, secret)
	case "get":
		id := flag.Arg(1)
		if id == "" {
			getUsage()
			return errInvalid
		}
		return get(logger, store, id)
	case "del":
		id := flag.Arg(1)
		if id == "" {
			delUsage()
			return errInvalid
		}
		return del(logger, store, id)
	case "sign":
		id := flag.Arg(1)
		json := flag.Arg(2)
		if id == "" || json == "" {
			signUsage()
			return errInvalid
		}
		var bytes []byte
		if strings.HasPrefix(json, "{") {
			bytes = []byte(json)
		} else {
			f, err := os.Open(json)
			if err != nil {
				logger.Errorf("Failed to open file '%s': %s", json, err.Error())
				return err
			}
			defer f.Close()
			bytes, err = ioutil.ReadAll(f)
			if err != nil {
				logger.Errorf("Failed to read file '%s': %s", json, err.Error())
				return err
			}
		}
		return sign(logger, store, id, bytes)
	case "verify":
		token := flag.Arg(1)
		if token == "" {
			verifyUsage()
			return errInvalid
		}
		var bytes []byte
		splits := strings.Split(token, ".")
		if len(splits) == 3 {
			bytes = []byte(token)
		} else {
			f, err := os.Open(token)
			if err != nil {
				logger.Errorf("Failed to open file '%s': %s", token, err.Error())
				return err
			}
			defer f.Close()
			bytes, err = ioutil.ReadAll(f)
			if err != nil {
				logger.Errorf("Failed to read file '%s': %s", token, err.Error())
				return err
			}
		}
		return verify(logger, store, bytes)
	default:
		if cmd != "" {
			logger.Errorf("Unknown command: %s", cmd)
		}
		flag.Usage()
		return errInvalid
	}
}
