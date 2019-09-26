package myLog

import (
	"log"
	"os"
)

var (
	debug *log.Logger
	info  *log.Logger
	err *log.Logger
	warn *log.Logger
)

func init() {
	debug = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	err = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	warn = log.New(os.Stderr, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Dubug (format string, args ...interface{}) {
	debug.Println(format,args)
}

func Info (format string, args ...interface{}) {
	info.Println(format,args)
}

func Error (format string, args ...interface{}) {
	err.Println(format,args)
}

func Warn (format string, args ...interface{}) {
	warn.Println(format,args)
}
