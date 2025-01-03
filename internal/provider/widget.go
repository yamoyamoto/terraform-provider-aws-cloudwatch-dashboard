package provider

import (
	"encoding/json"
)

type Widget interface {
	ToJSON() (string, error)
}

type TextWidget struct {
	Content string
}

func (tw TextWidget) ToJSON() (string, error) {
	widget := map[string]interface{}{
		"type": "text",
		"properties": map[string]interface{}{
			"content": tw.Content,
		},
	}
	jsonData, err := json.Marshal(widget)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
