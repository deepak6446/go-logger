## Logging in Golang 

logging library for go with following features.
1. Specify file size for each log file and Create new log file when file size is exceeded.
2. Each log can be created in async mode ( prefer using only when logging in the main goroutine ).
3. color code's for different level log file.

[example](/examples/console.png)

Set up is as simple as:
[example](/examples/main.go)

1. Install package </br>
go get github.com/deepak6446/go-logger/logger

2. import "logger" </br>
var Logger *logger.LoggerStack </br>
<pre>
Logger = &logger.LoggerStack{</br>
	Filename: "./logs/logs.json", 		// file name 
	Async: false,                       // files will be created asynchronous if set to true 
	</t>MaxSizeInBytes: 1000000,            // 1 MB 
}
</pre>

logger.Init(Logger)

3. log using</br>
logger.Info("Info level log");

```
NOTE: use Async: true only when using large logging in main thread. In thread using Async true will cause creation of a new thread for each log.
```
