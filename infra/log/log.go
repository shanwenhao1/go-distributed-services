// Package log provides syslog logging to a local or remote
// Syslog logger.  To specify a remote syslog host, set the
// "log.sysloghost" key in the Skynet configuration.  Specify
// the port with "log.syslogport".  If "log.sysloghost" is not provided,
// skynet will log to local syslog.
package log

import (
	"fmt"
	"github.com/alecthomas/log4go"
	"strings"
)

type LogLevel int8

var syslogHost string
var syslogPort int = 0

var minLevel LogLevel

const (
	TRACE = iota
	DEBUG
	INFO
	WARN
	ERROR
	CRITICAL
	PANIC
)

const (
	InitSer  = "Init Server"
	ReqParse = "Req Parsing"
)

// 将日志收集至指定的日志收集服务中,但由于syslog目前只支持unix.所以暂时弃用.
//var logger *syslog.Writer
//func Initialize() {
//	var e error
//	if len(syslogHost) > 0 {
//
//		logger, e = syslog.New(syslog.LOG_INFO|syslog.LOG_USER, "skynet")
//		if e != nil {
//			panic(e)
//		}
//	} else {
//		logger, e = syslog.Dial("tcp4", syslogHost+":"+strconv.Itoa(syslogPort), syslog.LOG_INFO|syslog.LOG_USER, "skynet")
//		if e != nil {
//			panic(e)
//		}
//	}
//}

var logger log4go.Logger

//记录基本日志
func Info(arg0 interface{}, args ...interface{}) {
	logger.Info(arg0, args...)
}

//记录调试日志
func Debug(arg0 interface{}, args ...interface{}) {
	logger.Debug(arg0, args...)
}

//记录警告日志
func Warn(arg0 interface{}, args ...interface{}) {
	logger.Warn(arg0, args...)
}

//记录错误日志
func Error(arg0 interface{}, args ...interface{}) {
	logger.Error(arg0, args...)
}

//记录崩溃日志
func Critical(arg0 interface{}, args ...interface{}) {
	logger.Critical(arg0, args)
}

//记录日志, 会具体记录哪种action操作
func LogWithTag(logType int, actionType string, arg2 interface{}, args ...interface{}) {
	var msg string
	switch first := arg2.(type) {
	case string:
		// Use the string as a format string
		msg = fmt.Sprintf(first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		msg = first()
	default:
		// Build a format string so that it will be similar to Sprint
		msg = fmt.Sprintf(fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
	}

	logMsg := "[" + actionType + "]: " + msg
	switch logType {
	case TRACE:
		logger.Trace(logMsg)
	case INFO:
		logger.Info(logMsg)
	case DEBUG:
		logger.Debug(logMsg)
	case WARN:
		logger.Warn(logMsg)
	case ERROR:
		logger.Error(logMsg)
	case CRITICAL:
		logger.Critical(logMsg)
	default:
		logger.Info(arg2, args)
	}
}

//日志框架初始化
// TODO 考虑将日志改为分布式收集
func init() {
	logger = make(log4go.Logger)
	logger.LoadConfiguration("config/log4go.xml")
	logger.Info("[" + InitSer + "]: " + "日志框架初始化完成")
}
