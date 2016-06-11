package log

import (
	"bufio"
	"log"
	"os"
	"runtime"
	"sync"
	"syscall"
	"time"
)

type LogLevel uint32

const (
	LEVEL_VERBOSE = 1 << 0
	LEVEL_DEBUG   = 1 << 1
	LEVEL_INFO    = 1 << 2
	LEVEL_WARN    = 1 << 3
	LEVEL_ERROR   = 1 << 4
)

var (
	out             *bufio.Writer
	logger          *log.Logger
	logLevel        LogLevel = 0xffff //enable all log by default
	mutex           *sync.Mutex
	logFile         *os.File
	_LOG_PATH_      string
	_LOGFILE_PREFIX string
)

func SetLogFilePrefix(prefix string) {
	_LOGFILE_PREFIX = prefix
}

func SetLogLevel(level LogLevel) {
	logLevel = 0
	switch level {
	case LEVEL_VERBOSE:
		logLevel |= LEVEL_VERBOSE
		fallthrough
	case LEVEL_DEBUG:
		logLevel |= LEVEL_DEBUG
		fallthrough
	case LEVEL_INFO:
		logLevel |= LEVEL_INFO
		fallthrough
	case LEVEL_WARN:
		logLevel |= LEVEL_WARN
		fallthrough
	case LEVEL_ERROR:
		logLevel |= LEVEL_ERROR
	}
}

func Setup(logDir string, isRedirect bool) {
	mutex = new(sync.Mutex)
	_LOG_PATH_ = logDir

	ensureLogDir()

	// 每天切换文件
	go func() {
		defer func() { recover() }()
		for {
			now := time.Now()
			switchLogFile(now, isRedirect)
			// 计算此刻到第二天零点的时间
			t := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, now.Nanosecond(), now.Location())
			duration := t.Sub(now)
			time.Sleep(duration)
		}
	}()

	runtime.Gosched()

	// 两秒刷新一次
	go func() {
		defer func() {
			if r := recover(); r != nil {
			}
		}()
		c := time.Tick(2 * time.Second)
		for _ = range c {
			Flush()
		}
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
			}
		}()
		c := time.Tick(10 * time.Second)
		memStat := new(runtime.MemStats)
		var (
			lastNumGC        uint32
			lastPauseTotalNs uint64
		)

		for _ = range c {

			runtime.ReadMemStats(memStat)
			justGC := memStat.NumGC - lastNumGC
			justPauseTotalNs := (memStat.PauseTotalNs - lastPauseTotalNs) / uint64(time.Millisecond)

			lastNumGC = memStat.NumGC
			lastPauseTotalNs = memStat.PauseTotalNs

			Infof("mem:%v goroutine:%v gc:%v pause:%vms", memStat.Alloc, runtime.NumGoroutine(), justGC, justPauseTotalNs)
		}
	}()
}

func ensureLogDir() {
	dir, _ := os.Stat(_LOG_PATH_)
	if dir == nil {
		os.Mkdir(_LOG_PATH_, 0777)
	}
}

func switchLogFile(now time.Time, isRedirect bool) {
	// file.Fd()==18446744073709551615为true  文件已经close 或者 没有打开
	// 目前正在查找是否有更加直观的写法
	mutex.Lock()
	defer mutex.Unlock()

	if logFile != nil {
		logFile.Close()
	}

	var err error
	logName := _LOG_PATH_ + "/" + _LOGFILE_PREFIX + now.Format("2006-01-02") + ".log"
	logFile, err = os.OpenFile(logName, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0777)
	if err != nil {
		panic(err)
	}
	if isRedirect {
		syscall.Dup2(int(logFile.Fd()), 1)
		syscall.Dup2(int(logFile.Fd()), 2)
	}

	out = bufio.NewWriterSize(logFile, 1024000)
	logger = log.New(out, "", log.Ldate|log.Ltime)
}

func Close() {
	// 永久锁定防止文件切换进程来操作
	mutex.Lock()
	out.Flush()
	logFile.Close()
}

func Flush() error {
	mutex.Lock()
	defer mutex.Unlock()
	return out.Flush()
}

func Info(format string) {
	if logLevel&LEVEL_INFO == 0 {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	logger.Printf("[I] - " + format)
}

func Warn(format string) {
	if logLevel&LEVEL_WARN == 0 {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	logger.Printf("[W] - " + format)
}

func Error(format string) {
	if logLevel&LEVEL_ERROR == 0 {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	logger.Printf("[E] - " + format)
}

func Debug(format string) {
	if logLevel&LEVEL_DEBUG == 0 {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	logger.Print("[D] - " + format)
}

func Verbose(format string) {
	if logLevel&LEVEL_VERBOSE == 0 {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	logger.Print("[V] - " + format)
}

func Infof(format string, v ...interface{}) {
	if logLevel&LEVEL_INFO == 0 {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	logger.Printf("[I] - "+format, v...)
}

func Warnf(format string, v ...interface{}) {
	if logLevel&LEVEL_WARN == 0 {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	logger.Printf("[W] - "+format, v...)
}

func Errorf(format string, v ...interface{}) {
	if logLevel&LEVEL_ERROR == 0 {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	logger.Printf("[E] - "+format, v...)
}

func Debugf(format string, v ...interface{}) {
	if logLevel&LEVEL_DEBUG == 0 {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	logger.Printf("[D] - "+format, v...)
}

func Verbosef(format string, v ...interface{}) {
	if logLevel&LEVEL_VERBOSE == 0 {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	logger.Printf("[V] - "+format, v...)
}
