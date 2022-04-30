package client

import (
	"encoding/json"
	"fmt"
)

type RapidproRestError struct {
	Status  int
	Details map[string]interface{}
}

func (e *RapidproRestError) Error() string {
	detailsJSON, _ := json.Marshal(e.Details)
	return fmt.Sprintf("Status: %d - Error: %s", e.Status, detailsJSON)
}
