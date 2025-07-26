package logger

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"
)

func PrintLog(message string, level ...string) {
	color := "\033[32m"
	logLevel := ""
	if len(level) > 0 {
		logLevel = level[0]
	}

	switch strings.ToUpper(logLevel) {
	case "ERROR":
		color = "\033[31m"
	case "WARN":
		color = "\033[33m"
	default:

	}

	pc, file, line, _ := runtime.Caller(1)

	now := time.Now()
	output := fmt.Sprintf(
		"\n[LOGGER] %s%d/%d/%d - %d:%d:%d %s %s:%d %s\x1b[0m\n",
		color,
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		runtime.FuncForPC(pc).Name(),
		path.Base(file),
		line,
		strings.ToUpper(message),
	)
	fmt.Printf("%s\n", output)
}
