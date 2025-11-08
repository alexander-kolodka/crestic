package healthchecks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/alexander-kolodka/crestic/internal/logger"
)

// Client sends signals to Healthchecks.io for monitoring command execution.
// Supports both formats: https://hc-ping.com/{uuid} or https://hc-ping.com/{uuid}/{slug}
type Client struct {
	http *http.Client
}

// NewClient creates a Healthchecks.io client.
func NewClient() *Client {
	const timeout = 3
	return &Client{
		http: &http.Client{
			Timeout: timeout * time.Second,
		},
	}
}

type Payload struct {
	JobName string `json:"jobName"`
	Err     string `json:"error,omitempty"`
}

// Start signals the beginning of task execution.
// Enables tracking of "hanging" tasks that started but never completed.
// url should be the full healthcheck URL (e.g., https://hc-ping.com/{uuid} or https://hc-ping.com/{uuid}/{slug})
// rid is the unique run ID for grouping signals.
func (c *Client) Start(ctx context.Context, url, rid string, p Payload) error {
	return c.post(ctx, url, "start", rid, p)
}

// Success reports successful task completion.
// Healthchecks.io automatically calculates duration between /start and this ping.
// url should be the full healthcheck URL (e.g., https://hc-ping.com/{uuid} or https://hc-ping.com/{uuid}/{slug})
// rid must match the one passed to Start.
func (c *Client) Success(ctx context.Context, url, rid string, p Payload) error {
	return c.post(ctx, url, "", rid, p)
}

// Fail reports task failure with error message.
// Healthchecks.io automatically calculates duration and stores the error message.
// errMsg is sent in request body and displayed in Healthchecks.io interface.
// url should be the full healthcheck URL (e.g., https://hc-ping.com/{uuid} or https://hc-ping.com/{uuid}/{slug})
// rid must match the one passed to Start.
func (c *Client) Fail(ctx context.Context, url, rid string, p Payload) error {
	return c.post(ctx, url, "fail", rid, p)
}

func (c *Client) post(ctx context.Context, url, endpoint, rid string, p Payload) error {
	url = buildURL(url, endpoint, rid)
	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("error marshalling Payload: %w", err)
	}

	return withRetry(ctx, func() error {
		return c.doPost(ctx, url, string(body))
	})
}

func buildURL(url, endpoint, rid string) string {
	url = strings.TrimSuffix(url, "/")
	if endpoint != "" {
		url = url + "/" + endpoint
	}

	if rid != "" {
		url = fmt.Sprintf("%s?rid=%s", url, rid)
	}

	return url
}

func (c *Client) doPost(ctx context.Context, url, body string) error {
	req := httptest.NewRequest(http.MethodPost, url, strings.NewReader(body))

	resp, err := c.http.Do(req)
	if err != nil {
		log := logger.FromContext(ctx)
		log.Warn().
			Err(err).
			Str("body", body).
			Str("url", url).
			Msg("Failed to send healthcheck request")
		return &retryableError{err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	respBody, _ := io.ReadAll(resp.Body)
	httpErr := fmt.Errorf("healthcheck error %d: %s", resp.StatusCode, string(respBody))

	log := logger.FromContext(ctx)
	log.Warn().
		Int("status_code", resp.StatusCode).
		Str("body", body).
		Str("url", url).
		Str("response", string(respBody)).
		Msg("Healthcheck request returned non-OK status")

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return httpErr
	}

	return &retryableError{err: httpErr}
}
