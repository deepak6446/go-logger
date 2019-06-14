package logger

import (
    
    "os"
    "sync"
	"time"
    "fmt"
    "log"
    "strings"
    "strconv"
    "io"

    js "encoding/json"
    
)

type LoggerStack struct {
    
    lock                sync.Mutex
    Filename            string                             
    fp                  *os.File
    Async               bool
    currentFileIndex    int 
    bytesLength         int64 
    MaxSizeInBytes      int64 

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
var currentFile string
func logToFile(json JsonLog) {

    bytes, _ := js.Marshal(json)
    Logger.fp.Write(bytes)
    Logger.fp.Write([]byte("\n"))
    
    rotate(int64(len(bytes)))

} 

func makeFile() (err error) {
  
    var FileName string
    
    preFix := Logger.Filename[:strings.LastIndex(Logger.Filename, ".json")]
    if Logger.currentFileIndex == 0 {
        FileName = preFix  + ".json"
    }else {
        FileName = preFix + strconv.Itoa(Logger.currentFileIndex) + ".json"
        return moveFile(FileName, preFix  + ".json")
    }

    fmt.Println("creating new log file: ",FileName)

    Logger.fp, err = os.OpenFile(FileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666) // read and write
    if err !=  nil {
        fmt.Println("error in openfile error: ", err)
    }else {
        Logger.bytesLength = fileSize(Logger.fp)
    }
    
    return err

}

func fileSize(fp *os.File) int64 {
    fi, err := fp.Stat()
    if err != nil {
      // Could not obtain stat, handle error
    }
    return fi.Size()
}

func moveFile(destFileName string, scrFileName string) (error error) {
    
    var destFile, srcFile *os.File

    Close()
    destFile, error = os.OpenFile(destFileName, os.O_CREATE|os.O_RDWR, 0666) // read and write
    srcFile, error = os.OpenFile(scrFileName, os.O_RDWR, 0666) // read and write
       
    if error != nil { 
        fmt.Println("error in create file srcFile: ", scrFileName, "destFile: ", destFileName)
        fmt.Println(" error: ", error)
        return error;      
    }

    if _, error = io.Copy(destFile, srcFile/*, scrFile*/); error != nil {
        fmt.Println("error in fileCopy src: ", scrFileName, "dest: ", destFileName)
        fmt.Println("error: ", error)
		return 
    }

    defer destFile.Close()
    srcFile.Truncate(0)
    
    error = os.Truncate(scrFileName, 0)
	if error != nil {
		fmt.Println("error in deleting file content, error: ", error)
	}

    Logger.fp = nil
    Logger.fp, error = os.OpenFile(scrFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666) // read and write
    if error !=  nil {
        fmt.Println("error in openfile error: ", error)
    }

    return 
}

// Perform the actual act of rotating and reopening file.
func rotate(len int64) (err error) {
    Logger.bytesLength = Logger.bytesLength + len
    
    if Logger.bytesLength > Logger.MaxSizeInBytes {
        Logger.currentFileIndex++
        Logger.bytesLength = 0

        if err := makeFile(); err != nil {
            fmt.Println("error in creating new file", err)
        }

    }
    
    return
}

func Close() (err error) {
    
    if Logger.fp != nil {
        
        err = Logger.fp.Close()
        Logger.fp = nil
        
        if err != nil {
            fmt.Println("Error closing log file error: ", err)
            return
        }

    }
    return
}