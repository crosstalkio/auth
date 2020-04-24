package goredis

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/crosstalkio/auth"
	"github.com/crosstalkio/log"
	"github.com/go-redis/redis/v7"
)

type factory struct{}

func (f *factory) CreateBlobStore(logger log.Logger, u *url.URL) (auth.BlobStore, error) {
	hit := false
	for _, scheme := range schemes {
		if u.Scheme == scheme {
			hit = true
			break
		}
	}
	if !hit {
		err := fmt.Errorf("Not supported scheme: %s", u.Scheme)
		log.NewSugar(logger).Errorf(err.Error())
		return nil, err
	}
	pass, _ := u.User.Password()
	opts := &Options{
		Options: redis.Options{
			Addr:     u.Host,
			Password: pass,
		},
		Prefix: strings.TrimPrefix(u.Path, "/"),
	}
	return NewBlobStore(logger, opts)
}
