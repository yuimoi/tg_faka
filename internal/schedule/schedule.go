package schedule

import (
	"fmt"
	"runtime/debug"
	"tg_go_faka/internal/services"
	"tg_go_faka/internal/utils/my_log"
	"tg_go_faka/internal/utils/tg_bot/tg_bot"
)

func clearPendingOrderScheduleFunc() {
	var err error

	defer func() {
		if r := recover(); r != nil {
			msgText := fmt.Sprintf("定时清理任务崩溃, %s", debug.Stack())
			my_log.LogError(msgText)
			services.HandlePanic(r)
		}
		if err != nil {
			msgText := fmt.Sprintf("定时清理任务出错")
			my_log.LogError(msgText)
			services.HandleError(err)
		}
	}()

	orders, err := services.ClearPendingOrder()
	for _, order := range orders {
		_ = tg_bot.DeleteMsg(order.TgID, order.MessageID)
	}
}
