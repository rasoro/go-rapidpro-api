package flows

import (
	"encoding/json"
	"net/url"
	"time"

	rapidpro "github.com/rasoro/rapidpro-api-go/client"
)

const SOURCE_PATH = "/v2/flows.json"

type ApiService struct {
	serviceURL     string
	requestHandler *rapidpro.RequestHandler
}

func NewService(requestHandler *rapidpro.RequestHandler, apiURL string) *ApiService {
	return &ApiService{
		requestHandler: requestHandler,
		serviceURL:     apiURL + SOURCE_PATH,
	}
}

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

type Response struct {
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []Flow      `json:"results"`
}

type QueryParams struct {
	UUID   string     `json:"uuid"`
	After  *time.Time `json:"after"`
	Before *time.Time `json:"before"`
}
