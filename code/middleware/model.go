package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type nonprofit struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	SubCategory string `json:"subcategory"`
}

type ScoreStat struct {
	Predictions []struct {
		Classes []string  `json:"classes"`
		Scores  []float64 `json:"scores"`
	} `json:"predictions"`
	DeployedModelID  string `json:"deployedModelId"`
	Model            string `json:"model"`
	ModelDisplayName string `json:"modelDisplayName"`
	ModelVersionID   string `json:"modelVersionId"`
}

// JSON marshalls the content of a non-profit to json.
func (t nonprofit) JSON() (string, error) {
	bytes, err := json.Marshal(t)
	if err != nil {
		return "", fmt.Errorf("could not marshal json for response: %s", err)
	}

	return string(bytes), nil
}

// JSONBytes marshalls the content of a elara to json as a byte array.
func (t nonprofit) JSONBytes() ([]byte, error) {
	bytes, err := json.Marshal(t)
	if err != nil {
		return []byte{}, fmt.Errorf("could not marshal json for response: %s", err)
	}

	return bytes, nil
}

// Key returns the id as a string.
func (t nonprofit) Key() string {
	return strconv.Itoa(t.ID)
}

type nonprofits []nonprofit

// JSON marshalls the content of a slice of nonprofit to json.
func (t nonprofits) JSON() (string, error) {
	bytes, err := json.Marshal(t)
	if err != nil {
		return "", fmt.Errorf("could not marshal json for response: %s", err)
	}

	return string(bytes), nil
}

// JSONBytes marshalls the content of a slice of nonprofit to json as a byte array.
func (t nonprofits) JSONBytes() ([]byte, error) {
	bytes, err := json.Marshal(t)
	if err != nil {
		return []byte{}, fmt.Errorf("could not marshal json for response: %s", err)
	}

	return bytes, nil
}
