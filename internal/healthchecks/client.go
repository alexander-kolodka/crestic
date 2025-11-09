package healthchecks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"path"
	"strings"
	"time"

	"github.com/alexander-kolodka/crestic/internal/logger"
)

type Client struct {
	http *http.Client
}

func NewClient() *Client {
	const timeout = 3 * time.Second
	return &Client{
		http: &http.Client{Timeout: timeout},
	}
}

type Payload struct {
	JobName string `json:"jobName"`
	Err     string `json:"error,omitempty"`
}

// Start signals the beginning of task execution.
// Enables tracking of "hanging" tasks that started but never completed.
// baseURL should be the full healthcheck URL (e.g., https://hc-ping.com/{uuid} or https://hc-ping.com/{uuid}/{slug})
// rid is the unique run ID for grouping signals.
func (c *Client) Start(ctx context.Context, baseURL, rid string, p Payload) error {
	return c.post(ctx, baseURL, "start", rid, p)
}

// Success reports successful task completion.
// Healthchecks.io automatically calculates duration between /start and this ping.
// baseURL should be the full healthcheck URL (e.g., https://hc-ping.com/{uuid} or https://hc-ping.com/{uuid}/{slug})
// rid must match the one passed to Start.
func (c *Client) Success(ctx context.Context, baseURL, rid string, p Payload) error {
	return c.post(ctx, baseURL, "", rid, p)
}

// Fail reports task failure with error message.
// Healthchecks.io automatically calculates duration and stores the error message.
// errMsg is sent in request body and displayed in Healthchecks.io interface.
// baseURL should be the full healthcheck URL (e.g., https://hc-ping.com/{uuid} or https://hc-ping.com/{uuid}/{slug})
// rid must match the one passed to Start.
func (c *Client) Fail(ctx context.Context, baseURL, rid string, p Payload) error {
	return c.post(ctx, baseURL, "fail", rid, p)
}

func (c *Client) post(ctx context.Context, baseURL, endpoint, rid string, p Payload) error {
	u, err := buildURL(baseURL, endpoint, rid)
	if err != nil {
		return fmt.Errorf("invalid healthcheck URL: %w", err)
	}

	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	return withRetry(ctx, func() error {
		return c.doPost(ctx, u, body)
	})
}

// buildURL builds https://hc-ping.com/{uuid}[/{slug}][/endpoint][?rid=...]
func buildURL(base, endpoint, rid string) (string, error) {
	base = strings.TrimSpace(base)
	if base == "" {
		return "", fmt.Errorf("empty base URL")
	}

	u, err := neturl.Parse(base)
	if err != nil {
		return "", err
	}

	if endpoint != "" {
		u.Path = path.Join(strings.TrimRight(u.Path, "/"), endpoint)
	}

	if rid != "" {
		q := u.Query()
		q.Set("rid", rid)
		u.RawQuery = q.Encode()
	}

	return u.String(), nil
}

func (c *Client) doPost(ctx context.Context, url string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		log := logger.FromContext(ctx)
		log.Warn().
			Err(err).
			RawJSON("body", body).
			Str("url", url).
			Msg("Failed to send healthcheck request")
		return &retryableError{err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	respBody, _ := io.ReadAll(resp.Body)
	httpErr := fmt.Errorf("healthcheck error %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))

	log := logger.FromContext(ctx)
	log.Warn().
		Int("status_code", resp.StatusCode).
		Str("url", url).
		RawJSON("body", body).
		Str("response", string(respBody)).
		Msg("Healthcheck request returned non-OK status")

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return httpErr
	}

	return &retryableError{err: httpErr}
}
