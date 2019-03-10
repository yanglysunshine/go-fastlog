package fastlog

func ExampleConsoleLogger_Info() {
	var flog FastLog = NewConsoleLogger("fastlog")
	flog.Info("hello fastlog!!!")
	flog.Warning("hello fastlog!!!")
	flog.Notice("hello fastlog!!!")
	//flog.Debug("hello fastlog!!!")
	//flog.Error("hello fastlog!!!")
	//flog.Fatal("hello fastlog!!!")

	// Output:
	// INFO: [fastlog] [info] 2019/03/09 22:33:55 hello fastlog!!!
	// INFO: [fastlog] [warning] 2019/03/09 22:33:55 hello fastlog!!!
	// INFO: [fastlog] [notice] 2019/03/09 22:33:55 hello fastlog!!!
}

func ExampleConsoleLogger_Infof() {
	var flog FastLog = NewConsoleLogger("fastlog")
	flog.Infof("%s\n", "hello fastlog!!!")
	flog.Warningf("%s\n", "hello fastlog!!!")
	flog.Noticef("%s\n", "hello fastlog!!!")
	//flog.Debug("%s\n", "hello fastlog!!!")
	//flog.Error("%s\n", "hello fastlog!!!")
	//flog.Fatal("%s\n", "hello fastlog!!!")

	// Output:
	// INFO: [fastlog] [info] 2019/03/09 22:33:55 hello fastlog!!!
	// INFO: [fastlog] [warning] 2019/03/09 22:33:55 hello fastlog!!!
	// INFO: [fastlog] [notice] 2019/03/09 22:33:55 hello fastlog!!!
}

func ExampleRotateLogger_Info() {
	var flog FastLog = NewRotateLogger("./fastlog", "fastlog", 1024*1024, 5)
	flog.Info("hello fastlog!!!")
	flog.Warningf("hello fastlog!!!")
	flog.Noticef("hello fastlog!!!")
}
