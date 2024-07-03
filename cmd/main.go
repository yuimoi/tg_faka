package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"tg_go_faka/internal/router"
	"tg_go_faka/internal/schedule"
	"tg_go_faka/internal/utils/config"
	"tg_go_faka/internal/utils/db"
	"tg_go_faka/internal/utils/tg_bot/tg_bot"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	port := flag.String("port", "8087", "http运行端口")
	flag.Parse()

	config.LoadAllConfig()
	db.InitDB()

	tg_bot.InitTGBot()

	schedule.StartSchedule()

	tgRouter := router.NewTgRouter()
	go tgRouter.Run()

	r := router.SetupGinRoutes()

	runningPath := fmt.Sprintf("0.0.0.0:%s", *port)
	fmt.Printf("http将运行在端口: %s\n", runningPath)

	r.Run(runningPath)
}
