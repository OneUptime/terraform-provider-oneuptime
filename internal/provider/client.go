package provider

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "regexp"
    "strings"
    "time"
)

// Client represents the API client for oneuptime
type Client struct {
    BaseURL    string
    HTTPClient *http.Client
    ApiKey     string
    UserAgent  string
}

// NewClient creates a new API client
func NewClient(oneuptimeUrl, apiKey, version string) (*Client, error) {
    // Ensure the oneuptimeUrl has the correct scheme
    if !strings.HasPrefix(oneuptimeUrl, "http://") && !strings.HasPrefix(oneuptimeUrl, "https://") {
        oneuptimeUrl = "https://" + oneuptimeUrl
    }

    // Append /api to the oneuptimeUrl
    if !strings.HasSuffix(oneuptimeUrl, "/api") {
        oneuptimeUrl = strings.TrimSuffix(oneuptimeUrl, "/") + "/api"
    }

    // Parse and validate the URL
    parsedURL, err := url.Parse(oneuptimeUrl)
    if err != nil {
        return nil, fmt.Errorf("invalid oneuptime_url: %w", err)
    }

    client := &Client{
        BaseURL: parsedURL.String(),
        HTTPClient: &http.Client{
            Timeout: time.Second * 60,
        },
        ApiKey:    apiKey,
        UserAgent: "terraform-provider-oneuptime/" + version,
    }

    return client, nil
}

// DoRequest performs an HTTP request. The context propagates Terraform's
// cancellation and deadlines into the HTTP layer.
func (c *Client) DoRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
    // Construct the full URL
    fullURL := c.BaseURL + path

    var jsonBody []byte
    if body != nil {
        var err error
        jsonBody, err = json.Marshal(body)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal request body: %w", err)
        }
    }

    buildRequest := func() (*http.Request, error) {
        var bodyReader io.Reader
        if jsonBody != nil {
            bodyReader = bytes.NewBuffer(jsonBody)
        }
        req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
        if err != nil {
            return nil, fmt.Errorf("failed to create request: %w", err)
        }
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Accept", "application/json")
        req.Header.Set("User-Agent", c.UserAgent)
        if c.ApiKey != "" {
            req.Header.Set("APIKey", c.ApiKey)
        }
        return req, nil
    }

    /*
     * Reads and deletes are safe to retry; creates and updates are not
     * (retrying a POST after an ambiguous failure could duplicate the
     * resource). Retry only on 429 and transient 5xx.
     */
    idempotent := method == "GET" || method == "DELETE" || strings.HasSuffix(path, "/get-item") || strings.HasSuffix(path, "/get-list") || strings.HasSuffix(path, "/count")
    attempts := 1
    if idempotent {
        attempts = 3
    }

    var resp *http.Response
    for attempt := 0; attempt < attempts; attempt++ {
        if attempt > 0 {
            select {
            case <-ctx.Done():
                return nil, ctx.Err()
            case <-time.After(time.Duration(500*(1<<attempt)) * time.Millisecond):
            }
        }

        req, err := buildRequest()
        if err != nil {
            return nil, err
        }

        resp, err = c.HTTPClient.Do(req)
        if err != nil {
            if attempt == attempts-1 {
                return nil, fmt.Errorf("failed to execute request: %w", err)
            }
            continue
        }

        if resp.StatusCode == http.StatusTooManyRequests ||
            resp.StatusCode == http.StatusBadGateway ||
            resp.StatusCode == http.StatusServiceUnavailable ||
            resp.StatusCode == http.StatusGatewayTimeout {
            if attempt < attempts-1 {
                resp.Body.Close()
                continue
            }
        }

        return resp, nil
    }

    return resp, nil
}

// Get performs a GET request
func (c *Client) Get(ctx context.Context, path string) (*http.Response, error) {
    return c.DoRequest(ctx, "GET", path, nil)
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, path string, body interface{}) (*http.Response, error) {
    return c.DoRequest(ctx, "POST", path, body)
}

// Put performs a PUT request
func (c *Client) Put(ctx context.Context, path string, body interface{}) (*http.Response, error) {
    return c.DoRequest(ctx, "PUT", path, body)
}

// Patch performs a PATCH request
func (c *Client) Patch(ctx context.Context, path string, body interface{}) (*http.Response, error) {
    return c.DoRequest(ctx, "PATCH", path, body)
}

// Delete performs a DELETE request
func (c *Client) Delete(ctx context.Context, path string) (*http.Response, error) {
    return c.DoRequest(ctx, "DELETE", path, nil)
}

// rejectedSelectColumn matches the server's column-permission error, e.g.
// "You do not have permissions to select on - serviceLanguage."
var rejectedSelectColumn = regexp.MustCompile(`select on - ([A-Za-z0-9_]+)`)

// PostWithSelect performs a POST request with a select parameter. When the
// server rejects a column in the select (permission-gated columns, or
// columns this server version does not know about yet), that column is
// dropped and the request retried — a version- or permission-skewed column
// must not fail the entire read.
func (c *Client) PostWithSelect(ctx context.Context, path string, selectParam interface{}) (*http.Response, error) {
    selectMap, _ := selectParam.(map[string]interface{})

    // One attempt per droppable column, bounded to keep worst cases sane.
    maxAttempts := 8
    for attempt := 0; attempt < maxAttempts; attempt++ {
        requestBody := map[string]interface{}{
            "select": selectParam,
        }
        resp, err := c.DoRequest(ctx, "POST", path, requestBody)
        if err != nil {
            return nil, err
        }
        if selectMap == nil ||
            (resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusUnprocessableEntity) {
            return resp, nil
        }

        body, readErr := io.ReadAll(resp.Body)
        resp.Body.Close()
        if readErr != nil {
            return nil, fmt.Errorf("failed to read response body: %w", readErr)
        }

        match := rejectedSelectColumn.FindSubmatch(body)
        rebuilt := func() *http.Response {
            resp.Body = io.NopCloser(bytes.NewReader(body))
            return resp
        }
        if match == nil {
            // Not a select-column rejection: surface the original error.
            return rebuilt(), nil
        }
        column := string(match[1])
        if _, present := selectMap[column]; !present {
            return rebuilt(), nil
        }
        delete(selectMap, column)
    }

    requestBody := map[string]interface{}{
        "select": selectParam,
    }
    return c.DoRequest(ctx, "POST", path, requestBody)
}

// apiErrorMessage extracts the server's human-readable error message from an
// error response body, falling back to the raw body.
func apiErrorMessage(body []byte) string {
    var parsed map[string]interface{}
    if err := json.Unmarshal(body, &parsed); err == nil {
        for _, key := range []string{"message", "error", "errorMessage"} {
            if msg, ok := parsed[key].(string); ok && msg != "" {
                return msg
            }
        }
    }
    trimmed := strings.TrimSpace(string(body))
    if trimmed == "" {
        return "(empty response body)"
    }
    return trimmed
}

// ParseResponse parses an HTTP response into a struct
func (c *Client) ParseResponse(resp *http.Response, target interface{}) error {
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("the OneUptime API returned status %d: %s", resp.StatusCode, apiErrorMessage(body))
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
