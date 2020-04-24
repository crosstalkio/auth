package goredis

import (
	"github.com/crosstalkio/auth"
	"github.com/crosstalkio/log"
	"github.com/go-redis/redis/v7"
)

type Options struct {
	redis.Options
	Prefix string
}

type blobStore struct {
	log.Sugar
	client *redis.Client
	prefix string
}

// NewBlobStore creates a redis-based blob store using go-redis client
func NewBlobStore(logger log.Logger, options *Options) (auth.BlobStore, error) {
	s := &blobStore{
		Sugar:  log.NewSugar(logger),
		client: redis.NewClient(&options.Options),
		prefix: options.Prefix,
	}
	err := s.client.Ping().Err()
	if err != nil {
		s.Errorf("Failed to ping redis: %s", err.Error())
		return nil, err
	}
	return s, nil
}

func (s *blobStore) Close() error {
	return s.client.Close()
}

func (s *blobStore) GetBlob(id string) ([]byte, error) {
	s.Debugf("Getting blob from redis: %s", id)
	val, err := s.client.Get(s.prefix + id).Bytes()
	if err != nil {
		if err == redis.Nil {
			s.Debugf("Blob not found in redis: %s", id)
			return nil, nil
		}
		s.Errorf("Failed to get blob from redis: %s", err.Error())
		return nil, err
	}
	return val, nil
}

func (s *blobStore) PutBlob(id string, val []byte) error {
	s.Debugf("Putting blob to redis: %s (%d bytes)", id, len(val))
	return s.client.Set(s.prefix+id, val, 0).Err()
}

func (s *blobStore) DelBlob(id string) (bool, error) {
	s.Debugf("Deleting blob from redis: %s", id)
	n, err := s.client.Del(s.prefix + id).Result()
	if err != nil {
		s.Errorf("Failed to delete redis: %s", err.Error())
		return false, err
	}
	return n > 0, nil
}

func (s *blobStore) ListBlobIDs() ([]string, error) {
	s.Debugf("Listing blob ID in redis")
	keys, err := s.client.Keys(s.prefix + "*").Result()
	if err != nil {
		s.Errorf("Failed to list blob ID in redis: %s", err.Error())
		return nil, err
	}
	ids := make([]string, len(keys))
	n := len(s.prefix)
	for i, key := range keys {
		ids[i] = key[n:]
	}
	return ids, nil
}
