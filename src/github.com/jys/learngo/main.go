package main

import (
	"errors"
	"fmt"
	"net/http"
)

type result struct {
	url    string
	status string
}

var errRequestFailed = errors.New("request failed")

func main() {
	results := make(map[string]string)
	channels := make(chan result)

	urls := []string{
		"https://www.baidu.com",
		"https://www.google.com",
		"https://www.amazon.com",
		"https://www.airbnb.com",
		"https://www.reddit.com",
		"https://www.soundcloude.com",
		"https://www.facebook.com",
		"https://www.instagram.com",
	}

	for _, url := range urls {
		go hitURL(url, channels)
	}

	for i := 0; i < len(urls); i++ {
		result := <-channels
		results[result.url] = result.status
	}

	for url, status := range results {
		fmt.Println(url, status)
	}
}

func hitURL(url string, c chan<- result) {
	fmt.Println("Checking:", url)
	resp, err := http.Get(url)
	status := "OK"
	if err != nil || resp.StatusCode >= 400 {
		status = "FAILED"
	}
	c <- result{url: url, status: status}
}
