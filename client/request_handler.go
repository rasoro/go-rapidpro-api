package client

import (
	"net/http"
	"net/url"
)

type RequestHandler struct {
	Client BaseClient
}

func NewRequestHandler(client BaseClient) *RequestHandler {
	return &RequestHandler{
		Client: client,
	}
}

func (c *RequestHandler) sendRequest(
	method string,
	rawURL string,
	queryParams url.Values,
	body interface{},
	headers map[string]interface{},
) (*http.Response, error) {
	return c.Client.SendRequest(method, rawURL, queryParams, body, headers)
}

func (c *RequestHandler) Post(
	path string,
	queryParams url.Values,
	body interface{},
	headers map[string]interface{},
) (*http.Response, error) {
	return c.sendRequest(http.MethodPost, path, queryParams, body, headers)
}

func (c *RequestHandler) Get(
	path string,
	queryParams url.Values,
	headers map[string]interface{},
) (*http.Response, error) {
	return c.sendRequest(http.MethodGet, path, queryParams, nil, headers)
}

func (c *RequestHandler) Delete(
	path string,
	queryParams url.Values,
	headers map[string]interface{},
) (*http.Response, error) {
	return c.sendRequest(http.MethodDelete, path, queryParams, nil, headers)
}
