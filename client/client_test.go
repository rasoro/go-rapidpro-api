package client_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	rapidpro "github.com/rasoro/rapidpro-api-go/client"
	"github.com/stretchr/testify/assert"
)

var mockServer *httptest.Server
var testClient *rapidpro.Client
var token = "token123"

func NewClient(token string) *rapidpro.Client {
	c := &rapidpro.Client{
		Credentials: &rapidpro.Credentials{Token: token},
		HTTPClient:  http.DefaultClient,
	}
	return c
}

func TestMain(m *testing.M) {
	mockServer = httptest.NewServer(http.HandlerFunc(
		func(writer http.ResponseWriter, r *http.Request) {
		}))
	defer mockServer.Close()
	testClient = NewClient("token123")
	os.Exit(m.Run())
}

func TestClient_SendRequestError(t *testing.T) {
	errorServer := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(400)
			_, _ = w.Write([]byte("{}"))
		}))
	defer errorServer.Close()

	resp, err := testClient.SendRequest(http.MethodGet, errorServer.URL, nil, nil, nil)
	rapidproErr := err.(*rapidpro.RapidproRestError)
	assert.Nil(t, resp)
	assert.Equal(t, 400, rapidproErr.Status)
}

func TestClient_SendRequestCreatesClient(t *testing.T) {
	c := &rapidpro.Client{
		Credentials: &rapidpro.Credentials{Token: token},
	}
	resp, err := c.SendRequest(http.MethodGet, mockServer.URL, nil, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
