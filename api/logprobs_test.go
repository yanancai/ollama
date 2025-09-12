package api

import (
	"encoding/json"
	"testing"
)

func TestTokenLogprobSerialization(t *testing.T) {
	// Test basic logprobs data structure
	logprobs := []TokenLogprob{
		{
			Token:   "hello",
			Logprob: -0.5108256,
			Bytes:   []int{104, 101, 108, 108, 111},
			TopLogprobs: []TopTokenLogprob{
				{Token: "hello", Logprob: -0.5108256, Bytes: []int{104, 101, 108, 108, 111}},
				{Token: "hi", Logprob: -1.2039728, Bytes: []int{104, 105}},
			},
		},
	}

	// Test marshaling
	data, err := json.Marshal(logprobs)
	if err != nil {
		t.Fatalf("Failed to marshal TokenLogprob: %v", err)
	}

	// Test unmarshaling
	var result []TokenLogprob
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal TokenLogprob: %v", err)
	}

	// Verify the structure
	if len(result) != 1 {
		t.Errorf("Expected 1 token, got %d", len(result))
	}

	token := result[0]
	if token.Token != "hello" {
		t.Errorf("Expected token 'hello', got '%s'", token.Token)
	}

	if token.Logprob != -0.5108256 {
		t.Errorf("Expected logprob -0.5108256, got %f", token.Logprob)
	}

	if len(token.TopLogprobs) != 2 {
		t.Errorf("Expected 2 top logprobs, got %d", len(token.TopLogprobs))
	}
}

func TestChatRequestOptions(t *testing.T) {
	// Test that ChatRequest correctly handles logprobs through options
	jsonData := `{
		"model": "test-model",
		"messages": [{"role": "user", "content": "Hello"}],
		"options": {
			"logprobs": true,
			"top_logprobs": 5
		}
	}`

	var req ChatRequest
	err := json.Unmarshal([]byte(jsonData), &req)
	if err != nil {
		t.Fatalf("Failed to unmarshal ChatRequest: %v", err)
	}

	if req.Options["logprobs"] != true {
		t.Error("Expected logprobs to be true in options")
	}

	topLogprobs, ok := req.Options["top_logprobs"].(float64) // JSON numbers are float64
	if !ok || int(topLogprobs) != 5 {
		t.Error("Expected top_logprobs to be 5 in options")
	}
}

func TestMessageWithLogprobs(t *testing.T) {
	// Test that Message correctly handles logprobs field
	message := Message{
		Role:    "assistant",
		Content: "Hello world!",
		Logprobs: []TokenLogprob{
			{Token: "Hello", Logprob: -0.1, Bytes: []int{72, 101, 108, 108, 111}},
			{Token: " world", Logprob: -0.2, Bytes: []int{32, 119, 111, 114, 108, 100}},
			{Token: "!", Logprob: -0.3, Bytes: []int{33}},
		},
	}

	// Test marshaling
	data, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("Failed to marshal Message: %v", err)
	}

	// Test unmarshaling
	var result Message
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal Message: %v", err)
	}

	if len(result.Logprobs) == 0 {
		t.Error("Expected logprobs to be present")
	}

	if len(result.Logprobs) != 3 {
		t.Errorf("Expected 3 tokens, got %d", len(result.Logprobs))
	}
}
