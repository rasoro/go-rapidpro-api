package contacts

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"

	rapidpro "github.com/rasoro/rapidpro-api-go/client"
)

const PATH = "/v2/contacts.json"

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

// Get makes a GET request to messages endpoint with *QueryParams and returns a Response.
func (s *ApiService) Get(params *QueryParams) (*Response, error) {
	data := url.Values{}

	if params != nil {
		if params.UUID != "" {
			data.Set("uuid", params.UUID)
		}
		if params.URN != "" {
			data.Set("urn", params.URN)
		}
		if params.Group != "" {
			data.Set("group", params.Group)
		}
		if params.Deleted != nil {
			data.Set("deleted", strconv.FormatBool(*params.Deleted))
		}
		if params.After != nil {
			data.Set("after", params.After.Format(time.RFC3339))
		}
		if params.Before != nil {
			data.Set("before", params.Before.Format(time.RFC3339))
		}
	}

	resp, err := s.requestHandler.Get(s.URL, data, nil)
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

// Post to add a new contact to your workspace passing PostBody and return a *Contact for the new created contact.
func (s *ApiService) Post(body PostBody, bodyPostBody string) (*Contact, error) {
	resp, err := s.requestHandler.Post(s.URL, url.Values{}, body, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	response := &Contact{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, err
	}
	return response, nil
}

// Contact represents a contact object.
type Contact struct {
	UUID       string                 `json:"uuid,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Language   interface{}            `json:"language,omitempty"`
	Urns       []string               `json:"urns,omitempty"`
	Groups     []interface{}          `json:"groups,omitempty"`
	Fields     map[string]interface{} `json:"fields,omitempty"`
	Blocked    bool                   `json:"blocked,omitempty"`
	Stopped    bool                   `json:"stopped,omitempty"`
	CreatedOn  *time.Time             `json:"created_on,omitempty"`
	ModifiedOn *time.Time             `json:"modified_on,omitempty"`
	LastSeenOn *time.Time             `json:"last_seen_on,omitempty"`
}

// Response represents the response of a request to contacts endpoint containing the contacts from an organization.
type Response struct {
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []Contact   `json:"results"`
}

// QueryParams represents the query parameters that can be passed in a GET request to contacts endpoint.
type QueryParams struct {
	UUID    string     `json:"uuid,omitempty"`
	URN     string     `json:"urn,omitempty"`
	Group   string     `json:"group,omitempty"`
	Deleted *bool      `json:"deleted,omitempty"`
	After   *time.Time `json:"after,omitempty"`
	Before  *time.Time `json:"before,omitempty"`
}

// PostBody represents the new contact to be created on a post request for contacts.
type PostBody struct {
	Name     string                 `json:"name,omitempty"`
	Language string                 `json:"language,omitempty"`
	URNS     []string               `json:"urns,omitempty"`
	Groups   []string               `json:"groups,omitempty"`
	Fields   map[string]interface{} `json:"fields,omitempty"`
}
