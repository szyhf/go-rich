package test

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/szyhf/go-rich"
)

type String struct {
	Hello string
	Word  int
}

func (this *String) MarshalBinary() (data []byte, err error) {
	return json.Marshal(this)
}

func (this *String) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, this)
}

func TestString(t *testing.T) {
	rich.SetLogger(func(level int, format string, v ...interface{}) {})
	sData := String{
		Hello: "World",
		Word:  10086,
	}
	sBytes, _ := sData.MarshalBinary()
	sStr := string(sBytes)

	Convey("Init", t, func() {
		richer := getRicher()
		sqs := richer.QueryString("HelloString")
		sqs = sqs.SetRebuildFunc(func() (interface{}, time.Duration) {
			return &sData, time.Hour
		})

		Convey("Get", func() {
			str, err := sqs.Get()
			So(err, ShouldBeNil)
			So(str, ShouldEqual, sStr)
		})

		Convey("Scan", func() {
			s := new(String)
			err := sqs.Scan(s)
			So(err, ShouldBeNil)
			So(reflect.DeepEqual(*s, sData), ShouldBeTrue)
		})
	})
}
