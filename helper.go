package vlog

import (
	"bufio"
	"fmt"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
	"time"
)

var log *logrus.Logger

// Init 初始化日志配置信息
// dir 日志文件保存在当前执行文件的哪个目录
// level 日志级别,设置为debug在控制台会有显示,设置为info,warn,error和其它只会打印在文件中
// day 文件默认保存多少天
func Init(dir, level string, day int) (*logrus.Logger, error) {
	exist, err := pathExists(dir)
	if err != nil {
		return nil, err
	}
	if !exist {
		// 创建文件夹
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	baseLogPath := path.Join(dir, "vlog")
	writer, err := rotatelogs.New(
		baseLogPath+"%Y-%m-%d %H:%M.log",
		//rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Duration(day)*24*time.Hour), // 文件最大保存时间(天)
		rotatelogs.WithRotationTime(24*time.Hour),              // 日志切割时间间隔
	)
	if err != nil {
		return nil, err
	}
	switch level {
	/*
		如果日志级别不是debug就不要打印日志到控制台了
	*/
	case "debug":
		log.SetLevel(logrus.DebugLevel)
		log.SetOutput(os.Stderr)
	case "info":
		setNull()
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		setNull()
		log.SetLevel(logrus.WarnLevel)
	case "error":
		setNull()
		log.SetLevel(logrus.ErrorLevel)
	default:
		setNull()
		log.SetLevel(logrus.InfoLevel)
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.TextFormatter{DisableColors: true, TimestampFormat: "2006-01-02 15:04:05.000"})
	log.AddHook(lfHook)
	return log, nil
}

func setNull() {
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	writer := bufio.NewWriter(src)
	log.SetOutput(writer)
}

// 打印debug日志
func Debug(format string, a ...interface{}) {
	log.WithFields(getStack()).Debugf(format, a...)
}

// 打印warn日志
//noinspection GoUnusedExportedFunction
func Warn(format string, a ...interface{}) {
	log.WithFields(getStack()).Warnf(format, a...)
}

// 打印error信息
func Error(format string, a ...interface{}) {
	log.WithFields(getStack()).Errorf(format, a...)
}

// 打印info信息
func Info(format string, a ...interface{}) {
	log.WithFields(getStack()).Infof(format, a...)
}

// getStack 获取堆栈信息
func getStack() logrus.Fields {
	pc, filePath, line, ok := runtime.Caller(2)
	if ok {
		return logrus.Fields{
			"file":   filePath,
			"line":   line,
			"method": runtime.FuncForPC(pc).Name(),
		}
	} else {
		return logrus.Fields{
			"file":   "?",
			"line":   "??",
			"method": "???",
		}
	}
}

// 判断文件夹是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
