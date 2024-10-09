package http

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

type HttpHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

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

func PostWithHeader(url string, headers []HttpHeader, params any, response any) error {

	body, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal params: %w", err)
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod("POST")
	req.SetRequestURI(url)

	req.SetBody(body)

	if headers != nil {
		for _, header := range headers {
			req.Header.Set(header.Key, header.Value)
		}
	}

	req.Header.SetContentType("application/json")

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := client.Do(req, resp); err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	body = resp.Body()

	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

func Post(url string, params any, response any) error {
	return PostWithHeader(url, nil, params, response)
}

func GetWithHeader(url string, headers []HttpHeader, response any) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod("GET")
	req.SetRequestURI(url)

	if headers != nil {
		for _, header := range headers {
			req.Header.Set(header.Key, header.Value)
		}
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := client.Do(req, resp); err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	body := resp.Body()

	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w body: %s", err, body)
	}

	return nil
}

func Get(url string, response any) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod("GET")
	req.SetRequestURI(url)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := client.Do(req, resp); err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	body := resp.Body()

	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w body: %s", err, body)
	}

	return nil
}
