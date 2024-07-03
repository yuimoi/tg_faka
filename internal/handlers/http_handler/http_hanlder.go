package http_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tg_go_faka/internal/services"
	"tg_go_faka/internal/utils/config"
)

func EpayNotify(c *gin.Context) {
	orderIDString := c.Query("out_trade_no")
	orderID, err := strconv.ParseInt(orderIDString, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "id格式错误")
		return
	}

	order, err := services.GetOrderByOrderID(orderID)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	epayConfig := config.EpayConfig
	epayKey := epayConfig.Key
	if err := services.EpayNotify(order, epayKey, c); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.String(http.StatusOK, "success")
}
