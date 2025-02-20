package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func getAITaskSuggestions(description string) (string, error) {
	apiKey := "your-openai-api-key"
	url := "https://api.openai.com/v1/completions"

	payload := map[string]interface{}{
		"model":      "text-davinci-003",
		"prompt":     "Generate 3 subtasks for: " + description,
		"max_tokens": 100,
	}

	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result["choices"].([]interface{})[0].(map[string]interface{})["text"].(string), nil
}
