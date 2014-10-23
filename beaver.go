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

var brushes = make(map[loglevel]brush)

func init() {
	brushes[DEBUG] = newBrush(cyan)
	brushes[INFO] = newBrush(green)
	brushes[NOTICE] = newBrush(green)
	brushes[WARNING] = newBrush(yellow)
	brushes[ERROR] = newBrush(red)
	brushes[CRITICAL] = newBrush(purple)

}

type StdoutLogger struct {
	level  loglevel
	color  bool
	writer *log.Logger
}

func messageFormatter(level loglevel, format string, color bool, text ...interface{}) string {
	var message string
	if format == "" {
		for _, i := range text {
			message += " " + fmt.Sprintf("%v", i)
		}
		message = levels[int(level)] + message
	} else {
		message = levels[int(level)] + " " + fmt.Sprintf(format, text...)
	}

	if color {
		message = brushes[level](message)
	}

	return message
}

func NewStdoutLogger(out io.Writer, prefix string, detail int, level loglevel) (*StdoutLogger, error) {
	if out == nil {
		out = os.Stdout
	}
	prefix = strings.TrimSpace(prefix) + " "

	return &StdoutLogger{
		level:  level,
		color:  false,
		writer: log.New(out, prefix, detail),
	}, nil
}

func (logger *StdoutLogger) SetLevel(level loglevel) {
	logger.level = level
}

func (logger *StdoutLogger) GetLevel() loglevel {
	return logger.level
}

func (logger *StdoutLogger) EnableColors(){
	logger.color = true
}

func (logger *StdoutLogger) DisableColors(){
	logger.color = false
}

func (logger *StdoutLogger) write(level loglevel, format string, text ...interface{}) {

	if level < logger.level {
		return
	}
	message := messageFormatter(level, format, logger.color, text...)
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
	logger.write(CRITICAL, format, v...)
}

func (logger *StdoutLogger) Critical(v ...interface{}) {
	logger.write(CRITICAL, "", v...)
}

type SyslogLogger struct {
	level  loglevel
	color  bool
	writer *syslog.Writer
}

func (logger *SyslogLogger) EnableColors(){
	logger.color = true
}

func (logger *SyslogLogger) DisableColors(){
	logger.color = false
}

func NewSyslogLogger(prefix string, level loglevel) (*SyslogLogger, error) {

	writer, err := syslog.New(0, prefix)
	if err != nil {
		return nil, err
	}

	return &SyslogLogger{
		level:  level,
		color:  false,
		writer: writer,
	}, nil
}

func (logger *SyslogLogger) SetLevel(level loglevel) {
	logger.level = level
}

func (logger *SyslogLogger) GetLevel() loglevel {
	return logger.level
}

func (logger *SyslogLogger) write(level loglevel, format string, text ...interface{}) {

	if level < logger.level {
		return
	}

	message := messageFormatter(level, format, logger.color, text...)

	switch level {
	case CRITICAL:
		_ = logger.writer.Crit(message)
	case ERROR:
		_ = logger.writer.Err(message)
	case WARNING:
		_ = logger.writer.Warning(message)
	case NOTICE:
		_ = logger.writer.Notice(message)
	case INFO:
		_ = logger.writer.Info(message)
	case DEBUG:
		_ = logger.writer.Debug(message)
	default:
	}
}

func (logger *SyslogLogger) Debug(v ...interface{}) {
	logger.write(DEBUG, "", v...)
}

func (logger *SyslogLogger) Debugf(format string, v ...interface{}) {
	logger.write(DEBUG, format, v...)
}

func (logger *SyslogLogger) Info(v ...interface{}) {
	logger.write(INFO, "", v...)
}

func (logger *SyslogLogger) Infof(format string, v ...interface{}) {
	logger.write(INFO, format, v...)
}

func (logger *SyslogLogger) Notice(v ...interface{}) {
	logger.write(NOTICE, "", v...)
}

func (logger *SyslogLogger) Noticef(format string, v ...interface{}) {
	logger.write(NOTICE, format, v...)
}

func (logger *SyslogLogger) Warn(v ...interface{}) {
	logger.write(WARNING, "", v...)
}

func (logger *SyslogLogger) Warnf(format string, v ...interface{}) {
	logger.write(WARNING, format, v...)
}

func (logger *SyslogLogger) Error(v ...interface{}) {
	logger.write(ERROR, "", v...)
}

func (logger *SyslogLogger) Errorf(format string, v ...interface{}) {
	logger.write(ERROR, format, v...)
}

func (logger *SyslogLogger) Critical(v ...interface{}) {
	logger.write(CRITICAL, "", v...)
}

func (logger *SyslogLogger) Criticalf(format string, v ...interface{}) {
	logger.write(CRITICAL, format, v...)
}
