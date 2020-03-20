package auth

import (
	"crypto/ecdsa"
	"crypto/x509"

	"github.com/crosstalkio/log"
	proto "github.com/golang/protobuf/proto"
)

func NewAPIKeyStore(logger log.Logger, store BlobStore) APIKeyStore {
	return &apiKeyBlobStore{
		Sugar: log.NewSugar(logger),
		store: store,
	}
}

type apiKeyBlobStore struct {
	log.Sugar
	store BlobStore
}

func (s *apiKeyBlobStore) PutAPIKey(key *APIKey) error {
	blob := &APIKeyBlob{
		Algorithm: string(key.Algorithm),
		Secret:    key.Secret,
	}
	var err error
	if key.ECDSAKey != nil {
		blob.ECDSAKey, err = x509.MarshalECPrivateKey(key.ECDSAKey)
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

func (s *apiKeyBlobStore) ListAPIKeyIDs() ([]string, error) {
	ids, err := s.store.ListBlobIDs()
	if err != nil {
		s.Errorf("Failed to list API key ID from blob store: %s", err.Error())
		return nil, err
	}
	return ids, nil
}

func (s *apiKeyBlobStore) GetAPIKey(id string) (*APIKey, error) {
	val, err := s.store.GetBlob(id)
	if err != nil {
		s.Errorf("Failed to get API key from blob store '%s': %s", id, err.Error())
		return nil, err
	}
	if val == nil {
		s.Debugf("API key not found in blob store: %s", id)
		return nil, nil
	}
	blob := &APIKeyBlob{}
	err = proto.Unmarshal(val, blob)
	if err != nil {
		s.Errorf("Failed to unmarshal API key '%s': %s", id, err.Error())
		return nil, err
	}
	key := NewAPIKey(s, id, blob.Secret)
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

func (s *apiKeyBlobStore) DelAPIKey(id string) error {
	err := s.store.DelBlob(id)
	if err != nil {
		s.Errorf("Failed to del API key from blob store: %s", err.Error())
		return err
	}
	return nil
}

func (s *apiKeyBlobStore) marshalECDSAKey(key *ecdsa.PrivateKey) ([]byte, error) {
	bytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		s.Errorf("Failed to marshal ECDSA key: %s", err.Error())
		return nil, err
	}
	return bytes, nil
}

func (s *apiKeyBlobStore) parseECDSAKey(bytes []byte) (*ecdsa.PrivateKey, error) {
	key, err := x509.ParseECPrivateKey(bytes)
	if err != nil {
		s.Errorf("Failed to parse ECDSA key: %s", err.Error())
		return nil, err
	}
	return key, nil
}
