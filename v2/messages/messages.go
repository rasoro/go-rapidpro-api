package messages

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"

	rapidpro "github.com/rasoro/rapidpro-api-go/client"
)

const PATH = "/v2/messages.json"

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

// Get makes a GET request to messages endpoint with *QueryParams and returns a Response.
func (s *ApiService) Get(params *QueryParams) (*Response, error) {
	data := url.Values{}
	headers := make(map[string]interface{})

	if params != nil {
		if params.ID != 0 {
			data.Set("id", strconv.Itoa(params.ID))
		}
		if params.Broadcast != 0 {
			data.Set("broadcast", strconv.Itoa(params.Broadcast))
		}
		if params.Contact != "" {
			data.Set("contact", params.Contact)
		}
		if params.Folder != "" {
			data.Set("folder", params.Folder)
		}
		if params.Label != "" {
			data.Set("label", params.Label)
		}
		if params.Before != nil {
			data.Set("before", params.Before.Format(time.RFC3339))
		}
		if params.After != nil {
			data.Set("after", params.After.Format(time.RFC3339))
		}
	}

	resp, err := s.requestHandler.Get(s.serviceURL, data, headers)
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

// Message represents a message objects
type Message struct {
	ID        int `json:"id,omitempty"`
	Broadcast int `json:"broadcast,omitempty"`
	Contact   struct {
		UUID string `json:"uuid,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"contact,omitempty"`
	Urn     string `json:"urn,omitempty"`
	Channel struct {
		UUID string `json:"uuid,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"channel,omitempty"`
	Direction   string `json:"direction,omitempty"`
	Type        string `json:"type,omitempty"`
	Status      string `json:"status,omitempty"`
	Visibility  string `json:"visibility,omitempty"`
	Text        string `json:"text,omitempty"`
	Attachments []struct {
		ContentType string `json:"content_type,omitempty"`
		URL         string `json:"url,omitempty"`
	} `json:"attachments,omitempty"`
	Labels []struct {
		Name string `json:"name,omitempty"`
		UUID string `json:"uuid,omitempty"`
	} `json:"labels,omitempty"`
	CreatedOn  *time.Time `json:"created_on,omitempty"`
	SentOn     *time.Time `json:"sent_on,omitempty"`
	ModifiedOn *time.Time `json:"modified_on,omitempty"`
}

// Response represents the response of a request in messages endpoint
type Response struct {
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []Message   `json:"results"`
}

// QueryParams is a struct that represents the query parameters that can be passed in a request to messages endpoint
type QueryParams struct {
	ID        int        `json:"id,omitempty"`
	Broadcast int        `json:"broadcast,omitempty"`
	Contact   string     `json:"contact,omitempty"`
	Folder    string     `json:"folder,omitempty"`
	Label     string     `json:"label,omitempty"`
	Before    *time.Time `json:"before,omitempty"`
	After     *time.Time `json:"after,omitempty"`
}
