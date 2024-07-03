package tg_handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
	"tg_go_faka/internal/models"
	"tg_go_faka/internal/services"
	"tg_go_faka/internal/utils/config"
	"tg_go_faka/internal/utils/db"
	"tg_go_faka/internal/utils/functions"
	"tg_go_faka/internal/utils/tg_bot/tg_bot"
	"tg_go_faka/internal/utils/tg_bot/tg_bot_router"
)

var ProductsCommandString = "/products"

func StartCommand(handlerData tg_bot_router.HandlerDataStruct) {
	var chatID int64
	if handlerData.Update.Message != nil {
		chatID = handlerData.Update.Message.From.ID
	} else if handlerData.Update.CallbackQuery != nil {
		chatID = handlerData.Update.CallbackQuery.Message.Chat.ID
	} else {
		return
	}

	msgText := fmt.Sprintf(`欢迎使用发卡机器人
	点击查看商品列表: %s
	`, ProductsCommandString)

	msg := tgbotapi.NewMessage(chatID, msgText)
	_, _ = tg_bot.Bot.Send(msg)
}

func ProductConfirmCallback(handlerData tg_bot_router.HandlerDataStruct) {
	chatID := handlerData.Update.CallbackQuery.Message.Chat.ID
	messageID := handlerData.Update.CallbackQuery.Message.MessageID

	productIDString := handlerData.Params["product_id"]
	productID, err := strconv.ParseInt(productIDString, 10, 64)
	if err != nil {
		_ = tg_bot.SendMsg(chatID, "id格式错误")
		return
	}

	product, err := services.GetValidProductByID(productID)
	if err != nil {
		_ = tg_bot.SendMsg(chatID, "没有该商品")
		return
	}

	productItemsInStockCount, _ := services.GetProductItemValidCounts(product.ID)

	productConfirmRow := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("确认支付", fmt.Sprintf("%s%d", "pay_product_", product.ID)),
		tgbotapi.NewInlineKeyboardButtonData("取消", "products_1"),
	}

	replyMarkup := tgbotapi.NewInlineKeyboardMarkup(productConfirmRow)

	text := fmt.Sprintf("%s\n%s\n库存:%d\n价格: %s¥", product.Name, product.Desc, productItemsInStockCount, product.Price)

	_ = tg_bot.SendEditMsg(chatID, messageID, text, replyMarkup)
}

func PayProductCallback(handlerData tg_bot_router.HandlerDataStruct) {
	chatID := handlerData.Update.CallbackQuery.Message.Chat.ID
	messageID := handlerData.Update.CallbackQuery.Message.MessageID

	productIDString := handlerData.Params["product_id"]
	productID, err := strconv.ParseInt(productIDString, 10, 64)
	if err != nil {
		_ = tg_bot.SendMsg(chatID, "id格式错误")
		return
	}

	product, err := services.GetValidProductByID(productID)
	if err != nil {
		_ = tg_bot.SendMsg(chatID, "没有该商品")
		return
	}

	// 检查是否有未支付的订单
	orders, err := services.GetUserPendingOrder(chatID)
	if len(orders) != 0 {
		_ = tg_bot.SendMsg(chatID, "有未支付的订单")
		return
	}

	newOrder, err := services.CreateOrder(product, chatID, messageID)
	if err != nil {
		_ = tg_bot.SendMsg(chatID, err.Error())
		return
	}

	// 构建订单url
	epayConfig := *config.EpayConfig
	payUrl := services.EpayUrl(fmt.Sprintf("%d", newOrder.ID), newOrder.Price, product.Name, epayConfig)

	sendText := fmt.Sprintf("<a href=\"%s\">点击支付</a>\n请在过期前支付\n过期时间: %s", payUrl, functions.TimestampToDatetime(newOrder.EndTime))

	productConfirmRow := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("主页", "start"),
	}
	replyMarkup := tgbotapi.NewInlineKeyboardMarkup(productConfirmRow)

	_ = tg_bot.SendEditMsg(chatID, messageID, sendText, replyMarkup)

	fmt.Println(chatID, messageID)
	//_ = tg_bot.DeleteMsg(chatID, messageID)
}
func ProductsCommand(handlerData tg_bot_router.HandlerDataStruct) {
	tgID := handlerData.Update.Message.From.ID

	products, pagination, err := services.GetProductsByPage(1)
	if err != nil {
		_ = tg_bot.SendMsg(tgID, "获取商品错误")
		return
	}

	replyMarkup := GetProductsPaginationMarkup(products, pagination)

	_ = tg_bot.SendMsg(tgID, "请选择商品", replyMarkup)

}

func ProductsCallback(handlerData tg_bot_router.HandlerDataStruct) {
	chatID := handlerData.Update.CallbackQuery.Message.Chat.ID
	messageID := handlerData.Update.CallbackQuery.Message.MessageID

	pageString := handlerData.Params["page"]
	page, err := strconv.Atoi(pageString)
	if err != nil {
		_ = tg_bot.SendMsg(chatID, "page错误")
		return
	}

	products, pagination, err := services.GetProductsByPage(page)
	if err != nil {
		_ = tg_bot.SendMsg(chatID, "获取商品错误")
		return
	}
	replyMarkup := GetProductsPaginationMarkup(products, pagination)

	tg_bot.SendEditMsg(chatID, messageID, "请选择商品", replyMarkup)

}

func DeleteCallback(handlerData tg_bot_router.HandlerDataStruct) {
	chatID := handlerData.Update.CallbackQuery.Message.Chat.ID
	messageID := handlerData.Update.CallbackQuery.Message.MessageID

	_ = tg_bot.DeleteMsg(chatID, messageID)
}

func AddProducts(handlerData tg_bot_router.HandlerDataStruct) {
	// 只能是admin
	siteConfig := config.GetSiteConfig()
	if siteConfig.AdminTGID != handlerData.Update.Message.From.ID {
		return
	}

	content := handlerData.Params["content"]

	tgID := handlerData.Update.Message.From.ID

	var products []*models.Product
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if functions.IsWhitespace(line) {
			continue
		}

		splitFunc := func(c rune) bool {
			return c == ',' || c == '\t' || c == ' '
		}
		parts := strings.FieldsFunc(line, splitFunc)

		var name string
		var desc string
		var price decimal.Decimal
		var priceString string
		var err error

		if len(parts) < 3 {
			_ = tg_bot.SendMsg(tgID, "格式错误")
			return
		}
		name = parts[0]
		desc = parts[1]
		priceString = parts[2]

		price, err = decimal.NewFromString(priceString)
		if err != nil {
			_ = tg_bot.SendMsg(tgID, "金额错误")
			return
		}

		product := &models.Product{
			Name:  name,
			Desc:  desc,
			Price: price,
		}
		products = append(products, product)
	}
	result := db.DB.Create(products)
	if result.Error != nil {
		_ = tg_bot.SendMsg(tgID, result.Error.Error())
		return
	}

	_ = tg_bot.SendMsg(tgID, "添加成功")
}
func AddProductItems(handlerData tg_bot_router.HandlerDataStruct) {
	// 只能是admin
	siteConfig := config.GetSiteConfig()
	if siteConfig.AdminTGID != handlerData.Update.Message.From.ID {
		return
	}

	tgID := handlerData.Update.Message.From.ID

	content, ok := handlerData.Params["content"]
	if !ok {
		_ = tg_bot.SendMsg(tgID, "没有内容")
		return
	}

	productIDString, ok := handlerData.Params["product_id"]
	if !ok {
		_ = tg_bot.SendMsg(tgID, "没有内容")
		return
	}

	productID, err := strconv.ParseInt(productIDString, 10, 64)
	if err != nil {
		_ = tg_bot.SendMsg(tgID, "id格式错误")
		return
	}

	var productItems []*models.ProductItem
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if functions.IsWhitespace(line) {
			continue
		}

		content := line

		productItem := &models.ProductItem{
			ProductID: productID,
			Content:   content,
		}
		productItems = append(productItems, productItem)
	}
	result := db.DB.Create(productItems)
	if result.Error != nil {
		_ = tg_bot.SendMsg(tgID, result.Error.Error())
		return
	}

	_ = tg_bot.SendMsg(tgID, "添加成功")
}

func ClearProductItems(handlerData tg_bot_router.HandlerDataStruct) {
	// 只能是admin
	siteConfig := config.GetSiteConfig()
	if siteConfig.AdminTGID != handlerData.Update.Message.From.ID {
		return
	}

	chatID := handlerData.Update.Message.Chat.ID

	// 发送文件
	items, _ := services.GetAllProductItems()
	if len(items) != 0 {
		excelBytes := services.GenerateExcelFromItems(items, []string{"content", "id", "product_id", "status", "create_time"}, nil)
		_ = tg_bot.SendTgFile(chatID, tgbotapi.FileBytes{Name: "product_items.xlsx", Bytes: excelBytes})
	}

	_ = services.DeleteAllProductItems()

	_ = tg_bot.SendMsg(chatID, "删除成功")
}

func ClearProducts(handlerData tg_bot_router.HandlerDataStruct) {
	// 只能是admin
	siteConfig := config.GetSiteConfig()
	if siteConfig.AdminTGID != handlerData.Update.Message.From.ID {
		return
	}

	chatID := handlerData.Update.Message.Chat.ID

	// 发送文件
	items, _ := services.GetAllProductItems()
	if len(items) != 0 {
		excelBytes := services.GenerateExcelFromItems(items, []string{"content", "id", "product_id", "status", "create_time"}, nil)
		_ = tg_bot.SendTgFile(chatID, tgbotapi.FileBytes{Name: "product_items.xlsx", Bytes: excelBytes})
	}

	_ = services.DeleteAllProductItems()
	_ = services.DeleteAllProducts()

	_ = tg_bot.SendMsg(chatID, "删除成功")
}

func ViewProducts(handlerData tg_bot_router.HandlerDataStruct) {
	// 只能是admin
	siteConfig := config.GetSiteConfig()
	if siteConfig.AdminTGID != handlerData.Update.Message.From.ID {
		return
	}

	chatID := handlerData.Update.Message.Chat.ID

	var msgText string
	products, err := services.GetAllProducts()
	if err != nil {
		_ = tg_bot.SendMsg(chatID, "获取商品失败")
		return
	}
	for _, product := range products {
		msgText = msgText + fmt.Sprintf("%v\n", functions.StructToMap(product))
	}

	_ = tg_bot.SendMsg(chatID, msgText)

}
