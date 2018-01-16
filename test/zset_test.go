package test

import (
	"sort"
	"testing"
	"time"

	"github.com/go-redis/redis"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/szyhf/go-rich"
)

func TestZSet(t *testing.T) {
	rich.SetLogger(func(level int, format string, v ...interface{}) {})
	Convey("init", t, func() {
		zqs := getRicher().QueryZSet("HelloZSet")
		zData := []redis.Z{
			redis.Z{Score: 1, Member: "A"},
			redis.Z{Score: 3, Member: "b"},
			redis.Z{Score: 2, Member: "C"},
			redis.Z{Score: 5, Member: "D"}}

		stringSlcCmd := richer.Querier().Keys("HelloZSet*")
		if stringSlcCmd.Err() == nil {
			richer.Querier().Del(stringSlcCmd.Val()...)
		}

		zDESCData := make([]redis.Z, len(zData))
		copy(zDESCData, zData)
		sort.SliceStable(zDESCData, func(i int, j int) bool {
			return zDESCData[i].Score > zDESCData[j].Score
		})
		zASCData := make([]redis.Z, len(zData))
		copy(zASCData, zData)
		sort.SliceStable(zASCData, func(i int, j int) bool {
			return zASCData[i].Score < zASCData[j].Score
		})
		_ = zData
		_ = zDESCData
		zqs = zqs.SetRebuildFunc(func() ([]redis.Z, time.Duration) {
			return zData, time.Hour
		})

		Convey("Members", func() {
			members, err := zqs.Members()
			So(err, ShouldBeNil)
			res := map[string]bool{}
			for _, m := range members {
				for _, mj := range zData {
					if mj.Member.(string) == m {
						res[m] = true
					}
				}
				So(res[m], ShouldBeTrue)
			}
		})

		Convey("IsMember", func() {
			ok, err := zqs.IsMember("b")
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)
		})

		Convey("Count", func() {
			count, err := zqs.Count()
			So(err, ShouldBeNil)
			So(count, ShouldEqual, len(zData))
		})

		Convey("Score", func() {
			for _, z := range zData {
				score, err := zqs.Score(z.Member.(string))
				So(err, ShouldBeNil)
				So(score, ShouldAlmostEqual, z.Score)
			}
		})

		Convey("RangeASC", func() {
			members, err := zqs.RangeASC(0, -1)
			So(err, ShouldBeNil)
			for i, zsd := range zASCData {
				So(members[i], ShouldEqual, zsd.Member.(string))
			}
		})

		Convey("RangeDESC", func() {
			members, err := zqs.RangeDESC(0, -1)
			So(err, ShouldBeNil)
			for i, zsd := range zDESCData {
				So(members[i], ShouldEqual, zsd.Member.(string))
			}
		})

		Convey("RangeASCWithScores", func() {
			zs, err := zqs.RangeASCWithScores(0, -1)
			So(err, ShouldBeNil)
			for i, zdd := range zASCData {
				So(zs[i].Score, ShouldAlmostEqual, zdd.Score)
				So(zs[i].Member, ShouldEqual, zdd.Member.(string))
			}
		})

		Convey("RangeDESCWithScores", func() {
			zs, err := zqs.RangeDESCWithScores(0, -1)
			So(err, ShouldBeNil)
			for i, zdd := range zDESCData {
				So(zs[i].Score, ShouldAlmostEqual, zdd.Score)
				So(zs[i].Member, ShouldEqual, zdd.Member.(string))
			}
		})

		Convey("RangeByScoreASC", func() {
			members, err := zqs.RangeByScoreASC("-inf", "+inf", 0, 100)
			So(err, ShouldBeNil)
			for i, zsd := range zASCData {
				So(members[i], ShouldEqual, zsd.Member.(string))
			}
		})

		Convey("RangeByScoreDESC", func() {
			members, err := zqs.RangeByScoreDESC("-inf", "+inf", 0, 100)
			So(err, ShouldBeNil)
			for i, zsd := range zDESCData {
				So(members[i], ShouldEqual, zsd.Member.(string))
			}
		})
	})
}
