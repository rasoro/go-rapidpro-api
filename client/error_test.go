package client_test

import (
	"testing"

	"github.com/rasoro/rapidpro-api-go/client"
	"github.com/stretchr/testify/assert"
)

const (
	errorMessage = "Bad request"
	errorStatus  = 400
)

func TestRapidproRestError_BadRequest(t *testing.T) {
	detail := make(map[string]interface{})
	detail["flow"] = []string{"this field is required."}
	err := &client.RapidproRestError{
		Status:  errorStatus,
		Details: detail,
	}
	expected := `Status: 400 - Error: {"flow":["this field is required."]}`
	assert.Equal(t, expected, err.Error())
}
