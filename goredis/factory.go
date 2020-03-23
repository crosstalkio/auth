package goredis

import (
	fmt "fmt"
	"net/url"
	"strings"

	"github.com/crosstalkio/auth"
	"github.com/go-redis/redis"
)

type Factory struct{}

func (f *Factory) CreateBlobStore(u *url.URL) (auth.BlobStore, error) {
	addr := fmt.Sprintf("%s:%s", u.Hostname(), u.Port())
	pass, _ := u.User.Password()
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return NewBlobStore(client, strings.TrimPrefix(u.Path, "/")), nil
}
