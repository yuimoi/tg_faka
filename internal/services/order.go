package services

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/url"
	"sort"
	"strings"
	"tg_go_faka/internal/models"
	"tg_go_faka/internal/utils/config"
	"tg_go_faka/internal/utils/db"
	"tg_go_faka/internal/utils/tg_bot/tg_bot"
	_type "tg_go_faka/internal/utils/type"

	"time"
)

func GetSuccessOrderByTgID(tgID int64) (*models.Order, error) {
	var order *models.Order
	result := db.DB.Model(models.Order{}).Where("tg_id=? and status=?", tgID, _type.OrderStatusSuccess).Find(&order)
	if result.Error != nil {
		return nil, errors.New("查询订单失败")
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("没有找到订单")
	}
	return order, nil
}

func EpayUrl(orderID string, price decimal.Decimal, productName string, epayConfig config.EpayConfigStruct) string {
	siteConfig := config.GetSiteConfig()
	notifyHost := siteConfig.Host
	submitData := map[string]string{
		"pid":          epayConfig.Pid,
		"type":         epayConfig.PayType,
		"out_trade_no": orderID,
		"notify_url":   fmt.Sprintf("%s%s", notifyHost, epayConfig.NotifyUrl),
		"return_url":   fmt.Sprintf("https://web.telegram.org"),
		"name":         productName,
		"money":        price.String(),
	}
	submitData["sign"] = EpaySign(submitData, epayConfig.Key)
	submitData["sign_type"] = "MD5"

	//生成url
	values := url.Values{}
	for key, value := range submitData {
		values.Add(key, value)
	}
	payUrl := fmt.Sprintf("%s?%s", epayConfig.Url, values.Encode())
	return payUrl
}

func EpaySign(mapInput map[string]string, epayKey string) string {
	//排序key获取排序后的key列表
	var keys []string
	for k := range mapInput {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var queryParts []string
	for _, key := range keys {
		key := key
		value := mapInput[key]
		if value == "" || key == "sign" || key == "sign_type" {
			continue
		}
		queryParts = append(queryParts, key+"="+value)
	}

	stringToSign := strings.Join(queryParts, "&")
	stringToSign = stringToSign + epayKey
	stringToSign, _ = url.QueryUnescape(stringToSign)

	//md5
	inputBytes := []byte(stringToSign)
	md5Hash := md5.Sum(inputBytes)
	md5String := fmt.Sprintf("%x", md5Hash)

	return md5String
}

func GetOrderByOrderID(orderID uuid.UUID) (*models.Order, error) {
	var order *models.Order
	if result := db.DB.Model(models.Order{}).Where("id = ?", orderID).Find(&order); result.RowsAffected == 0 {
		return nil, errors.New("没有找到订单")
	}
	return order, nil
}

func EpayNotify(order *models.Order, epayKey string, c *gin.Context) error {
	notifyData := map[string]string{
		"pid":          c.Query("pid"),
		"trade_no":     c.Query("trade_no"),
		"out_trade_no": c.Query("out_trade_no"),
		"type":         c.Query("type"),
		"name":         c.Query("name"),
		"money":        c.Query("money"),
		"trade_status": c.Query("trade_status"),
		"sign":         c.Query("sign"),
		"sign_type":    c.Query("sign_type"),
	}
	inputSign := c.Query("sign")
	calculateSign := EpaySign(notifyData, epayKey)

	if inputSign != calculateSign {
		return errors.New("签名错误")
	}

	// 更新订单状态，发送消息
	err := OrderSuccess(order)
	if err != nil {
		return err
	}

	return nil

}

func OrderSuccess(order *models.Order) error {
	result := db.DB.Model(models.Order{}).Where("id=? and status=?", order.ID, _type.OrderStatusPending).Updates(map[string]interface{}{
		"status": _type.OrderStatusSuccess,
	})
	if result.Error != nil {
		return errors.New("更新订单状态错误")
	}
	if result.RowsAffected == 0 {
		return errors.New("没有该订单")
	}

	result = db.DB.Model(models.ProductItem{}).Where("id=?", order.ProductItemID).Updates(map[string]interface{}{
		"status": _type.ProductStatusInvalid,
	})
	if result.Error != nil {
		//return errors.New("更新库存状态错误")
	}

	productItem, err := GetProductItemByID(order.ProductItemID)
	if err != nil {
		return errors.New("获取库存信息失败")
	}

	// 发送成功消息
	msgText := fmt.Sprintf("支付成功\n购买内容:\n%s", productItem.Content)
	msg := tgbotapi.NewMessage(order.TgID, msgText)
	_, _ = tg_bot.Bot.Send(msg)

	// 删除支付链接
	_ = tg_bot.DeleteMsg(order.TgID, order.MessageID)

	return nil
}

func ClearPendingOrder() ([]*models.Order, error) {
	nowTimestamp := time.Now().Unix()

	// 查询符合条件的订单
	var orders []*models.Order
	result := db.DB.Where("end_time < ? AND status=?", nowTimestamp, _type.OrderStatusPending).Find(&orders)
	if result.Error != nil {
		return nil, errors.New("查询进行订单失败")
	}

	_ = ReleaseOrders(orders)

	return orders, nil
}

func ReleaseOrders(orders []*models.Order) error {
	for _, order := range orders {
		result := db.DB.Model(models.Order{}).Where("id=?", order.ID).Updates(map[string]interface{}{
			"status": _type.OrderStatusTimeout,
		})
		if result.Error != nil {
			//return errors.New("更新订单错误")
		}

		db.DB.Model(models.ProductItem{}).Where("id=?", order.ProductItemID).Updates(map[string]interface{}{
			"status": _type.ProductItemStatusValid,
		})

	}

	return nil
}

func CreateOrder(product *models.Product, tgID int64, messageID int) (*models.Order, error) {
	var newOrder *models.Order
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var validProductItem *models.ProductItem
		result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(models.ProductItem{}).Where("status=? and product_id=?", _type.ProductItemStatusValid, product.ID).Find(&validProductItem)
		if result.RowsAffected == 0 {
			return errors.New("没有可用的库存")
		}
		// 锁定库存
		result = tx.Model(models.ProductItem{}).Where("id=?", validProductItem.ID).Updates(map[string]interface{}{
			"status": _type.ProductItemStatusPending,
		})
		if result.Error != nil {
			return errors.New("锁定失败")
		}

		siteConfig := config.GetSiteConfig()
		// 创建订单
		endTime := time.Now().Add(time.Minute * time.Duration(siteConfig.OrderDurationMinutes)).Unix()
		newOrder = models.NewOrder(product.Price, tgID, messageID, endTime, validProductItem.ID)
		result = tx.Create(&newOrder)
		if result.Error != nil {
			return errors.New("创建订单失败")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return newOrder, nil
}

func GetUserPendingOrder(tgID int64) ([]*models.Order, error) {
	var orders []*models.Order
	result := db.DB.Model(models.Order{}).Where("tg_id=? and status=?", tgID, _type.OrderStatusPending).Find(&orders)
	if result.Error != nil {
		return nil, errors.New("获取用户订单失败")
	}
	return orders, nil
}
