package main

import (
	"fmt"
	"io"
	"net/http"
)

type Bot struct {
	token        string
	webhook_host string
	webhook_port string
	webhook      WebhookHandler
}

func (bot *Bot) doMethod(methodName string) (string, error) {
	resp, err := http.Get("https://api.telegram.org/bot" + bot.token + "/" + methodName)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (bot *Bot) Start() {
	// TODO: отладочный вывод, надо удалить
	fmt.Println("Bot host:" + bot.webhook_host)
	fmt.Println("Bot token:" + bot.webhook_host)
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, "OK")
	}
	http.HandleFunc("/bwh", handler)
	http.ListenAndServe(":"+bot.webhook_port, nil)
}
