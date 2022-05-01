package rapidpro

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientWithParams(t *testing.T) {
	client := NewRestClientWithParams(ClientParams{
		Token: "token123",
	})
	assert.Equal(t, client.RequestHandler.Client.Token(), "token123")
}

func TestClientNewRestClient(t *testing.T) {
	os.Setenv("RAPIDPRO_API_GO_TOKEN", "token123")
	client := NewRestClient()
	assert.Equal(t, client.RequestHandler.Client.Token(), "token123")
}
