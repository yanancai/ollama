package llamarunner

import (
	"testing"

	"github.com/ollama/ollama/api"
	"github.com/ollama/ollama/llm"
)

func TestNewSequenceWithLogprobs(t *testing.T) {
	// Test that NewSequenceParams correctly handles logprobs configuration
	params := NewSequenceParams{
		numPredict:  10,
		logprobs:    true,
		topLogprobs: 5,
	}

	// Basic validation that the params structure works
	if !params.logprobs {
		t.Error("Expected logprobs to be true")
	}

	if params.topLogprobs != 5 {
		t.Errorf("Expected topLogprobs to be 5, got %d", params.topLogprobs)
	}
}

func TestResponseStructure(t *testing.T) {
	// Test the Response structure
	response := Response{
		Text: "Hello world",
		Logprobs: []api.TokenLogprob{
			{
				Token:   "Hello",
				Logprob: -0.1,
				Bytes:   []int{72, 101, 108, 108, 111},
				TopLogprobs: []api.TopTokenLogprob{
					{Token: "Hello", Logprob: -0.1, Bytes: []int{72, 101, 108, 108, 111}},
					{Token: "Hi", Logprob: -1.2, Bytes: []int{72, 105}},
				},
			},
			{
				Token:   " world",
				Logprob: -0.2,
				Bytes:   []int{32, 119, 111, 114, 108, 100},
			},
		},
	}

	if response.Text != "Hello world" {
		t.Errorf("Expected text 'Hello world', got '%s'", response.Text)
	}

	if len(response.Logprobs) != 2 {
		t.Errorf("Expected 2 logprobs, got %d", len(response.Logprobs))
	}

	// Check first token
	if response.Logprobs[0].Token != "Hello" {
		t.Errorf("Expected first token 'Hello', got '%s'", response.Logprobs[0].Token)
	}

	if len(response.Logprobs[0].TopLogprobs) != 2 {
		t.Errorf("Expected 2 top logprobs for first token, got %d", len(response.Logprobs[0].TopLogprobs))
	}
}

func TestCompletionRequestLogprobs(t *testing.T) {
	// Test the CompletionRequest structure in llm package
	req := llm.CompletionRequest{
		Prompt:      "Hello",
		Logprobs:    true,
		TopLogprobs: 10,
	}

	if !req.Logprobs {
		t.Error("Expected logprobs to be true")
	}

	if req.TopLogprobs != 10 {
		t.Errorf("Expected topLogprobs to be 10, got %d", req.TopLogprobs)
	}
}

func TestCompletionResponseLogprobs(t *testing.T) {
	// Test the CompletionResponse structure in llm package
	response := llm.CompletionResponse{
		Content: "Hello world",
		Logprobs: []api.TokenLogprob{
			{
				Token:   "Hello",
				Logprob: -0.5,
				Bytes:   []int{72, 101, 108, 108, 111},
			},
		},
		Done: true,
	}

	if response.Content != "Hello world" {
		t.Errorf("Expected content 'Hello world', got '%s'", response.Content)
	}

	if len(response.Logprobs) != 1 {
		t.Errorf("Expected 1 logprob, got %d", len(response.Logprobs))
	}

	if response.Logprobs[0].Token != "Hello" {
		t.Errorf("Expected token 'Hello', got '%s'", response.Logprobs[0].Token)
	}
}