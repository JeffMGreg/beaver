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
    DEBUG   Level = iota
    INFO
    WARNING
    ERROR
    FATAL
)

var levelBrushMap = make(map[Level]brush)

func init() {
    levelBrushMap[DEBUG]   = newBrush(colors[DEBUG])
    levelBrushMap[INFO]    = newBrush(colors[INFO])
    levelBrushMap[WARNING] = newBrush(colors[WARNING])
    levelBrushMap[ERROR]   = newBrush(colors[ERROR])
    levelBrushMap[FATAL]   = newBrush(colors[FATAL])
}

var (
    DevLog = NewLogger(os.Stdout, "").SetFlags(DEVFLAG).SetLevel(DEBUG)
    StdLog = NewLogger(os.Stderr, "").SetFlags(STDFLAG).SetLevel(INFO)
)

var levels = []string{"[DEBUG]", "[INFO]", "[WARN]", "[ERROR]", "[FATAL]"}

var mu = &sync.Mutex{}

type Logger struct {
    out         io.Writer
    level       Level
    writer      io.Writer
    flags       int
    prefix      string
    colorEnabled bool
}

func NewLogger(out io.Writer, prefix ...string) *Logger {
    if out == nil {
        out = os.Stdout
    }
    return &Logger{
        level:       INFO,
        writer:      out,
        colorEnabled: false,
        flags:       DEVFLAG,
        prefix:      strings.Join(prefix, " "),
    }
}

func (l *Logger) SetFlags(flag int) *Logger {
    l.flags = flag
    return l
}

func (l *Logger) GetFlags() int {
    return l.flags
}

func (l *Logger) SetLevel(level Level) *Logger {
    l.level = level
    return l
}

func (l *Logger) GetLevel() Level {
    return l.level
}

func (l *Logger) EnableColor() {
    l.colorEnabled = true
}

func (l *Logger) DisableColor() {
    l.colorEnabled = false
}

func (l *Logger) write(level Level, format string, a ...interface{}) (n int, err error) {
    if level < l.level {
        return
    }
    var levelName string = levels[int(level)]
    var sep = " "
    var prefix, outstr = l.prefix, ""

    if l.flags&DATETIME != 0 {
        now := time.Now()
        layout := ""
        if l.flags&DATE != 0 {
            layout += "2006/01/02"
        }
        if l.flags&TIME != 0 {
            layout += " 15:04:05"
        }
        layout = strings.TrimSpace(layout)
        prefix += now.Format(layout)
    }

    if l.flags&FILELINE != 0 {
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

    if l.colorEnabled && l.flags&COLOR != 0 {
        outstr = levelBrushMap[level](outstr)
    }

    mu.Lock()
    defer mu.Unlock()
    return l.writer.Write([]byte(prefix + sep + outstr))
}

func (l *Logger) Debug(v ...interface{}) {
    l.write(DEBUG, "", v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
    l.write(DEBUG, format, v...)
}

func (l *Logger) Info(v ...interface{}) {
    l.write(INFO, "", v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
    l.write(INFO, format, v...)
}

func (l *Logger) Warn(v ...interface{}) {
    l.write(WARNING, "", v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
    l.write(WARNING, format, v...)
}

func (l *Logger) Error(v ...interface{}) {
    l.write(ERROR, "", v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
    l.write(ERROR, format, v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
    l.write(FATAL, format, v...)
    os.Exit(1)
}

func (l *Logger) Fatal(v ...interface{}) {
    l.write(FATAL, "", v...)
    os.Exit(1)
}
