package logger

import (
	"fmt"
	"sync"
	"time"
)

var mu sync.Mutex

//ANSI color codes
const (
	yellow     = "\033[33m"    //timestamp
	brightBlue = "\033[1;34m"  //info
	green      = "\033[32m"    //success
	red        = "\033[31m"    //error
	reset      = "\033[0m"     //reset
)

//internal logging function
func log(prefix, color, format string, a ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	timestamp := time.Now().Format("15:04:05")
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("[%s%s%s] [%s%s%s] %s\n", yellow, timestamp, reset, color, prefix, reset, msg)
}

//Info logs general information (brightBlue)
func Info(format string, a ...interface{}) {
	log("*", brightBlue, format, a...)
}

//Success logs successful events (green)
func Success(format string, a ...interface{}) {
	log("+", green, format, a...)
}

//Error logs errors (red)
func Error(format string, a ...interface{}) {
	log("-", red, format, a...)
}