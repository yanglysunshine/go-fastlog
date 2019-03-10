# fastlog


- download
```bash
    go get -u -v "github.com/trunkszi/fastlog"
```


Branch    | style | Coverage
----------|-------|----------
master    | ![CircleCI](https://github.com/fastlog/go-fastlog/blob/master/style.png) | ![CircleCI](https://github.com/fastlog/go-fastlog/blob/master/cover.png)


- example ConsoleLogger
```go
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
```


- example RotateLogger
```go
    var flog FastLog = NewRotateLogger("./fastlog", "fastlog", 1024*1024, 5)
	flog.Info("hello fastlog!!!")
	flog.Warningf("hello fastlog!!!")
	flog.Noticef("hello fastlog!!!")
```

- LICENSE

DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE