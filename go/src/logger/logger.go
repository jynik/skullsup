// SPDX License Identifier: MIT
// Wrapper around the standard log package
package logger

import (
	"io"
	golog "log"
	"os"
	"strings"
)

type LogLevel int

const (
	Debug = iota
	Info
	Error
	Silent
)

type Logger struct {
	log   *golog.Logger
	level int
}

// Create a new Logger that writes to the specified file.
// Special values of filename are as follows:
//	- <empty>	Throw away log output
//	- stdout	Write to stdout instead of a file
//	- stderr	Write to stderr instead of a file
//
// The log level may be one of the following, listed
// inorder of decreasing verbosity.
//		Debug, Info, Error, Silent
//
func New(filename string, level string) (*Logger, error) {
	var iow io.Writer
	var err error

	logger := new(Logger)

	switch filename {
	case "":
		return logger, nil
	case "stdout":
		iow = os.Stdout
	case "stderr":
		iow = os.Stderr
	default:
		flags := os.O_WRONLY | os.O_CREATE | os.O_APPEND
		iow, err = os.OpenFile(filename, flags, 0600)
		if err != nil {
			return nil, err
		}
	}

	badLevel := false
	switch strings.ToLower(level) {
	case "debug":
		logger.level = Debug
	case "":
		logger.level = Info
	case "info":
		logger.level = Info
	case "error":
		logger.level = Error
	case "silent":
		logger.level = Silent
	default:
		logger.level = Info
		badLevel = true

	}
	logger.log = golog.New(iow, "", golog.LstdFlags)
	if badLevel {
		logger.Error("Invalid log level provided (%s). Defaulting to Info.\n", level)
	}
	return logger, nil
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= Debug {
		l.log.Printf("[D] "+format, v...)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= Info {
		l.log.Printf("[I] "+format, v...)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= Error {
		l.log.Printf("[E] "+format, v...)
	}
}
