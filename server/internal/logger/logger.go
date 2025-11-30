package logger

import (
	"log"
	"os"
)

type Logger interface {
	Info(msg string, args ...interface{})
	Er(msg string, err error, args ...interface{})
	Function(name string) Logger
}

type logger struct {
	name string
}

func New(name string) Logger {
	return &logger{name: name}
}

func (l *logger) Info(msg string, args ...interface{}) {
	logPrefix := "[" + l.name + "] INFO: "
	log.Printf(logPrefix+msg, args...)
}

func (l *logger) Er(msg string, err error, args ...interface{}) {
	logPrefix := "[" + l.name + "] ERROR: "
	if err != nil {
		log.Printf(logPrefix+msg+": %v", append(args, err)...)
	} else {
		log.Printf(logPrefix+msg, args...)
	}
}

func (l *logger) Function(name string) Logger {
	return &logger{name: l.name + "." + name}
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}