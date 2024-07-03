package services

import (
	"fmt"
	"runtime"
	"strings"
	"tg_go_faka/internal/utils/my_log"
)

// 处理定时任务中的panic
func HandlePanic(r interface{}) {
	var msg string
	threshold := 10 // 增加堆栈跟踪深度

	// 获取堆栈跟踪信息
	var stackTrace []string
	for skip := 0; ; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}

		funcName := runtime.FuncForPC(pc).Name()
		if !strings.Contains(funcName, "runtime.") && !strings.Contains(file, "handle_defender.go") {
			stackTrace = append(stackTrace, fmt.Sprintf("Function: %s, File: %s, Line: %d", funcName, file, line))
		}
		if len(stackTrace) >= threshold {
			break
		}
	}

	// 组合消息
	if len(stackTrace) == 0 {
		msg = "Unable to retrieve panic information."
	} else {
		msg = fmt.Sprintf("Panic occurred: %v\nStack trace:\n%s", r, strings.Join(stackTrace, "\n"))
	}

	// 日志记录
	my_log.LogError(msg)

}

func HandleError(err error) {
	my_log.LogError(fmt.Sprintf("发生错误: %s", err.Error()))

}
