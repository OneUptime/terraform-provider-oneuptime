package provider

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"
    "time"
)

// Client represents the API client for oneuptime
type Client struct {
    BaseURL    string
    HTTPClient *http.Client
    ApiKey     string
}

// NewClient creates a new API client
func NewClient(host, apiKey string) (*Client, error) {
    // Ensure the host has the correct scheme
    if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
        host = "https://" + host
    }

    // Append /api to the host
    if !strings.HasSuffix(host, "/api") {
        host = strings.TrimSuffix(host, "/") + "/api"
    }

    // Parse and validate the URL
    parsedURL, err := url.Parse(host)
    if err != nil {
        return nil, fmt.Errorf("invalid host URL: %w", err)
    }

    client := &Client{
        BaseURL: parsedURL.String(),
        HTTPClient: &http.Client{
            Timeout: time.Second * 30,
        },
        ApiKey: apiKey,
    }

    return client, nil
}

// DoRequest performs an HTTP request
func (c *Client) DoRequest(method, path string, body interface{}) (*http.Response, error) {
    // Construct the full URL
    fullURL := c.BaseURL + path

    var bodyReader io.Reader
    if body != nil {
        jsonBody, err := json.Marshal(body)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal request body: %w", err)
        }
        bodyReader = bytes.NewBuffer(jsonBody)
    }

    req, err := http.NewRequest(method, fullURL, bodyReader)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    // Set headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")

    // Set authentication
    if c.ApiKey != "" {
        req.Header.Set("APIKey", c.ApiKey)
    }

    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to execute request: %w", err)
    }

    return resp, nil
}

// Get performs a GET request
func (c *Client) Get(path string) (*http.Response, error) {
    return c.DoRequest("GET", path, nil)
}

// Post performs a POST request
func (c *Client) Post(path string, body interface{}) (*http.Response, error) {
    return c.DoRequest("POST", path, body)
}

// Put performs a PUT request
func (c *Client) Put(path string, body interface{}) (*http.Response, error) {
    return c.DoRequest("PUT", path, body)
}

// Patch performs a PATCH request
func (c *Client) Patch(path string, body interface{}) (*http.Response, error) {
    return c.DoRequest("PATCH", path, body)
}

// Delete performs a DELETE request
func (c *Client) Delete(path string) (*http.Response, error) {
    return c.DoRequest("DELETE", path, nil)
}

// PostWithSelect performs a POST request with select parameter to fetch full object
func (c *Client) PostWithSelect(path string, selectParam interface{}) (*http.Response, error) {
    requestBody := map[string]interface{}{
        "select": selectParam,
    }
    return c.DoRequest("POST", path, requestBody)
}

// ParseResponse parses an HTTP response into a struct
func (c *Client) ParseResponse(resp *http.Response, target interface{}) error {
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
    }

    if target == nil {
        return nil
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("failed to read response body: %w", err)
    }

    if len(body) == 0 {
        return nil
    }

    err = json.Unmarshal(body, target)
    if err != nil {
        return fmt.Errorf("failed to unmarshal response: %w", err)
    }

    return nil
}
