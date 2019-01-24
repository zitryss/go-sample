package log

import (
	"io/ioutil"
	"log"
	"os"
)

const (
	debug = iota
	info
	warn
	error_
	critical
)

var (
	l = &logger{
		Logger: log.New(ioutil.Discard, "", log.LstdFlags),
		level:  info,
	}
)

type logger struct {
	*log.Logger
	level int
}

func Setup() {
	l.SetOutput(os.Stderr)
	Info("logging initialized")
}

func Debug(v interface{}) {
	if l.level <= debug {
		l.Printf("debug: %v\n", v)
	}
}

func Info(v interface{}) {
	if l.level <= info {
		l.Printf("info: %v\n", v)
	}
}

func Warn(v interface{}) {
	if l.level <= warn {
		l.Printf("warn: %v\n", v)
	}
}

func Error(v interface{}) {
	if l.level <= error_ {
		l.Printf("error: %+v\n", v)
	}
}

func Critical(v interface{}) {
	if l.level <= critical {
		l.Printf("critical: %+v\n", v)
	}
}
