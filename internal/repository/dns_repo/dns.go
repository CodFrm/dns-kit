package dns_repo

import (
	"context"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/dns-kit/internal/model/entity/dns_entity"
)

type DnsRepo interface {
	Find(ctx context.Context, id int64) (*dns_entity.Dns, error)
	FindPage(ctx context.Context, page httputils.PageRequest) ([]*dns_entity.Dns, int64, error)
	Create(ctx context.Context, dns *dns_entity.Dns) error
	Update(ctx context.Context, dns *dns_entity.Dns) error
	Delete(ctx context.Context, id int64) error
}

var defaultDns DnsRepo

func Dns() DnsRepo {
	return defaultDns
}

func RegisterDns(i DnsRepo) {
	defaultDns = i
}

type dnsRepo struct {
}

func NewDns() DnsRepo {
	return &dnsRepo{}
}

func (u *dnsRepo) Find(ctx context.Context, id int64) (*dns_entity.Dns, error) {
	ret := &dns_entity.Dns{}
	if err := db.Ctx(ctx).Where("id=? and status=?", id, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *dnsRepo) Create(ctx context.Context, dns *dns_entity.Dns) error {
	return db.Ctx(ctx).Create(dns).Error
}

func (u *dnsRepo) Update(ctx context.Context, dns *dns_entity.Dns) error {
	return db.Ctx(ctx).Updates(dns).Error
}

func (u *dnsRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Model(&dns_entity.Dns{}).Where("id=?", id).Update("status", consts.DELETE).Error
}

func (u *dnsRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*dns_entity.Dns, int64, error) {
	var list []*dns_entity.Dns
	var count int64
	find := db.Ctx(ctx).Model(&dns_entity.Dns{}).Where("status=?", consts.ACTIVE)
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := find.Order("createtime desc").Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}
