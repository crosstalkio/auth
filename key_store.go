package auth

import (
	"crypto/ecdsa"
	"crypto/x509"
	"net/url"

	"github.com/crosstalkio/log"
	"google.golang.org/protobuf/proto"
)

func NewKeyStore(logger log.Logger, url *url.URL) (KeyStore, error) {
	store, err := NewBlobStore(url)
	if err != nil {
		return nil, err
	}
	return NewKeyStoreFromBobStore(logger, store), nil
}

func NewKeyStoreFromBobStore(logger log.Logger, store BlobStore) KeyStore {
	return &keyBlobStore{
		Sugar: log.NewSugar(logger),
		store: store,
	}
}

type keyBlobStore struct {
	log.Sugar
	store BlobStore
}

func (s *keyBlobStore) PutKey(key *Key) error {
	blob := &Blob{
		Algorithm: string(key.Algorithm),
		Secret:    key.Secret,
	}
	var err error
	if key.ECDSAKey != nil {
		blob.ECDSAKey, err = s.marshalECDSAKey(key.ECDSAKey)
		if err != nil {
			return err
		}
	}
	val, err := proto.Marshal(blob)
	if err != nil {
		s.Errorf("Failed to marshal API key '%s': %s", key.ID, err.Error())
		return err
	}
	err = s.store.PutBlob(key.ID, val)
	if err != nil {
		s.Errorf("Failed to put API key to blob store: %s", err.Error())
		return err
	}
	return nil
}

func (s *keyBlobStore) ListKeyIDs() ([]string, error) {
	ids, err := s.store.ListBlobIDs()
	if err != nil {
		s.Errorf("Failed to list API key ID from blob store: %s", err.Error())
		return nil, err
	}
	return ids, nil
}

func (s *keyBlobStore) GetKey(id string) (*Key, error) {
	val, err := s.store.GetBlob(id)
	if err != nil {
		s.Errorf("Failed to get API key from blob store '%s': %s", id, err.Error())
		return nil, err
	}
	if val == nil {
		s.Debugf("API key not found in blob store: %s", id)
		return nil, nil
	}
	blob := &Blob{}
	err = proto.Unmarshal(val, blob)
	if err != nil {
		s.Errorf("Failed to unmarshal API key '%s': %s", id, err.Error())
		return nil, err
	}
	key := NewKey(id, blob.Secret)
	if blob.Algorithm == "" {
		s.Warningf("Using HS256 as default algorithm: %s", id)
		key.Algorithm = HS256
	} else {
		key.Algorithm = Algorithm(blob.Algorithm)
	}
	if blob.ECDSAKey != nil {
		key.ECDSAKey, err = s.parseECDSAKey(blob.ECDSAKey)
		if err != nil {
			return nil, err
		}
	}
	return key, nil
}

func (s *keyBlobStore) DelKey(id string) error {
	err := s.store.DelBlob(id)
	if err != nil {
		s.Errorf("Failed to del API key from blob store: %s", err.Error())
		return err
	}
	return nil
}

func (s *keyBlobStore) marshalECDSAKey(key *ecdsa.PrivateKey) ([]byte, error) {
	bytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		s.Errorf("Failed to marshal ECDSA key: %s", err.Error())
		return nil, err
	}
	return bytes, nil
}

func (s *keyBlobStore) parseECDSAKey(bytes []byte) (*ecdsa.PrivateKey, error) {
	key, err := x509.ParseECPrivateKey(bytes)
	if err != nil {
		s.Errorf("Failed to parse ECDSA key: %s", err.Error())
		return nil, err
	}
	return key, nil
}
