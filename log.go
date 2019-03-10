package fastlog

type FastLog interface {

	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Infoln(v ...interface{})

	Notice(v ...interface{})
	Noticef(format string, v ...interface{})
	Noticeln(v ...interface{})

	Warning(v ...interface{})
	Warningf(format string, v ...interface{})
	Warningln(v ...interface{})

	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Errorln(v ...interface{})

	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Debugln(v ...interface{})
}