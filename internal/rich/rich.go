package rich

import (
	"fmt"

	"github.com/go-redis/redis"

	"github.com/szyhf/go-rich/internal/query"
	"github.com/szyhf/go-rich/internal/queryset"
	richTypes "github.com/szyhf/go-rich/types"
)

func NewRicher() richTypes.Richer {
	cli, ok := GetRedisClient("default")
	if !ok {
		panic("client `default` not registered.")
	}
	return newRich(cli)
}

type rich struct {
	query *query.Query
}

func newRich(cli *redis.Client) *rich {
	return &rich{
		query: query.NewQuery(cli),
	}
}

func (this *rich) Using(alias string) error {
	cli, ok := GetRedisClient(alias)
	if !ok {
		return fmt.Errorf("alias = %s not reigsted.", alias)
	}
	this.query.Client = cli
	return nil
}

func (this *rich) QueryZSet(i interface{}) richTypes.ZSetQuerySeter {
	switch m := i.(type) {
	case string:
		return queryset.NewZSet(m, this.query)
	case *string:
		panic("Not imp")
	case interface {
		ZSetName() string
	}:
		panic("Not imp")
	default:
		panic("Not imp")
	}
}

func (this *rich) QueryString(i interface{}) richTypes.StringQuerySeter {
	switch m := i.(type) {
	case string:
		return queryset.NewString(m, this.query)
	default:
		panic("Not imp")
	}
}

func (this *rich) Querier() richTypes.Querier {
	return this.query
}
