package acme

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/codfrm/dns-kit/pkg/jws"
)

var DefaultDirectoryUrl = "https://acme-v02.api.letsencrypt.org/directory"

var DefaultClient = &Client{
	options: newClientOptions(),
}

type Client struct {
	options *ClientOptions
}

func NewClient(opts ...ClientOption) (*Client, error) {
	options := newClientOptions(opts...)
	client := &Client{
		options: options,
	}
	if options.directory == nil {
		_, err := client.Directory()
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

func (c *Client) GetDirectory() *Directory {
	return c.options.directory
}

func (c *Client) SetDirectory(directory *Directory) {
	c.options.directory = directory
}
func (c *Client) Directory() (*Directory, error) {
	// 请求目录
	req, err := http.NewRequest(http.MethodGet, c.options.directoryUrl, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.options.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	directory := &Directory{}
	// 解析目录
	if err := json.NewDecoder(resp.Body).Decode(directory); err != nil {
		return nil, err
	}
	c.options.directory = directory
	return directory, nil
}

func (c *Client) NewNonce() (string, error) {
	req, err := http.NewRequest(http.MethodHead, c.options.directory.NewNonce, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	nonce := resp.Header.Get("Replay-Nonce")
	if nonce == "" {
		return "", fmt.Errorf("Replay-Nonce not found")
	}
	return nonce, nil
}

func (c *Client) GetPrivateKey() *ecdsa.PrivateKey {
	return c.options.privateKey
}

func (c *Client) SetPrivateKey(privateKey *ecdsa.PrivateKey) {
	c.options.privateKey = privateKey
}

func (c *Client) SetKid(kid string) {
	c.options.kid = kid
}

func (c *Client) GetKid() string {
	return c.options.kid
}

var ErrPrivateKeyRequired = errors.New("private key required")

func (c *Client) newRequest(url string, payload any) (*http.Request, error) {
	nonce, err := c.NewNonce()
	if err != nil {
		return nil, err
	}
	// 注册账户需要签名
	if c.options.privateKey == nil {
		return nil, ErrPrivateKeyRequired
	}
	var header *jws.Header
	// 如果有kid则使用kid签名
	if c.options.kid != "" {
		header = jws.NewHeader(newEs256(c.options.kid, c.options.privateKey))
	} else {
		header = jws.NewHeader(jws.ES256(c.options.privateKey))
	}
	data, err := jws.Encode(header.Set("nonce", nonce).Set("url", url),
		payload, jws.WithSerialization(jws.JSONSerialization))
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/jose+json")
	return req, nil
}

func (c *Client) do(url string, payload any) ([]byte, *http.Response, error) {
	req, err := c.newRequest(url, payload)
	if err != nil {
		return nil, nil, err
	}
	resp, err := c.options.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return body, resp, nil
}

func (c *Client) NewAccount(contact []string) (string, error) {
	body, resp, err := c.do(c.options.directory.NewAccount, map[string]interface{}{
		"termsOfServiceAgreed": true,
		"contact":              contact,
	})
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("NewAccount failed: %s", body)
	}
	if resp.Header.Get("Location") == "" {
		return "", fmt.Errorf("location not found: %s", body)
	}
	return resp.Header.Get("Location"), nil
}

type Identifiers struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type NewOrderResponse struct {
	Status      string    `json:"status"`
	Expires     time.Time `json:"expires"`
	NotBefore   time.Time `json:"notBefore"`
	NotAfter    time.Time `json:"notAfter"`
	Identifiers []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"identifiers"`
	Authorizations []string `json:"authorizations"`
	Finalize       string   `json:"finalize"`
}

func (c *Client) NewOrder(identifiers []Identifiers) (*NewOrderResponse, error) {
	body, resp, err := c.do(c.options.directory.NewOrder, map[string]interface{}{
		"identifiers": identifiers,
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("NewOrder failed: %s", body)
	}
	order := &NewOrderResponse{}
	if err := json.Unmarshal(body, order); err != nil {
		return nil, err
	}
	return order, nil
}

type AuthorizationResponse struct {
	Identifier struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"identifier"`
	Status     string    `json:"status"`
	Expires    time.Time `json:"expires"`
	Challenges []struct {
		Type   string `json:"type"`
		Status string `json:"status"`
		Url    string `json:"url"`
		Token  string `json:"token"`
	} `json:"challenges"`
}

func (c *Client) GetAuthorization(url string) (*AuthorizationResponse, error) {
	body, resp, err := c.do(url, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetAuthorization failed: %s", body)
	}
	auth := &AuthorizationResponse{}
	if err := json.Unmarshal(body, auth); err != nil {
		return nil, err
	}
	return auth, nil
}

type ChallengeResponse struct {
	Type             string `json:"type"`
	Status           string `json:"status"`
	Url              string `json:"url"`
	Token            string `json:"token"`
	ValidationRecord []struct {
		Hostname      string   `json:"hostname"`
		ResolverAddrs []string `json:"resolverAddrs"`
	} `json:"validationRecord"`
	Validated time.Time `json:"validated"`
}

// GetChallenge 获取挑战
func (c *Client) GetChallenge(url string) (*ChallengeResponse, error) {
	body, resp, err := c.do(url, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetChanllenge failed: %s", body)
	}
	challenge := &ChallengeResponse{}
	if err := json.Unmarshal(body, challenge); err != nil {
		return nil, err
	}
	return challenge, nil
}

// RequestChallenge 请求挑战
// 当你当http-01/dns-01记录准备好后，调用此接口
// 然后使用GetChallenge或者GetAuthorization轮询查看状态
func (c *Client) RequestChallenge(url string) (*ChallengeResponse, error) {
	body, resp, err := c.do(url, "{}")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RequestChallenge failed: %s", body)
	}
	challenge := &ChallengeResponse{}
	if err := json.Unmarshal(body, challenge); err != nil {
		return nil, err
	}
	return challenge, nil
}

func (c *Client) thumbprint() string {
	sha256Bytes := sha256.Sum256([]byte(jws.ES256Jwk(c.options.privateKey.PublicKey)))
	return base64.RawURLEncoding.EncodeToString(sha256Bytes[:])
}

func (c *Client) keyAuthorization(token string) string {
	return token + "." + c.thumbprint()
}

func (c *Client) ChallengeRecord(token string) string {
	hash := sha256.Sum256([]byte(c.keyAuthorization(token)))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

type FinalizeResponse struct {
	Status      string    `json:"status"`
	Expires     time.Time `json:"expires"`
	Identifiers []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"identifiers"`
	Authorizations []string `json:"authorizations"`
	Finalize       string   `json:"finalize"`
	Certificate    string   `json:"certificate"`
}

func (c *Client) Finalize(url string, csr []byte) (*FinalizeResponse, error) {
	body, resp, err := c.do(url, map[string]interface{}{
		"csr": base64.RawURLEncoding.EncodeToString(csr),
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FinalizeOrder failed: %s", body)
	}
	finalize := &FinalizeResponse{}
	if err := json.Unmarshal(body, finalize); err != nil {
		return nil, err
	}
	return finalize, nil
}

func (c *Client) GetCertificate(url string) ([]byte, error) {
	body, resp, err := c.do(url, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetCertificate failed: %s", body)
	}
	return body, nil
}
