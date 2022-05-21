package main

import (
	"net/url"
	"os"
	"strconv"
)

var token = os.Getenv("BOT_TOKEN")

func main() {
	handler := func(update Update, bot *Bot) {
		params := make(url.Values)
		params.Set("chat_id", strconv.FormatUint(update.Message.From.Id, 10))
		params.Set("text", "Привет, путник!\n\nЭто бот самопроверки по комплексной безопасности, вдохновленный одним из заданий хакатона DemHack. Пока он находится в стадии активной разработки и умеет отвечать только таким вот сообщением, но в будущем сможет больше.")
		bot.doMethod("sendMessage", params)
	}
	bot := Bot{token, os.Getenv("HOST"), os.Getenv("PORT"), handler}
	bot.Start()
}
