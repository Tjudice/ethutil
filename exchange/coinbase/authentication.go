package coinbase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gfx.cafe/open/ghost/hexutil"
	"gopkg.in/yaml.v3"
)

type Transport struct {
	Auth         Authenticator
	PrefixString string
	Base         http.RoundTripper
}

// RoundTrip authorizes and authenticates the request with an
// access token from Transport's Source.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqBodyClosed := false
	if req.Body != nil {
		defer func() {
			if !reqBodyClosed {
				req.Body.Close()
			}
		}()
	}
	req2 := cloneRequest(req) // per RoundTripper contract
	t.Auth.SignRequest(strings.TrimPrefix(req.URL.Path, t.PrefixString), req2)
	// req.Body is assumed to be closed by the base RoundTripper.
	reqBodyClosed = true
	return t.base().RoundTrip(req2)
}

func (t *Transport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}

func cloneRequest(r *http.Request) *http.Request {
	r2 := new(http.Request)
	*r2 = *r
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}

type AuthenticationType int

var (
	ExchangeAuth      AuthenticationType = 1
	AdvancedTradeAuth AuthenticationType = 2
	OAUTH2            AuthenticationType = 3
)

type Authenticator interface {
	SignRequest(requestPath string, r *http.Request) error
	SignWebsocketRequest(channels []string, products []string) (*SignedMessage, error)
}

type ExchangeAccountAuth struct {
	API_KEY        string `json:"api_key" yaml:"api_key"`
	API_PASSPHRASE string `json:"api_passphrase" yaml:"api_passphrase"`
	API_SECRET     string `json:"api_secret" yaml:"api_secret"`
}

func LoadAccount(accountType AuthenticationType, filepath string) (Authenticator, error) {
	bts, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var out Authenticator
	switch accountType {
	case ExchangeAuth:
		out = &ExchangeAccountAuth{}
	case AdvancedTradeAuth:
		out = &AdvancedTradeAuthenticator{}
	case OAUTH2:
		panic("not implemented")
	}
	err = yaml.Unmarshal(bts, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (a *ExchangeAccountAuth) SignRequest(requestPath string, r *http.Request) error {
	access_timestamp := time.Now().Unix()
	body, err := encodeBody(r)
	if err != nil {
		return err
	}
	message := strconv.FormatInt(access_timestamp, 10) + r.Method + requestPath + body
	decoded, err := base64.StdEncoding.DecodeString(a.API_SECRET)
	if err != nil {
		return err
	}
	secretHmac := hmac.New(sha256.New, decoded)
	_, err = secretHmac.Write([]byte(message))
	if err != nil {
		return err
	}
	sig := secretHmac.Sum(make([]byte, 0, secretHmac.Size()))
	sig64 := base64.StdEncoding.EncodeToString(sig)
	r.Header.Add("CB-ACCESS-KEY", a.API_KEY)
	r.Header.Add("CB-ACCESS-TIMESTAMP", strconv.FormatInt(access_timestamp, 10))
	r.Header.Add("CB-ACCESS-PASSPHRASE", a.API_PASSPHRASE)
	r.Header.Add("CB-ACCESS-SIGN", string(sig64))
	return nil
}

type SignedMessage struct {
	Key        string `json:"key"`
	Timestamp  string `json:"timestamp"`
	Passphrase string `json:"passphrase,omitempty"`
	Sig        string `json:"signature"`
}

func (a *ExchangeAccountAuth) SignWebsocketRequest(channels []string, products []string) (*SignedMessage, error) {
	access_timestamp := time.Now().Unix()
	message := strconv.FormatInt(access_timestamp, 10) + http.MethodGet + "/users/self/verify"
	decoded, err := base64.StdEncoding.DecodeString(a.API_SECRET)
	if err != nil {
		return nil, err
	}
	secretHmac := hmac.New(sha256.New, decoded)
	_, err = secretHmac.Write([]byte(message))
	if err != nil {
		return nil, err
	}
	sig := secretHmac.Sum(make([]byte, 0, secretHmac.Size()))
	sig64 := base64.StdEncoding.EncodeToString(sig)
	r := &SignedMessage{
		Key:        a.API_KEY,
		Timestamp:  strconv.FormatInt(access_timestamp, 10),
		Passphrase: a.API_PASSPHRASE,
		Sig:        sig64,
	}
	return r, nil
}

type AdvancedTradeAuthenticator struct {
	API_KEY    string `json:"api_key" yaml:"api_key"`
	API_SECRET string `json:"api_secret" yaml:"api_secret"`
}

func (a *AdvancedTradeAuthenticator) SignRequest(requestPath string, r *http.Request) error {
	access_timestamp := time.Now().Unix()
	body, err := encodeBody(r)
	if err != nil {
		return err
	}
	message := strconv.FormatInt(access_timestamp, 10) + r.Method + requestPath + body
	secretHmac := hmac.New(sha256.New, []byte(a.API_SECRET))
	_, err = secretHmac.Write([]byte(message))
	if err != nil {
		return err
	}
	sig := secretHmac.Sum(make([]byte, 0, secretHmac.Size()))
	r.Header.Add("CB-ACCESS-KEY", a.API_KEY)
	r.Header.Add("CB-ACCESS-TIMESTAMP", strconv.FormatInt(access_timestamp, 10))
	r.Header.Add("CB-ACCESS-SIGN", hexutil.Bytes(sig).String()[2:])
	return nil
}

func (a *AdvancedTradeAuthenticator) SignWebsocketRequest(channels []string, products []string) (*SignedMessage, error) {
	access_timestamp := time.Now().Unix()
	channelStr := strings.Join(channels, "") + strings.Join(products, ",")
	message := strconv.FormatInt(access_timestamp, 10) + channelStr
	secretHmac := hmac.New(sha256.New, []byte(a.API_SECRET))
	_, err := secretHmac.Write([]byte(message))
	if err != nil {
		return nil, err
	}
	sig := secretHmac.Sum(make([]byte, 0, secretHmac.Size()))
	r := &SignedMessage{
		Key:       a.API_KEY,
		Timestamp: strconv.FormatInt(access_timestamp, 10),
		Sig:       hexutil.Bytes(sig).String()[2:],
	}
	return r, nil
}

func encodeBody(r *http.Request) (string, error) {
	if r.Body == nil {
		return "", nil
	}
	body, err := r.GetBody()
	if err != nil {
		return "", err
	}
	bodyBts, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}
	if len(bodyBts) == 0 {
		return "", nil
	}
	marshalled, err := json.Marshal(bodyBts)
	if err != nil {
		return "", err
	}
	return string(marshalled), nil
}
