package client_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestClient_SendRequestWithData(t *testing.T) {
	dataServer := httptest.NewServer(http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			_ = request.ParseForm()

			var value string

			switch request.Method {
			case http.MethodGet:
				value = request.FormValue("foo")
			case http.MethodPost:
				var body map[string]string
				decoder := json.NewDecoder(request.Body)
				err := decoder.Decode(&body)
				if err != nil {
					t.Error(err)
				}
				value = body["foo"]
			}
			assert.Equal(t, "bar", value)
			d := map[string]interface{}{
				"response": "ok",
			}
			encoder := json.NewEncoder(writer)
			err := encoder.Encode(&d)
			if err != nil {
				t.Error(err)
			}
		}))
	defer dataServer.Close()

	methods := []string{http.MethodGet, http.MethodPost}
	for _, tcm := range methods {
		t.Run(tcm, func(t *testing.T) {
			var resp *http.Response
			var err error
			switch tcm {
			case http.MethodGet:
				data := url.Values{}
				data.Set("foo", "bar")
				resp, err = testClient.SendRequest(tcm, dataServer.URL, data, nil, nil)
			case http.MethodPost:
				body := map[string]string{"foo": "bar"}
				resp, err = testClient.SendRequest(tcm, dataServer.URL, nil, body, nil)
			}
			assert.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)
		})
	}
}
