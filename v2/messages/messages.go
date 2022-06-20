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

func NewSErvice(requestHandler *rapidpro.RequestHandler, apiURL string) *ApiService {
	return &ApiService{
		requestHandler: requestHandler,
		serviceURL:     apiURL + PATH,
	}
}

// Get makes a GET request to messages endpoint with *QueryParams and returns a Response
func (s *ApiService) Get(params *QueryParams) (*Response, error) {
	data := url.Values{}
	headers := make(map[string]interface{})

	if params != nil {
		if params.ID != 0 {
			data.Set("id", strconv.FormatInt(params.ID, 10))
		}
		if params.Broadcast != "" {
			data.Set("broadcast", params.Broadcast)
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

type Message struct {
	ID        int64 `json:"id"`
	Broadcast int   `json:"broadcast"`
	Contact   struct {
		UUID string `json:"uuid"`
		Name string `json:"name"`
	} `json:"contact"`
	Urn     string `json:"urn"`
	Channel struct {
		UUID string `json:"uuid"`
		Name string `json:"name"`
	} `json:"channel"`
	Direction   string `json:"direction"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Visibility  string `json:"visibility"`
	Text        string `json:"text"`
	Attachments []struct {
		ContentType string `json:"content_type"`
		URL         string `json:"url"`
	} `json:"attachments"`
	Labels []struct {
		Name string `json:"name"`
		UUID string `json:"uuid"`
	} `json:"labels"`
	CreatedOn  time.Time `json:"created_on"`
	SentOn     time.Time `json:"sent_on"`
	ModifiedOn time.Time `json:"modified_on"`
}

type Response struct {
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []Message   `json:"message"`
}

// QueryParams is a struct that represents the query parameters that can be passed in a request to messages endpoint
type QueryParams struct {
	ID        int64      `json:"id"`
	Broadcast string     `json:"broadcast"`
	Contact   string     `json:"contact"`
	Folder    string     `json:"folder"`
	Label     string     `json:"label"`
	Before    *time.Time `json:"before"`
	After     *time.Time `json:"after"`
}
