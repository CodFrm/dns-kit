package acme

import (
	"crypto/ecdsa"
	"net/http"
)

type ClientOptions struct {
	directoryUrl string
	httpClient   *http.Client
	directory    *Directory
	privateKey   *ecdsa.PrivateKey
	kid          string
}

type ClientOption func(*ClientOptions)

func newClientOptions(options ...ClientOption) *ClientOptions {
	opts := &ClientOptions{
		directoryUrl: DefaultDirectoryUrl,
		httpClient:   http.DefaultClient,
	}
	for _, o := range options {
		o(opts)
	}
	return opts
}

func WithDirectoryUrl(directoryUrl string) ClientOption {
	return func(options *ClientOptions) {
		options.directoryUrl = directoryUrl
	}
}

func WithHttpClient(httpClient *http.Client) ClientOption {
	return func(options *ClientOptions) {
		options.httpClient = httpClient
	}
}

func WithPrivateKey(privateKey *ecdsa.PrivateKey) ClientOption {
	return func(options *ClientOptions) {
		options.privateKey = privateKey
	}
}

func WithDirectory(directory *Directory) ClientOption {
	return func(options *ClientOptions) {
		options.directory = directory
	}
}

func WithKid(kid string) ClientOption {
	return func(options *ClientOptions) {
		options.kid = kid
	}
}

type Option func(*Options)

type Options struct {
	client *Client
}

func WithClient(client *Client) Option {
	return func(options *Options) {
		options.client = client
	}
}

func newOptions(opts ...Option) (*Options, error) {
	opt := &Options{}
	for _, o := range opts {
		o(opt)
	}
	if opt.client == nil {
		client, err := NewClient()
		if err != nil {
			return nil, err
		}
		opt.client = client
	}
	return opt, nil
}
