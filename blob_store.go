package auth

import (
	"fmt"
	"io"
	"net/url"

	"github.com/crosstalkio/log"
)

type BlobStore interface {
	io.Closer
	GetBlob(id string) ([]byte, error)
	PutBlob(id string, val []byte) error
	DelBlob(id string) (bool, error)
	ListBlobIDs() ([]string, error)
}

type BlobStoreFactory interface {
	CreateBlobStore(logger log.Logger, url *url.URL) (BlobStore, error)
}

var blobStoreFactories = map[string]BlobStoreFactory{}

func RegisterBlobStore(proto string, factory BlobStoreFactory) {
	blobStoreFactories[proto] = factory
}

func NewBlobStore(logger log.Logger, url_ string) (BlobStore, error) {
	s := log.NewSugar(logger)
	url, err := url.Parse(url_)
	if err != nil {
		s.Errorf("Invalid URL: %s (%s)", url_, err.Error())
		return nil, err
	}
	factory := blobStoreFactories[url.Scheme]
	if factory == nil {
		err = fmt.Errorf("Factory not registered: %s", url.Scheme)
		s.Errorf(err.Error())
		return nil, err
	}
	return factory.CreateBlobStore(logger, url)
}
