package google

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
	"golang.org/x/oauth2/google"
)

// Option is an option for Google ID token.
type Option interface {
	Apply(*settings)
}

// WithHTTPClient returns a Option that specifies the HTTP client to use
func WithHTTPClient(client *http.Client) Option {
	return withHTTPClient{client}
}

// WithCredentialsFile returns a Option with credientials file
func WithCredentialsFile(filename string) Option {
	return withCredFile(filename)
}

// WithCredentialsJSON returns a Option with credientials JSON
func WithCredentialsJSON(p []byte) Option {
	return withCredentialsJSON(p)
}

// GetIDToken returns the ID token for given service URL
func GetIDToken(serviceURL string, options ...Option) (string, error) {
	s := &settings{
		HTTPClient:         http.DefaultClient,
		IDTokenExchangeURL: "https://www.googleapis.com/oauth2/v4/token",
	}
	for _, o := range options {
		o.Apply(s)
	}
	jwt, err := s.signJWT(serviceURL)
	if err != nil {
		return "", err
	}
	return s.exchangeJWTForIDToken(string(jwt))
}

type settings struct {
	CredentialsFile    string
	CredentialsJSON    []byte
	HTTPClient         *http.Client
	IDTokenExchangeURL string
}

type withHTTPClient struct{ client *http.Client }

func (w withHTTPClient) Apply(o *settings) {
	o.HTTPClient = w.client
}

type withCredFile string

func (w withCredFile) Apply(o *settings) {
	o.CredentialsFile = string(w)
}

type withCredentialsJSON []byte

func (w withCredentialsJSON) Apply(o *settings) {
	o.CredentialsJSON = make([]byte, len(w))
	copy(o.CredentialsJSON, w)
}

func (s *settings) signJWT(serviceURL string) ([]byte, error) {
	if s.CredentialsFile != "" {
		if s.CredentialsJSON != nil {
			return nil, fmt.Errorf("Can only accept single credentials")
		}
		var err error
		s.CredentialsJSON, err = ioutil.ReadFile(s.CredentialsFile)
		if err != nil {
			return nil, err
		}
	}
	if s.CredentialsJSON == nil {
		return nil, fmt.Errorf("No credentials")
	}
	cfg, err := google.JWTConfigFromJSON(s.CredentialsJSON)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	jot := &struct {
		*jwt.Payload
		TargetAudience string `json:"target_audience"`
	}{
		Payload: &jwt.Payload{
			Issuer:         cfg.Email,
			Subject:        cfg.Email,
			Audience:       jwt.Audience{s.IDTokenExchangeURL},
			IssuedAt:       jwt.NumericDate(now),
			ExpirationTime: jwt.NumericDate(now.Add(time.Hour)),
		},
		TargetAudience: serviceURL,
	}
	block, _ := pem.Decode(cfg.PrivateKey)
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch key := key.(type) {
	case *rsa.PrivateKey:
		return jwt.Sign(jot, jwt.NewRS256(jwt.RSAPrivateKey(key)))
	default:
		return nil, fmt.Errorf("Not a RSA key")
	}
}

func (s *settings) exchangeJWTForIDToken(jwt string) (string, error) {
	res, err := s.HTTPClient.PostForm(s.IDTokenExchangeURL, url.Values{
		"grant_type": []string{"urn:ietf:params:oauth:grant-type:jwt-bearer"},
		"assertion":  []string{jwt},
	})
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("%s", res.Status)
	}
	body := &struct {
		IDToken string `json:"id_token"`
	}{}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(body)
	if err != nil {
		return "", err
	}
	return body.IDToken, nil
}
