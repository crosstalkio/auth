package auth

import (
	"time"

	"github.com/crosstalkio/log"
	"github.com/gbrlsnchs/jwt/v3"
)

type Manager struct {
	log.Sugar
	store APIKeyStore
}

func NewManager(logger log.Logger, store APIKeyStore) *Manager {
	return &Manager{
		Sugar: log.NewSugar(logger),
		store: store,
	}
}

func (m *Manager) VerifyToken(token []byte, payload interface{}) (*APIKey, error) {
	now := time.Now()
	jot := &jwt.Payload{}
	hdr, err := jwt.Verify(token, jwt.None(), jot)
	if err != nil {
		m.Errorf("Invalid JWT: %s (%s)", err.Error(), token)
		return nil, nil
	}
	if jot.IssuedAt != nil && jot.IssuedAt.Time.After(now) {
		m.Errorf("JWT issued in the future: %v", jot.IssuedAt.Time)
		return nil, nil
	}
	if jot.ExpirationTime != nil && jot.ExpirationTime.Time.Before(now) {
		m.Errorf("JWT expired: %v", jot.ExpirationTime.Time)
		return nil, nil
	}
	if jot.NotBefore != nil && jot.NotBefore.Time.After(now) {
		m.Errorf("JWT not yet valid: %v", jot.NotBefore.Time)
		return nil, nil
	}
	var key *APIKey
	if jot.Issuer == "" {
		m.Errorf("Missing JWT isser")
		return nil, nil
	} else {
		key, err = m.store.GetAPIKey(jot.Issuer)
		if err != nil {
			return nil, err
		}
		if key == nil {
			m.Errorf("JWT isser not found: %s", jot.Issuer)
			return nil, nil
		}
	}
	if hdr.Algorithm != string(key.Algorithm) {
		m.Errorf("Algorithm mismatch: expect=%s, actual=%s", key.Algorithm, hdr.Algorithm)
		return nil, nil
	}
	err = key.ParseToken(token, payload)
	if err != nil {
		return nil, nil
	}
	return key, nil
}
