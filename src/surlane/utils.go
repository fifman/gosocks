package surlane

import (
	"log"
	"os"
	"context"
	"time"
	"net"
	"encoding/binary"
	"strings"
	"strconv"
	"fmt"
	"github.com/pkg/errors"
)

const (
	FLAG = log.Ldate | log.Ltime | log.Lshortfile | log.Lmicroseconds
	LevelPanic = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

var (
	DebugMark = 0
	currentRoutineId = 0
	RootContext = LocalContext{
		context.Background(),
		currentRoutineId,
		"root",
		nil,
		newContextLogger("root", currentRoutineId),
	}
	callDepth = 20
)

type contextLogger struct {
	panic *log.Logger
	error *log.Logger
	warn  *log.Logger
	info  *log.Logger
	debug *log.Logger
	trace *log.Logger
	level int
}

func (logger *contextLogger) logError(err error, msg string) {
	logger.Errorf(msg + " reason: %+v", errors.WithStack(err))
}

func (logger *contextLogger) Level(level int) {
	logger.level = level
}

func (logger *contextLogger) Errorf(format string, v ...interface{}) {
	if logger.level >= LevelError {
		logger.error.Output(callDepth, fmt.Sprintf(format, v))
	}
}

func (logger *contextLogger) Panicf(format string, v ...interface{}) {
	if logger.level >= LevelPanic {
		logger.panic.Output(callDepth, fmt.Sprintf(format, v))
	}
}

func (logger *contextLogger) Infof(format string, v ...interface{}) {
	if logger.level >= LevelInfo {
		logger.info.Output(callDepth, fmt.Sprintf(format, v))
	}
}

func (logger *contextLogger) Warnf(format string, v ...interface{}) {
	if logger.level >= LevelWarn {
		logger.warn.Output(callDepth, fmt.Sprintf(format, v))
	}
}

func (logger *contextLogger) Debugf(format string, v ...interface{}) {
	if logger.level >= LevelDebug {
		logger.debug.Output(callDepth, fmt.Sprintf(format, v))
	}
}

func (logger *contextLogger) Tracef(format string, v ...interface{}) {
	if logger.level >= LevelTrace {
		logger.trace.Output(callDepth, fmt.Sprintf(format, v))
	}
}

func (logger *contextLogger) Error(v ...interface{}) {
	if logger.level >= LevelError {
		logger.error.Output(callDepth, fmt.Sprint(v))
	}
}

func (logger *contextLogger) Panic(v ...interface{}) {
	if logger.level >= LevelPanic {
		logger.panic.Output(callDepth, fmt.Sprint(v))
	}
}

func (logger *contextLogger) Info(v ...interface{}) {
	if logger.level >= LevelInfo {
		logger.info.Output(callDepth, fmt.Sprint(v))
	}
}

func (logger *contextLogger) Warn(v ...interface{}) {
	if logger.level >= LevelWarn {
		logger.warn.Output(callDepth, fmt.Sprint(v))
	}
}

func (logger *contextLogger) Debug(v ...interface{}) {
	if logger.level >= LevelDebug {
		logger.debug.Output(callDepth, fmt.Sprint(v))
	}
}

func (logger *contextLogger) Trace(v ...interface{}) {
	if logger.level >= LevelTrace {
		logger.trace.Output(callDepth, fmt.Sprint(v))
	}
}

func newContextLogger(name string, id int) *contextLogger {
	return &contextLogger{
		log.New(os.Stderr, "PANIC{" + name + "$" + strconv.Itoa(id) + "}:  ", FLAG),
		log.New(os.Stderr, "ERROR{" + name + "$" + strconv.Itoa(id) + "}:  ", FLAG),
		log.New(os.Stdout, "WARN{" + name + "$" + strconv.Itoa(id) + "}:  ", FLAG),
		log.New(os.Stdout, "INFO{" + name + "$" + strconv.Itoa(id) + "}:  ", FLAG),
		log.New(os.Stdout, "DEBUG{" + name + "$" + strconv.Itoa(id) + "}:  ", FLAG),
		log.New(os.Stdout, "TRACE{" + name + "$" + strconv.Itoa(id) + "} >>>>>  ", FLAG),
		LevelError,
	}
}

type AppLogger interface {
	Error(...interface{})
	Panic(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Debug(...interface{})
	Trace(...interface{})
	Errorf(string, ...interface{})
	Panicf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Debugf(string, ...interface{})
	Tracef(string, ...interface{})
	Level(int)
	logError(error, string)
}

type LocalContext struct {
	context.Context
	Id int
	Name string
	Cancel context.CancelFunc
	AppLogger
}

func NewContext(parent *LocalContext, name string) *LocalContext {
	currentRoutineId += 1
	sub, cancel := context.WithCancel(parent.Context)
	return &LocalContext{
		sub,
		currentRoutineId,
		name,
		cancel,
		newContextLogger(name, currentRoutineId),
	}
}

func NewContextWithDeadline(parent *LocalContext, name string, duration time.Duration) *LocalContext {
	currentRoutineId += 1
	sub, cancel := context.WithDeadline(parent.Context, time.Now().Add(duration))
	return &LocalContext{
		sub,
		currentRoutineId,
		name,
		cancel,
		newContextLogger(name, currentRoutineId),
	}
}

func (ctx *LocalContext) GetId() int {
	return ctx.Id
}

func ParseRawAddr(rawAddr []byte) (address string, port uint16) {
	addrType := rawAddr[0]
	addrLen := len(rawAddr)-2
	switch addrType {
	case 0, 1:
		address = net.IP(rawAddr[1:addrLen]).String()
	default:
		address = string(rawAddr[1:addrLen])
	}
	port = binary.BigEndian.Uint16(rawAddr[addrLen:addrLen+2])
	return
}

func GenRawAddr(address string, port uint16) []byte {
	var (
		addrBytes []byte
		addrType  byte
		ip net.IP
	)
	addrType = byte(len(address))
	if ip = net.ParseIP(address); ip != nil {
		addrType = 0
		if strings.Contains(address, ":") {
			addrType = 1
		}
	}
	switch addrType {
	case 0:
		addrBytes = ip.To4()
	case 1:
		addrBytes = ip.To16()
	default:
		addrBytes = []byte(address)
	}
	addrLen := len(addrBytes)
	result := make([]byte, addrLen+3)
	result[0] = addrType
	copy(result[1:addrLen+1], addrBytes)
	binary.BigEndian.PutUint16(result[addrLen+1:addrLen+3], port)
	return result
}