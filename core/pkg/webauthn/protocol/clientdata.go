package protocol

import (
	"encoding/json"
	"log"
)

func ParseClientData(clientData string) (*CollectedClientData, error) {
	var cd CollectedClientData
	err := json.Unmarshal([]byte(clientData), &cd)
	if err != nil {
		log.Printf("failed to unmarshal client data, %v", err)
		return nil, err
	}

	return &cd, nil
}
