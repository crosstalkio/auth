package auth

import (
	"io"

	"github.com/crosstalkio/log"
)

type KeyStore interface {
	io.Closer
	PutKey(key *Key) error
	GetKey(id string) (*Key, error)
	DelKey(id string) error
	ListKeyIDs() ([]string, error)
}

func NewKeyStore(logger log.Logger, url string) (KeyStore, error) {
	store, err := NewBlobStore(logger, url)
	if err != nil {
		return nil, err
	}
	return NewBlobKeyStore(logger, store), nil
}
