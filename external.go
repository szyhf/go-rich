package rich

import (
	"github.com/go-redis/redis"
	"github.com/szyhf/go-rich/internal/log"
	"github.com/szyhf/go-rich/internal/rich"
	"github.com/szyhf/go-rich/types"
)

func RegisterRedisClient(alias string, cli *redis.Client) error {
	return rich.RegisterRedisClient(alias, cli)
}

func NewRicher() Richer {
	return rich.NewRicher()
}

type Richer = richTypes.Richer

type LogFunc = log.LogFunc

func SetLogger(l LogFunc) {
	log.Logf = l
}
