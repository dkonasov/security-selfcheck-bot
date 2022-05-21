package main

import (
	"os"
)

var token = os.Getenv("BOT_TOKEN")

func main() {
	bot := Bot{token, os.Getenv("HOST"), os.Getenv("PORT"), nil}
	bot.Start()
}
