#!/usr/bin/env python3
"""
Integration test script to verify logprobs implementation works correctly
with both generation and chat APIs using the new options-based configuration.
"""

import json
import requests
import sys
from typing import Dict, Any

def test_generate_api_with_logprobs():
    """Test /api/generate endpoint with logprobs via options"""
    url = "http://localhost:11434/api/generate"
    
    # Test with logprobs enabled via options
    data = {
        "model": "qwen2.5:0.5b",  # Use a small model for testing
        "prompt": "Say hello",
        "stream": False,
        "options": {
            "logprobs": True,
            "top_logprobs": 3
        }
    }
    
    print("Testing Generate API with logprobs via options...")
    print(f"Request: {json.dumps(data, indent=2)}")
    
    try:
        response = requests.post(url, json=data, timeout=30)
        response.raise_for_status()
        
        result = response.json()
        print(f"Response: {json.dumps(result, indent=2)}")
        
        # Verify logprobs structure
        if 'logprobs' in result and result['logprobs']:
            print("✅ Generate API: Logprobs found in response")
            
            # Verify structure: should be array of TokenLogprob
            logprobs = result['logprobs']
            if isinstance(logprobs, list) and len(logprobs) > 0:
                first_token = logprobs[0]
                required_fields = ['token', 'logprob']
                if all(field in first_token for field in required_fields):
                    print("✅ Generate API: Logprobs structure is correct")
                    
                    # Check for top_logprobs if requested
                    if 'top_logprobs' in first_token and first_token['top_logprobs']:
                        print("✅ Generate API: Top logprobs found")
                    else:
                        print("⚠️  Generate API: Top logprobs not found (may be normal)")
                else:
                    print("❌ Generate API: Logprobs structure missing required fields")
                    return False
            else:
                print("❌ Generate API: Logprobs is not a proper array")
                return False
        else:
            print("❌ Generate API: No logprobs in response")
            return False
            
        return True
        
    except requests.exceptions.RequestException as e:
        print(f"❌ Generate API request failed: {e}")
        return False
    except json.JSONDecodeError as e:
        print(f"❌ Generate API response parsing failed: {e}")
        return False

def test_chat_api_with_logprobs():
    """Test /api/chat endpoint with logprobs via options"""
    url = "http://localhost:11434/api/chat"
    
    # Test with logprobs enabled via options
    data = {
        "model": "qwen2.5:0.5b",  # Use a small model for testing
        "messages": [{"role": "user", "content": "Say hello"}],
        "stream": False,
        "options": {
            "logprobs": True,
            "top_logprobs": 3
        }
    }
    
    print("\nTesting Chat API with logprobs via options...")
    print(f"Request: {json.dumps(data, indent=2)}")
    
    try:
        response = requests.post(url, json=data, timeout=30)
        response.raise_for_status()
        
        result = response.json()
        print(f"Response: {json.dumps(result, indent=2)}")
        
        # Verify logprobs structure in message
        if 'message' in result and 'logprobs' in result['message'] and result['message']['logprobs']:
            print("✅ Chat API: Logprobs found in message")
            
            # Verify structure: should be array of TokenLogprob
            logprobs = result['message']['logprobs']
            if isinstance(logprobs, list) and len(logprobs) > 0:
                first_token = logprobs[0]
                required_fields = ['token', 'logprob']
                if all(field in first_token for field in required_fields):
                    print("✅ Chat API: Logprobs structure is correct")
                    
                    # Check for top_logprobs if requested
                    if 'top_logprobs' in first_token and first_token['top_logprobs']:
                        print("✅ Chat API: Top logprobs found")
                    else:
                        print("⚠️  Chat API: Top logprobs not found (may be normal)")
                else:
                    print("❌ Chat API: Logprobs structure missing required fields")
                    return False
            else:
                print("❌ Chat API: Logprobs is not a proper array")
                return False
        else:
            print("❌ Chat API: No logprobs in message")
            return False
            
        return True
        
    except requests.exceptions.RequestException as e:
        print(f"❌ Chat API request failed: {e}")
        return False
    except json.JSONDecodeError as e:
        print(f"❌ Chat API response parsing failed: {e}")
        return False

def test_openai_compatibility():
    """Test OpenAI-compatible endpoint with logprobs"""
    url = "http://localhost:11434/v1/chat/completions"
    
    data = {
        "model": "qwen2.5:0.5b",
        "messages": [
            {"role": "user", "content": "Say hello"}
        ],
        "logprobs": True,
        "top_logprobs": 3,
        "stream": False
    }
    
    print("\nTesting OpenAI-compatible API with logprobs...")
    print(f"Request: {json.dumps(data, indent=2)}")
    
    try:
        response = requests.post(url, json=data, timeout=30)
        response.raise_for_status()
        
        result = response.json()
        print(f"Response: {json.dumps(result, indent=2)}")
        
        # Verify OpenAI-style response structure
        if 'choices' in result and len(result['choices']) > 0:
            choice = result['choices'][0]
            if 'logprobs' in choice and choice['logprobs']:
                print("✅ OpenAI API: Logprobs found in choice")
                return True
            else:
                print("❌ OpenAI API: No logprobs in choice")
                return False
        else:
            print("❌ OpenAI API: Invalid response structure")
            return False
            
    except requests.exceptions.RequestException as e:
        print(f"❌ OpenAI API request failed: {e}")
        return False
    except json.JSONDecodeError as e:
        print(f"❌ OpenAI API response parsing failed: {e}")
        return False

def main():
    """Run all integration tests"""
    print("🧪 Running Logprobs Integration Tests")
    print("=" * 50)
    
    tests = [
        ("Generate API", test_generate_api_with_logprobs),
        ("Chat API", test_chat_api_with_logprobs),  
        ("OpenAI Compatibility", test_openai_compatibility)
    ]
    
    results = {}
    for test_name, test_func in tests:
        print(f"\n📋 Running {test_name} test...")
        try:
            results[test_name] = test_func()
        except Exception as e:
            print(f"❌ {test_name} test failed with exception: {e}")
            results[test_name] = False
    
    print("\n" + "=" * 50)
    print("📊 Test Results Summary:")
    
    all_passed = True
    for test_name, passed in results.items():
        status = "✅ PASS" if passed else "❌ FAIL"
        print(f"  {test_name}: {status}")
        if not passed:
            all_passed = False
    
    if all_passed:
        print("\n🎉 All tests passed! Logprobs implementation is working correctly.")
        sys.exit(0)
    else:
        print("\n⚠️  Some tests failed. Please check the implementation.")
        sys.exit(1)

if __name__ == "__main__":
    main()