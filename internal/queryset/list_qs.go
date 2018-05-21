package queryset

import (
	"context"
	"strings"
	"time"

	"github.com/szyhf/go-rich/internal/log"
	"github.com/szyhf/go-rich/internal/query"
	richTypes "github.com/szyhf/go-rich/types"
)

type ListQuerySet struct {
	*querySet
	ctx         context.Context
	listQuery   *query.ListQuery
	rebuildFunc func(ctx context.Context) ([]interface{}, time.Duration, error)
}

func NewList(key string, q *query.Query) *ListQuerySet {
	return &ListQuerySet{
		querySet:  New(key, q),
		listQuery: query.NewListQuery(q),
	}
}

func (this *ListQuerySet) ScanSlice(slc interface{}, start, stop int64) error {
	intCmd := this.listQuery.Exists(this.Key())
	if intCmd.Val() != 0 {
		slcCmd := this.listQuery.LRange(this.Key(), start, stop)
		if slcCmd.Err() == nil {
			return slcCmd.ScanSlice(slc)
		} else if strings.HasPrefix(slcCmd.Err().Error(), "WRONGTYPE") {
			// 数据库保护产生的空键
			return nil
		}
	}

	if err := this.rebuildingProcess(this); err == nil {
		return this.ScanSlice(slc, start, stop)
	} else {
		return err
	}
}

// ========= 连贯操作接口 =========
// 保护数据库
func (this ListQuerySet) Protect(expire time.Duration) richTypes.ListQuerySeter {
	this.isProtectDB = true
	this.protectExpire = expire
	return &this
}

// 重构List的方法（基于RPush）
func (this ListQuerySet) SetRebuildFunc(f func(ctx context.Context) ([]interface{}, time.Duration, error)) richTypes.ListQuerySeter {
	this.rebuildFunc = f
	return &this
}

func (this ListQuerySet) WithContext(ctx context.Context) richTypes.ListQuerySeter {
	this.ctx = ctx
	return &this
}

// ========= 辅助方法 =========

func (this *ListQuerySet) Rebuilding() error {
	// 重建缓存
	log.Debug("ListQuerySet.rebuild(", this.Key(), ")")
	if valueList, expire, err := this.callRebuildFunc(); err == nil && len(valueList) > 0 {
		intCmd := this.listQuery.RPush(this.Key(), valueList...)
		if intCmd.Err() != nil {
			return intCmd.Err()
		}
		boolCmd := this.listQuery.Expire(this.Key(), expire)
		return boolCmd.Err()
	} else {
		return richTypes.ErrorRebuildNil
	}
}

func (this *ListQuerySet) callRebuildFunc() ([]interface{}, time.Duration, error) {
	if this.rebuildFunc == nil {
		return nil, -1, richTypes.ErrorRebuildNil
	}
	return this.rebuildFunc(this.ctx)
}
