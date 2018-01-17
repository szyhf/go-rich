package queryset

import (
	"context"
	"time"

	"github.com/go-redis/redis"

	"github.com/szyhf/go-rich/internal/log"
	"github.com/szyhf/go-rich/internal/query"
	richTypes "github.com/szyhf/go-rich/types"
)

type ZSetQuerySet struct {
	*querySet
	ctx         context.Context
	zsetQuery   *query.ZSetQuery
	rebuildFunc func(ctx context.Context) ([]redis.Z, time.Duration)
}

func NewZSet(key string, q *query.Query) *ZSetQuerySet {
	return &ZSetQuerySet{
		querySet:  New(key, q),
		zsetQuery: query.NewZSetQuery(q),
	}
}

// ========= 查询接口 =========

func (this *ZSetQuerySet) Count() (int64, error) {
	// 尝试直接从缓存拿
	count, err := this.zsetQuery.ZCardIfExist(this.Key())
	if err == nil {
		return count, nil
	}

	// 重建缓存
	if err := this.rebuildingProcess(this); err == nil {
		// 重建成功则重新获取
		return this.Count()
	} else {
		return 0, err
	}
}

func (this *ZSetQuerySet) Score(member string) (float64, error) {
	score, err := this.zsetQuery.ZScoreIfExist(this.Key(), member)
	if err == nil {
		return score, nil
	}

	if err := this.rebuildingProcess(this); err == nil {
		return this.Score(member)
	} else {
		return 0, err
	}
}

func (this *ZSetQuerySet) IsMember(member string) (bool, error) {
	// 尝试直接从缓存拿
	exist, err := this.zsetQuery.ZIsMemberIfExist(this.Key(), member)
	if err == nil {
		return exist, nil
	}

	// 重建缓存
	if err := this.rebuildingProcess(this); err == nil {
		return this.IsMember(member)
	} else {
		return false, err
	}
}

func (this *ZSetQuerySet) RangeASC(start, stop int64) ([]string, error) {
	// 尝试直接从缓存拿
	members, err := this.zsetQuery.ZRangeIfExist(this.Key(), start, stop)
	if err == nil {
		return members, nil
	}

	// 缓存获取失败尝试重构缓存
	if err := this.rebuildingProcess(this); err == nil {
		return this.RangeASC(start, stop)
	} else {
		return nil, err
	}
}

func (this *ZSetQuerySet) RangeDESC(start, stop int64) ([]string, error) {
	// 尝试直接从缓存拿
	members, err := this.zsetQuery.ZRevRangeIfExist(this.Key(), start, stop)
	if err == nil {
		return members, nil
	}

	// 缓存获取失败尝试重构缓存
	if err := this.rebuildingProcess(this); err == nil {
		return this.RangeDESC(start, stop)
	} else {
		return nil, err
	}
}

func (this *ZSetQuerySet) Members() ([]string, error) {
	// 利用Range的负数参数指向倒数的元素的特性实现
	return this.RangeASC(0, -1)
}

func (this *ZSetQuerySet) RangeByScoreASC(min, max string, offset, count int64) ([]string, error) {
	members, err := this.zsetQuery.ZRangeByScoreIfExist(this.Key(), redis.ZRangeBy{
		Max:    max,
		Min:    min,
		Offset: offset,
		Count:  count,
	})
	if err == nil {
		return members, nil
	}

	if err := this.rebuildingProcess(this); err == nil {
		return this.RangeByScoreASC(min, max, offset, count)
	} else {
		return nil, err
	}
}

func (this *ZSetQuerySet) RangeByScoreDESC(min, max string, offset, count int64) ([]string, error) {
	members, err := this.zsetQuery.ZRevRangeByScoreIfExist(this.Key(), redis.ZRangeBy{
		Max:    max,
		Min:    min,
		Offset: offset,
		Count:  count,
	})
	if err == nil {
		return members, nil
	}

	if err := this.rebuildingProcess(this); err == nil {
		return this.RangeByScoreDESC(min, max, offset, count)
	} else {
		return nil, err
	}
}

func (this *ZSetQuerySet) RangeASCWithScores(start, stop int64) ([]redis.Z, error) {
	members, err := this.zsetQuery.ZRangeWithScoresIfExist(this.Key(), start, stop)
	if err == nil {
		return members, nil
	}

	if err := this.rebuildingProcess(this); err == nil {
		return this.RangeASCWithScores(start, stop)
	} else {
		return nil, err
	}
}

func (this *ZSetQuerySet) RangeDESCWithScores(start, stop int64) ([]redis.Z, error) {
	members, err := this.zsetQuery.ZRevRangeWithScoresIfExist(this.Key(), start, stop)
	if err == nil {
		return members, nil
	}

	if err := this.rebuildingProcess(this); err == nil {
		return this.RangeDESCWithScores(start, stop)
	} else {
		return nil, err
	}
}

// ========= 写入接口 =========

func (this *ZSetQuerySet) AddExpire(member interface{}, score float64, expire time.Duration) (int64, error) {
	// 如果不增加过期方法，可能会创建一个不会过期的集合
	num, err := this.zsetQuery.
		ZAddExpireIfExist(this.Key(),
			[]redis.Z{
				redis.Z{
					Member: member,
					Score:  score,
				}},
			expire)
	if err == nil {
		return num, nil
	}

	if err := this.rebuildingProcess(this); err == nil {
		return this.AddExpire(member, score, expire)
	} else {
		return 0, err
	}
}

func (this *ZSetQuerySet) Rem(member ...interface{}) error {
	cmd := this.zsetQuery.ZRem(this.Key(), member...)
	return cmd.Err()
}

// ============= 连贯操作 =============

// 防止频繁重建
// expire 保护有效时间
func (this ZSetQuerySet) Protect(expire time.Duration) richTypes.ZSetQuerySeter {
	this.isProtectDB = true
	this.protectExpire = expire
	return &this
}

func (this ZSetQuerySet) SetRebuildFunc(rebuildFunc func(context.Context) ([]redis.Z, time.Duration)) richTypes.ZSetQuerySeter {
	this.rebuildFunc = rebuildFunc
	return &this
}

func (this ZSetQuerySet) WithContext(ctx context.Context) richTypes.ZSetQuerySeter {
	this.ctx = ctx
	return &this
}

func (this ZSetQuerySet) Querier() richTypes.Querier {
	return this.zsetQuery
}

// ============= 辅助方法 =============

func (this *ZSetQuerySet) Rebuilding() error {
	// 重建缓存
	log.Notice("zsetQuerySet.rebuild(", this.Key(), ")")
	// 见 issue#1，移除可能存在的保护键
	cmd := this.Querier().Del(this.Key())

	if members, expire := this.callRebuildFunc(); len(members) > 0 {
		if cmd.Err() == nil {
			_, err := this.zsetQuery.ZAddExpire(this.Key(), members, expire)
			return err
		}
		return cmd.Err()
	}
	return richTypes.ErrorCanNotRebuild
}

func (this *ZSetQuerySet) callRebuildFunc() ([]redis.Z, time.Duration) {
	if this.rebuildFunc == nil {
		return []redis.Z{}, -1
	}
	return this.rebuildFunc(this.ctx)
}
