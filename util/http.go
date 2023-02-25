package util

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"go-micro.dev/v4/util/log"
	"time"
)

type FastHttp struct {
	client  *fasthttp.Client
	timeout time.Duration
}

func NewClient(readTimeout, writeTimeout, maxIdleConnDuration, timeout time.Duration) *FastHttp {
	client := &fasthttp.Client{
		ReadTimeout:                   readTimeout,
		WriteTimeout:                  writeTimeout,
		MaxIdleConnDuration:           maxIdleConnDuration,
		NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
		DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
		DisablePathNormalizing:        true,
		// increase DNS cache time to an hour instead of default minute
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}
	return &FastHttp{
		client:  client,
		timeout: timeout,
	}
}

// SendGetRequest send a get request return []byte, error
func (f *FastHttp) SendGetRequest(uri string) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(uri)
	req.Header.SetMethod(fasthttp.MethodGet)
	resp := fasthttp.AcquireResponse()
	err := f.client.DoTimeout(req, resp, f.timeout)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		log.Errorf("ERR Connection error: %v\n", err)
		return nil, err
	}
	body := resp.Body()
	fasthttp.ReleaseResponse(resp)
	return body, err
}

func (f *FastHttp) SendPostRequest(reqBody interface{}, uri string, contentType string) ([]byte, error) {
	reqEntityBytes, _ := json.Marshal(reqBody)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(uri)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentTypeBytes([]byte(contentType))
	req.SetBodyRaw(reqEntityBytes)

	resp := fasthttp.AcquireResponse()
	err := f.client.DoTimeout(req, resp, f.timeout)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		log.Errorf("ERR Connection error: %v\n", err)
		return nil, err
	}
	body := resp.Body()
	fasthttp.ReleaseResponse(resp)
	return body, nil
}
