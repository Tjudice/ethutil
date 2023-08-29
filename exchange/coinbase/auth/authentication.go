package auth

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
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

type SignedMessage struct {
	Key        string `json:"key"`
	Timestamp  string `json:"timestamp"`
	Passphrase string `json:"passphrase,omitempty"`
	Sig        string `json:"signature"`
}

type Authenticator interface {
	SignRequest(requestPath string, r *http.Request) error
	SignWebsocketRequest(channels []string, products []string) (*SignedMessage, error)
}

func EncodeBody(r *http.Request) (string, error) {
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
