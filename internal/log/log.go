package log

import (
	"fmt"
	"io"
	"os"
	"strings"
)

var logger io.Writer = os.Stdout

func SetLogger(l io.Writer) {
	logger = l
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

func Logf(level LevelID, format string, v ...interface{}) {
	if logger != nil {
		msg := fmt.Sprintf(format, v...)
		logger.Write([]byte(fmt.Sprintf("[%s] %s\n", level, msg)))
	}
}

// 调试使用的消息
func Debug(v ...interface{}) {
	Logf(Level.Debug, strings.Repeat("%v ", len(v)), v...)
}

func Debugf(fmtStr string, v ...interface{}) {
	Logf(Level.Debug, fmtStr, v...)
}

// 业务需要记录的消息
func Info(v ...interface{}) {
	Logf(Level.Info, strings.Repeat("%v ", len(v)), v...)
}

func Infof(fmtStr string, v ...interface{}) {
	// Logf(LevelInformational, fmtStr, v...)
	Logf(Level.Info, fmtStr, v...)
}

// 运维需要知道的消息
func Notice(v ...interface{}) {
	Logf(Level.Notice, strings.Repeat("%v ", len(v)), v...)
}

func Noticef(fmtStr string, v ...interface{}) {
	Logf(Level.Notice, fmtStr, v...)
}

// 运维需要关注的消息
func Warn(v ...interface{}) {
	Logf(Level.Warn, strings.Repeat("%v ", len(v)), v...)
}

func Warnf(fmtStr string, v ...interface{}) {
	Logf(Level.Warn, fmtStr, v...)
}

// 开发需要尽快处理的消息
func Error(v ...interface{}) {
	Logf(Level.Error, strings.Repeat("%v ", len(v)), v...)
}

func Errorf(fmtStr string, v ...interface{}) {
	Logf(Level.Error, fmtStr, v...)
}

// 开发需要马上处理的消息
func Emergency(v ...interface{}) {
	Logf(Level.Emergency, strings.Repeat("%v ", len(v)), v...)
}

func Emergencyf(fmtStr string, v ...interface{}) {
	Logf(Level.Emergency, fmtStr, v...)
}

func Alert(v ...interface{}) {
	Logf(Level.Alert, strings.Repeat("%v ", len(v)), v...)
}

func Alertf(fmtStr string, v ...interface{}) {
	Logf(Level.Alert, fmtStr, v...)
}
