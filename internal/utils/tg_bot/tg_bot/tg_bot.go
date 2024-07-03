package tg_bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/net/proxy"
	"net/http"
	"net/url"
	"runtime"
	"tg_go_faka/internal/utils/config"
)

var Bot *tgbotapi.BotAPI

func InitTGBot() {
	client := &http.Client{}

	myProxy := config.GetSiteConfig().Proxy
	if myProxy.EnableProxy == true && runtime.GOOS == "windows" {
		tgProxyURL, err := url.Parse(fmt.Sprintf("%s://%s:%d", myProxy.Protocol, myProxy.Host, myProxy.Port))
		if err != nil {
			panic(fmt.Sprintf("Failed to parse proxy: %s\n", err))
		}

		tgDialer, err := proxy.FromURL(tgProxyURL, proxy.Direct)
		if err != nil {
			panic(fmt.Sprintf("Failed to obtain proxy dialer: %s\n", err))
		}
		tgTransport := &http.Transport{
			Dial: tgDialer.Dial,
		}
		client.Transport = tgTransport
	}

	fmt.Println("正在连接TG")
	var err error
	Bot, err = tgbotapi.NewBotAPIWithClient(config.GetSiteConfig().TgBotToken, "https://api.telegram.org/bot%s/%s", client)
	if err != nil {
		panic(err)
	}
	fmt.Println("TG连接成功")

	//Bot.Debug = config.GetSiteConfig().EnableTGBotDebug
}

func SendMsg(tgID int64, msgText string, opts ...interface{}) error {
	msg := tgbotapi.NewMessage(tgID, msgText)
	msg.DisableWebPagePreview = true
	msg.ParseMode = "HTML"

	// 检查是否有传入匿名参数且是否为 tgbotapi.InlineKeyboardMarkup 类型
	if len(opts) > 0 {
		if keyboard, ok := opts[0].(tgbotapi.InlineKeyboardMarkup); ok {
			msg.ReplyMarkup = keyboard
		}
	}

	_, err := Bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func SendEditMsg(chatID int64, messageID int, msgText string, opts ...interface{}) error {
	newMsg := tgbotapi.NewEditMessageText(chatID, messageID, msgText)
	newMsg.DisableWebPagePreview = true
	newMsg.ParseMode = "HTML"

	// 检查是否有传入匿名参数且是否为 tgbotapi.InlineKeyboardMarkup 类型
	if len(opts) > 0 {
		if keyboard, ok := opts[0].(tgbotapi.InlineKeyboardMarkup); ok {
			newMsg.ReplyMarkup = &keyboard
		}
	}

	_, err := Bot.Send(newMsg)

	if err != nil {
		return err
	}
	return nil
}

func SendCallback(callbackID string, callbackText string) error {
	callback := tgbotapi.NewCallback(callbackID, callbackText)
	_, err := Bot.Request(callback)
	if err != nil {
		return err
	}
	return nil
}

func SendTgFile(chatID int64, tgFileBytes tgbotapi.FileBytes) error {
	req := tgbotapi.NewDocument(chatID, tgFileBytes)
	_, err := Bot.Send(req)
	if err != nil {
		return err
	}
	return nil
}
func DeleteMsg(chatID int64, msgID int) error {
	deleteConfig := tgbotapi.DeleteMessageConfig{
		ChatID:    chatID,
		MessageID: msgID,
	}
	resp, err := Bot.Request(deleteConfig)
	if err != nil {
		fmt.Println(resp)
		return err
	}
	return nil
}

type PaginationData struct {
	Page  int64
	limit int64
	Items []*PaginationItemsData
}
type PaginationItemsData struct {
	Name         string
	CallbackData string
}

func PaginationMarkup() {

}
