package rich

import (
	"io"

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

func SetLogger(l io.Writer) {
	log.SetLogger(l)
}
