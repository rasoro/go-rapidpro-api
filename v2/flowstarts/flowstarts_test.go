package flowstarts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rasoro/rapidpro-api-go/client"
	"github.com/stretchr/testify/assert"
)

type FlowStartsTestCase struct {
	Label        string
	Method       string
	QueryParams  *QueryParams
	Status       int
	Headers      map[string]string
	PostBody     *PostBody
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
		ResponseBody: testDataGet,
		ResultsCount: 1,
		Error:        nil,
	},
	{
		Label:        "Test Get flow starts with error",
		Method:       "GET",
		QueryParams:  &QueryParams{ID: "asd"},
		Status:       400,
		ResponseBody: `{"detail": "Value for id must be an integer"}`,
		ResultsCount: 0,
		Error: &client.RapidproRestError{
			Status:  400,
			Details: map[string]interface{}{"detail": "Value for id must be an integer"},
		},
	},
	{
		Label:        "Test Get flow starts with QueryParam ID",
		Method:       "GET",
		QueryParams:  &QueryParams{ID: "1"},
		Status:       200,
		ResponseBody: testDataGetWithParams,
		ResultsCount: 0,
		Error:        nil,
	},
	{
		Label:  "Test Get flow starts with QueryParam After & Before",
		Method: "GET",
		QueryParams: &QueryParams{
			After:  timeToPointer(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
			Before: timeToPointer(time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)),
		},
		Status:       200,
		ResponseBody: testDataGetWithParams,
		ResultsCount: 0,
		Error:        nil,
	},
	{
		Label:       "Test Post flow starts",
		Method:      "POST",
		QueryParams: nil,
		Status:      200,
		PostBody: &PostBody{
			Flow: "d6efc9ff-cf7d-4a5c-b4b3-46eda997d461",
			URNs: []string{"telegram:938623661"},
		},
		ResponseBody: testDataPost,
		ResultsCount: 1,
		Error:        nil,
	},
	{
		Label:       "Test Post flow starts with error",
		Method:      "POST",
		QueryParams: nil,
		Status:      400,
		PostBody: &PostBody{
			Flow: "f5901b62-ba76-4003-9c62-72fdacc1b7b7",
			URNs: []string{"telegram:938623661"},
		},
		ResponseBody: testDataPostWithError,
		Error: &client.RapidproRestError{
			Status:  400,
			Details: map[string]interface{}{"flow": []interface{}{"No such object: f5901b62-ba76-4003-9c62-72fdacc1b7b7"}},
		},
	},
}

func Sp(str interface{}) *string { asStr := fmt.Sprintf("%s", str); return &asStr }

func timeToPointer(t time.Time) *time.Time { return &t }

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
		if tc.Method == "GET" {
			resp, err := service.Get(tc.QueryParams)
			assert.Equal(t, tc.Error, err)
			if err == nil {
				assert.Equal(t, tc.ResultsCount, len(resp.Results))
			}
		} else {
			flowStartResponse, err := service.Post(*tc.PostBody)
			assert.Equal(t, tc.Error, err)

			if err == nil {
				var responseBodyFlowStart FlowStart
				err = json.Unmarshal([]byte(tc.ResponseBody), &responseBodyFlowStart)
				assert.NoError(t, err)
				assert.Equal(t, flowStartResponse.UUID, responseBodyFlowStart.UUID)
			}
		}
	}
}

var (
	testDataGet = `
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
	testDataGetWithParams = `
	{
    "next": null,
    "previous": null,
    "results": []
	}
	`
)

var (
	testDataPost = `
	{
		"id": 92776060,
		"uuid": "6846356a-b25b-4a2c-b999-b55c53880dd0",
		"flow": {
			"uuid": "d6efc9ff-cf7d-4a5c-b4b3-46eda997d461",
			"name": "dummy flow"
		},
		"status": "pending",
		"groups": [],
		"contacts": [],
		"restart_participants": true,
		"exclude_active": false,
		"extra": null,
		"params": null,
		"created_on": "2022-05-29T15:08:21.238970Z",
		"modified_on": "2022-05-29T15:08:21.238978Z"
	}`
	testDataPostWithError = `
	{
    "flow": [
        "No such object: f5901b62-ba76-4003-9c62-72fdacc1b7b7"
    ]
	}
	`
)
