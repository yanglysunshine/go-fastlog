package fastlog

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
	LdebugFlags   = Ldate | Ltime | Llongfile
)

var (
	layout string = "2006_01_02"
)

var (
	info    = fmt.Sprint(aurora.Bold(aurora.Green("[info] ")))
	notice  = fmt.Sprint(aurora.Bold(aurora.Blue("[notice] ")))
	warning = fmt.Sprint(aurora.Bold(aurora.Red("[warning] ")))
	fatal   = fmt.Sprint(aurora.Bold(aurora.Magenta("[fatal] ")))
	errFlag = fmt.Sprint(aurora.Bold(aurora.Magenta("[error] ")))
	debug   = fmt.Sprint(aurora.Bold(aurora.Cyan("[debug] ")))
)

type RotateLogger struct {
	mu                             sync.Mutex     // ensures atomic writes; protects the following fields
	prefix                         [6]string      // prefix to write at beginning of each line
	flag                           [8]int         // properties
	out                            io.WriteCloser // destination for output
	buf                            []byte         // for accumulating text to write
	fileName, dir                  string         // rotate file name
	fileSize                       int64          // rotate file size,rotate file Count, Fixed time rotation every day 0-23
	curRotate, duration, fileCount int            //
}

func NewRotateLogger(fileName string, module string, fileSize int64, fileCount int, args ...interface{}) *RotateLogger {
	if fileSize <= 0 || fileCount <= 0 || fileName == "" {
		panic("fileSize or fileCount cannot be less than or equal to 0")
	}
	var (
		err                         error
		fd                          *os.File
		logger                      *RotateLogger
		duration, curRotate         int
		absPath, dir, file, fmtTime string
	)
	fmtTime = time.Now().Format(layout)
	if len(args) == 0 {
		duration = -1
	}
	absPath, err = filepath.Abs(fileName)
	if err != nil {
		panic(err)
	}

	if !isExist(absPath) {
		dir, file = filepath.Split(absPath)
		if !isExist(dir) {
			err = os.MkdirAll(dir, 0777)
			if err != nil {
				panic(err)
			}
			if fd, err = os.Create(dir + string(filepath.Separator) + file + "-" + fmtTime); err != nil {
				panic(err)
			}

		} else {
			fileName = dir + string(filepath.Separator) + file + "-" + fmtTime
			if isExist(fileName) {
				if fd, err = os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0666); err != nil {
					panic(err)
				}
			} else {
				if fd, err = os.Create(dir + string(filepath.Separator) + file + "-" + fmtTime); err != nil {
					panic(err)
				}
			}

		}
	} else {
		dir, file = filepath.Split(absPath)
		if ss := strings.Split(file, "-"); len(ss) != 0 {
			file = ss[0]
		}
		if fd, err = os.OpenFile(absPath, os.O_RDWR|os.O_APPEND, 0666); err != nil {
			panic(err)
		}
	}
	layout = fmtTime
	logger = &RotateLogger{
		out: fd,
		prefix: [6]string{
			"[" + module + "] " + info,
			"[" + module + "] " + notice,
			"[" + module + "] " + warning,
			"[" + module + "] " + fatal,
			"[" + module + "] " + errFlag,
			"[" + module + "] " + debug},
		flag:      [8]int{Ldate, Ltime, Lmicroseconds, Llongfile, Lshortfile, LUTC, LstdFlags, LdebugFlags},
		fileSize:  fileSize,
		fileCount: fileCount,
		duration:  duration,
		fileName:  file,
		dir:       dir,
		curRotate: curRotate,
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGABRT)
	go logger.loop(quit)
	return logger
}

func isExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

func fileSize(fileName string) (size int64) {
	info, err := os.Stat(fileName)
	if err != nil {
		panic(err)
	}
	return info.Size()
}

func calcNextTime(hour int) <-chan time.Time {
	now := time.Now()
	next := now.Add(time.Hour * 24)
	next = time.Date(next.Year(), next.Month(), next.Day(), hour, 0, 0, 0, next.Location())
	return time.NewTimer(next.Sub(now)).C
}

func (l *RotateLogger) switchFile() {
	layout = time.Now().Format(layout)
	var (
		curName   string
		fd        *os.File
		err       error
		curRotate int
	)
	if l.curRotate == 0 {
		curRotate = l.curRotate + 1
	}

	for i := curRotate; i < l.fileCount; i++ {
		curName = l.dir + l.fileName + "-" + layout + "." + strconv.Itoa(curRotate)
		if isExist(curName) {
			continue
		}
		l.curRotate = curRotate
		break
	}
	if isExist(curName) {
		curName = l.dir + l.fileName + "-" + layout
		if isExist(curName) {
			fd, err = os.OpenFile(curName, os.O_RDWR, 0666)
		} else {
			fd, err = os.OpenFile(curName, os.O_RDWR|os.O_CREATE, 0666)
		}
		l.curRotate = 0
		l.SetOutput(fd)
		return
	}

	fd, err = os.Create(curName)
	if err != nil {
		panic(err)
	}
	l.SetOutput(fd)
}

func (l *RotateLogger) check() bool {
	curName := l.dir + l.fileName + "-" + layout + "." + strconv.Itoa(l.curRotate)
	curSize := fileSize(curName)
	if curSize >= l.fileSize {
		return true
	}
	return false
}

func (l *RotateLogger) loop(quit chan os.Signal) {
	var t <-chan time.Time
	if l.duration != -1 {
		t = calcNextTime(l.duration)
	}

	for {
		select {
		case <-t: //If t is empty then it will never be selected
			t = calcNextTime(l.duration)
			l.switchFile()
		case <-time.NewTicker(time.Duration(time.Second * 5)).C:
			if l.check() {
				l.switchFile()
			}
		case sig := <-quit:
			l.Infoln("Received signal: ", sig)
			l.out.Close()
			return
		}
	}
}

func (l *RotateLogger) SetOutput(w io.WriteCloser) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

func (l *RotateLogger) formatHeader(buf *[]byte, t time.Time, file string, line, prefixFlag, properties int) {
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

func (l *RotateLogger) Output(calldepth int, s string, flag, properties int) error {
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

func (l *RotateLogger) Info(v ...interface{}) { l.Output(2, fmt.Sprint(v...), 0, 6) }
func (l *RotateLogger) Infof(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf(format, v...), 0, 6)
}
func (l *RotateLogger) Infoln(v ...interface{}) { l.Output(2, fmt.Sprintln(v...), 0, 6) }

func (l *RotateLogger) Notice(v ...interface{}) { l.Output(2, fmt.Sprint(v...), 1, 6) }
func (l *RotateLogger) Noticef(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf(format, v...), 1, 6)
}
func (l *RotateLogger) Noticeln(v ...interface{}) { l.Output(2, fmt.Sprintln(v...), 1, 6) }

func (l *RotateLogger) Warning(v ...interface{}) { l.Output(2, fmt.Sprint(v...), 2, 6) }
func (l *RotateLogger) Warningf(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf(format, v...), 2, 6)
}
func (l *RotateLogger) Warningln(v ...interface{}) { l.Output(2, fmt.Sprintln(v...), 2, 6) }

func (l *RotateLogger) Fatal(v ...interface{}) { l.Output(2, fmt.Sprint(v...), 3, 6); os.Exit(1) }
func (l *RotateLogger) Fatalf(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf(format, v...), 3, 6)
	os.Exit(1)
}
func (l *RotateLogger) Fatalln(v ...interface{}) { l.Output(2, fmt.Sprintln(v...), 3, 6); os.Exit(1) }

func (l *RotateLogger) Error(v ...interface{}) { s := fmt.Sprint(v...); l.Output(2, s, 4, 6); panic(s) }
func (l *RotateLogger) Errorf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.Output(2, fmt.Sprintf(format, v...), 4, 6)
	panic(s)
}
func (l *RotateLogger) Errorln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	l.Output(2, s, 4, 6)
	panic(s)
}

func (l *RotateLogger) Debug(v ...interface{}) { l.Output(2, fmt.Sprint(v...), 5, 7) }
func (l *RotateLogger) Debugf(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf(format, v...), 5, 7)
}
func (l *RotateLogger) Debugln(v ...interface{}) { l.Output(2, fmt.Sprintln(v...), 5, 7) }

func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}
