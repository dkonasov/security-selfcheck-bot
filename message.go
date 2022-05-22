package main

type Message struct {
	Text string `json:"text"`
	From User   `json:"from"`
}
