package main

type GetWebhookResponse struct {
	Result struct {
		Url string `json:"url"`
	} `json:"result"`
}
