package schedule

import "time"

func StartSchedule() {
	go startClearPendingOrderSchedule()

}

func startClearPendingOrderSchedule() {
	for {
		clearPendingOrderScheduleFunc()

		time.Sleep(time.Hour * 24)
	}
}
