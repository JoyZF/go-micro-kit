package util

import (
	"fmt"
	"testing"
	"time"
)

func TestFastHttp_SendGetRequest(t *testing.T) {
	bytes, err := NewClient(10*time.Second, 10*time.Second, 10*time.Second, 10*time.Second).
		SendGetRequest("https://www.baidu.com")
	fmt.Println(string(bytes))
	fmt.Println(err)
}

func TestFastHttp_SendPostRequest(t *testing.T) {
	bytes, err := NewClient(10*time.Second, 10*time.Second, 10*time.Second, 10*time.Second).
		SendPostRequest("123", "https://www.baidu.com", "")
	fmt.Println(string(bytes))
	fmt.Println(err)
}
