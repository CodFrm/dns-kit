package acme_repo

import (
	"context"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/dns-kit/internal/model/entity/acme_entity"
)

type AcmeRepo interface {
	Find(ctx context.Context, id int64) (*acme_entity.Acme, error)
	FindPage(ctx context.Context, page httputils.PageRequest) ([]*acme_entity.Acme, int64, error)
	Create(ctx context.Context, acme *acme_entity.Acme) error
	Update(ctx context.Context, acme *acme_entity.Acme) error
	Delete(ctx context.Context, id int64) error

	FindByEmail(ctx context.Context, email string) (*acme_entity.Acme, error)
}

var defaultAcme AcmeRepo

func Acme() AcmeRepo {
	return defaultAcme
}

func RegisterAcme(i AcmeRepo) {
	defaultAcme = i
}

type acmeRepo struct {
}

func NewAcme() AcmeRepo {
	return &acmeRepo{}
}

func (u *acmeRepo) Find(ctx context.Context, id int64) (*acme_entity.Acme, error) {
	ret := &acme_entity.Acme{}
	if err := db.Ctx(ctx).Where("id=? and status=?", id, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *acmeRepo) Create(ctx context.Context, acme *acme_entity.Acme) error {
	return db.Ctx(ctx).Create(acme).Error
}

func (u *acmeRepo) Update(ctx context.Context, acme *acme_entity.Acme) error {
	return db.Ctx(ctx).Updates(acme).Error
}

func (u *acmeRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Model(&acme_entity.Acme{}).Where("id=?", id).Update("status", consts.DELETE).Error
}

func (u *acmeRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*acme_entity.Acme, int64, error) {
	var list []*acme_entity.Acme
	var count int64
	find := db.Ctx(ctx).Model(&acme_entity.Acme{}).Where("status=?", consts.ACTIVE)
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := find.Order("createtime desc").Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (u *acmeRepo) FindByEmail(ctx context.Context, email string) (*acme_entity.Acme, error) {
	ret := &acme_entity.Acme{}
	if err := db.Ctx(ctx).Where("email=? and status=?", email, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}
