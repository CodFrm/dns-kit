package cert_repo

import (
	"context"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/dns-kit/internal/model/entity/cert_entity"
)

type CertRepo interface {
	Find(ctx context.Context, id int64) (*cert_entity.Cert, error)
	FindPage(ctx context.Context, page httputils.PageRequest) ([]*cert_entity.Cert, int64, error)
	Create(ctx context.Context, cert *cert_entity.Cert) error
	Update(ctx context.Context, cert *cert_entity.Cert) error
	UpdateStatus(ctx context.Context, id int64, status int32) error
	Delete(ctx context.Context, id int64) error

	FindByStatus(ctx context.Context, apply int) ([]*cert_entity.Cert, error)
}

var defaultCert CertRepo

func Cert() CertRepo {
	return defaultCert
}

func RegisterCert(i CertRepo) {
	defaultCert = i
}

type certRepo struct {
}

func NewCert() CertRepo {
	return &certRepo{}
}

func (u *certRepo) Find(ctx context.Context, id int64) (*cert_entity.Cert, error) {
	ret := &cert_entity.Cert{}
	if err := db.Ctx(ctx).Where("id=? and status!=?", id, consts.DELETE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *certRepo) Create(ctx context.Context, cert *cert_entity.Cert) error {
	return db.Ctx(ctx).Create(cert).Error
}

func (u *certRepo) Update(ctx context.Context, cert *cert_entity.Cert) error {
	return db.Ctx(ctx).Updates(cert).Error
}

func (u *certRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Model(&cert_entity.Cert{}).Where("id=?", id).Update("status", consts.DELETE).Error
}

func (u *certRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*cert_entity.Cert, int64, error) {
	var list []*cert_entity.Cert
	var count int64
	find := db.Ctx(ctx).Model(&cert_entity.Cert{}).Where("status!=?", consts.DELETE)
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := find.Order("createtime desc").Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (u *certRepo) FindByStatus(ctx context.Context, apply int) ([]*cert_entity.Cert, error) {
	var list []*cert_entity.Cert
	if err := db.Ctx(ctx).Model(&cert_entity.Cert{}).Where("status=?", apply).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (u *certRepo) UpdateStatus(ctx context.Context, id int64, status int32) error {
	return db.Ctx(ctx).Model(&cert_entity.Cert{}).Where("id=?", id).Update("status", status).Error
}
