package test

import (
	"sort"
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
	res, err := zqs.RangeASCWithScores(0, -1)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	isASCSorted := sort.SliceIsSorted(res, func(i int, j int) bool {
		return res[i].Score > res[j].Score
	})
	if !isASCSorted {
		t.Errorf("not asc sorted: %+v", res)
	}
}
