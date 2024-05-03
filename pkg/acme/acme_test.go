package acme

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/codfrm/cago/pkg/utils"
	"github.com/codfrm/dns-kit/pkg/jws"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestNewAcme(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodGet, DefaultDirectoryUrl, httpmock.NewStringResponder(http.StatusOK, `{
    "ErmLfw9nUsc": "https://community.letsencrypt.org/t/adding-random-entries-to-the-directory/33417",
    "keyChange": "https://acme-v02.api.letsencrypt.org/acme/key-change",
    "meta": {
        "caaIdentities": [
            "letsencrypt.org"
        ],
        "termsOfService": "https://letsencrypt.org/documents/LE-SA-v1.3-September-21-2022.pdf",
        "website": "https://letsencrypt.org"
    },
    "newAccount": "https://acme-v02.api.letsencrypt.org/acme/new-acct",
    "newNonce": "https://acme-v02.api.letsencrypt.org/acme/new-nonce",
    "newOrder": "https://acme-v02.api.letsencrypt.org/acme/new-order",
    "renewalInfo": "https://acme-v02.api.letsencrypt.org/draft-ietf-acme-ari-02/renewalInfo/",
    "revokeCert": "https://acme-v02.api.letsencrypt.org/acme/revoke-cert"
}`))
	httpmock.RegisterResponder(http.MethodHead, "https://acme-v02.api.letsencrypt.org/acme/new-nonce", func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			Header: map[string][]string{"Replay-Nonce": {utils.RandString(32, utils.Mix)}},
		}, nil
	})
	ctx := context.Background()
	acme, err := NewAcme("yz@ggnb.top")
	assert.NoError(t, err)
	assert.NotNil(t, acme)
	var cacheJws map[string]any
	httpmock.RegisterResponder(http.MethodPost, "https://acme-v02.api.letsencrypt.org/acme/new-acct", func(request *http.Request) (*http.Response, error) {
		defer request.Body.Close()
		body, err := io.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		data := struct {
			Contact              []string `json:"contact"`
			TermsOfServiceAgreed bool     `json:"termsOfServiceAgreed"`
		}{}
		header := jws.NewHeader(jws.ES256(nil))
		if err := jws.Decode(string(body), header, &data, jws.WithUnmarshaler(jws.JSONUnmarshaler)); err != nil {
			return nil, err
		}
		assert.Equal(t, "mailto:yz@ggnb.top", data.Contact[0])
		cacheJws = header.Get("jwk").(map[string]any)
		return &http.Response{
			StatusCode: http.StatusCreated,
			Header: map[string][]string{
				"Location": {"https://acme-v02.api.letsencrypt.org/acme/acct/1"},
			},
		}, nil
	})
	httpmock.RegisterResponder(http.MethodPost, "https://acme-v02.api.letsencrypt.org/acme/new-order", func(request *http.Request) (*http.Response, error) {
		defer request.Body.Close()
		body, err := io.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		data := struct {
			Identifiers []Identifier `json:"identifiers"`
		}{}
		header := jws.NewHeader(&es256{
			Algorithm: jws.ES256(nil),
			jwk:       cacheJws,
		})
		if err := jws.Decode(string(body), header, &data, jws.WithUnmarshaler(jws.JSONUnmarshaler)); err != nil {
			return nil, err
		}
		assert.Equal(t, "dns", data.Identifiers[0].Type)
		assert.Equal(t, "test2.ggnb.top", data.Identifiers[0].Value)
		respData, err := json.Marshal(NewOrderResponse{
			Status:      "pending",
			Expires:     time.Time{},
			NotBefore:   time.Time{},
			NotAfter:    time.Time{},
			Identifiers: data.Identifiers,
			Authorizations: []string{
				"https://acme-v02.api.letsencrypt.org/acme/authz-v3/1",
			},
			Finalize: "https://acme-v02.api.letsencrypt.org/acme/finalize/1",
		})
		if err != nil {
			return nil, err
		}
		return &http.Response{
			StatusCode: http.StatusCreated,
			Header: map[string][]string{
				"Location": {"https://acme-v02.api.letsencrypt.org/acme/order/1"},
			},
			Body: io.NopCloser(bytes.NewReader(respData)),
		}, nil
	})
	httpmock.RegisterResponder(http.MethodPost, "https://acme-v02.api.letsencrypt.org/acme/authz-v3/1", func(request *http.Request) (*http.Response, error) {
		respData, err := json.Marshal(AuthorizationResponse{
			Identifier: Identifier{
				Type:  "dns",
				Value: "test2.ggnb.top",
			},
			Status:  "pending",
			Expires: time.Time{},
			Challenges: []AuthorizationChallenge{
				{
					Type:  "dns-01",
					Url:   "https://acme-v02.api.letsencrypt.org/acme/chall-v3/1",
					Token: "test",
				}},
		})
		if err != nil {
			return nil, err
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(respData)),
		}, nil
	})
	challenges, err := acme.GetChallenge(ctx, []string{"test2.ggnb.top"})
	assert.NoError(t, err)
	assert.NotNil(t, challenges)
	assert.Equal(t, acme.options.client.DNS01ChallengeRecord("test"), challenges[0].Record)
	// 设置dns操作
	// 等待acme验证
	httpmock.RegisterResponder(http.MethodPost, "https://acme-v02.api.letsencrypt.org/acme/chall-v3/1", func(request *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		var data string
		header := jws.NewHeader(&es256{
			Algorithm: jws.ES256(nil),
			jwk:       cacheJws,
		})
		if err := jws.Decode(string(body), header, &data, jws.WithUnmarshaler(jws.JSONUnmarshaler)); err != nil {
			return nil, err
		}
		assert.Equal(t, "{}", data)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`{"status":"pending"}`))),
		}, nil
	})
	httpmock.RegisterResponder(http.MethodPost, "https://acme-v02.api.letsencrypt.org/acme/authz-v3/1", func(request *http.Request) (*http.Response, error) {
		respData, err := json.Marshal(AuthorizationResponse{
			Identifier: Identifier{
				Type:  "dns",
				Value: "test2.ggnb.top",
			},
			Status:  "valid",
			Expires: time.Time{},
			Challenges: []AuthorizationChallenge{
				{
					Type:   "dns-01",
					Url:    "https://acme-v02.api.letsencrypt.org/acme/chall-v3/1",
					Token:  "test",
					Status: "valid",
				}},
		})
		if err != nil {
			return nil, err
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(respData)),
		}, nil
	})
	err = acme.WaitChallenge(ctx)
	assert.NoError(t, err)
	// 生成证书
	httpmock.RegisterResponder(http.MethodPost, "https://acme-v02.api.letsencrypt.org/acme/finalize/1", func(request *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		data := struct {
			Csr string `json:"csr"`
		}{}
		header := jws.NewHeader(&es256{
			Algorithm: jws.ES256(nil),
			jwk:       cacheJws,
		})
		if err := jws.Decode(string(body), header, &data, jws.WithUnmarshaler(jws.JSONUnmarshaler)); err != nil {
			return nil, err
		}
		assert.NotEmpty(t, data.Csr)
		respData, err := json.Marshal(FinalizeResponse{
			Status:         "valid",
			Expires:        time.Time{},
			Identifiers:    nil,
			Authorizations: nil,
			Finalize:       "",
			Certificate:    "https://acme-v02.api.letsencrypt.org/acme/cert/1",
		})
		if err != nil {
			return nil, err
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(respData)),
		}, nil
	})
	httpmock.RegisterResponder(http.MethodPost, "https://acme-v02.api.letsencrypt.org/acme/cert/1", httpmock.NewStringResponder(http.StatusOK, "test"))
	data, err := acme.GetCertificate(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, data)
}
