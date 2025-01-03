package provider

import (
	"encoding/json"
)

type Widget interface {
	ToJSON() (string, error)
}

type TextWidget struct {
	Text string `json:"text"`
}

func (tw TextWidget) ToJSON() (string, error) {
	jsonData, err := json.Marshal(tw)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
