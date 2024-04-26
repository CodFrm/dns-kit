package domain_repo

import (
	"context"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/dns-kit/internal/model/entity/domain_entity"
)

type DomainRepo interface {
	Find(ctx context.Context, id int64) (*domain_entity.Domain, error)
	FindPage(ctx context.Context, page httputils.PageRequest) ([]*domain_entity.Domain, int64, error)
	Create(ctx context.Context, domain *domain_entity.Domain) error
	Update(ctx context.Context, domain *domain_entity.Domain) error
	Delete(ctx context.Context, id int64) error
}

var defaultDomain DomainRepo

func Domain() DomainRepo {
	return defaultDomain
}

func RegisterDomain(i DomainRepo) {
	defaultDomain = i
}

type domainRepo struct {
}

func NewDomain() DomainRepo {
	return &domainRepo{}
}

func (u *domainRepo) Find(ctx context.Context, id int64) (*domain_entity.Domain, error) {
	ret := &domain_entity.Domain{}
	if err := db.Ctx(ctx).Where("id=? and status=?", id, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *domainRepo) Create(ctx context.Context, domain *domain_entity.Domain) error {
	return db.Ctx(ctx).Create(domain).Error
}

func (u *domainRepo) Update(ctx context.Context, domain *domain_entity.Domain) error {
	return db.Ctx(ctx).Updates(domain).Error
}

func (u *domainRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Model(&domain_entity.Domain{}).Where("id=?", id).Update("status", consts.DELETE).Error
}

func (u *domainRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*domain_entity.Domain, int64, error) {
	var list []*domain_entity.Domain
	var count int64
	find := db.Ctx(ctx).Model(&domain_entity.Domain{}).Where("status=?", consts.ACTIVE)
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := find.Order("createtime desc").Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}
