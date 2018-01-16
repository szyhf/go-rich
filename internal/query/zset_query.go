package query

import (
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/szyhf/go-rich/internal/log"
	"github.com/szyhf/go-rich/types"
)

type ZSetQuery struct {
	*Query
}

func NewZSetQuery(q *Query) *ZSetQuery {
	return &ZSetQuery{
		Query: q,
	}
}

// ========== 写入操作 ==============

// 使用pipline实现的带过期时间的ZAdd
func (this *ZSetQuery) ZAddExpire(key string, members []redis.Z, expire time.Duration) (int64, error) {
	log.Notice("[Redis ZAddExpire]", key, members, expire)
	cmds, err := this.ExecPipeline(func(pipe redis.Pipeliner) error {
		pipe.ZAdd(key, members...)
		pipe.Expire(key, expire)
		return nil
	})
	if err == nil {
		return cmds[0].(*redis.IntCmd).Val(), nil
	}

	return 0, err
}

// 使用pipline实现的带过期时间的ZAdd（仅当key存在时添加）
func (this *ZSetQuery) ZAddExpireIfExist(key string, members []redis.Z, expire time.Duration) (int64, error) {
	log.Notice("[Redis ZAddExpireIfExist]", key, members, expire)
	cmds, _ := this.ExecPipeline(func(pipe redis.Pipeliner) error {
		pipe.Exists(key)
		pipe.ZAdd(key, members...)
		pipe.Expire(key, expire)
		return nil
	})
	// Pipeline默认返回的是最后一个err，所以这里的判定方式要做调整
	if cmds[0].Err() != nil {
		return 0, cmds[0].Err()
	}
	if cmds[0].(*redis.IntCmd).Val() == 1 {
		if cmds[1].Err() == nil {
			return cmds[1].(*redis.IntCmd).Val(), nil
		} else {
			return 0, cmds[1].Err()
		}
	} else {
		return 0, richTypes.ErrorKeyNotExist
	}
}

// ========== 读取操作 ==============

// 使用Pipline实现的优先检查存在性的ZCard
func (this *ZSetQuery) ZCardIfExist(key string) (int64, error) {
	log.Notice("[Redis ZCardIfExist]", key)
	cmds, _ := this.ExecPipeline(func(pipe redis.Pipeliner) error {
		pipe.Exists(key)
		pipe.ZCard(key)
		return nil
	})
	// Pipeline默认返回的是最后一个err，所以这里的判定方式要做调整
	if cmds[0].Err() != nil {
		return 0, cmds[0].Err()
	}
	if cmds[0].(*redis.IntCmd).Val() == 1 {
		// utils.Display("cmd1", cmds[1])
		if cmds[1].Err() == nil {
			return cmds[1].(*redis.IntCmd).Val(), nil
		} else if strings.HasPrefix(cmds[1].Err().Error(), "WRONGTYPE") {
			// 数据库保护产生的空键
			return 0, nil
		} else {
			return 0, cmds[1].Err()
		}
	} else {
		return 0, richTypes.ErrorKeyNotExist
	}
}

// 判定Key是否存在，如果存在则返回指定排序区间的成员（正序）
func (this *ZSetQuery) ZRangeIfExist(key string, start, stop int64) ([]string, error) {
	log.Notice("[Redis ZRangeIfExist]", key, start, stop)
	cmds, _ := this.ExecPipeline(func(pipe redis.Pipeliner) error {
		pipe.Exists(key)
		pipe.ZRange(key, start, stop)
		return nil
	})
	// Pipeline默认返回的是最后一个err，所以这里的判定方式要做调整
	if cmds[0].Err() != nil {
		return nil, cmds[0].Err()
	}
	if cmds[0].(*redis.IntCmd).Val() == 1 {
		if cmds[1].Err() == nil {
			return cmds[1].(*redis.StringSliceCmd).Val(), nil
		} else if strings.HasPrefix(cmds[1].Err().Error(), "WRONGTYPE") {
			// 数据库保护产生的空键
			return nil, nil
		} else {
			return nil, cmds[1].Err()
		}
	} else {
		return nil, richTypes.ErrorKeyNotExist
	}
}

// 判定Key是否存在，如果存在则返回指定排序区间的成员（逆序）
func (this *ZSetQuery) ZRevRangeIfExist(key string, start, stop int64) ([]string, error) {
	log.Notice("[Redis ZRevRangeIfExist]", key, start, stop)
	cmds, _ := this.ExecPipeline(func(pipe redis.Pipeliner) error {
		pipe.Exists(key)
		pipe.ZRevRange(key, start, stop)
		return nil
	})
	// Pipeline默认返回的是最后一个err，所以这里的判定方式要做调整
	if cmds[0].Err() != nil {
		return nil, cmds[0].Err()
	}
	if cmds[0].(*redis.IntCmd).Val() == 1 {
		if cmds[1].Err() == nil {
			return cmds[1].(*redis.StringSliceCmd).Val(), nil
		} else if strings.HasPrefix(cmds[1].Err().Error(), "WRONGTYPE") {
			// 数据库保护产生的空键
			return nil, nil
		} else {
			return nil, cmds[1].Err()
		}
	} else {
		return nil, richTypes.ErrorKeyNotExist
	}
}

func (this *ZSetQuery) ZRangeWithScoresIfExist(key string, start, stop int64) ([]redis.Z, error) {
	log.Notice("[Redis ZRangeWithScoresIfExist]", key, start, stop)
	cmds, _ := this.ExecPipeline(func(pipe redis.Pipeliner) error {
		pipe.Exists(key)
		pipe.ZRangeWithScores(key, start, stop)
		return nil
	})
	// Pipeline默认返回的是最后一个err，所以这里的判定方式要做调整
	if cmds[0].Err() != nil {
		return nil, cmds[0].Err()
	}
	if cmds[0].(*redis.IntCmd).Val() == 1 {
		if cmds[1].Err() == nil {
			return cmds[1].(*redis.ZSliceCmd).Val(), nil
		} else if strings.HasPrefix(cmds[1].Err().Error(), "WRONGTYPE") {
			// 数据库保护产生的空键
			return nil, nil
		} else {
			return nil, cmds[1].Err()
		}
	} else {
		return nil, richTypes.ErrorKeyNotExist
	}
}

func (this *ZSetQuery) ZRevRangeWithScoresIfExist(key string, start, stop int64) ([]redis.Z, error) {
	log.Notice("[Redis ZRevRangeWithScoresIfExist]", key, start, stop)
	cmds, _ := this.ExecPipeline(func(pipe redis.Pipeliner) error {
		pipe.Exists(key)
		pipe.ZRevRangeWithScores(key, start, stop)
		return nil
	})
	// Pipeline默认返回的是最后一个err，所以这里的判定方式要做调整
	if cmds[0].Err() != nil {
		return nil, cmds[0].Err()
	}
	if cmds[0].(*redis.IntCmd).Val() == 1 {
		if cmds[1].Err() == nil {
			return cmds[1].(*redis.ZSliceCmd).Val(), nil
		} else if strings.HasPrefix(cmds[1].Err().Error(), "WRONGTYPE") {
			// 数据库保护产生的空键
			return nil, nil
		} else {
			return nil, cmds[1].Err()
		}
	} else {
		return nil, richTypes.ErrorKeyNotExist
	}
}

func (this *ZSetQuery) ZRangeByScoreIfExist(key string, opt redis.ZRangeBy) ([]string, error) {
	log.Notice("[Redis ZRangeByScoreIfExist]", key, opt)

	cmds, _ := this.ExecPipeline(func(pipe redis.Pipeliner) error {
		pipe.Exists(key)
		pipe.ZRangeByScore(key, opt)
		return nil
	})
	// Pipeline默认返回的是最后一个err，所以这里的判定方式要做调整
	if cmds[0].Err() != nil {
		return nil, cmds[0].Err()
	}
	if cmds[0].(*redis.IntCmd).Val() == 1 {
		if cmds[1].Err() == nil {
			return cmds[1].(*redis.StringSliceCmd).Val(), nil
		} else if strings.HasPrefix(cmds[1].Err().Error(), "WRONGTYPE") {
			// 数据库保护产生的空键
			return nil, nil
		} else {
			return nil, cmds[1].Err()
		}
	} else {
		return nil, richTypes.ErrorKeyNotExist
	}
}
func (this *ZSetQuery) ZRevRangeByScoreIfExist(key string, opt redis.ZRangeBy) ([]string, error) {
	log.Notice("[Redis ZRevRangeByScoreIfExist]", key, opt)
	cmds, _ := this.ExecPipeline(func(pipe redis.Pipeliner) error {
		pipe.Exists(key)
		pipe.ZRevRangeByScore(key, opt)
		return nil
	})
	// Pipeline默认返回的是最后一个err，所以这里的判定方式要做调整
	if cmds[0].Err() != nil {
		return nil, cmds[0].Err()
	}
	if cmds[0].(*redis.IntCmd).Val() == 1 {
		if cmds[1].Err() == nil {
			return cmds[1].(*redis.StringSliceCmd).Val(), nil
		} else if strings.HasPrefix(cmds[1].Err().Error(), "WRONGTYPE") {
			// 数据库保护产生的空键
			return nil, nil
		} else {
			return nil, cmds[1].Err()
		}
	} else {
		return nil, richTypes.ErrorKeyNotExist
	}
}

// 判定Key是否存在，如果存在则检查member是否在集合中
func (this *ZSetQuery) ZIsMemberIfExist(key string, member string) (bool, error) {
	log.Notice("[Redis ZIsMemberIfExist]", key, member)
	// 通过ZRank间接实现存在性判断
	// ZScore返回member在ZSet中的Index
	cmds, _ := this.ExecPipeline(func(pipe redis.Pipeliner) error {
		pipe.Exists(key)
		pipe.ZScore(key, member)
		return nil
	})

	// Pipeline默认返回的是最后一个err，所以这里的判定方式要做调整
	if cmds[0].Err() != nil {
		return false, cmds[0].Err()
	}
	if cmds[0].(*redis.IntCmd).Val() == 1 {
		// 如果member不存在，则会返回error=redis.Nil
		if cmds[1].Err() == nil {
			// member存在
			return true, nil
		} else if cmds[1].Err() == redis.Nil {
			// member不存在，虽然有err但属于正常情况
			return false, nil
		} else if strings.HasPrefix(cmds[1].Err().Error(), "WRONGTYPE") {
			// 数据库保护产生的空键
			return false, nil
		} else {
			// err!=redis.Nil����说明是其他异常，要返回异常
			return false, cmds[1].Err()
		}
	} else {
		return false, richTypes.ErrorKeyNotExist
	}
}

func (this *ZSetQuery) ZScoreIfExist(key string, member string) (float64, error) {
	log.Notice("[Redis ZIsMemberIfExist]", key, member)
	// 通过ZRank间接实现存在性判断
	// ZScore返回member在ZSet中的Index
	cmds, _ := this.ExecPipeline(func(pipe redis.Pipeliner) error {
		pipe.Exists(key)
		pipe.ZScore(key, member)
		return nil
	})

	// Pipeline默认返回的是最后一个err，所以这里的判定方式要做调整
	if cmds[0].Err() != nil {
		return 0, cmds[0].Err()
	}
	if cmds[0].(*redis.IntCmd).Val() == 1 {
		// 如果member不存在，则会返回error=redis.Nil
		if cmds[1].Err() == nil {
			// member存在
			return cmds[1].(*redis.FloatCmd).Val(), nil
		} else {
			return 0, richTypes.ErrorMemberNotExist
		}
	} else {
		return 0, richTypes.ErrorKeyNotExist
	}
}
