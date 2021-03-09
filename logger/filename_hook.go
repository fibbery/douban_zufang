package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

type FilenameHook struct {
	Field  string
	Skip   int
	levels []logrus.Level
}

func NewFilenameHook(levels ...logrus.Level) *FilenameHook {
	hook := &FilenameHook{
		Field:  "source",
		Skip:   10,
		levels: levels,
	}
	if len(levels) == 0 {
		hook.levels = logrus.AllLevels
	}
	return hook
}

func (f *FilenameHook) Levels() []logrus.Level {
	return f.levels
}

func (f *FilenameHook) Fire(entry *logrus.Entry) error {
	entry.Data[f.Field] = findCaller(f.Skip)
	return nil
}

func findCaller(skip int) string {
	file := ""
	line := 0
	for i := 0; i < 10; i++ {
		file, line = getCaller(i + skip)
		if !strings.HasPrefix(file, "logrus") {
			break
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func getCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0
	}
	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file, line
}
