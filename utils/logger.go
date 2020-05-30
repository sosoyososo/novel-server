package utils

import (
	"fmt"
	"io"
	"os"
	"time"
)

type loggerLevel int

const (
	loggerLevelDebug loggerLevel = iota
	loggerLevelTest
	loggerLevelInfo
	loggerLevelError
	loggerLevelPanic
	loggerLevelBusiness
	loggerLevelWebServer
	loggerLevelDB
)

var (
	DebugLogger     = Logger{loggerLevelDebug, nil, "info.log"}
	TestLogger      = Logger{loggerLevelTest, nil, "info.log"}
	InfoLogger      = Logger{loggerLevelInfo, nil, "info.log"}
	ErrorLogger     = Logger{loggerLevelError, nil, "err.log"}
	PanicLogger     = Logger{loggerLevelPanic, nil, "err.log"}
	BusinessLogger  = Logger{loggerLevelBusiness, nil, "info.log"}
	WebServerLogger = Logger{loggerLevelWebServer, nil, "gin.log"}
	DBLogger        = Logger{loggerLevelBusiness, nil, "db.log"}
)

var (
	DebugLoggerEnabled     = true
	InfoLoggerEnabled      = true
	ErrorLoggerEnabled     = true
	PanicLoggerEnabled     = true
	TestLoggerEnabled      = true
	BusinessLoggerEnabled  = true
	WebServerLoggerEnabled = true
	DBLoggerEnabled        = true
)

func InitLogger() {
	loggers := []*Logger{
		&InfoLogger,
		&DebugLogger,
		&TestLogger,
		&ErrorLogger,
		&PanicLogger,
		&BusinessLogger,
		&WebServerLogger,
		&DBLogger,
	}
	fileMapper := map[string]*os.File{}
	for i, l := range loggers {
		f := fileMapper[l.baseName]
		if nil == f {
			f = loggerCreateFile(l.baseName)
		}
		loggers[i].w = f
	}
	// f := loggerCreateFile("gin.log")
	// InfoLogger.w = f
	// DebugLogger.w = f
	// TestLogger.w = f
	// ErrorLogger.w = f
	// PanicLogger.w = f
	// BusinessLogger.w = f
	// WebServerLogger.w = f
	// DBLogger.w = f
	// gin.DefaultWriter = f
	// gin.DefaultErrorWriter = f
}

func loggerCreateFile(fileName string) *os.File {
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
	return f
}

type Logger struct {
	level    loggerLevel
	w        io.Writer
	baseName string
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
		BusinessLoggerEnabled,
		WebServerLoggerEnabled,
		DBLoggerEnabled,
	}[l.level]
}

func (l *Logger) prefix() string {
	return []string{
		"[ -- Debug -- ]",
		"[ -- Test --]",
		"[ -- Info -- ]",
		"[ -- Error -- ]",
		"[ -- Fetal -- ]",
		"[ -- Business -- ]",
		"",
		""}[l.level]
}

func (l *Logger) Logf(format string, a ...interface{}) {
	CallFuncInNewRecoveryRoutine(func() {
		if l.Enabled() {
			fmt.Fprintf(l.w, "%v %v ", l.prefix(), FormatTime(time.Now()))
			fmt.Fprintf(l.w, format+"\n", a...)
		}
	})
}

func (l *Logger) Logln(a ...interface{}) {
	CallFuncInNewRecoveryRoutine(func() {
		if l.Enabled() {
			fmt.Fprintf(l.w, "%v %v ", l.prefix(), FormatTime(time.Now()))
			fmt.Fprintln(l.w, a...)
		}
	})
}

func (l Logger) Write(p []byte) (n int, err error) {
	return l.w.Write(p)
}
