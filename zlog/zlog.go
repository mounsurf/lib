package zlog

import (
	"fmt"
	"io"
	"log"
	"os"
)

//对外开放，string类型，支持哪些级别
const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelPanic = "panic"
	LevelFatal = "fatal"
)

//对内开放，int类型，用以快速索引
const (
	levelDebug = iota
	levelInfo
	levelWarn
	levelError
	levelPanic
	levelFatal
)

var levelName = []string{
	"[DEBUG]",
	"[INFO] ",
	"[WARN] ",
	"[ERROR]",
	"[PANIC]",
	"[FATAL]",
}

const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if lDate or lTime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime // initial values for the standard log
)

var logIns *logger

func init() {
	SetConfig("info", "")
}

func newLogger(w io.Writer) *logger {
	return &logger{
		l: log.New(w, "", 0),
	}
}

type logger struct {
	level int
	l     *log.Logger
}

func (self *logger) setFlags(flag int) *logger {
	self.l.SetFlags(flag)
	return self
}

func (self *logger) setLevel(level int) *logger {
	self.level = level
	return self
}

func (self *logger) doLog(level int, v ...interface{}) bool {
	if level < self.level {
		return false
	}
	self.l.Output(3, levelName[level]+" "+fmt.Sprintln(v...))
	return true
}

func (self *logger) doLogf(level int, format string, v ...interface{}) bool {
	if level < self.level {
		return false
	}
	self.l.Output(3, levelName[level]+" "+fmt.Sprintln(fmt.Sprintf(format, v...)))
	return true
}

func SetConfig(level, file string) {
	var (
		levelMap = map[string]int{
			LevelDebug: levelDebug,
			LevelInfo:  levelInfo,
			LevelWarn:  levelWarn,
			LevelError: levelError,
			LevelPanic: levelPanic,
			LevelFatal: levelFatal,
		}
	)

	if file != "" {
		f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			logIns = newLogger(os.Stdout).setFlags(Lshortfile | LstdFlags)
			Error(err)
		} else {
			logIns = newLogger(f).setFlags(Lshortfile | LstdFlags)
		}
	} else {
		logIns = newLogger(os.Stdout).setFlags(Lshortfile | LstdFlags)
	}

	if l, ok := levelMap[level]; ok {
		logIns.setLevel(l)
	} else {
		logIns.setLevel(levelError)
		Errorf(`unsupported log level: %q, default level "error" seted`, level)
	}
}

func Debug(v ...interface{}) {
	logIns.doLog(levelDebug, v...)
}

func Info(v ...interface{}) {
	logIns.doLog(levelInfo, v...)
}

func Warn(v ...interface{}) {
	logIns.doLog(levelWarn, v...)
}

func Error(v ...interface{}) {
	logIns.doLog(levelError, v...)
}

func Panic(v ...interface{}) {
	if logIns.doLog(levelPanic, v...) {
		panic(fmt.Sprintln(v...))
	}
}

func Fatal(v ...interface{}) {
	if logIns.doLog(levelFatal, v...) {
		os.Exit(1)
	}
}

func Debugf(format string, v ...interface{}) {
	logIns.doLogf(levelDebug, format, v...)
}

func Infof(format string, v ...interface{}) {
	logIns.doLogf(levelInfo, format, v...)
}

func Warnf(format string, v ...interface{}) {
	logIns.doLogf(levelWarn, format, v...)
}

func Errorf(format string, v ...interface{}) {
	logIns.doLogf(levelError, format, v...)
}

func Panicf(format string, v ...interface{}) {
	if logIns.doLogf(levelPanic, format, v...) {
		panic(fmt.Sprintf(format, v...))
	}
}

func Fatalf(format string, v ...interface{}) {
	if logIns.doLogf(levelFatal, format, v...) {
		os.Exit(1)
	}
}
