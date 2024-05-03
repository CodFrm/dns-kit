package cdn_repo

import (
	"context"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/dns-kit/internal/model/entity/cdn_entity"
)

type CdnRepo interface {
	Find(ctx context.Context, id int64) (*cdn_entity.Cdn, error)
	FindPage(ctx context.Context, page httputils.PageRequest) ([]*cdn_entity.Cdn, int64, error)
	Create(ctx context.Context, cdn *cdn_entity.Cdn) error
	Update(ctx context.Context, cdn *cdn_entity.Cdn) error
	Delete(ctx context.Context, id int64) error
}

var defaultCdn CdnRepo

func Cdn() CdnRepo {
	return defaultCdn
}

func RegisterCdn(i CdnRepo) {
	defaultCdn = i
}

type cdnRepo struct {
}

func NewCdn() CdnRepo {
	return &cdnRepo{}
}

func (u *cdnRepo) Find(ctx context.Context, id int64) (*cdn_entity.Cdn, error) {
	ret := &cdn_entity.Cdn{}
	if err := db.Ctx(ctx).Where("id=? and status=?", id, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *cdnRepo) Create(ctx context.Context, cdn *cdn_entity.Cdn) error {
	return db.Ctx(ctx).Create(cdn).Error
}

func (u *cdnRepo) Update(ctx context.Context, cdn *cdn_entity.Cdn) error {
	return db.Ctx(ctx).Updates(cdn).Error
}

func (u *cdnRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Model(&cdn_entity.Cdn{}).Where("id=?", id).Update("status", consts.DELETE).Error
}

func (u *cdnRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*cdn_entity.Cdn, int64, error) {
	var list []*cdn_entity.Cdn
	var count int64
	find := db.Ctx(ctx).Model(&cdn_entity.Cdn{}).Where("status=?", consts.ACTIVE)
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := find.Order("createtime desc").Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}
