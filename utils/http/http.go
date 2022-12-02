package http

import (
	"io"
	"log"
	"net/http"
)

// Get 请求网络资源，HTTP : GET
func Get(url string) []byte {
	log.Printf("HTTP REQUEST [GET] : %s", url)
	resp, err := http.Get(url)
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
