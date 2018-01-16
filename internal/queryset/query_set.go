package queryset

import (
	"time"

	"github.com/szyhf/go-rich/internal/log"
	"github.com/szyhf/go-rich/internal/query"
	richTypes "github.com/szyhf/go-rich/types"
)

type querySet struct {
	*query.Query

	key           string
	isProtectDB   bool
	protectExpire time.Duration

	isRebuilding bool
}

func New(key string, q *query.Query) *querySet {
	return &querySet{
		Query: q,
		key:   key,
	}
}

func (this querySet) Protect(expire time.Duration) *querySet {
	this.isProtectDB = true
	this.protectExpire = expire
	return &this
}

func (this *querySet) Key() string {
	return this.key
}

func (this *querySet) Rebuilding() error {
	panic(`Should implement method "Rebuilding" in sub queryset`)
}

func (this *querySet) Querier() richTypes.Querier {
	return this.Query
}

func (this *querySet) tryGetRebuildLock(key string) bool {
	log.Notice("tryGetRebuildLock:", key)
	// 通过setNX设置锁，同设置超时，防止del失败
	if cmd := this.Querier().SetNX(key+":mutex", "", 30*time.Second); cmd.Err() == nil {
		return cmd.Val()
	} else {
		log.Warn("querySet.TryGetRebuildLock(", key, ") failed: ", cmd.Err())
	}
	return false
}

func (this *querySet) tryReleaseRebuildLock(key string) bool {
	log.Notice("tryReleaseRebuildLock:", key)
	if cmd := this.Querier().Del(key + ":mutex"); cmd.Err() == nil {
		return true
	} else {
		log.Warn("querySet.TryReleaseRebuildLock(", key, ") failed: ", cmd.Err())
	}

	return false
}

func (this *querySet) tryProtectDB(key string) bool {
	cmd := this.Querier().Set(key, nil, this.protectExpire)
	log.Notice("tryProtectDB:", key, "for", this.protectExpire, "seconds.")
	return cmd.Err() == nil
}

func (this *querySet) rebuildingProcess(qs richTypes.QuerySeter) bool {
	if this.isRebuilding {
		log.Warn("Rebuilding break for dead loop.")
		// 防止重构缓存失败陷入死循环
		return false
	}

	// 获取缓存重建锁
	if this.tryGetRebuildLock(this.Key()) {
		this.isRebuilding = true
		defer this.tryReleaseRebuildLock(this.Key())
		if err := qs.Rebuilding(); err != nil {
			// 失败了，建立缓存保护盾保护DB
			if this.isProtectDB {
				this.tryProtectDB(this.Key())
			}
		} else {
			return true
		}
	}
	return false
}
