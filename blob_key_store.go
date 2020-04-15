package auth

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"fmt"

	"github.com/crosstalkio/log"
	"google.golang.org/protobuf/proto"
)

func NewBlobKeyStore(logger log.Logger, store BlobStore) KeyStore {
	return &blobKeyStore{
		Sugar: log.NewSugar(logger),
		store: store,
	}
}

type blobKeyStore struct {
	log.Sugar
	store BlobStore
}

func (s *blobKeyStore) Close() error {
	return s.store.Close()
}

func (s *blobKeyStore) PutKey(key *Key) error {
	if key.ID == "" {
		err := fmt.Errorf("Key ID is empty")
		s.Errorf("%s", err.Error())
		return err
	}
	_, pb, err := s.getKey(key.ID)
	if err != nil {
		return err
	}
	blob := &Blob{
		Algorithm: string(key.Algorithm),
		Secret:    key.Secret,
	}
	if key.ECDSAKey != nil {
		blob.EcdsaKey, err = s.marshalECDSAKey(key.ECDSAKey)
		if err != nil {
			return err
		}
	}
	var val []byte
	if pb {
		val, err = proto.Marshal(blob)
	} else {
		val, err = json.Marshal(blob)
	}
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

func (s *blobKeyStore) ListKeyIDs() ([]string, error) {
	ids, err := s.store.ListBlobIDs()
	if err != nil {
		s.Errorf("Failed to list API key ID from blob store: %s", err.Error())
		return nil, err
	}
	return ids, nil
}

func (s *blobKeyStore) GetKey(id string) (*Key, error) {
	key, _, err := s.getKey(id)
	return key, err
}

func (s *blobKeyStore) getKey(id string) (*Key, bool, error) {
	pb := true
	val, err := s.store.GetBlob(id)
	if err != nil {
		s.Errorf("Failed to get API key from blob store '%s': %s", id, err.Error())
		return nil, pb, err
	}
	if val == nil {
		s.Debugf("API key not found in blob store: %s", id)
		return nil, pb, nil
	}
	blob := &Blob{}
	err = proto.Unmarshal(val, blob)
	if err != nil {
		err = json.Unmarshal(val, blob)
		if err != nil {
			s.Errorf("Failed to unmarshal API key '%s': %s", id, err.Error())
			return nil, pb, err
		}
		pb = false
	}
	key := NewKey(id, blob.Secret)
	if blob.Algorithm == "" {
		s.Warningf("Using HS256 as default algorithm: %s", id)
		key.Algorithm = HS256
	} else {
		key.Algorithm = Algorithm(blob.Algorithm)
	}
	if blob.EcdsaKey != nil {
		key.ECDSAKey, err = s.parseECDSAKey(blob.EcdsaKey)
		if err != nil {
			return nil, pb, err
		}
	}
	return key, pb, nil
}

func (s *blobKeyStore) DelKey(id string) (bool, error) {
	deleted, err := s.store.DelBlob(id)
	if err != nil {
		s.Errorf("Failed to del API key from blob store: %s", err.Error())
		return deleted, err
	}
	return deleted, nil
}

func (s *blobKeyStore) marshalECDSAKey(key *ecdsa.PrivateKey) ([]byte, error) {
	bytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		s.Errorf("Failed to marshal ECDSA key: %s", err.Error())
		return nil, err
	}
	return bytes, nil
}

func (s *blobKeyStore) parseECDSAKey(bytes []byte) (*ecdsa.PrivateKey, error) {
	key, err := x509.ParseECPrivateKey(bytes)
	if err != nil {
		s.Errorf("Failed to parse ECDSA key: %s", err.Error())
		return nil, err
	}
	return key, nil
}
