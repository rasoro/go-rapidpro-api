package messages

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/rasoro/rapidpro-api-go/client"
	"github.com/stretchr/testify/assert"
)

type MessagesTestCase struct {
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
var uuid = "7f2df3ad-e419-499a-8d6a-a6493ea6f163"

var testCases = []MessagesTestCase{
	{
		Label:        "Test Get messages",
		Method:       "GET",
		QueryParams:  nil,
		Status:       200,
		ResponseBody: testData[0],
		ResultsCount: 2,
		Error:        nil,
	},
	{
		Label:        "Test Get messages with status 500 error",
		Method:       "GET",
		QueryParams:  nil,
		Status:       500,
		ResponseBody: "{}",
		ResultsCount: 0,
		Error: &client.RapidproRestError{
			Status:  500,
			Details: map[string]interface{}{},
		},
	},
	{
		Label:        "Test Get messages with response error",
		Method:       "GET",
		Status:       200,
		ResponseBody: "{",
		Error:        errors.New("unexpected EOF"),
	},
	{
		Label:  "Test Get messages with params",
		Method: "GET",
		QueryParams: &QueryParams{
			ID:        123,
			Broadcast: 123,
			Contact:   "f8c17c2d-18e5-432b-938d-1144b56a3d32",
			Folder:    "inbox",
			Label:     "e74c6be1-dcfb-42ba-96d5-692d3d19b63d",
			After:     timeToPointer(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
			Before:    timeToPointer(time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)),
		},
		Status:       200,
		ResponseBody: testData[1],
		ResultsCount: 0,
		Error:        nil,
	},
}

func timeToPointer(t time.Time) *time.Time { return &t }

func TestMessages(t *testing.T) {
	for _, tc := range testCases {
		mockServer := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(tc.Status)
				w.Write([]byte(tc.ResponseBody))
			}))
		defer mockServer.Close()

		t.Run(tc.Label, func(t *testing.T) {
			defaultClient := &client.Client{
				Credentials: &client.Credentials{Token: "token123"},
			}
			requestHandler := client.NewRequestHandler(defaultClient)
			service := NewService(requestHandler, mockServer.URL)
			resp, err := service.Get(tc.QueryParams)
			if err != nil {
				assert.Equal(t, tc.Error.Error(), err.Error())
			}
			if tc.ResultsCount > 0 {
				assert.Equal(t, tc.ResultsCount, len(resp.Results))
			}
		})
	}
}

var testData = []string{
	`{
		"next": null,
		"previous": null,
		"results": [
			{
					"id": 4105426,
					"broadcast": 2690007,
					"contact": {"uuid": "d33e9ad5-5c35-414c-abd4-e7451c69ff1d", "name": "Bob McFlow"},
					"urn": "twitter:textitin",
					"channel": {"uuid": "9a8b001e-a913-486c-80f4-1356e23f582e", "name": "Vonage"},
					"direction": "out",
					"type": "inbox",
					"status": "wired",
					"visibility": "visible",
					"text": "How are you?",
					"attachments": [{"content_type": "audio/wav", "url": "http://domain.com/recording.wav"}],
					"labels": [{"name": "Important", "uuid": "5a4eb79e-1b1f-4ae3-8700-09384cca385f"}],
					"created_on": "2016-01-06T15:33:00.813162Z",
					"sent_on": "2016-01-06T15:35:03.675716Z",
					"modified_on": "2016-01-06T15:35:03.675716Z"
			},
			{
					"id": 5216537,
					"broadcast": 2690008,
					"contact": {"uuid": "d33e9ad5-5c35-414c-abd4-e7451c69ff1d", "name": "Bob McFlow"},
					"urn": "twitter:textitin",
					"channel": {"uuid": "9a8b001e-a913-486c-80f4-1356e23f582e", "name": "Vonage"},
					"direction": "out",
					"type": "inbox",
					"status": "wired",
					"visibility": "visible",
					"text": "How are you?",
					"attachments": [{"content_type": "audio/wav", "url": "http://domain.com/recording.wav"}],
					"labels": [{"name": "Important", "uuid": "5a4eb79e-1b1f-4ae3-8700-09384cca385f"}],
					"created_on": "2016-01-06T15:33:00.813162Z",
					"sent_on": "2016-01-06T15:35:03.675716Z",
					"modified_on": "2016-01-06T15:35:03.675716Z"
			}
		]
	}`,
	`{
		"next": null,
		"previous": null,
		"results": []
	}`,
}
