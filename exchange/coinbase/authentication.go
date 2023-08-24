package coinbase

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Transport struct {
	Auth         Authenticator
	PrefixString string
	Base         http.RoundTripper
}

// Stolen from OAUATH2
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqBodyClosed := false
	if req.Body != nil {
		defer func() {
			if !reqBodyClosed {
				req.Body.Close()
			}
		}()
	}
	req2 := cloneRequest(req)
	t.Auth.SignRequest(strings.Split(strings.TrimPrefix(req.URL.Path, t.PrefixString), "?")[0], req2)
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
