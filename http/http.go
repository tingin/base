package http

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

var client = &fasthttp.Client{
	MaxConnsPerHost:     100,
	MaxIdleConnDuration: 100 * time.Second,
	TLSConfig: &tls.Config{
		InsecureSkipVerify: true,
	},
}

var clientMutex sync.Mutex

func ReInit(maxConnsPerHost int) {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	client = &fasthttp.Client{
		MaxConnsPerHost: maxConnsPerHost,
	}
}

func Post(url string, body any, response any) error {
	result, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod("POST")
	req.SetRequestURI(url)
	req.SetBody(result)
	req.Header.SetContentType("application/json")

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := client.Do(req, resp); err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return nil
}
