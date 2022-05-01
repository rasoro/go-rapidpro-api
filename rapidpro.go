package rapidpro

import (
	"os"

	"github.com/rasoro/rapidpro-api-go/client"
)

type RestClient struct {
	*client.RequestHandler
}

type ClientParams struct {
	Client client.BaseClient
	Token  string
}

func NewRestClient() *RestClient {
	return NewRestClientWithParams(ClientParams{})
}

func NewRestClientWithParams(params ClientParams) *RestClient {
	requestHandler := client.NewRequestHandler(params.Client)

	if params.Client == nil {
		token := params.Token
		if token == "" {
			token = os.Getenv("RAPIDPRO_API_GO_TOKEN")
		}
		defaultClient := &client.Client{
			Credentials: &client.Credentials{Token: token},
		}
		requestHandler = client.NewRequestHandler(defaultClient)
	}
	c := &RestClient{
		RequestHandler: requestHandler,
	}
	return c
}
