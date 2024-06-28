package http

import (
	"bytes"
	"net/http"
	"time"
)

type CustomHTTPClient struct {
	Client  *http.Client
	BaseURL string
}

func NewCustomHTTPClient() *CustomHTTPClient {
	transport := &defaultTransport{
		Transport: http.DefaultTransport,
	}

	return &CustomHTTPClient{
		Client: &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second,
		},
		BaseURL: "",
	}
}

func (c *CustomHTTPClient) SetBaseURL(baseURL string) {
	c.BaseURL = baseURL
}

func (c *CustomHTTPClient) SetDefaultHeaders(headers map[string]string) {
	transport, ok := c.Client.Transport.(*defaultTransport)
	if !ok {
		transport = &defaultTransport{
			Transport: c.Client.Transport,
		}
		c.Client.Transport = transport
	}

	transport.Headers = headers
}

func (c *CustomHTTPClient) SetDefaultContentType(contentType string) {
	transport, ok := c.Client.Transport.(*defaultTransport)
	if !ok {
		transport = &defaultTransport{
			Transport: c.Client.Transport,
		}
		c.Client.Transport = transport
	}

	transport.Headers["Content-Type"] = contentType
}

func (c *CustomHTTPClient) SetBearerAuthorization(token string) {
	transport, ok := c.Client.Transport.(*defaultTransport)
	if !ok {
		transport = &defaultTransport{
			Transport: c.Client.Transport,
		}
		c.Client.Transport = transport
	}

	transport.Authorization = "Bearer " + token
}

func (c *CustomHTTPClient) SetTimeout(timeout time.Duration) {
	c.Client.Timeout = timeout
}

type defaultTransport struct {
	Transport     http.RoundTripper
	BaseURL       string
	Headers       map[string]string
	Authorization string
}

func (dt *defaultTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range dt.Headers {
		req.Header.Set(key, value)
	}

	if dt.Authorization != "" {
		req.Header.Set("Authorization", dt.Authorization)
	}

	if dt.BaseURL != "" {
		req.URL.Path = dt.BaseURL + req.URL.Path
	}

	return dt.Transport.RoundTrip(req)
}

func (c *CustomHTTPClient) Get(path string) (*http.Response, error) {
	fullURL := c.BaseURL + path

	return c.Client.Get(fullURL)
}

func (c *CustomHTTPClient) Post(path string, body []byte) (*http.Response, error) {
	fullURL := c.BaseURL + path

	var contentType string

	if transport, ok := c.Client.Transport.(*defaultTransport); ok {
		contentType = transport.Headers["Content-Type"]
	}

	return c.Client.Post(fullURL, contentType, bytes.NewBuffer(body))
}

func (c *CustomHTTPClient) Put(path string, body []byte) (*http.Response, error) {
	fullURL := c.BaseURL + path

	req, err := http.NewRequest(http.MethodPut, fullURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	if transport, ok := c.Client.Transport.(*defaultTransport); ok {
		for key, value := range transport.Headers {
			req.Header.Set(key, value)
		}
	}

	if transport, ok := c.Client.Transport.(*defaultTransport); ok && transport.Authorization != "" {
		req.Header.Set("Authorization", transport.Authorization)
	}

	return c.Client.Do(req)
}

func (c *CustomHTTPClient) Patch(path string, body []byte) (*http.Response, error) {
	fullURL := c.BaseURL + path

	req, err := http.NewRequest(http.MethodPatch, fullURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	if transport, ok := c.Client.Transport.(*defaultTransport); ok {
		for key, value := range transport.Headers {
			req.Header.Set(key, value)
		}
	}

	if transport, ok := c.Client.Transport.(*defaultTransport); ok && transport.Authorization != "" {
		req.Header.Set("Authorization", transport.Authorization)
	}

	return c.Client.Do(req)
}

func (c *CustomHTTPClient) Delete(path string) (*http.Response, error) {
	fullURL := c.BaseURL + path

	req, err := http.NewRequest(http.MethodDelete, fullURL, nil)
	if err != nil {
		return nil, err
	}

	if transport, ok := c.Client.Transport.(*defaultTransport); ok {
		for key, value := range transport.Headers {
			req.Header.Set(key, value)
		}
	}

	if transport, ok := c.Client.Transport.(*defaultTransport); ok && transport.Authorization != "" {
		req.Header.Set("Authorization", transport.Authorization)
	}

	return c.Client.Do(req)
}
