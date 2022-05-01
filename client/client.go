package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Credentials struct {
	Token string
}

func NewCredentials(token string) *Credentials {
	return &Credentials{Token: token}
}

type Client struct {
	*Credentials
	HTTPClient *http.Client
}

func defaultHTTPClient() *http.Client {
	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: time.Second * 15,
	}
}

func (c *Client) Token() string {
	return c.Credentials.Token
}

func (c *Client) SetToken(token string) {
	c.Credentials = NewCredentials(token)
}

func (c *Client) SetTimeout(timeout time.Duration) {
	if c.HTTPClient == nil {
		c.HTTPClient = defaultHTTPClient()
	}
	c.HTTPClient.Timeout = timeout
}

func (c *Client) SendRequest(
	method string,
	rawURL string,
	queryParams url.Values,
	body interface{},
	headers map[string]interface{},
) (*http.Response, error) {
	reader := &strings.Reader{}
	goVersion := runtime.Version()

	if method == http.MethodPost {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = strings.NewReader(string(jsonBody))
	}

	req, err := http.NewRequest(method, rawURL, reader)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for k, v := range queryParams {
		for _, vv := range v {
			q.Add(k, vv)
		}
	}
	req.URL.RawQuery = q.Encode()

	token := fmt.Sprintf("Token %s", c.Credentials.Token)
	req.Header.Set("Authorization", token)

	userAgent := fmt.Sprintf("rapidro-api-go (%s %s) go/%s", runtime.GOOS, runtime.GOARCH, goVersion)
	req.Header.Add("User-Agent", userAgent)

	if method == http.MethodPost {
		req.Header.Add("Content-Type", "aplication/json")
	}

	for k, v := range headers {
		req.Header.Add(k, fmt.Sprint(v))
	}

	return c.do(req)
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	client := c.HTTPClient

	if client == nil {
		client = defaultHTTPClient()
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 400 {
		details := make(map[string]interface{})

		if decodeErr := json.NewDecoder(res.Body).Decode(&details); decodeErr != nil {
			err = errors.Wrap(decodeErr, "error decoding response for HTTP error code: "+strconv.Itoa(res.StatusCode))
			return nil, err
		}

		err = &RapidproRestError{
			Status:  res.StatusCode,
			Details: details,
		}

		return nil, err
	}
	return res, nil
}
