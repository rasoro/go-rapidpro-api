package flowstarts

import (
	"encoding/json"
	"net/url"
	"time"

	rapidpro "github.com/rasoro/rapidpro-api-go/client"
)

const PATH = "/v2/flow_starts.json"

type ApiService struct {
	URL            string
	requestHandler *rapidpro.RequestHandler
}

func NewService(requestHandler *rapidpro.RequestHandler, apiURL string) *ApiService {
	return &ApiService{
		requestHandler: requestHandler,
		URL:            apiURL + PATH,
	}
}

func (s *ApiService) Get(params *QueryParams) (*Response, error) {
	data := url.Values{}
	headers := make(map[string]interface{})

	if params != nil {
		if params.ID != "" {
			data.Set("id", params.ID)
		}
		if params.After != nil {
			data.Set("after", params.After.Format(time.RFC3339))
		}
		if params.Before != nil {
			data.Set("before", params.Before.Format(time.RFC3339))
		}
	}

	resp, err := s.requestHandler.Get(s.URL, data, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := &Response{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, err
	}
	return response, nil
}

func (s *ApiService) Post(body PostBody) (*FlowStart, error) {
	queryParams := url.Values{}
	resp, err := s.requestHandler.Post(s.URL, queryParams, body, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	response := &FlowStart{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, err
	}
	return response, nil
}

type FlowStart struct {
	UUID string `json:"uuid,omitempty"`
	Flow struct {
		UUID string `json:"uuid,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"flow,omitempty"`
	Groups []struct {
		UUID string `json:"uuid,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"groups,omitempty"`
	Contacts []struct {
		UUID string `json:"uuid,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"contacts,omitempty"`
	RestartParticipants bool   `json:"restart_participants,omitempty"`
	ExcludeActive       bool   `json:"exclude_active,omitempty"`
	Status              string `json:"status,omitempty"`
	Params              struct {
		FirstName string `json:"first_name,omitempty"`
		LastName  string `json:"last_name,omitempty"`
	} `json:"params,omitempty"`
	CreatedOn  *time.Time `json:"created_on,omitempty"`
	ModifiedOn *time.Time `json:"modified_on,omitempty"`
}

type Response struct {
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []FlowStart `json:"results"`
}

type QueryParams struct {
	ID     string     `json:"id,omitempty"`
	After  *time.Time `json:"after,omitempty"`
	Before *time.Time `json:"before,omitempty"`
}

type PostBody struct {
	Flow     string             `json:"flow,omitempty"`
	Groups   []string           `json:"groups,omitempty"`
	Contacts []string           `json:"contacts,omitempty"`
	URNs     []string           `json:"urns,omitempty"`
	Params   *map[string]string `json:"params,omitempty"`
}
