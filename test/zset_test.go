package test

import (
	"testing"
	"time"

	"github.com/go-redis/redis"
)

func TestZSetAddMember(t *testing.T) {
	zqs := getRicher().QueryZSet("HelloZSet")

	zqs = zqs.SetRebuildFunc(func() ([]redis.Z, time.Duration) {
		return []redis.Z{redis.Z{Score: 1, Member: "A"},
			redis.Z{Score: 3, Member: "b"},
			redis.Z{Score: 2, Member: "C"},
			redis.Z{Score: 5, Member: "D"}}, time.Hour
	})
	a, b := zqs.RangeASCWithScores(0, -1)
}
