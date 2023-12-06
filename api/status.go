package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

const (
	endpoint string = "https://status.vultr.com"
)

type Client struct {
	client      *http.Client
	rateLimiter *rate.Limiter
}

func NewClient() *Client {
	return &Client{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		rateLimiter: rate.NewLimiter(rate.Every(time.Second), 5),
	}
}
func (c *Client) Do(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
	rqst, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		msg := "unable to create Vultr status client"
		slog.Info(msg)
		return nil, fmt.Errorf(msg)
	}

	// Add Content-Type header
	rqst.Header.Set(
		"Content-Type",
		"application/json",
	)

	// Add Accept header
	rqst.Header.Set(
		"Accept",
		"application/json",
	)

	// Apply rate limiter to request
	if err := c.rateLimiter.Wait(ctx); err != nil {
		msg := "API rate limit exceeed"
		slog.Info(msg)
		return nil, fmt.Errorf(msg)
	}

	resp, err := c.client.Do(rqst)
	if err != nil {
		msg := "unable to perform HTTP request"
		slog.Info(msg,
			"method", method,
			"url", url,
		)
		return nil, fmt.Errorf("%s\n%+v", msg, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		msg := "unable to read response body"
		slog.Info(msg)
		return nil, fmt.Errorf(msg)
	}

	if resp.StatusCode != http.StatusOK {
		msg := "unsuccessful"
		slog.Info(msg,
			"status", resp.StatusCode,
		)
		return nil, fmt.Errorf("%s [%d]", msg, resp.StatusCode)
	}

	return respBody, nil
}
func (c *Client) Alerts() (*Alerts, error) {
	ctx := context.Background()
	method := http.MethodGet
	url := fmt.Sprintf("%s/alerts.json", endpoint)

	alerts := &Alerts{}

	resp, err := c.Do(ctx, method, url, nil)
	if err != nil {
		msg := "unable to get alerts"
		slog.Info(msg)
		return alerts, fmt.Errorf(msg)
	}

	if err := json.Unmarshal(resp, alerts); err != nil {
		msg := "unable to unmarshal result as Alerts"
		slog.Info(msg,
			"response", resp,
		)
		return alerts, fmt.Errorf(msg)
	}

	slog.Info("Returning",
		"alerts", len(alerts.ServiceAlerts),
	)

	return alerts, nil
}
func (c *Client) Status() (*Status, error) {
	ctx := context.Background()
	method := http.MethodGet
	url := fmt.Sprintf("%s/status.json", endpoint)

	status := &Status{}

	resp, err := c.Do(ctx, method, url, nil)
	if err != nil {
		msg := "unable to get status"
		slog.Info(msg)
		return status, fmt.Errorf(msg)
	}

	if err := json.Unmarshal(resp, status); err != nil {
		msg := "unable to unmarshal result as Status"
		slog.Info(msg,
			"response", resp,
		)
		return status, fmt.Errorf(msg)
	}

	slog.Info("Returning",
		"alerts", len(status.ServiceAlerts),
		"regions", len(status.Regions),
	)

	return status, nil
}
