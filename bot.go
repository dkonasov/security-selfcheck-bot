package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Bot struct {
	token        string
	webhook_host string
	webhook_port string
	webhook      WebhookHandler
}

func (bot *Bot) doMethod(methodName string, params url.Values) (string, error) {
	method_url := "https://api.telegram.org/bot" + bot.token + "/" + methodName
	if len(params) > 0 {
		method_url += ("?" + params.Encode())
	}
	resp, err := http.Get(method_url)
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

func (bot *Bot) Start() error {
	get_webhook_response, err := bot.doMethod("getWebhookInfo", make(url.Values))
	if err != nil {
		return err
	}
	_ = get_webhook_response
	var webhook GetWebhookResponse
	err = json.Unmarshal([]byte(get_webhook_response), &webhook)
	if err != nil {
		return err
	}
	webhook_url := bot.webhook_host + "/bwh"
	if webhook.Result.Url != webhook_url {
		set_webhook_params := make(url.Values)
		set_webhook_params.Set("url", webhook_url)
		response, err := bot.doMethod("setWebhook", set_webhook_params)
		_ = response
		if err != nil {
			return err
		}
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		res, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprint(w, "Bad request")
		} else {
			var update Update
			err = json.Unmarshal(res, &update)
			if err != nil {
				w.WriteHeader(400)
				fmt.Fprint(w, "Bad request")
			} else {
				bot.webhook(update, bot)
				w.WriteHeader(200)
				fmt.Fprint(w, "OK")
			}
		}
	}
	http.HandleFunc("/bwh", handler)
	return http.ListenAndServe(":"+bot.webhook_port, nil)
}
