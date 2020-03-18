package goredis

import (
	"github.com/crosstalkio/auth"
)

func init() {
	factory := &Factory{}
	auth.RegisterBlobStore("redis", factory)
	auth.RegisterBlobStore("goredis", factory)
}
