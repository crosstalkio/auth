package goredis

import (
	"fmt"

	"github.com/go-redis/redis"
)

type BlobStore struct {
	client *redis.Client
	prefix string
}

func NewBlobStore(client *redis.Client, prefix string) *BlobStore {
	return &BlobStore{client: client, prefix: prefix}
}

func (s *BlobStore) GetBlob(id string) ([]byte, error) {
	val, err := s.client.Get(s.prefix + id).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return val, nil
}

func (s *BlobStore) PutBlob(id string, val []byte) error {
	return s.client.Set(s.prefix+id, val, 0).Err()
}

func (s *BlobStore) DelBlob(id string) error {
	n, err := s.client.Del(s.prefix + id).Result()
	if err != nil {
		return err
	}
	if n <= 0 {
		return fmt.Errorf("Not found: %s", s.prefix+id)
	}
	return nil
}
