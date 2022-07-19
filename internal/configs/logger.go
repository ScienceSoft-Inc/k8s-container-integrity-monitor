package config

import (
	"bytes"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(config *LoggerConfig) *logrus.Logger {
	l := logrus.New()
	l.Level = logrus.Level(config.Level)
	l.SetReportCaller(true)

	logfile := &lumberjack.Logger{
		Filename:   "./logs/sha256sum.log",
		MaxSize:    5, // MB
		MaxBackups: 10,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
	}

	l.SetOutput(io.MultiWriter(logfile, os.Stdout))
	l.Formatter = &formatter{"[sha256sum]"}
	return l
}

// Formatter implements logrus.Formatter interface.
type formatter struct {
	prefix string
}

// Format building log message.
func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var sb bytes.Buffer

	var newLine = "\n"
	if runtime.GOOS == "linux" {
		newLine = "\r\n"
	}

	sb.WriteString(strings.ToUpper(entry.Level.String()) + " " + entry.Time.Format(time.RFC3339) + " " + f.prefix + " " + entry.Message + " ")
	file, ok := entry.Data["file"].(string)
	if ok {
		sb.WriteString("file:" + file)
	}
	line, ok := entry.Data["line"].(int)
	if ok {
		sb.WriteString(":" + strconv.Itoa(line))
	}
	function, ok := entry.Data["function"].(string)
	if ok {
		sb.WriteString(" " + "func:" + function)
	}
	sb.WriteString(newLine)

	return sb.Bytes(), nil
}
