package log

import (
	"fmt"
	goLog "log"
)

func Format(level, str string) string {
	return fmt.Sprintf("[%s] %s", level, str)
}

func Debug(str string) {
	goLog.Println(Format("DEBUG", str))
}

func Info(str string) {
	goLog.Println(Format("INFO", str))
}
