package chapter4

import (
	"fmt"
	"net/http"
)

func PrintError() {
	checkStatus := func(
		done <-chan interface{},
		urls ...string,
	) <-chan *http.Response {
		ch := make(chan *http.Response)

		go func() {
			defer close(ch)
			for _, url := range urls {
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println(err)
					continue
				}

				select {
				case <-done:
					return
				case ch <- resp:
				}
			}
		}()

		return ch
	}

	done := make(chan interface{})
	defer close(done)

	urls := []string{
		"https://www.google.com",
		"https://badhost",
	}

	for resp := range checkStatus(done, urls...) {
		fmt.Printf("Response: %v\n", resp.Status)
	}
}

func HandleError() {
	type Result struct {
		Error    error
		Response *http.Response
	}

	get := func(
		done <-chan interface{},
		urls ...string,
	) <-chan Result {
		result := make(chan Result)

		go func() {
			defer close(result)

			for _, url := range urls {
				resp, err := http.Get(url)

				select {
				case <-done:
					return
				case result <- Result{Response: resp, Error: err}:
				}
			}
		}()

		return result
	}

	done := make(chan interface{})
	defer close(done)

	urls := []string{"https://www.google.com", "https:badhost"}
	for result := range get(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v", result.Error)
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}
