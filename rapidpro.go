package rapidpro

import (
	"log"

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
	// TODO: implement
	log.Fatal("Not implemented")
	return nil
}
