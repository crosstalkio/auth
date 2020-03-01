package auth

import (
	fmt "fmt"

	"github.com/crosstalkio/log"
	"github.com/gbrlsnchs/jwt/v3"
)

type Algorithm string

const (
	HS256 = "HS256"
	HS384 = "HS384"
	HS512 = "HS512"
)

type APIKeyStore interface {
	PutAPIKey(key *APIKey) error
	GetAPIKey(id string) (*APIKey, error)
	DelAPIKey(id string) error
}

type APIKey struct {
	log.Sugar
	ID     string
	Secret []byte
}

func NewAPIKey(logger log.Logger, id string, secret []byte) *APIKey {
	return &APIKey{
		Sugar:  log.NewSugar(logger),
		ID:     id,
		Secret: secret,
	}
}

func (k *APIKey) CreateToken(algorithm Algorithm, payload interface{}) ([]byte, error) {
	algo, err := k.algorithm(algorithm)
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

func (k *APIKey) VerifyToken(algorithm Algorithm, token []byte, payload interface{}) error {
	algo, err := k.algorithm(algorithm)
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
	switch algo {
	case HS256:
		return jwt.NewHS256(k.Secret), nil
	case HS384:
		return jwt.NewHS384(k.Secret), nil
	case HS512:
		return jwt.NewHS512(k.Secret), nil
	default:
		err := fmt.Errorf("Unsupported algorithm: %s", algo)
		k.Errorf("%s", err.Error())
		return nil, err
	}
}
