package queryset

import (
	"context"
	"time"

	"github.com/szyhf/go-rich/internal/log"
	"github.com/szyhf/go-rich/internal/query"
	richTypes "github.com/szyhf/go-rich/types"
)

type StringQuerySet struct {
	*querySet
	ctx         context.Context
	stringQuery *query.StringQuery
	rebuildFunc func(ctx context.Context) (interface{}, time.Duration)
}

func NewString(key string, q *query.Query) *StringQuerySet {
	return &StringQuerySet{
		querySet:    New(key, q),
		stringQuery: query.NewStringQuery(q),
	}
}

// ======== 读取接口 ========

func (this *StringQuerySet) Get() (string, error) {
	// 尝试直接从缓存获取
	cmd := this.stringQuery.Get(this.Key())
	if cmd.Err() == nil {
		return cmd.Val(), nil
	}

	// 尝试重建缓存
	if err := this.rebuildingProcess(this); err == nil {
		return this.Get()
	} else {
		return "", err
	}
}

func (this *StringQuerySet) Scan(value interface{}) error {
	// 尝试直接从缓存获取
	cmd := this.stringQuery.Get(this.Key())
	if cmd.Err() == nil {
		return cmd.Scan(value)
	}

	if err := this.rebuildingProcess(this); err == nil {
		return this.Scan(value)
	} else {
		return err
	}
}

// ========= 写入接口 =========
// 写入当前key
func (this *StringQuerySet) Set(value interface{}, expire time.Duration) error {
	cmd := this.stringQuery.Set(this.Key(), value, expire)
	return cmd.Err()
}

// 移除当前key
func (this *StringQuerySet) Del() error {
	cmd := this.stringQuery.Del(this.Key())
	return cmd.Err()
}

// 如果key存在，则给当前key增长指定的值
func (this *StringQuerySet) IncrBy(incr int64) (int64, error) {
	val, err := this.stringQuery.
		IncrByIfExist(this.Key(), incr)
	if err == nil {
		return val, nil
	}

	if err := this.rebuildingProcess(this); err == nil {
		return this.IncrBy(incr)
	} else {
		return 0, err
	}
}

// ========= 连贯操作接口 =========
// 保护数据库
func (this StringQuerySet) Protect(expire time.Duration) richTypes.StringQuerySeter {
	this.isProtectDB = true
	this.protectExpire = expire
	return &this
}

// 重构String的方法
func (this StringQuerySet) SetRebuildFunc(f func(ctx context.Context) (interface{}, time.Duration)) richTypes.StringQuerySeter {
	this.rebuildFunc = f
	return &this
}

func (this StringQuerySet) WithContext(ctx context.Context) richTypes.StringQuerySeter {
	this.ctx = ctx
	return &this
}

// ========= 辅助方法 =========

func (this *StringQuerySet) Rebuilding() error {
	// 重建缓存
	log.Debug("StringQuerySet.rebuild(", this.Key(), ")")
	if value, expire := this.callRebuildFunc(); value != nil {
		cmd := this.stringQuery.Set(this.Key(), value, expire)
		return cmd.Err()
	}
	return richTypes.ErrorRebuildNil
}

func (this *StringQuerySet) callRebuildFunc() (interface{}, time.Duration) {
	if this.rebuildFunc == nil {
		return nil, -1
	}
	return this.rebuildFunc(this.ctx)
}
