package llama

import (
	"math"
	"testing"
)

func TestSoftmax(t *testing.T) {
	tests := []struct {
		name      string
		input     []float32
		tolerance float32
	}{
		{
			name:      "simple case",
			input:     []float32{1.0, 2.0, 3.0},
			tolerance: 0.001,
		},
		{
			name:      "large values",
			input:     []float32{100.0, 101.0, 102.0},
			tolerance: 0.001,
		},
		{
			name:      "negative values",
			input:     []float32{-1.0, -2.0, -3.0},
			tolerance: 0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Softmax(tt.input)
			
			// Check length
			if len(result) != len(tt.input) {
				t.Errorf("Expected length %d, got %d", len(tt.input), len(result))
			}
			
			// Check sum to 1.0
			sum := float32(0.0)
			for _, v := range result {
				sum += v
			}
			if math.Abs(float64(sum-1.0)) > float64(tt.tolerance) {
				t.Errorf("Expected sum to be 1.0, got %f", sum)
			}
			
			// Check all values are positive
			for i, v := range result {
				if v < 0 {
					t.Errorf("Expected positive value at index %d, got %f", i, v)
				}
			}
			
			// Check ordering (higher input -> higher output for softmax)
			for i := 1; i < len(result); i++ {
				if tt.input[i] > tt.input[i-1] {
					if result[i] < result[i-1] {
						t.Errorf("Softmax should preserve ordering: input[%d]=%f > input[%d]=%f, but result[%d]=%f < result[%d]=%f", 
							i, tt.input[i], i-1, tt.input[i-1], i, result[i], i-1, result[i-1])
					}
				}
			}
		})
	}
}

func TestSelectTopN(t *testing.T) {
	probs := []float32{0.1, 0.3, 0.05, 0.4, 0.15}
	
	tests := []struct {
		name        string
		n           int
		expectedLen int
		firstToken  int32 // Should be index of highest probability
	}{
		{
			name:        "top 1",
			n:           1,
			expectedLen: 1,
			firstToken:  3, // index of 0.4
		},
		{
			name:        "top 3",
			n:           3,
			expectedLen: 3,
			firstToken:  3, // index of 0.4
		},
		{
			name:        "top all",
			n:           5,
			expectedLen: 5,
			firstToken:  3, // index of 0.4
		},
		{
			name:        "more than available",
			n:           10,
			expectedLen: 5, // all available
			firstToken:  3, // index of 0.4
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SelectTopN(probs, tt.n)
			
			// Check length
			if len(result) != tt.expectedLen {
				t.Errorf("Expected length %d, got %d", tt.expectedLen, len(result))
			}
			
			// Check first token (highest probability)
			if len(result) > 0 && result[0].ID != tt.firstToken {
				t.Errorf("Expected first token ID %d, got %d", tt.firstToken, result[0].ID)
			}
			
			// Check descending order
			for i := 1; i < len(result); i++ {
				if result[i-1].Logit < result[i].Logit {
					t.Errorf("Expected descending order of logits: result[%d].Logit=%f >= result[%d].Logit=%f", 
						i-1, result[i-1].Logit, i, result[i].Logit)
				}
			}
		})
	}
}

func TestTokenToBytes(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected []int
	}{
		{
			name:     "ascii",
			token:    "hello",
			expected: []int{104, 101, 108, 108, 111},
		},
		{
			name:     "unicode",
			token:    "café",
			expected: []int{99, 97, 102, 195, 169},
		},
		{
			name:     "empty",
			token:    "",
			expected: []int{},
		},
		{
			name:     "space",
			token:    " ",
			expected: []int{32},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TokenToBytes(tt.token)
			
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
			}
			
			for i, expected := range tt.expected {
				if i < len(result) && result[i] != expected {
					t.Errorf("Expected byte[%d]=%d, got %d", i, expected, result[i])
				}
			}
		})
	}
}