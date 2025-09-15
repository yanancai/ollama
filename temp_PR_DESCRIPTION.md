# api: add logprobs support for token probability output

## Overview

This PR adds logprobs support to Ollama's API, allowing users to receive log probability information for generated tokens. This feature enables users to understand the model's confidence in its token choices and explore alternative token candidates.

This is a long-requested feature by the community (see [issue #2415](https://github.com/ollama/ollama/issues/2415)) that brings Ollama's API closer to OpenAI's standard and provides essential functionality for applications requiring probability-aware text generation.

## What is Logprobs?

Logprobs (log probabilities) provide insight into the model's token selection process by returning:
- **Log probability** for each generated token (indicating confidence)
- **Top alternative tokens** with their probabilities (when `top_logprobs` > 0)

This is useful for:
- Understanding model confidence in generated text
- Analyzing alternative token choices
- Building applications that need probability-aware text generation
- Debugging and improving prompts based on model uncertainty

## Features Added

### Request Parameters

Both parameters are configured through the `options` object:

- **`logprobs`** (boolean): Enable log probability output
- **`top_logprobs`** (integer, 0-20): Number of alternative tokens to include with probabilities

### Response Format

The response format varies by API endpoint:

#### Generate API (`/api/generate`)

```json
{
  "model": "qwen2.5:0.5b",
  "response": "Hello! How can I help you today?",
  "done": true,
  "logprobs": [
    {
      "token": "Hello",
      "logprob": -0.0234361093949319,
      "bytes": [72, 101, 108, 108, 111],
      "top_logprobs": [
        {
          "token": "Hello",
          "logprob": -0.023436108604073524,
          "bytes": [72, 101, 108, 108, 111]
        },
        {
          "token": "Hi",
          "logprob": -5.994608402252197,
          "bytes": [72, 105]
        }
      ]
    }
  ]
}
```

#### Chat API (`/api/chat`)

```json
{
  "model": "qwen2.5:0.5b",
  "created_at": "2025-09-15T22:51:43.641524894Z",
  "message": {
    "role": "assistant",
    "content": "Hello! It's nice to meet you.",
    "logprobs": [
      {
        "token": "Hello",
        "logprob": -0.0234361093949319,
        "bytes": [72, 101, 108, 108, 111],
        "top_logprobs": [
          {
            "token": "Hello",
            "logprob": -0.023436108604073524,
            "bytes": [72, 101, 108, 108, 111]
          },
          {
            "token": "Hi",
            "logprob": -5.994608402252197,
            "bytes": [72, 105]
          }
        ]
      }
    ]
  },
  "done": true
}
```

#### OpenAI Compatible API (`/v1/chat/completions`)

```json
{
  "id": "chatcmpl-916",
  "object": "chat.completion",
  "created": 1757976716,
  "model": "qwen2.5:0.5b",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "Hello there! How can I assist you today?"
      },
      "finish_reason": "stop",
      "logprobs": {
        "content": [
          {
            "token": "Hello",
            "logprob": -0.0234361093949319,
            "bytes": [72, 101, 108, 108, 111],
            "top_logprobs": [
              {
                "token": "Hello",
                "logprob": -0.023436108604073524,
                "bytes": [72, 101, 108, 108, 111]
              },
              {
                "token": "Hi",
                "logprob": -5.994608402252197,
                "bytes": [72, 105]
              }
            ]
          }
        ]
      }
    }
  ],
  "usage": {
    "prompt_tokens": 31,
    "completion_tokens": 10,
    "total_tokens": 41
  }
}
```

## API Usage

### Native Ollama API

```bash
curl -X POST http://localhost:11434/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "model": "llama3.2:1b",
    "prompt": "Hello world",
    "options": {
      "logprobs": true,
      "top_logprobs": 3
    }
  }'
```

### Chat API

```bash
curl -X POST http://localhost:11434/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "model": "llama3.2:1b",
    "messages": [{"role": "user", "content": "Hello"}],
    "options": {
      "logprobs": true,
      "top_logprobs": 2
    }
  }'
```

### OpenAI-Compatible API

```bash
curl -X POST http://localhost:11434/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "llama3.2:1b",
    "messages": [{"role": "user", "content": "Hello"}],
    "logprobs": true,
    "top_logprobs": 3
  }'
```

## Configuration

Logprobs settings can be:
- Set per-request in the `options` object
- Configured in model files for default behavior
- Inherited from parent configurations

Example in a Modelfile:
```
FROM llama3.2:1b
PARAMETER logprobs true
PARAMETER top_logprobs 2
```

## Validation

- `top_logprobs` requires `logprobs` to be `true`
- `top_logprobs` must be between 0 and 20
- Invalid configurations return clear error messages

## Backward Compatibility

This change is **fully backward compatible**:
- Existing requests continue to work unchanged
- New fields are optional and default to disabled
- No breaking changes to existing API contracts
- OpenAI API compatibility maintained

## Implementation

- **`api/types.go`**: Added logprobs configuration to `Options` struct and response types
- **`server/routes.go`**: Added validation and processing logic for both Generate and Chat handlers  
- **`openai/openai.go`**: Added middleware mapping for OpenAI API compatibility
- **Tests**: Comprehensive test coverage for all API endpoints and edge cases

## Testing

- ✅ Unit tests for logprobs data structures and serialization
- ✅ Integration tests for Generate and Chat APIs  
- ✅ OpenAI compatibility tests
- ✅ Validation and error handling tests
- ✅ Backward compatibility verification