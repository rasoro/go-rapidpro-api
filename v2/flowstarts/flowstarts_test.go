package flowstarts

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rasoro/rapidpro-api-go/client"
	"github.com/stretchr/testify/assert"
)

type FlowStartsTestCase struct {
	Label        string
	Method       string
	QueryParams  *QueryParams
	Status       int
	ResponseBody string
	ResultsCount int
	Error        error
}

var testCases = []FlowStartsTestCase{
	{
		Label:        "Test Get flow starts",
		Method:       "GET",
		QueryParams:  nil,
		Status:       200,
		ResponseBody: testData,
		ResultsCount: 1,
		Error:        nil,
	},
}

func TestFlowStarts(t *testing.T) {
	for _, tc := range testCases {
		mockServer := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.Status)
				_, _ = w.Write([]byte(tc.ResponseBody))
			}))
		defer mockServer.Close()

		defaultClient := &client.Client{
			Credentials: &client.Credentials{Token: "token123"},
		}

		requestHandler := client.NewRequestHandler(defaultClient)
		service := NewService(requestHandler, mockServer.URL)
		resp, err := service.Get(tc.QueryParams)
		assert.Equal(t, tc.Error, err)
		assert.Equal(t, tc.ResultsCount, len(resp.Results))
	}
}

var (
	testData = `
	{
    "next": "http://example.com/api/v2/flow_starts.json?cursor=cD0yMDE1LTExLTExKzExJTNBM40NjQlMkIwMCUzRv",
    "previous": null,
    "results": [
			{
				"uuid": "09d23a05-47fe-11e4-bfe9-b8f6b119e9ab",
				"flow": {"uuid": "f5901b62-ba76-4003-9c62-72fdacc1b7b7", "name": "Thrift Shop"},
				"groups": [
							{"uuid": "f5901b62-ba76-4003-9c62-72fdacc1b7b7", "name": "Ryan & Macklemore"}
				],
				"contacts": [
							{"uuid": "f5901b62-ba76-4003-9c62-fjjajdsi15553", "name": "Wanz"}
				],
				"restart_participants": true,
				"exclude_active": false,
				"status": "complete",
				"params": {
						"first_name": "Ryan",
						"last_name": "Lewis"
				},
				"created_on": "2013-08-19T19:11:21.082Z",
				"modified_on": "2013-08-19T19:11:21.082Z"
			}
    ]
	}
	`
)
