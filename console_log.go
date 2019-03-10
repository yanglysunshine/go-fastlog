package fastlog

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

type ConsoleLogger struct {
	mu     sync.Mutex // ensures atomic writes; protects the following fields
	prefix [6]string  // prefix to write at beginning of each line
	flag   [8]int     // properties
	out    io.Writer  // destination for output
	buf    []byte     // for accumulating text to write
}

func NewConsoleLogger(module string) *ConsoleLogger {
	return &ConsoleLogger{
		out: os.Stderr,
		prefix: [6]string{
			"[" + module + "] " + info,
			"[" + module + "] " + notice,
			"[" + module + "] " + warning,
			"[" + module + "] " + fatal,
			"[" + module + "] " + errFlag,
			"[" + module + "] " + debug},
		flag:      [8]int{Ldate, Ltime, Lmicroseconds, Llongfile, Lshortfile, LUTC, LstdFlags, LdebugFlags},
	}
}

func (l *ConsoleLogger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

func (l *ConsoleLogger) formatHeader(buf *[]byte, t time.Time, file string, line, prefixFlag, properties int) {
	*buf = append(*buf, l.prefix[prefixFlag]...)
	if l.flag[properties]&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.flag[properties]&LUTC != 0 {
			t = t.UTC()
		}
		if l.flag[properties]&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flag[properties]&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag[properties]&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}
	if l.flag[properties]&(Lshortfile|Llongfile) != 0 {
		if l.flag[properties]&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
}

func (l *ConsoleLogger) Output(calldepth int, s string, flag, properties int) error {
	now := time.Now() // get this early.
	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag[properties]&(Lshortfile|Llongfile) != 0 {
		// Release lock while getting caller info - it's expensive.
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}
	l.buf = l.buf[:0]
	l.formatHeader(&l.buf, now, file, line, flag, properties)
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, err := l.out.Write(l.buf)
	return err
}

func (l *ConsoleLogger) Info(v ...interface{})                 { l.Output(2, fmt.Sprint(v...), 0, 6) }
func (l *ConsoleLogger) Infof(format string, v ...interface{}) { l.Output(2, fmt.Sprintf(format, v...), 0, 6) }
func (l *ConsoleLogger) Infoln(v ...interface{})               { l.Output(2, fmt.Sprintln(v...), 0, 6) }

func (l *ConsoleLogger) Notice(v ...interface{})                 { l.Output(2, fmt.Sprint(v...), 1, 6) }
func (l *ConsoleLogger) Noticef(format string, v ...interface{}) { l.Output(2, fmt.Sprintf(format, v...), 1, 6) }
func (l *ConsoleLogger) Noticeln(v ...interface{})               { l.Output(2, fmt.Sprintln(v...), 1, 6) }

func (l *ConsoleLogger) Warning(v ...interface{})                 { l.Output(2, fmt.Sprint(v...), 2, 6) }
func (l *ConsoleLogger) Warningf(format string, v ...interface{}) { l.Output(2, fmt.Sprintf(format, v...), 2, 6) }
func (l *ConsoleLogger) Warningln(v ...interface{})               { l.Output(2, fmt.Sprintln(v...), 2, 6) }

func (l *ConsoleLogger) Fatal(v ...interface{})                 { l.Output(2, fmt.Sprint(v...), 3, 6); os.Exit(1) }
func (l *ConsoleLogger) Fatalf(format string, v ...interface{}) { l.Output(2, fmt.Sprintf(format, v...), 3, 6); os.Exit(1) }
func (l *ConsoleLogger) Fatalln(v ...interface{})               { l.Output(2, fmt.Sprintln(v...), 3, 6); os.Exit(1) }

func (l *ConsoleLogger) Error(v ...interface{})                 { s := fmt.Sprint(v...); l.Output(2, s, 4, 6); panic(s) }
func (l *ConsoleLogger) Errorf(format string, v ...interface{}) { s := fmt.Sprintf(format, v...); l.Output(2, fmt.Sprintf(format, v...), 4, 6); panic(s) }
func (l *ConsoleLogger) Errorln(v ...interface{})               { s := fmt.Sprintln(v...); l.Output(2, s, 4, 6); panic(s) }

func (l *ConsoleLogger) Debug(v ...interface{})                 { l.Output(2, fmt.Sprint(v...), 5, 7) }
func (l *ConsoleLogger) Debugf(format string, v ...interface{}) { l.Output(2, fmt.Sprintf(format, v...), 5, 7) }
func (l *ConsoleLogger) Debugln(v ...interface{})               { l.Output(2, fmt.Sprintln(v...), 5, 7) }
