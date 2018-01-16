package query

import (
	"github.com/go-redis/redis"
)

// 操作RedisClient的代理
// 封装一些组合操作。

type Query struct {
	*redis.Client
}

func NewQuery(cli *redis.Client) *Query {
	return &Query{
		Client: cli,
	}
}

func (this *Query) ExecPipeline(fun func(pipe redis.Pipeliner) error) ([]redis.Cmder, error) {
	pipe := this.Client.TxPipeline()
	defer pipe.Close()
	if err := fun(pipe); err == nil {
		cmderAry, err := pipe.Exec()
		return cmderAry, err
	} else {
		pipe.Discard()
		return nil, err
	}
}
