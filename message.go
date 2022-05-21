package main

type Message struct {
	Text string `json:"message"`
	From User   `json:"from"`
}
