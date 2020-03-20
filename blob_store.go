package auth

import (
	"fmt"
	"net/url"
)

type BlobStore interface {
	GetBlob(id string) ([]byte, error)
	PutBlob(id string, val []byte) error
	DelBlob(id string) error
	ListBlobIDs() ([]string, error)
}

type BlobStoreFactory interface {
	CreateBlobStore(url *url.URL) (BlobStore, error)
}

var blobStoreFactories = map[string]BlobStoreFactory{}

func RegisterBlobStore(proto string, factory BlobStoreFactory) {
	blobStoreFactories[proto] = factory
}

func NewBlobStore(url *url.URL) (BlobStore, error) {
	factory := blobStoreFactories[url.Scheme]
	if factory == nil {
		return nil, fmt.Errorf("Scheme not supported: %s", url.Scheme)
	}
	return factory.CreateBlobStore(url)
}
