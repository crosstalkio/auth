package goredis

import (
	"github.com/crosstalkio/auth"
)

var schemes = []string{
	"redis",
	"goredis",
}

func init() {
	factory := &factory{}
	for _, scheme := range schemes {
		auth.RegisterBlobStore(scheme, factory)
	}
}
