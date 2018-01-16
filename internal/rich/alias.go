package rich

import (
	"fmt"
	"sync"

	"github.com/go-redis/redis"
)

var redisRegistry = sync.Map{}

func RegisterRedisClient(alias string, cli *redis.Client) error {
	if cmd := cli.Ping(); cmd.Err() != nil {
		return fmt.Errorf("ping redis-client failed: %s", cmd.Err().Error())
	}
	redisRegistry.Store(alias, cli)
	return nil
}

func GetRedisClient(alias string) (*redis.Client, bool) {
	if cli, ok := redisRegistry.Load(alias); ok {
		if c, ok := cli.(*redis.Client); ok {
			return c, true
		} else {
			return nil, false
		}
	}
	return nil, false
}
