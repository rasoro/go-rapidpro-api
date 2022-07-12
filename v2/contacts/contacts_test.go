package contacts

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rasoro/rapidpro-api-go/client"
	"github.com/stretchr/testify/assert"
)

type ContactsTestCase struct {
	Label        string
	Method       string
	Params       *QueryParams
	Status       int
	PostBody     *PostBody
	ResponseBody string
	ResultsCount int
	Error        error
}

var testCases = []ContactsTestCase{
	{
		Label:        "Test Get contacts",
		Method:       "GET",
		Params:       nil,
		Status:       200,
		ResponseBody: testDataGet[0],
		ResultsCount: 1,
		Error:        nil,
	},
	{
		Label:  "Test Get contacts passing uuid param",
		Method: "GET",
		Params: &QueryParams{
			UUID: "21043846-e18f-476f-beb2-820b36ea27a1",
		},
		Status:       200,
		ResponseBody: testDataGet[1],
		ResultsCount: 0,
		Error:        nil,
	},
	{
		Label:  "Test Get contacts passing urn param",
		Method: "GET",
		Params: &QueryParams{
			URN: "tel:+250788123123",
		},
		Status:       200,
		ResponseBody: testDataGet[1],
		ResultsCount: 0,
		Error:        nil,
	},
	{
		Label:  "Test Get contacts passing group param",
		Method: "GET",
		Params: &QueryParams{
			Group: "508b9353-1f26-475a-ba91-e7ab1c0aa753",
		},
		Status:       200,
		ResponseBody: testDataGet[1],
		ResultsCount: 0,
		Error:        nil,
	},
	{
		Label:  "Test Get contacts passing params for deleted, after and before",
		Method: "GET",
		Params: &QueryParams{
			Deleted: Bp(true),
			After:   &after,
			Before:  &before,
		},
		Status:       200,
		ResponseBody: testDataGet[1],
		ResultsCount: 0,
		Error:        nil,
	},
	{
		Label:        "Test Get contacts with status 500 error",
		Method:       "GET",
		Params:       nil,
		Status:       500,
		ResponseBody: "{}",
		Error: &client.RapidproRestError{
			Status:  500,
			Details: map[string]interface{}{},
		},
	},
	{
		Label:        "Test Get contacts with response error",
		Method:       "GET",
		Params:       nil,
		Status:       200,
		ResponseBody: "{",
		Error:        errors.New("unexpected EOF"),
	},
	{
		Label:  "Test Post contact",
		Method: "POST",
		Params: nil,
		Status: 200,
		PostBody: &PostBody{
			Name:     "Dummy Contact",
			URNs:     []string{"tel:+250788123123", "twitter:dummy"},
			Language: "eng",
			Groups:   []string{"6685e933-26e1-4363-a468-8f7268ab63a9"},
			Fields: map[string]interface{}{
				"nickname":  "Dummy Person",
				"side_kick": "Another Dummy",
			},
		},
		ResponseBody: testDataPost[0],
		ResultsCount: 1,
		Error:        nil,
	},
	{
		Label:  "Test Post contact with error 500",
		Method: "POST",
		Params: nil,
		Status: 500,
		PostBody: &PostBody{
			Name:     "Dummy Contact",
			URNs:     []string{"tel:+250788123123", "twitter:dummy"},
			Language: "eng",
			Groups:   []string{"6685e933-26e1-4363-a468-8f7268ab63a9"},
			Fields: map[string]interface{}{
				"nickname":  "Dummy Person",
				"side_kick": "Another Dummy",
			},
		},
		ResponseBody: "{}",
		Error: &client.RapidproRestError{
			Status:  500,
			Details: map[string]interface{}{},
		},
	},
	{
		Label:  "Test Post contact with response error",
		Method: "POST",
		Params: nil,
		Status: 200,
		PostBody: &PostBody{
			Name:     "Dummy Contact",
			URNs:     []string{"tel:+250788123123", "twitter:dummy"},
			Language: "eng",
			Groups:   []string{"6685e933-26e1-4363-a468-8f7268ab63a9"},
			Fields: map[string]interface{}{
				"nickname":  "Dummy Person",
				"side_kick": "Another Dummy",
			},
		},
		ResponseBody: "{",
		Error:        errors.New("unexpected EOF"),
	},
	{
		Label:  "Test Post contact to update querying by UUID",
		Method: "POST",
		Status: 200,
		Params: &QueryParams{
			UUID: "18d69114-ad53-45ac-a948-bd72288977d9",
		},
		PostBody: &PostBody{
			Fields: map[string]interface{}{
				"nickname": "The Dummy",
			},
		},
		ResponseBody: testDataPost[1],
	},
	{
		Label:  "Test Post contact to update querying by URN",
		Method: "POST",
		Status: 200,
		Params: &QueryParams{
			URN: "tel:+250788123123",
		},
		PostBody: &PostBody{
			Fields: map[string]interface{}{
				"nickname": "The Dummy",
			},
		},
		ResponseBody: testDataPost[1],
	},
}

var after = time.Date(2022, 06, 01, 0, 0, 0, 0, time.UTC)
var before = time.Date(2022, 06, 02, 0, 0, 0, 0, time.UTC)

func Sp(str interface{}) *string { asStr := fmt.Sprintf("%s", str); return &asStr }

func Bp(v bool) *bool { return &v }

func TestContacts(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Label, func(t *testing.T) {
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
			contacts := NewService(requestHandler, mockServer.URL)

			if tc.Method == "GET" {
				resp, err := contacts.Get(tc.Params)
				assert.Equal(t, tc.Error, err)
				if err == nil {
					assert.Equal(t, tc.ResultsCount, len(resp.Results))
				}

			} else {
				responseContact, err := contacts.Post(*tc.PostBody, tc.Params)
				assert.Equal(t, tc.Error, err)
				if err == nil {
					assert.NoError(t, err)
					if tc.PostBody.Name != "" {
						assert.Equal(t, tc.PostBody.Name, responseContact.Name)
					}
					if tc.PostBody.Language != "" {
						assert.Equal(t, tc.PostBody.Language, responseContact.Language)
					}
					if len(tc.PostBody.URNs) > 0 {
						assert.Equal(t, tc.PostBody.URNs, responseContact.URNs)
					}
					if len(tc.PostBody.Groups) > 0 {
						for _, pgroup := range tc.PostBody.Groups {
							hasGroup := false
							for _, rgroup := range responseContact.Groups {
								if pgroup == rgroup["uuid"] {
									hasGroup = true
								}
							}
							assert.True(t, hasGroup)
						}
					}
					if len(tc.PostBody.Fields) > 0 {
						assert.Equal(t, tc.PostBody.Fields, responseContact.Fields)
					}
				}
			}
		})
	}
}

var testDataGet = []string{
	`{
    "next": null,
    "previous": null,
    "results": [
    {
        "uuid": "09d23a05-47fe-11e4-bfe9-b8f6b119e9ab",
        "name": "Ben Haggerty",
        "language": null,
        "urns": ["tel:+250788123123"],
        "groups": [{"name": "Customers", "uuid": "5a4eb79e-1b1f-4ae3-8700-09384cca385f"}],
        "fields": {
          "nickname": "Macklemore",
          "side_kick": "Ryan Lewis"
        },
        "blocked": false,
        "stopped": false,
        "created_on": "2015-11-11T13:05:57.457742Z",
        "modified_on": "2020-08-11T13:05:57.576056Z",
        "last_seen_on": "2020-07-11T13:05:57.576056Z"
    }]
	}`,
	`{
		"next": null,
		"previous": null,
		"results": []
	}`,
}

var testDataPost = []string{
	`{
    "uuid": "09d23a05-47fe-11e4-bfe9-b8f6b119e9ab",
    "name": "Dummy Contact",
    "language": "eng",
    "urns": ["tel:+250788123123", "twitter:dummy"],
    "groups": [{"name": "Customers", "uuid": "5a4eb79e-1b1f-4ae3-8700-09384cca385f"},{"name":"dummy-group", "uuid": "6685e933-26e1-4363-a468-8f7268ab63a9"}],
    "fields": {
      "nickname": "Dummy Person",
      "side_kick": "Another Dummy"
    },
    "blocked": false,
    "stopped": false,
    "created_on": "2015-11-11T13:05:57.457742Z",
    "modified_on": "2015-11-11T13:05:57.576056Z",
    "last_seen_on": null
	}`,
	`{
    "uuid": "09d23a05-47fe-11e4-bfe9-b8f6b119e9ab",
    "name": "Dummy Contact",
    "language": "eng",
    "urns": ["tel:+250788123123", "twitter:dummy"],
    "groups": [{"name": "Customers", "uuid": "5a4eb79e-1b1f-4ae3-8700-09384cca385f"}],
    "fields": {
      "nickname": "The Dummy"
    },
    "blocked": false,
    "stopped": false,
    "created_on": "2015-11-11T13:05:57.457742Z",
    "modified_on": "2015-11-11T13:05:57.576056Z",
    "last_seen_on": null
	}`,
}
