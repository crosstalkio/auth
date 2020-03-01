package auth

import (
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
	val, err := proto.Marshal(&APIKeyBlob{Secret: key.Secret})
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
	return NewAPIKey(s, id, blob.Secret), nil
}

func (s *apiKeyBlobStore) DelAPIKey(id string) error {
	err := s.store.DelBlob(id)
	if err != nil {
		s.Errorf("Failed to del API key from blob store: %s", err.Error())
		return err
	}
	return nil
}
