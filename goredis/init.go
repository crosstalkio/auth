package goredis

import (
	"github.com/crosstalkio/auth"
)

func init() {
	factory := &factory{}
	auth.RegisterBlobStore("redis", factory)
	auth.RegisterBlobStore("goredis", factory)
}
