package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Schema struct {
	URL            string `json:"url"`
	Command        string `json:"command"`
	HTTPS          bool   `json:"https"`
	ExpectedOutput string `json:"expected_output"`
	StatusCode     int    `json:"status_code"`
}

// Run is the function that will get called to run an instance of a check
func Run(ctx context.Context, config string) error {
	// Define a new Schema
	schema := Schema{}

	// Unmarshal the config to the Schema
	err := json.Unmarshal([]byte(config), &schema)
	if err != nil {
		return err
	}

	var requestType string

	switch strings.ToUpper(schema.Command) {
	case "GET":
		requestType = http.MethodGet
	case "POST":
		requestType = http.MethodPost
	case "PUT":
		requestType = http.MethodPut
	case "DELETE":
		requestType = http.MethodDelete
	case "PATCH":
		requestType = http.MethodPatch
	case "HEAD":
		requestType = http.MethodHead
	case "OPTIONS":
		requestType = http.MethodOptions
	case "CONNECT":
		requestType = http.MethodConnect
	case "TRACE":
		requestType = http.MethodTrace
	default:
		return fmt.Errorf("provided invalid command/http verb: \"%v\"" + schema.Command)
	}

	req, err := http.NewRequestWithContext(ctx, requestType, schema.URL, nil)
	if err != nil {
		return fmt.Errorf("encounted error while creating request: %v", err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("encounted error while making request: %v", err.Error())
	}
	defer resp.Body.Close()

	if schema.StatusCode != 0 && resp.StatusCode != schema.StatusCode {
		return fmt.Errorf("expected status code: %v, got: %v", schema.StatusCode, resp.StatusCode)
	}

	if schema.ExpectedOutput != "" {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("encounted error while reading response body: %v", err.Error())
		}

		body := string(bodyBytes)

		if !strings.Contains(body, schema.ExpectedOutput) {
			return fmt.Errorf("expected output: \"%v\" not found in response body", schema.ExpectedOutput)
		}
	}

	return nil
}
