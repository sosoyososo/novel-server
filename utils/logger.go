package utils

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type loggerLevel int

const (
	loggerLevelDebug loggerLevel = iota
	loggerLevelTest
	loggerLevelInfo
	loggerLevelError
	loggerLevelPanic
)

var (
	DebugLogger = Logger{loggerLevelDebug, nil}
	TestLogger  = Logger{loggerLevelTest, nil}
	InfoLogger  = Logger{loggerLevelInfo, nil}
	ErrorLogger = Logger{loggerLevelError, nil}
	PanicLogger = Logger{loggerLevelPanic, nil}
)

var (
	DebugLoggerEnabled = true
	InfoLoggerEnabled  = true
	ErrorLoggerEnabled = true
	PanicLoggerEnabled = true
	TestLoggerEnabled  = true
)

func InitLogger() {
	fileName := "gin.log"
	logFilePath := GetPathRelativeToProjRoot(fileName)
	st, err := os.Stat(logFilePath)
	var f *os.File
	if os.IsNotExist(err) {
		f, err = os.Create(logFilePath)
	} else if err == nil && st.IsDir() == false && st.Size() > 1024*1024*10 {
		err = os.Rename(
			logFilePath,
			GetPathRelativeToProjRoot(FormatTime(time.Now())+"-"+fileName),
		)
		if nil != err {
			InfoLogger.Logf("日志文件转移错误 %v\n", err)
			panic(err)
		}
		f, err = os.Create(logFilePath)
	} else if nil != err {
		InfoLogger.Logf("日志文件信息读取错误 %v\n", err)
		panic(err)
	}

	if nil != err {
		InfoLogger.Logf("日志文件创建错误 %v\n", err)
		panic(err)
	}

	f, err = os.OpenFile(logFilePath, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	DebugLogger.w = f
	InfoLogger.w = f
	ErrorLogger.w = f
	PanicLogger.w = f
	TestLogger.w = f

	gin.DefaultWriter = f
	gin.DefaultErrorWriter = f
}

type Logger struct {
	level loggerLevel
	w     io.Writer
}

func (l *Logger) SetWriter(w io.Writer) {
	l.w = w
}

func (l *Logger) Enabled() bool {
	return []bool{
		DebugLoggerEnabled,
		InfoLoggerEnabled,
		ErrorLoggerEnabled,
		PanicLoggerEnabled,
		TestLoggerEnabled,
	}[l.level]
}

func (l *Logger) prefix() string {
	return []string{"[ -- Debug -- ]", "[ -- Test --]", "[ -- Info -- ]", "[ -- Error -- ]", "[ -- Fetal -- ]"}[l.level]
}

func (l *Logger) Logf(format string, a ...interface{}) {
	if l.Enabled() {
		fmt.Fprintf(l.w, "%v %v ", l.prefix(), FormatTime(time.Now()))
		fmt.Fprintf(l.w, format+"\n", a...)
	}
}

func (l *Logger) Logln(a ...interface{}) {
	if l.Enabled() {
		fmt.Fprintf(l.w, "%v %v ", l.prefix(), FormatTime(time.Now()))
		fmt.Fprintln(l.w, a...)
	}
}

func (l Logger) Write(p []byte) (n int, err error) {
	return l.w.Write(p)
}
