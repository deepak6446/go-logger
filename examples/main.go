package main

import "../logger"

func main() {

	var Logger *logger.LoggerStack
	Logger = &logger.LoggerStack{
		Filename: "./logs/logs.json", 		// file name 
		Async: false,                       // files will be created asynchronous if set to true
		MaxSizeInBytes: 1000000,            // 1 MB
	}

	logger.Init(Logger)
		
	logger.Log("hola!, working");
	logger.Info("Info  level log");
	logger.Warn("Warn  level log");
	logger.Error("Error level log");

}