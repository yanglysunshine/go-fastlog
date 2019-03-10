package fastlog

import "testing"

var flog FastLog = NewConsoleLogger("testing")
var rlog FastLog = NewRotateLogger("./fastLog", "testing", 1024*1024, 5)

func fastlog_test() bool {
	flog.Info("hello FastLog!!!")
	flog.Infof("%s", "hello FastLog!!!")
	flog.Infoln("hello FastLog!!!")

	flog.Warning("hello FastLog!!!")
	flog.Warningf("%s", "hello FastLog!!!")
	flog.Warningln("hello FastLog!!!")

	flog.Notice("hello FastLog!!!")
	flog.Noticef("%s", "hello FastLog!!!")
	flog.Noticeln("hello FastLog!!!")

	flog.Debug("hello FastLog!!!")
	flog.Debugf("%s", "hello FastLog!!!")
	flog.Debugln("hello FastLog!!!")

	rlog.Info("hello FastLog!!!")
	rlog.Infof("%s", "hello FastLog!!!")
	rlog.Infoln("hello FastLog!!!")

	rlog.Warning("hello FastLog!!!")
	rlog.Warningf("%s", "hello FastLog!!!")
	rlog.Warningln("hello FastLog!!!")

	rlog.Notice("hello FastLog!!!")
	rlog.Noticef("%s", "hello FastLog!!!")
	rlog.Noticeln("hello FastLog!!!")

	rlog.Debug("hello FastLog!!!")
	rlog.Debugf("%s", "hello FastLog!!!")
	rlog.Debugln("hello FastLog!!!")
	return true
}

func TestConsoleLogger_Info(t *testing.T) {
	tests := []struct {
		success bool
	}{
		{true},
		{true},
		{true},
		{true},
		{true},
		{true},
		{true},
		{true},
		{true},
	}
	for _, tt := range tests {
		actual := fastlog_test()
		if actual != tt.success {
			t.Error("got ", actual, "expected ", tt.success)
		}
	}
}

// go test -bench .
// go test -bench . -cpuprofile cpu.out
func BenchmarkConsoleLogger_Info(b *testing.B) {
	for i := 0; i < b.N; i++ { //b.N是自动算出来的一个测试次数
		actual := fastlog_test()
		if actual != true {
			b.Error("got ", actual, "expected ", true)
		}
	}
}
