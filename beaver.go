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
)

const (
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
	"[WARN]",
	"[ERROR]",
	"[CRITICAL]",
	"[FATAL]",
	"[PANIC]",
}

type Logger struct {
	level  loglevel
	writer *log.Logger
}

func NewLogger(out io.Writer, prefix string, detail int) (*Logger, error) {
	if out == nil {
		out = os.Stdout
	}
	prefix = strings.TrimSpace(prefix) + " "

	return &Logger{
		level:  DEBUG,
		writer: log.New(out, prefix, detail),
	}, nil
}

func (logger *Logger) SetLevel(level loglevel) {
	logger.level = level
}

func (logger *Logger) GetLevel() loglevel {
	return logger.level
}

func (logger *Logger) write(level loglevel, format string, text ...interface{}) {

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

	if level < FATAL {
		logger.writer.Print(message)
	} else if level < PANIC {
		logger.writer.Fatal(message)
	} else {
		logger.writer.Panic(message)
	}
}

func (logger *Logger) Debug(v ...interface{}) {
	logger.write(DEBUG, "", v...)
}

func (logger *Logger) Debugf(format string, v ...interface{}) {
	logger.write(DEBUG, format, v...)
}

func (logger *Logger) Info(v ...interface{}) {
	logger.write(INFO, "", v...)
}

func (logger *Logger) Infof(format string, v ...interface{}) {
	logger.write(INFO, format, v...)
}

func (logger *Logger) Warn(v ...interface{}) {
	logger.write(WARNING, "", v...)
}

func (logger *Logger) Warnf(format string, v ...interface{}) {
	logger.write(WARNING, format, v...)
}

func (logger *Logger) Error(v ...interface{}) {
	logger.write(ERROR, "", v...)
}

func (logger *Logger) Errorf(format string, v ...interface{}) {
	logger.write(ERROR, format, v...)
}

func (logger *Logger) Fatalf(format string, v ...interface{}) {
	logger.write(FATAL, format, v...)
}

func (logger *Logger) Fatal(v ...interface{}) {
	logger.write(FATAL, "", v...)
}

func (logger *Logger) Panicf(format string, v ...interface{}) {
	logger.write(PANIC, format, v...)
}

func (logger *Logger) Panic(v ...interface{}) {
	logger.write(PANIC, "", v...)
}

type Sysger struct {
	level  loglevel
	writer *syslog.Writer
}

func NewSyslog(prefix string) (*Sysger, error) {

	writer, err := syslog.New(0, prefix)
	if err != nil {
		return nil, err
	}

	return &Sysger{
		level:  DEBUG,
		writer: writer,
	}, nil
}

func (logger *Sysger) write(level loglevel, format string, text ...interface{}) error {

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

func (logger *Sysger) Debug(v ...interface{}) {
	_ = logger.write(DEBUG, "", v...)
}

func (logger *Sysger) Debugf(format string, v ...interface{}) {
	_ = logger.write(DEBUG, format, v...)
}

func (logger *Sysger) Info(v ...interface{}) {
	_ = logger.write(INFO, "", v...)
}

func (logger *Sysger) Infof(format string, v ...interface{}) {
	_ = logger.write(INFO, format, v...)
}

func (logger *Sysger) Notice(v ...interface{}) {
	_ = logger.write(NOTICE, "", v...)
}

func (logger *Sysger) Noticef(format string, v ...interface{}) {
	_ = logger.write(NOTICE, format, v...)
}

func (logger *Sysger) Warn(v ...interface{}) {
	_ = logger.write(WARNING, "", v...)
}

func (logger *Sysger) Warnf(format string, v ...interface{}) {
	_ = logger.write(WARNING, format, v...)
}

func (logger *Sysger) Error(v ...interface{}) {
	_ = logger.write(ERROR, "", v...)
}

func (logger *Sysger) Errorf(format string, v ...interface{}) {
	_ = logger.write(ERROR, format, v...)
}

func (logger *Sysger) Critical(v ...interface{}) {
	_ = logger.write(CRITICAL, "", v...)
}

func (logger *Sysger) Criticalf(format string, v ...interface{}) {
	_ = logger.write(CRITICAL, format, v...)
}
