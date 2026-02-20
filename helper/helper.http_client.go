package helper

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/bytedance/sonic"
)

type ExternalJSONRequestOptions struct {
	Headers     map[string]string
	Query       map[string]string
	QueryValues url.Values
	JSONBody    interface{}
	FormBody    url.Values
	Body        io.Reader
	ContentType string
}

type ExternalJSONResponseError struct {
	StatusCode int
	Body       string
}

type ExternalResponse struct {
	StatusCode int
	Body       []byte
	Headers    map[string]string
	Success    *bool
}

func (e *ExternalJSONResponseError) Error() string {
	if e == nil {
		return "external json response error"
	}
	if e.Body == "" {
		return fmt.Sprintf("external request failed with status %d", e.StatusCode)
	}
	return fmt.Sprintf("external request failed with status %d: %s", e.StatusCode, e.Body)
}

func DoExternalJSONRequest(ctx context.Context, client *http.Client, method, endpoint string, options ExternalJSONRequestOptions) (*ExternalResponse, error) {
	if client == nil {
		client = http.DefaultClient
	}

	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	query := parsedURL.Query()
	for key, values := range options.QueryValues {
		for _, value := range values {
			query.Add(key, value)
		}
	}
	for key, value := range options.Query {
		query.Set(key, value)
	}
	parsedURL.RawQuery = query.Encode()

	requestBody := options.Body
	contentType := options.ContentType

	if options.FormBody != nil {
		requestBody = strings.NewReader(options.FormBody.Encode())
		if contentType == "" {
			contentType = "application/x-www-form-urlencoded"
		}
	}

	if options.JSONBody != nil {
		jsonBody, marshalErr := sonic.Marshal(options.JSONBody)
		if marshalErr != nil {
			return nil, marshalErr
		}
		requestBody = bytes.NewReader(jsonBody)
		if contentType == "" {
			contentType = "application/json"
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, parsedURL.String(), requestBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &ExternalResponse{
		StatusCode: resp.StatusCode,
		Body:       bodyBytes,
		Headers:    map[string]string{},
	}

	for key, values := range resp.Header {
		if len(values) > 0 {
			result.Headers[key] = values[0]
		}
	}

	var envelope struct {
		Success bool `json:"success"`
	}
	if len(bodyBytes) > 0 && sonic.Unmarshal(bodyBytes, &envelope) == nil {
		success := envelope.Success
		result.Success = &success
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return result, &ExternalJSONResponseError{
			StatusCode: resp.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	return result, nil
}
