package flows

import (
	"encoding/json"
	"net/url"
	"time"

	rapidpro "github.com/rasoro/rapidpro-api-go/client"
)

const PATH = "/v2/flows.json"

type ApiService struct {
	serviceURL     string
	requestHandler *rapidpro.RequestHandler
}

func NewService(requestHandler *rapidpro.RequestHandler, apiURL string) *ApiService {
	return &ApiService{
		requestHandler: requestHandler,
		serviceURL:     apiURL + PATH,
	}
}

// Get makes a GET request to flows endpoint with *QueryParams and returns a Response
func (s *ApiService) Get(params *QueryParams) (*Response, error) {
	data := url.Values{}
	headers := make(map[string]interface{})

	if params != nil {
		if params.UUID != "" {
			data.Set("uuid", params.UUID)
		}
		if params.After != nil {
			data.Set("after", params.After.Format(time.RFC3339))
		}
		if params.Before != nil {
			data.Set("before", params.Before.Format(time.RFC3339))
		}
	}

	resp, err := s.requestHandler.Get(s.serviceURL, data, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	flowResponse := &Response{}
	if err := json.NewDecoder(resp.Body).Decode(flowResponse); err != nil {
		return nil, err
	}
	return flowResponse, err
}

// Flow is a struct that represents a flow object
type Flow struct {
	UUID     string        `json:"uuid"`
	Name     string        `json:"name"`
	Type     string        `json:"type"`
	Archived bool          `json:"archived"`
	Labels   []interface{} `json:"labels"`
	Expires  int           `json:"expires"`
	Runs     struct {
		Active      int `json:"active"`
		Completed   int `json:"completed"`
		Interrupted int `json:"interrupted"`
		Expired     int `json:"expired"`
	} `json:"runs"`
	Results []struct {
		Key        string   `json:"key"`
		Name       string   `json:"name"`
		Categories []string `json:"categories"`
		NodeUUIDS  []string `json:"node_uuids"`
	} `json:"results"`
	ParentRefs []interface{} `json:"parent_refs"`
	CreatedOn  time.Time     `json:"created_on"`
	ModifiedOn time.Time     `json:"modified_on"`
}

// Response is a struct that represents the response of a request in flows endpoint
type Response struct {
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []Flow      `json:"results"`
}

// QueryParams is a struct that represents the query parameters that can be passed in a request to flows endpoint
type QueryParams struct {
	UUID   string     `json:"uuid"`
	After  *time.Time `json:"after"`
	Before *time.Time `json:"before"`
}
