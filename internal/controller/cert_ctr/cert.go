package cert_ctr

import (
	"archive/zip"
	"context"

	"github.com/codfrm/dns-kit/internal/pkg/utils"
	"github.com/gin-gonic/gin"

	api "github.com/codfrm/dns-kit/internal/api/cert"
	"github.com/codfrm/dns-kit/internal/service/cert_svc"
)

type Cert struct {
}

func NewCert() *Cert {
	return &Cert{}
}

// List 获取证书列表
func (c *Cert) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	return cert_svc.Cert().List(ctx, req)
}

// Create 创建证书
func (c *Cert) Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error) {
	return cert_svc.Cert().Create(ctx, req)
}

// Download 下载证书
func (c *Cert) Download(ctx *gin.Context, req *api.DownloadRequest) error {
	resp, err := cert_svc.Cert().Download(ctx, req)
	if err != nil {
		return err
	}
	// .csr .crt, .key, .pem 文件作为一个zip包
	cert, err := utils.DecodeCertPEM(resp.Cert)
	if err != nil {
		return err
	}
	ctx.Header("Content-Type", "application/zip")
	name := cert.Subject.CommonName
	ctx.Header("Content-Disposition", "attachment; filename="+name+"_cert.zip")
	w := zip.NewWriter(ctx.Writer)
	defer w.Close()
	_ = writeZipFile(w, name+".csr", resp.CSR)
	_ = writeZipFile(w, name+"_bundle.crt", resp.Cert)
	_ = writeZipFile(w, name+".key", resp.Key)
	_ = writeZipFile(w, name+"_bundle.pem", resp.Cert)
	return nil
}

func writeZipFile(w *zip.Writer, name, content string) error {
	file, err := w.Create(name)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(content))
	return err
}

// Delete 删除证书
func (c *Cert) Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error) {
	return cert_svc.Cert().Delete(ctx, req)
}
