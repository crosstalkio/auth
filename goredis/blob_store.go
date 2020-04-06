package goredis

import (
	"fmt"

	"github.com/crosstalkio/auth"
	"github.com/go-redis/redis"
)

type blobStore struct {
	client *redis.Client
	prefix string
}

// NewBlobStore creates a redis-based blob store using go-redis client
func NewBlobStore(client *redis.Client, prefix string) auth.BlobStore {
	return &blobStore{client: client, prefix: prefix}
}

func (s *blobStore) Close() error {
	return s.client.Close()
}

func (s *blobStore) GetBlob(id string) ([]byte, error) {
	val, err := s.client.Get(s.prefix + id).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return val, nil
}

func (s *blobStore) PutBlob(id string, val []byte) error {
	return s.client.Set(s.prefix+id, val, 0).Err()
}

func (s *blobStore) DelBlob(id string) error {
	n, err := s.client.Del(s.prefix + id).Result()
	if err != nil {
		return err
	}
	if n <= 0 {
		return fmt.Errorf("Not found: %s", s.prefix+id)
	}
	return nil
}

func (s *blobStore) ListBlobIDs() ([]string, error) {
	keys, err := s.client.Keys(s.prefix + "*").Result()
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(keys))
	n := len(s.prefix)
	for i, key := range keys {
		ids[i] = key[n:]
	}
	return ids, nil
}
