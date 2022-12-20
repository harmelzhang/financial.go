package http

import (
	"financial/config"
	"io"
	"log"
	"math/rand"
	"net/http"
)

// Get 请求网络资源，HTTP : GET
func Get(url string) []byte {
	log.Printf("HTTP REQUEST [GET] : %s", url)

	// resp, err := http.Get(url)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", config.HttpAccept)
	req.Header.Add("User-Agent", config.UserAgent[rand.Intn(len(config.UserAgent))])
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("执行出错 : %s", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		log.Fatalf("网络请求出错，Status Code: %d", resp.StatusCode)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("执行出错 : %s", err)
	}
	return bytes
}
