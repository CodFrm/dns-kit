package cert_hosting_repo

import (
	"context"
	"github.com/codfrm/cago/pkg/logger"
	"go.uber.org/zap"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/dns-kit/internal/model/entity/cert_hosting_entity"
)

type CertHostingRepo interface {
	Find(ctx context.Context, id int64) (*cert_hosting_entity.CertHosting, error)
	FindPage(ctx context.Context, page httputils.PageRequest) ([]*cert_hosting_entity.CertHosting, int64, error)
	Create(ctx context.Context, certHosting *cert_hosting_entity.CertHosting) error
	Update(ctx context.Context, certHosting *cert_hosting_entity.CertHosting) error
	Delete(ctx context.Context, id int64) error

	UpdateStatus(ctx context.Context, id int64, status cert_hosting_entity.CertHostingStatus) error
	FindByCDN(ctx context.Context, cdnID int64) ([]*cert_hosting_entity.CertHosting, error)
	FindByCert(ctx context.Context, certID int64) ([]*cert_hosting_entity.CertHosting, error)
}

var defaultCertHosting CertHostingRepo

func CertHosting() CertHostingRepo {
	return defaultCertHosting
}

func RegisterCertHosting(i CertHostingRepo) {
	defaultCertHosting = i
}

type certHostingRepo struct {
}

func NewCertHosting() CertHostingRepo {
	return &certHostingRepo{}
}

func (u *certHostingRepo) Find(ctx context.Context, id int64) (*cert_hosting_entity.CertHosting, error) {
	ret := &cert_hosting_entity.CertHosting{}
	if err := db.Ctx(ctx).Where("id=? and status=?", id, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *certHostingRepo) Create(ctx context.Context, certHosting *cert_hosting_entity.CertHosting) error {
	return db.Ctx(ctx).Create(certHosting).Error
}

func (u *certHostingRepo) Update(ctx context.Context, certHosting *cert_hosting_entity.CertHosting) error {
	return db.Ctx(ctx).Updates(certHosting).Error
}

func (u *certHostingRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Model(&cert_hosting_entity.CertHosting{}).Where("id=?", id).Update("status", consts.DELETE).Error
}

func (u *certHostingRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*cert_hosting_entity.CertHosting, int64, error) {
	var list []*cert_hosting_entity.CertHosting
	var count int64
	find := db.Ctx(ctx).Model(&cert_hosting_entity.CertHosting{}).Where("status=?", consts.ACTIVE)
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := find.Order("createtime desc").Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (u *certHostingRepo) FindByCDN(ctx context.Context, cdnID int64) ([]*cert_hosting_entity.CertHosting, error) {
	var list []*cert_hosting_entity.CertHosting
	if err := db.Ctx(ctx).Where("cdn_id=? and status!=?", cdnID, consts.DELETE).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (u *certHostingRepo) FindByCert(ctx context.Context, certID int64) ([]*cert_hosting_entity.CertHosting, error) {
	var list []*cert_hosting_entity.CertHosting
	if err := db.Ctx(ctx).Where("cert_id=? and status!=?", certID, consts.DELETE).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (u *certHostingRepo) UpdateStatus(ctx context.Context, id int64, status cert_hosting_entity.CertHostingStatus) error {
	if err := db.Ctx(ctx).Model(&cert_hosting_entity.CertHosting{}).Where("id=?", id).Update("status", status).Error; err != nil {
		logger.Ctx(ctx).Error("update status error", zap.Int64("id", id), zap.Error(err))
		return err
	}
	return nil
}
