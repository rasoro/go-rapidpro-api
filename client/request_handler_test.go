package client_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/rasoro/rapidpro-api-go/client"
	"github.com/stretchr/testify/assert"
)

func NewRequestHandler() *client.RequestHandler {
	c := NewClient(token)
	return client.NewRequestHandler(c)
}

func TestRequestHandler(t *testing.T) {
	h := NewRequestHandler()
	tcs := []string{http.MethodGet, http.MethodPost, http.MethodDelete}
	for _, tc := range tcs {
		var resp *http.Response
		err := errors.New("")
		switch tc {
		case http.MethodGet:
			resp, err = h.Get(mockServer.URL, nil, nil)
		case http.MethodPost:
			resp, err = h.Post(mockServer.URL, nil, nil, nil)
		case http.MethodDelete:
			resp, err = h.Delete(mockServer.URL, nil, nil)
		}
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	}
}
