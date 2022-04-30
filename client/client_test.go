package client_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	rapidpro "github.com/rasoro/rapidpro-api-go/client"
)

var mockServer *httptest.Server
var testClient *rapidpro.Client

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
