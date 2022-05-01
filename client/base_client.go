package client

import (
	"net/http"
	"net/url"
	"time"
)

type BaseClient interface {
	SetTimeout(timeout time.Duration)
	SendRequest(
		method string,
		rawURL string,
		queryParams url.Values,
		body interface{},
		headers map[string]interface{},
	) (*http.Response, error)
	Token() string
}
