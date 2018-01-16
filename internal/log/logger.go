package log

import (
	"fmt"
	"time"
)

type LogFunc func(level LevelID, format string, v ...interface{})

var Logf LogFunc = func(level LevelID, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fmt.Printf("[%s][%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), level, msg)
}

type LevelID int

var Level = struct {
	Emergency LevelID
	Alert     LevelID
	Critical  LevelID
	Error     LevelID
	Warn      LevelID
	Notice    LevelID
	Info      LevelID
	Debug     LevelID
}{
	Emergency: 0,
	Alert:     1,
	Critical:  2,
	Error:     3,
	Warn:      4,
	Notice:    5,
	Info:      6,
	Debug:     7,
}

var strMap = map[LevelID]string{
	Level.Emergency: "M",
	Level.Alert:     "A",
	Level.Critical:  "C",
	Level.Error:     "E",
	Level.Warn:      "W",
	Level.Notice:    "N",
	Level.Info:      "I",
	Level.Debug:     "D",
}

func (id LevelID) String() string {
	return strMap[id]
}

func (id LevelID) Int() int {
	return int(id)
}
