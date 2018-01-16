package query

import (
	"github.com/go-redis/redis"
	"github.com/szyhf/go-rich/internal/log"
	"github.com/szyhf/go-rich/types"
)

type StringQuery struct {
	*Query
}

func NewStringQuery(q *Query) *StringQuery {
	return &StringQuery{
		Query: q,
	}
}

func (r *StringQuery) IncrByIfExist(key string, incr int64) (int64, error) {
	log.Notice("[Redis IncrByIfExist]", key)
	cmds, err := r.ExecPipeline(func(pipe redis.Pipeliner) error {
		pipe.Exists(key)
		pipe.IncrBy(key, incr)
		return nil
	})

	if err != nil {
		return 0, err
	}
	if cmds[0].(*redis.BoolCmd).Val() {
		return cmds[1].(*redis.IntCmd).Val(), nil
	} else {
		return 0, richTypes.ErrorKeyNotExist
	}
}
