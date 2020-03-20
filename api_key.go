package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	fmt "fmt"

	"github.com/crosstalkio/log"
	"github.com/gbrlsnchs/jwt/v3"
)

type Algorithm string

const (
	HS256 = "HS256"
	HS384 = "HS384"
	HS512 = "HS512"
	ES256 = "ES256"
	ES384 = "ES384"
	ES512 = "ES512"
)

type APIKeyStore interface {
	PutAPIKey(key *APIKey) error
	GetAPIKey(id string) (*APIKey, error)
	DelAPIKey(id string) error
	ListAPIKeyIDs() ([]string, error)
}

type APIKey struct {
	log.Sugar
	ID        string
	Algorithm Algorithm
	Secret    []byte
	ECDSAKey  *ecdsa.PrivateKey
}

func NewAPIKey(logger log.Logger, id string, secret []byte) *APIKey {
	return &APIKey{
		Sugar:     log.NewSugar(logger),
		ID:        id,
		Secret:    secret,
		Algorithm: HS256,
	}
}

func (k *APIKey) SetAlgorithm(algo Algorithm) error {
	var err error
	switch algo {
	case HS256:
		k.Secret = make([]byte, 32)
		_, err = rand.Read(k.Secret)
	case HS384:
		k.Secret = make([]byte, 48)
		_, err = rand.Read(k.Secret)
	case HS512:
		k.Secret = make([]byte, 64)
		_, err = rand.Read(k.Secret)
	case ES256:
		k.ECDSAKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case ES384:
		k.ECDSAKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case ES512:
		k.ECDSAKey, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		err = fmt.Errorf("Unknown algorithm: %s", algo)
		k.Errorf("%s", err.Error())
		return err
	}
	if err != nil {
		k.Errorf("Failed to generete key: %s", err.Error())
		return err
	}
	k.Algorithm = algo
	return nil
}

func (k *APIKey) CreateToken(payload interface{}) ([]byte, error) {
	algo, err := k.algorithm(k.Algorithm)
	if err != nil {
		return nil, err
	}
	token, err := jwt.Sign(payload, algo)
	if err != nil {
		k.Errorf("Failed to sign JWT: %s", err.Error())
		return nil, err
	}
	return token, nil
}

func (k *APIKey) ParseToken(token []byte, payload interface{}) error {
	algo, err := k.algorithm(k.Algorithm)
	if err != nil {
		return err
	}
	_, err = jwt.Verify(token, algo, payload)
	if err != nil {
		k.Errorf("Invalid signature: %s", err.Error())
		return err
	}
	return nil
}

func (k *APIKey) algorithm(algo Algorithm) (jwt.Algorithm, error) {
	var err error
	switch algo {
	case HS256:
		return jwt.NewHS256(k.Secret), nil
	case HS384:
		return jwt.NewHS384(k.Secret), nil
	case HS512:
		return jwt.NewHS512(k.Secret), nil
	case ES256:
		if k.ECDSAKey == nil {
			err = fmt.Errorf("Missing ECDSA key for: %s", algo)
		} else {
			return jwt.NewES256(jwt.ECDSAPrivateKey(k.ECDSAKey)), nil
		}
	case ES384:
		if k.ECDSAKey == nil {
			err = fmt.Errorf("Missing ECDSA key for: %s", algo)
		} else {
			return jwt.NewES384(jwt.ECDSAPrivateKey(k.ECDSAKey)), nil
		}
	case ES512:
		if k.ECDSAKey == nil {
			err = fmt.Errorf("Missing ECDSA key for: %s", algo)
		} else {
			return jwt.NewES512(jwt.ECDSAPrivateKey(k.ECDSAKey)), nil
		}
	default:
		err = fmt.Errorf("Unsupported algorithm: %s", algo)
	}
	k.Errorf("%s", err.Error())
	return nil, err
}
