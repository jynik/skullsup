// SPDX License Identifier: MIT
// Wrapper around the standard log package
package logger

import (
	"io"
	golog "log"
	"os"
)

type Logger struct {
	log     *golog.Logger
	verbose bool
}

// Create a new Logger that writes to the specified file.
// Special values of filename are as follows:
//	- <empty>	Throw away log output
//	- stdout	Write to stdout instead of a file
//	- stderr	Write to stderr instead of a file
//
// If verbose is true, the Logger functions prefixed with
// "V" will be output. Otherwise, they will be supressed.
func New(filename string, verbose bool) (*Logger, error) {
	var iow io.Writer
	var err error

	logger := new(Logger)
	logger.verbose = verbose

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

	logger.log = golog.New(iow, "", golog.LstdFlags)
	return logger, nil
}

func Fatal(v ...interface{}) {
	golog.Fatal(v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	if l.log != nil {
		l.log.Fatal(v...)
	} else {
		golog.Fatal(v...)
	}
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.log.Printf(format, v...)
}

func (l *Logger) Print(v ...interface{}) {
	l.log.Print(v...)
}

func (l *Logger) Println(v ...interface{}) {
	l.log.Println(v...)
}

func (l *Logger) VPrintf(format string, v ...interface{}) {
	if l.verbose {
		l.log.Printf(format, v...)
	}
}

func (l *Logger) VPrint(v ...interface{}) {
	if l.verbose {
		l.log.Print(v...)
	}
}

func (l *Logger) VPrintln(v ...interface{}) {
	if l.verbose {
		l.log.Println(v...)
	}
}
