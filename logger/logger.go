package logger

import (
    
    "os"
    "sync"
	"time"
    "fmt"
    "log"
    "strings"
    "strconv"

    js "encoding/json"
    
)

type LoggerStack struct {
    
    lock                sync.Mutex
    Filename            string                             
    fp                  *os.File
    Async               bool
    currentFileIndex    int 
    bytesLength         int 
    MaxSizeInBytes      int 

}

type JsonLog map[string]string
var Logger *LoggerStack
var json JsonLog 

// func Init(Filename string, Async bool, maxSize int) {
func Init(info *LoggerStack) {
    
    Logger = info
    Logger.currentFileIndex = 0
    
    Logger.lock.Lock()
    defer Logger.lock.Unlock()

    json = make(JsonLog, 0)
    if err := makeFile(); err != nil {
        log.Fatal("Error in init logger error: ", err)
    }

}

// Write info level logs to file
func Info(v ...interface{}) {
    
    output := fmt.Sprint(v...)
    if Logger.Async == true { 
        go logFile(output, "INFO") 
    }else {
        logFile(output, "INFO")
    }
	
}

// Write Error level logs to file
func Error(v ...interface{}) {
    
    output := fmt.Sprint(v...)
    if Logger.Async == true {
        go logFile(output, "ERROR") 
    }else {
        logFile(output, "ERROR")
    }
	
}

// Write Log level logs to file
func Log(v ...interface{}) {
    
    output := fmt.Sprint(v...)
    if Logger.Async == true {
        go logFile(output, "LOG") 
    }else {
        logFile(output, "LOG")
    }

}

// Write Warn level logs to file
func Warn(v ...interface{}) {
    
    output := fmt.Sprint(v...)
    if Logger.Async == true {
        go logFile(output, "WARN") 
    }else {
        logFile(output, "WARN")
    }
	
}

func Fatal(v ...interface{}) {
    log.Fatal(v...)
}

func logFile(message string, level string) {

	Logger.lock.Lock()
	defer Logger.lock.Unlock()

    col := getColor(level)
	fmt.Println(col, message)
	
	json["message"] = message
	json["time"] = time.Now().String()
    json["level"] = level
    logToFile(json)
    
}

func getColor(level string) (string) {
    
    switch level {
	case "INFO":
		return "\x1b[32;1m";   
	case "ERROR":
        return "\x1b[31;1m";
    case "LOG":
        return "\x1b[0m";
    case "WARN":
        return "\x1b[33;1m";
	default:
		return ""
    }
    
}

func logToFile(json JsonLog) {

    bytes, _ := js.Marshal(json)
    Logger.fp.Write(bytes)
    Logger.fp.Write([]byte("\n"))
    
    rotate(len(bytes))

} 

func makeFile() (err error) {
  
    var FileName string
    
    preFix := Logger.Filename[:strings.LastIndex(Logger.Filename, ".json")]
    if Logger.currentFileIndex == 0 {
        FileName = preFix  + ".json"
    }else {
        FileName = preFix + strconv.Itoa(Logger.currentFileIndex) + ".json"
    }
    
    fmt.Println("creating new log file: ",FileName)
    Logger.fp, err = os.OpenFile(FileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600) // read and write
    
    return

}

// Perform the actual act of rotating and reopening file.
func rotate(len int) (err error) {
    
    Logger.bytesLength = Logger.bytesLength + len
    
    if Logger.bytesLength > Logger.MaxSizeInBytes {
        Logger.currentFileIndex++
        Logger.bytesLength = 0

        makeFile()

    }
    
    return
}

func Close() (err error) {
    
    if Logger.fp != nil {
        
        fmt.Println("Closing log file at: ", Logger.Filename)
        err = Logger.fp.Close()
        Logger.fp = nil
        
        if err != nil {
            fmt.Println("Error closing log file error: ", err)
            return
        }

    }
    return
}