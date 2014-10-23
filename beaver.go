package beaver

import (
	"fmt"
	"io"
	"log"
	"log/syslog"
	"os"
	"strings"
)

type loglevel int
type logflags int

const (
	DEBUG loglevel = iota
	INFO
	NOTICE
	WARNING
	ERROR
	CRITICAL
	FATAL
	PANIC

	Ldate         = log.Ldate
	Ltime         = log.Ltime
	Lmicroseconds = log.Lmicroseconds
	Llongfile     = log.Llongfile
	Lshortfile    = log.Lshortfile
	LstdFlags     = Ldate | Ltime
)

var levels = []string{
	"[DEBUG]",
	"[INFO]",
	"[NOTICE]",
	"[WARNING]",
	"[ERROR]",
	"[CRITICAL]",
}

//type Beaver interface {
//	Debug(string, ...interface {})
//	Debugf(...interface {})
//	Info(string, ...interface {})
//	Infof(...interface {})
//	Notice(string, ...interface {})
//	Noticef(...interface {})
//	Error(string, ...interface {})
//	Errorf(...interface {})
//	Critical(string, ...interface {})
//	Criticalf(...interface {})
//}


type StdoutLogger struct {
	level  loglevel
	writer *log.Logger
}

func NewStdoutLogger(out io.Writer, prefix string, detail int) (*StdoutLogger, error) {
	if out == nil {
		out = os.Stdout
	}
	prefix = strings.TrimSpace(prefix) + " "

	return &StdoutLogger{
		level:  DEBUG,
		writer: log.New(out, prefix, detail),
	}, nil
}

func (logger *StdoutLogger) SetLevel(level loglevel) {
	logger.level = level
}

func (logger *StdoutLogger) GetLevel() loglevel {
	return logger.level
}

func (logger *StdoutLogger) write(level loglevel, format string, text ...interface{}) {

	if level < logger.level {
		return
	}

	var message string
	if format == "" {
		for _, i := range text {
			message += " " + fmt.Sprintf("%v", i)
		}
		message = levels[int(level)] + message
	} else {
		message = levels[int(level)] + " " + fmt.Sprintf(format, text...)
	}

	switch level {
	case CRITICAL:
		logger.writer.Panic(message)
	case ERROR:
		logger.writer.Print(message)
	case WARNING:
		logger.writer.Print(message)
	case NOTICE:
		logger.writer.Print(message)
	case INFO:
		logger.writer.Print(message)
	case DEBUG:
		logger.writer.Print(message)
	default:
	}
	panic("unhandled log level")
}

func (logger *StdoutLogger) Debug(v ...interface{}) {
	logger.write(DEBUG, "", v...)
}

func (logger *StdoutLogger) Debugf(format string, v ...interface{}) {
	logger.write(DEBUG, format, v...)
}

func (logger *StdoutLogger) Info(v ...interface{}) {
	logger.write(INFO, "", v...)
}

func (logger *StdoutLogger) Infof(format string, v ...interface{}) {
	logger.write(INFO, format, v...)
}

func (logger *StdoutLogger) Noticef(format string, v ...interface{}) {
	logger.write(NOTICE, format, v...)
}

func (logger *StdoutLogger) Notice(v ...interface{}) {
	logger.write(NOTICE, "", v...)
}

func (logger *StdoutLogger) Warn(v ...interface{}) {
	logger.write(WARNING, "", v...)
}

func (logger *StdoutLogger) Warnf(format string, v ...interface{}) {
	logger.write(WARNING, format, v...)
}

func (logger *StdoutLogger) Error(v ...interface{}) {
	logger.write(ERROR, "", v...)
}

func (logger *StdoutLogger) Errorf(format string, v ...interface{}) {
	logger.write(ERROR, format, v...)
}

func (logger *StdoutLogger) Criticalf(format string, v ...interface{}) {
	logger.write(PANIC, format, v...)
}

func (logger *StdoutLogger) Critical(v ...interface{}) {
	logger.write(PANIC, "", v...)
}

type SyslogLoggger struct {
	level  loglevel
	writer *syslog.Writer
}

func NewSyslogLogger(prefix string) (*SyslogLoggger, error) {

	writer, err := syslog.New(0, prefix)
	if err != nil {
		return nil, err
	}

	return &SyslogLoggger{
		level:  DEBUG,
		writer: writer,
	}, nil
}

func (logger *SyslogLogger) SetLevel(level loglevel) {
	logger.level = level
}

func (logger *SyslogLogger) GetLevel() loglevel {
	return logger.level
}

func (logger *SyslogLogger) write(level loglevel, format string, text ...interface{}) error {

	if level < logger.level {
		return nil
	}

	var message string
	if format == "" {
		for _, i := range text {
			message += " " + fmt.Sprintf("%v", i)
		}
		message = levels[int(level)] + message
	} else {
		message = levels[int(level)] + " " + fmt.Sprintf(format, text...)
	}

	switch level {
	case CRITICAL:
		return logger.writer.Crit(message)
	case ERROR:
		return logger.writer.Err(message)
	case WARNING:
		return logger.writer.Warning(message)
	case NOTICE:
		return logger.writer.Notice(message)
	case INFO:
		return logger.writer.Info(message)
	case DEBUG:
		return logger.writer.Debug(message)
	default:
	}
	panic("unhandled log level")

}

func (logger *SyslogLogger) Debug(v ...interface{}) {
	_ = logger.write(DEBUG, "", v...)
}

func (logger *SyslogLogger) Debugf(format string, v ...interface{}) {
	_ = logger.write(DEBUG, format, v...)
}

func (logger *SyslogLogger) Info(v ...interface{}) {
	_ = logger.write(INFO, "", v...)
}

func (logger *SyslogLogger) Infof(format string, v ...interface{}) {
	_ = logger.write(INFO, format, v...)
}

func (logger *SyslogLogger) Notice(v ...interface{}) {
	_ = logger.write(NOTICE, "", v...)
}

func (logger *SyslogLogger) Noticef(format string, v ...interface{}) {
	_ = logger.write(NOTICE, format, v...)
}

func (logger *SyslogLogger) Warn(v ...interface{}) {
	_ = logger.write(WARNING, "", v...)
}

func (logger *SyslogLogger) Warnf(format string, v ...interface{}) {
	_ = logger.write(WARNING, format, v...)
}

func (logger *SyslogLogger) Error(v ...interface{}) {
	_ = logger.write(ERROR, "", v...)
}

func (logger *SyslogLogger) Errorf(format string, v ...interface{}) {
	_ = logger.write(ERROR, format, v...)
}

func (logger *SyslogLogger) Critical(v ...interface{}) {
	_ = logger.write(CRITICAL, "", v...)
}

func (logger *SyslogLogger) Criticalf(format string, v ...interface{}) {
	_ = logger.write(CRITICAL, format, v...)
}
