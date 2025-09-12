package openai

import (
	"encoding/json"
	"testing"

	"github.com/ollama/ollama/api"
)

func TestChatCompletionRequestLogprobs(t *testing.T) {
	// Test that ChatCompletionRequest correctly handles logprobs fields
	jsonData := `{
		"model": "test-model",
		"messages": [{"role": "user", "content": "Hello"}],
		"logprobs": true,
		"top_logprobs": 5
	}`

	var req ChatCompletionRequest
	err := json.Unmarshal([]byte(jsonData), &req)
	if err != nil {
		t.Fatalf("Failed to unmarshal ChatCompletionRequest: %v", err)
	}

	if req.Logprobs == nil || !*req.Logprobs {
		t.Error("Expected logprobs to be true")
	}

	if req.TopLogprobs == nil || *req.TopLogprobs != 5 {
		t.Error("Expected top_logprobs to be 5")
	}
}

func TestCompletionRequestLogprobs(t *testing.T) {
	// Test that CompletionRequest correctly handles logprobs fields
	// CompletionRequest expects `logprobs` as an integer (0-5)
	jsonData := `{
		"model": "test-model",
		"prompt": "Hello",
		"logprobs": 1,
		"top_logprobs": 10
	}`

	var req CompletionRequest
	err := json.Unmarshal([]byte(jsonData), &req)
	if err != nil {
		t.Fatalf("Failed to unmarshal CompletionRequest: %v", err)
	}

	if req.Logprobs == nil || *req.Logprobs != 1 {
		t.Error("Expected logprobs to be 1")
	}

	if req.TopLogprobs == nil || *req.TopLogprobs != 10 {
		t.Error("Expected top_logprobs to be 10")
	}
}

func TestChoiceWithLogprobs(t *testing.T) {
	// Test that Choice correctly handles logprobs field
	choice := Choice{
		Index:   0,
		Message: Message{Role: "assistant", Content: "Hello!"},
		Logprobs: &ChatLogprobs{
			Content: []api.TokenLogprob{
				{Token: "Hello", Logprob: -0.1, Bytes: []int{72, 101, 108, 108, 111}},
				{Token: "!", Logprob: -0.2, Bytes: []int{33}},
			},
		},
	}

	// Test marshaling
	data, err := json.Marshal(choice)
	if err != nil {
		t.Fatalf("Failed to marshal Choice: %v", err)
	}

	// Test unmarshaling
	var result Choice
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal Choice: %v", err)
	}

	if result.Logprobs == nil || len(result.Logprobs.Content) == 0 {
		t.Error("Expected logprobs to be present")
	}

	if len(result.Logprobs.Content) != 2 {
		t.Errorf("Expected 2 tokens, got %d", len(result.Logprobs.Content))
	}
}

func TestChunkChoiceWithLogprobs(t *testing.T) {
	// Test that ChunkChoice correctly handles logprobs field
	chunk := ChunkChoice{
		Index: 0,
		Delta: Message{Content: "Hello"},
		Logprobs: &ChatLogprobs{
			Content: []api.TokenLogprob{
				{Token: "Hello", Logprob: -0.15, Bytes: []int{72, 101, 108, 108, 111}},
			},
		},
	}

	// Test marshaling
	data, err := json.Marshal(chunk)
	if err != nil {
		t.Fatalf("Failed to marshal ChunkChoice: %v", err)
	}

	// Test unmarshaling
	var result ChunkChoice
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal ChunkChoice: %v", err)
	}

	if result.Logprobs == nil || len(result.Logprobs.Content) == 0 {
		t.Error("Expected logprobs to be present")
	}

	if len(result.Logprobs.Content) != 1 {
		t.Errorf("Expected 1 token, got %d", len(result.Logprobs.Content))
	}
}
