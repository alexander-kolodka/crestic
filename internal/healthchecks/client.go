package healthchecks

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"path"
	"strings"
	"time"

	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/alexander-kolodka/crestic/internal/logger"
)

type JobsList struct {
	Jobs []string `json:"jobs"`
}

func NewJobsList(jobs []string) *JobsList {
	return &JobsList{Jobs: jobs}
}

type Client struct {
	http    *http.Client
	baseURL string
}

func NewClient(baseURL string) (*Client, error) {
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		return nil, errors.New("empty base URL")
	}

	const timeout = 3 * time.Second
	return &Client{
		http:    &http.Client{Timeout: timeout},
		baseURL: baseURL,
	}, nil
}

// Start signals the beginning of task execution.
// Enables tracking of "hanging" tasks that started but never completed.
// baseURL should be the full healthcheck URL (e.g., https://hc-ping.com/{uuid} or https://hc-ping.com/{uuid}/{slug})
// rid is the unique run ID for grouping signals.
func (c *Client) Start(ctx context.Context, rid string, j *JobsList) error {
	return c.post(ctx, "start", rid, j)
}

// Success reports successful task completion.
// Healthchecks.io automatically calculates duration between /start and this ping.
// baseURL should be the full healthcheck URL (e.g., https://hc-ping.com/{uuid} or https://hc-ping.com/{uuid}/{slug})
// rid must match the one passed to Start.
func (c *Client) Success(ctx context.Context, rid string, r *entity.JobResults) error {
	return c.post(ctx, "", rid, r)
}

// Fail reports task failure with error message.
// Healthchecks.io automatically calculates duration and stores the error message.
// errMsg is sent in request body and displayed in Healthchecks.io interface.
// baseURL should be the full healthcheck URL (e.g., https://hc-ping.com/{uuid} or https://hc-ping.com/{uuid}/{slug})
// rid must match the one passed to Start.
func (c *Client) Fail(ctx context.Context, rid string, r *entity.JobResults) error {
	return c.post(ctx, "fail", rid, r)
}

func (c *Client) post(ctx context.Context, endpoint, rid string, p any) error {
	u, err := buildURL(c.baseURL, endpoint, rid)
	if err != nil {
		log := logger.FromContext(ctx)
		log.Error().Err(err).Msg("failed to build url")
		return fmt.Errorf("invalid healthcheck URL: %w", err)
	}

	body, err := json.Marshal(p)
	if err != nil {
		log := logger.FromContext(ctx)
		log.Error().Err(err).Msg("failed to marshal payload")
		return fmt.Errorf("marshal payload: %w", err)
	}

	err = withRetry(ctx, func() error {
		return c.doPost(ctx, u, body)
	})
	if err != nil {
		log := logger.FromContext(ctx)
		log.Error().Err(err).Msg("failed to do post")
		return err
	}

	return nil
}

// buildURL builds https://hc-ping.com/{uuid}[/{slug}][/endpoint][?rid=...]
func buildURL(base, endpoint, rid string) (string, error) {
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
