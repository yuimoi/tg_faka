package tg_bot_router

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"strings"
	"tg_go_faka/internal/services"
	"tg_go_faka/internal/utils/tg_bot/tg_bot"
)

type HandlerDataStruct struct {
	Bot    *tgbotapi.BotAPI
	Update tgbotapi.Update
	Params map[string]string
}

type HandlerType func(HandlerDataStruct)

type route struct {
	pattern *regexp.Regexp
	handler HandlerType
}

type TgRouter struct {
	commandHandlers  []route
	messageHandlers  []route
	callbackHandlers []route
	defaultHandler   HandlerType
}

// 用于将各种路径注册到路由中
func (r *TgRouter) Default(handler HandlerType) {
	r.defaultHandler = handler
}
func (r *TgRouter) Command(pattern string, handler HandlerType) {
	regexPattern := createRegexPattern(pattern)
	r.commandHandlers = append(r.commandHandlers, route{pattern: regexPattern, handler: handler})
}
func (r *TgRouter) Message(pattern string, handler HandlerType) {
	regexPattern := createRegexPattern(pattern)
	r.messageHandlers = append(r.messageHandlers, route{pattern: regexPattern, handler: handler})
}
func (r *TgRouter) Callback(pattern string, handler HandlerType) {
	regexPattern := createRegexPattern(pattern)
	r.callbackHandlers = append(r.callbackHandlers, route{pattern: regexPattern, handler: handler})
}

func (r *TgRouter) Run() {
	bot := tg_bot.Bot
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		go handleUpdate(r, bot, update)
	}
}

func handleUpdate(r *TgRouter, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// 处理路由，根据已注册的路由
	defer func() {
		if rec := recover(); rec != nil {
			msgText := fmt.Sprintf("机器人处理消息崩溃, Error: %v", rec)
			fmt.Println(msgText)
			services.HandlePanic(rec)
		}
	}()

	var commandHandlers []route
	if update.Message != nil {
		var patternText string
		if update.Message.IsCommand() {
			// 匹配命令
			//patternText = update.Message.Command()
			patternText = strings.Trim(update.Message.Text, "/")
			commandHandlers = r.commandHandlers
		} else {
			// 匹配文本
			patternText = update.Message.Text
			commandHandlers = r.messageHandlers
		}

		for _, route := range commandHandlers {
			if matches := route.pattern.FindStringSubmatch(patternText); matches != nil {
				params := extractParams(route.pattern, matches)
				route.handler(HandlerDataStruct{Bot: bot, Update: update, Params: params})
				return
			}
		}

		// 如果没有匹配到任何命令或消息路由，调用默认处理函数，仅处理message的
		if r.defaultHandler != nil {
			r.defaultHandler(HandlerDataStruct{Bot: bot, Update: update, Params: nil})
		}

		return
	}

	if update.CallbackQuery != nil {
		// 及时响应，不然一直转圈圈
		_ = tg_bot.SendCallback(update.CallbackQuery.ID, "")

		callbackData := update.CallbackQuery.Data
		for _, route := range r.callbackHandlers {
			if matches := route.pattern.FindStringSubmatch(callbackData); matches != nil {
				params := extractParams(route.pattern, matches)
				route.handler(HandlerDataStruct{Bot: bot, Update: update, Params: params})
				return
			}
		}

		return
	}

}

func createRegexPattern(pattern string) *regexp.Regexp {
	regexPattern := "^" + pattern
	// 使用正则表达式替换格式正确的占位符
	regexPattern = regexp.MustCompile(`\{\:(\w+)\}`).ReplaceAllString(regexPattern, `(?P<$1>[^/]+)`) // 这里创建pattern的时候，匹配变量的内容不能包含斜杠/，因此使用斜杠作为分隔符，变量匹配的时候遇到分隔符会自动停止
	regexPattern += "$"

	return regexp.MustCompile(regexPattern)
}

func extractParams(pattern *regexp.Regexp, matches []string) map[string]string {
	params := make(map[string]string)
	for i, name := range pattern.SubexpNames() {
		if i != 0 && name != "" {
			params[name] = matches[i]
		}
	}
	return params
}
