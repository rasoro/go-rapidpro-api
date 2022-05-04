package flows

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rasoro/rapidpro-api-go/client"
	"github.com/stretchr/testify/assert"
)

type FlowsTestCase struct {
	Label        string
	Method       string
	QueryParams  *QueryParams
	Status       int
	ResponseBody string
	ResultsCount int
	Error        error
}

var after = time.Date(2016, time.January, 6, 15, 0, 0, 0, time.UTC)
var before = time.Date(2016, time.January, 10, 0, 0, 0, 0, time.UTC)
var uuid = "8ecb3849-226b-431f-a790-37f879418a6b"

var testCases = []FlowsTestCase{
	{
		Label:        "Test Get flows",
		Method:       "GET",
		QueryParams:  nil,
		Status:       200,
		ResponseBody: testData,
		ResultsCount: 2,
		Error:        nil,
	},
	{
		Label:  "Test Get flows with params after and before",
		Method: "GET",
		QueryParams: &QueryParams{
			After:  &after,
			Before: &before,
		},
		Status:       200,
		ResponseBody: testData,
		ResultsCount: 2,
		Error:        nil,
	},
	{
		Label:        "Test Get flows with param uuid",
		Method:       "GET",
		QueryParams:  &QueryParams{UUID: uuid},
		Status:       200,
		ResponseBody: testData2,
		ResultsCount: 0,
		Error:        nil,
	},
}

func TestFlows(t *testing.T) {
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

var testData = `
{
	"next": null,
	"previous": null,
	"results": [
		{
			"uuid": "5f05311e-8f81-4a67-a5b5-1501b6d6496a",
			"name": "Survey1",
			"type": "message",
			"archived": false,
			"labels": [{"name": "Important", "uuid": "5a4eb79e-1b1f-4ae3-8700-09384cca385f"}],
			"expires": 600,
			"runs": {
				"active": 47,
				"completed": 123,
				"interrupted": 2,
				"expired": 34
			},
			"results": [
				{
					"key": "has_water",
					"name": "Has Water",
					"categories": ["Yes", "No", "Other"],
					"node_uuids": ["99afcda7-f928-4d4a-ae83-c90c96deb76d"]
				}
			],
			"parent_refs": [],
			"created_on": "2016-01-06T15:33:00.813162Z",
			"modified_on": "2017-01-07T13:14:00.453567Z"
		},
		{
			"uuid": "9d9dba87-6e91-4e08-85db-fabeadffac02",
			"name": "Survey2",
			"type": "message",
			"archived": false,
			"labels": [{"name": "Important", "uuid": "728fa0b1-3dd0-47a8-b6ef-17b071b5a280"}],
			"expires": 600,
			"runs": {
				"active": 47,
				"completed": 123,
				"interrupted": 2,
				"expired": 34
			},
			"results": [
				{
					"key": "has_water",
					"name": "Has Water",
					"categories": ["Yes", "No", "Other"],
					"node_uuids": ["ec6f2bde-50fa-4589-a93b-c8c9eba93c58"]
				}
			],
			"parent_refs": [],
			"created_on": "2016-01-06T15:33:00.813162Z",
			"modified_on": "2017-01-07T13:14:00.453567Z"
		}
	]
}
`
var testData2 = `
{
	"next": null,
	"previous": null,
	"results": []
}
`
