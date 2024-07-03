package tg_handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg_go_faka/internal/models"
	_type "tg_go_faka/internal/utils/type"
)

func GetProductsPaginationMarkup(products []*models.Product, pagination _type.PaginationQueryDataStruct) tgbotapi.InlineKeyboardMarkup {
	itemCallbackPrefix := "product_"
	paginateCallbackPrefix := "products_"

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, product := range products {
		buttonText := fmt.Sprintf("%s 价格:%s¥ 库存:%d", product.Name, product.Price, product.InStockCount)

		row := []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(buttonText, fmt.Sprintf("%s%d", itemCallbackPrefix, product.ID))}
		rows = append(rows, row)
	}

	var paginationRow []tgbotapi.InlineKeyboardButton
	if pagination.Page > 1 {
		paginationRow = append(paginationRow, tgbotapi.NewInlineKeyboardButtonData("上一页", paginateCallbackPrefix+fmt.Sprintf("%d", pagination.Page-1)))
	}
	if pagination.Page < pagination.TotalPage {
		paginationRow = append(paginationRow, tgbotapi.NewInlineKeyboardButtonData("下一页", paginateCallbackPrefix+fmt.Sprintf("%d", pagination.Page+1)))
	}
	if len(paginationRow) != 0 {
		rows = append(rows, paginationRow)
	}
	rows = append(rows, deleteMsgRow())
	replyMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return replyMarkup
}

func deleteMsgRow() []tgbotapi.InlineKeyboardButton {
	var paymentSelectRow []tgbotapi.InlineKeyboardButton
	paymentSelectRow = append(paymentSelectRow, tgbotapi.NewInlineKeyboardButtonData("关闭", "delete"))

	return paymentSelectRow
}
