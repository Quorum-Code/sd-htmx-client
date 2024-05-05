package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func (wsc *WSConfig) eventHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/event-stream")
	resp.Header().Set("Cache-Control", "no-cache")
	resp.Header().Set("Connection", "keep-alive")

	dataCh := make(chan string)

	_, cancel := context.WithCancel(req.Context())
	defer cancel()

	go func() {
		for data := range dataCh {
			fmt.Fprintf(resp, "data: %s\n\n", data)
			resp.(http.Flusher).Flush()
		}
	}()

	x := 0
	for {
		dataCh <- strconv.Itoa(x)
		time.Sleep(1 * time.Second)
		x++
	}
}

func (wsc WSConfig) redirectHandler(resp http.ResponseWriter, req *http.Request) {
	time.Sleep(3 * time.Second)

	resp.Header().Set("Content-Type", "text/event-stream")
	resp.Header().Set("Cache-Control", "no-cache")
	resp.Header().Set("Connection", "keep-alive")

	fmt.Fprintf(resp, "event: redirect\n")
	data, _ := json.Marshal(map[string]string{"redirectTo": "/logged-in"})
	fmt.Fprintf(resp, "data: %s\n\n", data)

	resp.(http.Flusher).Flush()
}
