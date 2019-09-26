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
	if len(args) <= 0 {
		debug.Printf(format)
		return
	}
	debug.Printf(format,args...)
}

func Info (format string, args ...interface{}) {
	if len(args) <= 0 {
		info.Printf(format)
		return
	}
	info.Printf(format,args...)
}

func Error (format string, args ...interface{}) {
	if len(args) <= 0 {
		err.Printf(format)
		return
	}
	err.Printf(format,args...)
}

func Warn (format string, args ...interface{}) {
	if len(args) <= 0 {
		warn.Printf(format)
		return
	}
	warn.Printf(format,args...)
}
