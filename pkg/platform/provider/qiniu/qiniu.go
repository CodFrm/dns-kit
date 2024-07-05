package qiniu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/codfrm/dns-kit/pkg/platform"
	"github.com/qiniu/go-sdk/v7/auth"
)

type Qiniu struct {
	mac       *auth.Credentials
	accessKey string
}

func NewQiniu(accessKey, secretKey string) (platform.CDNManager, error) {
	mac := auth.New(accessKey, secretKey)
	return &Qiniu{
		mac:       mac,
		accessKey: accessKey,
	}, nil
}

type requestOption func(*requestOptions)

type requestOptions struct {
	host   string
	method string
}

func newRequestOptions(opts ...requestOption) *requestOptions {
	r := &requestOptions{
		host:   "api.qiniu.com",
		method: http.MethodGet,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (q *Qiniu) request(ctx context.Context, path string, body, resp any, options ...requestOption) (any, error) {
	opts := newRequestOptions(options...)

	urlStr := fmt.Sprintf("%s%s", "https://"+opts.host, path)

	method := opts.method
	var reqData []byte
	if body != nil {
		reqData, _ = json.Marshal(body)
		if opts.method == http.MethodGet {
			method = http.MethodPost
		}
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, urlStr, bytes.NewReader(reqData))
	if err != nil {
		return nil, err
	}

	accessToken, err := q.mac.SignRequest(httpReq)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Add("Authorization", "QBox "+accessToken)
	httpReq.Header.Add("Content-Type", "application/json")

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	resData, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d, response: %s", httpResp.StatusCode, string(resData))
	}

	if resp != nil {
		err = json.Unmarshal(resData, resp)
		if err != nil {
			return nil, err
		}
	}

	return httpResp, nil
}

func (q *Qiniu) GetCDNList(ctx context.Context) ([]*platform.CDNItem, error) {
	resp := &CDNDomainListResponse{}
	_, err := q.request(ctx, "/domain", nil, resp)
	if err != nil {
		return nil, err
	}
	ret := make([]*platform.CDNItem, 0)
	for _, item := range resp.Domains {
		ret = append(ret, &platform.CDNItem{
			ID:     item.Name,
			Domain: item.Name,
		})
	}
	return ret, nil
}

func (q *Qiniu) GetCDNDetail(ctx context.Context, domain *platform.CDNItem) (*platform.CDNItem, error) {
	resp := &CDNDomainDetailResponse{}
	_, err := q.request(ctx, "/domain/"+domain.ID, nil, resp)
	if err != nil {
		return nil, err
	}
	return domain, nil
}

func (q *Qiniu) SetCDNHttpsCert(ctx context.Context, domain *platform.CDNItem, cert, key string) error {
	// 获取域名信息
	domainDetailResp := &CDNDomainDetailResponse{}
	_, err := q.request(ctx, "/domain/"+domain.ID, nil, domainDetailResp)
	if err != nil {
		return err
	}
	// 上传证书
	resp := &UploadCertResponse{}
	_, err = q.request(ctx, "/sslcert", &UploadCertRequest{
		Name:       "dns-kit-cert-" + strconv.FormatInt(time.Now().Unix(), 10),
		CommonName: domain.ID,
		Pri:        key,
		Ca:         cert,
	}, resp)
	if err != nil {
		return err
	}
	// 设置证书
	_, err = q.request(ctx, "/domain/"+domain.ID+"/sslize", &SetCDNHttpsRequest{
		Certid:      resp.CertID,
		ForceHttps:  domainDetailResp.Https.ForceHttps,
		Http2Enable: domainDetailResp.Https.Http2Enable,
	}, nil, func(options *requestOptions) {
		options.method = http.MethodPut
	})
	return err
}

func (q *Qiniu) UserDetails(ctx context.Context) (*platform.User, error) {
	_, err := q.GetCDNList(ctx)
	if err != nil {
		return nil, err
	}
	return &platform.User{
		ID:       q.accessKey[:8],
		Username: q.accessKey[:8],
	}, nil
}
