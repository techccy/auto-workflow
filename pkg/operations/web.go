package operations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// WebHandler handles web automation operations
type WebHandler struct {
	client *http.Client
}

// NewWebHandler creates a new web handler
func NewWebHandler() *WebHandler {
	return &WebHandler{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetTimeout sets the HTTP client timeout
func (h *WebHandler) SetTimeout(timeout time.Duration) {
	h.client.Timeout = timeout
}

// GetType returns the operation type
func (h *WebHandler) GetType() string {
	return "web"
}

// Validate checks if the parameters are valid
func (h *WebHandler) Validate(params map[string]interface{}) error {
	action, ok := params["action"].(string)
	if !ok {
		return fmt.Errorf("action parameter is required")
	}

	switch action {
	case "get", "post", "put", "patch", "delete":
		if _, ok := params["url"]; !ok {
			return fmt.Errorf("url parameter is required for %s action", action)
		}
	case "request":
		if _, ok := params["url"]; !ok {
			return fmt.Errorf("url parameter is required for request action")
		}
		if _, ok := params["method"]; !ok {
			return fmt.Errorf("method parameter is required for request action")
		}
	}

	return nil
}

// Execute performs the web operation
func (h *WebHandler) Execute(ctx context.Context, params map[string]interface{}, execCtx interface{}) (map[string]interface{}, error) {
	action := params["action"].(string)

	var result map[string]interface{}
	var err error

	switch action {
	case "get":
		result, err = h.executeGet(ctx, params)
	case "post":
		result, err = h.executePost(ctx, params)
	case "put":
		result, err = h.executePut(ctx, params)
	case "patch":
		result, err = h.executePatch(ctx, params)
	case "delete":
		result, err = h.executeDelete(ctx, params)
	case "request":
		result, err = h.executeRequest(ctx, params)
	default:
		return nil, fmt.Errorf("unsupported web action: %s", action)
	}

	return result, err
}

func (h *WebHandler) executeGet(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	url := getStringParam(params, "url")
	headers := getHeaders(params)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	startTime := time.Now()
	resp, err := h.client.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %v", err)
	}

	return map[string]interface{}{
		"success":     resp.StatusCode >= 200 && resp.StatusCode < 300,
		"status_code": resp.StatusCode,
		"body":        string(body),
		"headers":     resp.Header,
		"duration":    duration.Milliseconds(),
	}, nil
}

func (h *WebHandler) executePost(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	url := getStringParam(params, "url")
	headers := getHeaders(params)

	var body io.Reader
	if content, ok := params["body"]; ok {
		switch v := content.(type) {
		case string:
			body = bytes.NewBufferString(v)
		case map[string]interface{}:
			jsonData, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("marshal body failed: %v", err)
			}
			body = bytes.NewBuffer(jsonData)
			if headers["Content-Type"] == "" {
				headers["Content-Type"] = "application/json"
			}
		}
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	startTime := time.Now()
	resp, err := h.client.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %v", err)
	}

	return map[string]interface{}{
		"success":     resp.StatusCode >= 200 && resp.StatusCode < 300,
		"status_code": resp.StatusCode,
		"body":        string(responseBody),
		"headers":     resp.Header,
		"duration":    duration.Milliseconds(),
	}, nil
}

func (h *WebHandler) executePut(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	url := getStringParam(params, "url")
	headers := getHeaders(params)

	var body io.Reader
	if content, ok := params["body"]; ok {
		switch v := content.(type) {
		case string:
			body = bytes.NewBufferString(v)
		case map[string]interface{}:
			jsonData, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("marshal body failed: %v", err)
			}
			body = bytes.NewBuffer(jsonData)
			if headers["Content-Type"] == "" {
				headers["Content-Type"] = "application/json"
			}
		}
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", url, body)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	startTime := time.Now()
	resp, err := h.client.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %v", err)
	}

	return map[string]interface{}{
		"success":     resp.StatusCode >= 200 && resp.StatusCode < 300,
		"status_code": resp.StatusCode,
		"body":        string(responseBody),
		"headers":     resp.Header,
		"duration":    duration.Milliseconds(),
	}, nil
}

func (h *WebHandler) executePatch(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	url := getStringParam(params, "url")
	headers := getHeaders(params)

	var body io.Reader
	if content, ok := params["body"]; ok {
		switch v := content.(type) {
		case string:
			body = bytes.NewBufferString(v)
		case map[string]interface{}:
			jsonData, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("marshal body failed: %v", err)
			}
			body = bytes.NewBuffer(jsonData)
			if headers["Content-Type"] == "" {
				headers["Content-Type"] = "application/json"
			}
		}
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, body)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	startTime := time.Now()
	resp, err := h.client.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %v", err)
	}

	return map[string]interface{}{
		"success":     resp.StatusCode >= 200 && resp.StatusCode < 300,
		"status_code": resp.StatusCode,
		"body":        string(responseBody),
		"headers":     resp.Header,
		"duration":    duration.Milliseconds(),
	}, nil
}

func (h *WebHandler) executeDelete(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	url := getStringParam(params, "url")
	headers := getHeaders(params)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	startTime := time.Now()
	resp, err := h.client.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %v", err)
	}

	return map[string]interface{}{
		"success":     resp.StatusCode >= 200 && resp.StatusCode < 300,
		"status_code": resp.StatusCode,
		"body":        string(body),
		"headers":     resp.Header,
		"duration":    duration.Milliseconds(),
	}, nil
}

func (h *WebHandler) executeRequest(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	url := getStringParam(params, "url")
	method := getStringParam(params, "method")
	headers := getHeaders(params)

	var body io.Reader
	if content, ok := params["body"]; ok {
		switch v := content.(type) {
		case string:
			body = bytes.NewBufferString(v)
		case map[string]interface{}:
			jsonData, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("marshal body failed: %v", err)
			}
			body = bytes.NewBuffer(jsonData)
			if headers["Content-Type"] == "" {
				headers["Content-Type"] = "application/json"
			}
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	startTime := time.Now()
	resp, err := h.client.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %v", err)
	}

	return map[string]interface{}{
		"success":     resp.StatusCode >= 200 && resp.StatusCode < 300,
		"status_code": resp.StatusCode,
		"body":        string(responseBody),
		"headers":     resp.Header,
		"duration":    duration.Milliseconds(),
	}, nil
}

// Helper functions
func getHeaders(params map[string]interface{}) map[string]string {
	headers := make(map[string]string)
	if h, ok := params["headers"].(map[string]interface{}); ok {
		for k, v := range h {
			if str, ok := v.(string); ok {
				headers[k] = str
			}
		}
	}
	return headers
}
