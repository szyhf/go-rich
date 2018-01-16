package test

import (
	"fmt"
	"testing"

	"github.com/go-redis/redis"

	rich "github.com/szyhf/go-rich"
	"github.com/szyhf/go-rich/types"
)

var cli = redis.NewClient(&redis.Options{
	Addr: "127.0.0.1:6379",
})

var richer richTypes.Richer

func getRicher() richTypes.Richer {
	if richer == nil {
		alias := "default"
		err := rich.RegisterRedisClient(alias, cli)
		if err != nil {
			panic(fmt.Errorf("can not register client alias `%s`: %s", alias, err))
		}
		// cli, ok := rich.GetRedisClient(alias)
		// if !ok {
		// 	panic(fmt.Errorf("client alias `%s` not exist", alias))
		// }
		// if cmd := cli.Ping(); cmd.Err() != nil {
		// 	panic(cmd.Err())
		// }
		richer = rich.NewRicher()
	}
	return richer
}

func TestRegister(t *testing.T) {
	alias := "default"
	err := rich.RegisterRedisClient(alias, cli)
	if err != nil {
		t.Errorf("can not register client alias `%s`: %s", alias, err)
		t.FailNow()
	}
	// cli, ok := rich.GetRedisClient(alias)
	// if !ok {
	// 	t.Errorf("client alias `%s` not exist", alias)
	// 	t.FailNow()
	// }
	// if cmd := cli.Ping(); cmd.Err() != nil {
	// 	t.Error(cmd.Err())
	// 	t.FailNow()
	// }
}
