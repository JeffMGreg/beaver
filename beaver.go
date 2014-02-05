package beaver

import (
    "fmt"
    "io"
    "os"
    "runtime"
    "strings"
    "sync"
    "time"
)

// Colors ======================================================================
var (
    DarkGray Paint = `1;30`
    Red      Paint = `1;31`
    Green    Paint = `1;32`
    Yellow   Paint = `1;33`
    Blue     Paint = `1;34`
    Purple   Paint = `1;35`
    Cyan     Paint = `1;36`
    White    Paint = `1;37`

    colors = []Paint{Cyan, Green, Yellow, Red, Purple}
)

type brush func(string) string

type Paint string

type Style struct {
    color Paint
    code  string
}

func (s Style) brush() brush {
    return func(text string) string {
        return s.code + text + "\033[0m"
    }
}

func newBrush(color Paint) brush {
    return Style{color, "\033[" + string(color) + "m" + ``}.brush()
}

// =============================================================================

type Level int

const (
    FILELINE = 1 << iota // show filename:lineno
    DATE
    TIME
    COLOR

    DATETIME = DATE | TIME
    DEVFLAG  = DATETIME | FILELINE | COLOR // for develop use
    STDFLAG  = DATETIME | COLOR
)

const (
    DEBUG Level = iota
    INFO
    WARNING
    ERROR
    FATAL
)

var levelBrushMap = make(map[Level]brush)

func init() {
    levelBrushMap[DEBUG] = newBrush(colors[DEBUG])
    levelBrushMap[INFO] = newBrush(colors[INFO])
    levelBrushMap[WARNING] = newBrush(colors[WARNING])
    levelBrushMap[ERROR] = newBrush(colors[ERROR])
    levelBrushMap[FATAL] = newBrush(colors[FATAL])
}

var levels = []string{"[DEBUG]", "[INFO]", "[WARN]", "[ERROR]", "[FATAL]"}

var mutex = &sync.Mutex{}

type Logger struct {
    out          io.Writer
    level        Level
    writer       io.Writer
    flags        int
    prefix       string
    colorEnabled bool
}

func NewLogger(out io.Writer, prefix ...string) *Logger {
    if out == nil {
        out = os.Stdout
    }
    return &Logger{
        level:        INFO,
        writer:       out,
        colorEnabled: false,
        flags:        DEVFLAG,
        prefix:       strings.Join(prefix, " "),
    }
}

func (logger *Logger) SetFlags(flag int){
    logger.flags = flag
}

func (logger *Logger) GetFlags() int {
    return logger.flags
}

func (logger *Logger) SetLevel(level Level) {
    logger.level = level
}

func (logger *Logger) GetLevel() Level {
    return logger.level
}

func (logger *Logger) EnableColor() {
    logger.colorEnabled = true
}

func (logger *Logger) DisableColor() {
    logger.colorEnabled = false
}

func (logger *Logger) write(level Level, format string, a ...interface{}) (n int, err error) {
    if level < logger.level {
        return
    }
    var levelName string = levels[int(level)]
    var sep = " "
    var prefix, outstr = logger.prefix, ""

    if logger.flags&DATETIME != 0 {
        now := time.Now()
        layout := ""
        if logger.flags&DATE != 0 {
            layout += "2006/01/02"
        }
        if logger.flags&TIME != 0 {
            layout += " 15:04:05"
        }
        layout = strings.TrimSpace(layout)
        prefix += now.Format(layout)
    }

    if logger.flags&FILELINE != 0 {
        // Retrieve the stack infos
        _, file, line, ok := runtime.Caller(2)
        if !ok {
            file = "<unknown>"
            line = -1
        } else {
            file = file[strings.LastIndex(file, "/")+1:]
        }
        prefix = fmt.Sprintf("%s %s:%d", prefix, file, line)
    }

    outstr += levelName

    if format == "" {
        for _, i := range a {
            outstr += sep + fmt.Sprintf("%v", i)
        }
    } else {
        outstr = outstr + sep + fmt.Sprintf(format, a...)
    }
    if !strings.HasSuffix(outstr, "\n") {
        outstr += "\n"
    }

    if logger.colorEnabled && logger.flags&COLOR != 0 {
        outstr = levelBrushMap[level](outstr)
    }

    mutex.Lock()
    defer mutex.Unlock()
    return logger.writer.Write([]byte(prefix + sep + outstr))
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
    os.Exit(1)
}

func (logger *Logger) Fatal(v ...interface{}) {
    logger.write(FATAL, "", v...)
    os.Exit(1)
}
