package router

import (
	"github.com/gin-gonic/gin"
	"tg_go_faka/internal/handlers/http_handler"
	"tg_go_faka/internal/handlers/tg_handler"
	"tg_go_faka/internal/utils/tg_bot/tg_bot_router"
)

func SetupGinRoutes() *gin.Engine {
	r := gin.Default()

	r.GET("/api/epay_notify", http_handler.EpayNotify)

	return r
}

func NewTgRouter() *tg_bot_router.TgRouter {
	tgRouter := &tg_bot_router.TgRouter{}
	tgRouter.Command("start", tg_handler.StartCommand)
	tgRouter.Default(tg_handler.StartCommand)
	tgRouter.Callback("start", tg_handler.StartCommand)

	tgRouter.Command("products", tg_handler.ProductsCommand)
	tgRouter.Callback("products_{:page}", tg_handler.ProductsCallback)
	tgRouter.Callback("product_{:product_id}", tg_handler.ProductConfirmCallback)
	tgRouter.Callback("pay_product_{:product_id}", tg_handler.PayProductCallback)

	tgRouter.Callback("delete", tg_handler.DeleteCallback)

	tgRouter.Command("add_products\n{:content}", tg_handler.AddProducts)
	tgRouter.Command("add_product_items /{:product_id}/\n{:content}", tg_handler.AddProductItems)
	tgRouter.Command("clear_products", tg_handler.ClearProducts)
	tgRouter.Command("clear_product_items", tg_handler.ClearProductItems)
	tgRouter.Command("view_products", tg_handler.ViewProducts)

	return tgRouter
}
