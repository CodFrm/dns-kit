package handler

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/codfrm/dns-kit/internal/service/cert_svc"

	"github.com/codfrm/dns-kit/internal/model/entity/domain_entity"

	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/dns-kit/internal/model/entity/cert_entity"
	"github.com/codfrm/dns-kit/internal/pkg/utils"
	"github.com/codfrm/dns-kit/internal/repository/cert_repo"
	"github.com/codfrm/dns-kit/internal/repository/domain_repo"
	"github.com/codfrm/dns-kit/internal/task/queue"
	"github.com/codfrm/dns-kit/internal/task/queue/message"
	"github.com/codfrm/dns-kit/pkg/platform"
	"go.uber.org/zap"
)

type CertHandler struct {
}

func (c *CertHandler) Register(ctx context.Context) error {
	if err := queue.SubscribeCertCreate(ctx, c.CreateCert); err != nil {
		return err
	}
	return nil
}

func (c *CertHandler) CreateCert(ctx context.Context, msg *message.CreateCertMessage) error {
	cert, err := cert_repo.Cert().Find(ctx, msg.ID)
	if err != nil {
		return err
	}
	logger := logger.Ctx(ctx).With(zap.Int64("cert_id", cert.ID))
	if err := cert.Check(ctx); err != nil {
		logger.Error("cert check failed", zap.Error(err))
		return err
	}
	if cert.Status != cert_entity.CertStatusApply {
		logger.Error("cert status is not apply", zap.Int32("status", int32(cert.Status)))
		return errors.New("cert status is not apply")
	}
	// 后续的错误都更新为申请失败
	defer func() {
		msg := &message.CreateCertAfterMessage{
			ID:      msg.ID,
			Success: true,
		}
		if err != nil {
			msg.Success = false
			if err := cert_repo.Cert().UpdateStatus(ctx, cert.ID, cert_entity.CertStatusApplyFail); err != nil {
				logger.Error("update cert failed", zap.Error(err))
			}
		}
		if err := queue.PublishCertCreateAfter(ctx, msg); err != nil {
			logger.Error("publish cert create after failed", zap.Error(err))
		}
	}()
	acmeInstance, err := cert_svc.Cert().NewACME(ctx, cert.Email)
	if err != nil {
		logger.Error("new acme failed", zap.Error(err))
		return err
	}
	challenges, err := acmeInstance.GetChallenge(ctx, cert.Domains)
	if err != nil {
		logger.Error("get challenges failed", zap.Error(err))
		return err
	}
	logger.Info("get challenges success", zap.Any("challenges", challenges))
	// 设置dns记录
	wg := sync.WaitGroup{}
	for _, v := range challenges {
		logger := logger.With(zap.String("domain", v.Domain))
		// 获取tld
		var tld string
		tld, err = utils.GetTLD(v.Domain)
		if err != nil {
			logger.Error("get tld failed", zap.Error(err))
			return err
		}
		// 设置dns解析
		var domain *domain_entity.Domain
		domain, err = domain_repo.Domain().FindByDomain(ctx, tld)
		if err != nil {
			logger.Error("find domain failed", zap.Error(err))
			return err
		}
		if err = domain.Check(ctx); err != nil {
			logger.Error("domain check failed", zap.Error(err))
			return err
		}
		var manager platform.DNSManager
		manager, err = domain.DnsManager(ctx)
		if err != nil {
			logger.Error("domain factory failed", zap.Error(err))
			return err
		}
		record := &platform.Record{
			Type:  "TXT",
			Value: v.Record,
		}
		// 获取记录名
		if v.Domain == tld {
			record.Name = "_acme-challenge"
		} else {
			record.Name = "_acme-challenge." + strings.TrimSuffix(v.Domain, "."+tld)
		}
		// 删除老的记录
		var recordList []*platform.Record
		recordList, err = manager.GetRecordList(ctx)
		if err != nil {
			logger.Error("get record list failed", zap.Error(err))
			return err
		}
		for _, v := range recordList {
			if v.Name == record.Name {
				if err = manager.DelRecord(ctx, v.ID); err != nil {
					logger.Error("del record failed", zap.String("value", v.Value), zap.Error(err))
					return err
				}
			}
		}
		if err = manager.AddRecord(ctx, record); err != nil {
			logger.Error("add record failed", zap.Error(err))
			return err
		}
		defer func() {
			if err := manager.DelRecord(ctx, record.ID); err != nil {
				logger.Error("del record failed", zap.String("value", record.Value), zap.Error(err))
			}
		}()
		wg.Add(1)
		// 等待dns记录更新
		go func() {
			defer wg.Done()
			sctx, cancel := context.WithTimeout(ctx, time.Minute*5)
			defer cancel()
			equalNum := 0
			for {
				select {
				case <-sctx.Done():
					return
				default:
					time.Sleep(time.Second * 5)
					list, err := net.LookupTXT(record.Name + "." + tld)
					if err != nil {
						equalNum = 0
						logger.Error("lookup txt failed", zap.Error(err))
						continue
					}
					// 判断是否有记录且只有一条
					if len(list) != 1 {
						equalNum = 0
						continue
					}
					// 连续3次记录相等
					if list[0] == v.Record {
						if equalNum++; equalNum == 3 {
							return
						}
					} else {
						equalNum = 0
					}
				}
			}
		}()
	}
	wg.Wait()
	// 等待申请完成
	var cancel context.CancelFunc
	sctx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()
	logger.Info("wait challenge")
	if err = acmeInstance.WaitChallenge(sctx); err != nil {
		logger.Error("wait challenge failed", zap.Error(err))
		return err
	}
	// 获取证书
	logger.Info("get certificate")
	certData, err := acmeInstance.GetCertificate(ctx)
	if err != nil {
		logger.Error("get certificate failed", zap.Error(err))
		return err
	}
	logger.Info("get certificate success")
	// 解析证书获取到期时间
	certInfo, err := utils.DecodeCertPEM(string(certData))
	if err != nil {
		logger.Error("parse certificate failed", zap.Error(err))
		return err
	}
	cert.Expiretime = certInfo.NotAfter.Unix()
	// 保存证书与更新状态
	cert.Certificate = string(certData)
	buf := bytes.NewBuffer(nil)
	if err = pem.Encode(buf, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: acmeInstance.GetCSR()}); err != nil {
		logger.Error("encode csr failed", zap.Error(err))
		return err
	}
	cert.CertificateRequest = buf.String()
	privateKey, err := x509.MarshalECPrivateKey(acmeInstance.GetPrivateKey())
	if err != nil {
		logger.Error("marshal private key failed", zap.Error(err))
		return err
	}
	buf.Reset()
	if err = pem.Encode(buf, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privateKey}); err != nil {
		logger.Error("encode private key failed", zap.Error(err))
		return err
	}
	cert.PrivateKey = buf.String()
	cert.Status = cert_entity.CertStatusActive
	cert.Updatetime = time.Now().Unix()
	if err = cert_repo.Cert().Update(ctx, cert); err != nil {
		logger.Error("update cert failed", zap.Error(err))
		return err
	}
	return nil
}
