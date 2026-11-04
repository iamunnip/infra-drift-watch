package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// request never returns response bodies or URLs in errors: both may contain secrets.
func (e *Engine) request(ctx context.Context, method, address, auth string, payload, output any, retries int) error {
	var data []byte
	if payload != nil {
		var err error
		data, err = json.Marshal(payload)
		if err != nil {
			return errors.New("cannot encode delivery payload")
		}
	}
	client := http.Client{Timeout: 30 * time.Second}
	if e.Client != nil {
		client = *e.Client
	}
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return errors.New("redirects are disabled") }
	for attempt := 0; ; attempt++ {
		req, err := http.NewRequestWithContext(ctx, method, address, bytes.NewReader(data))
		if err != nil {
			return errors.New("invalid delivery request")
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", "infra-drift-watch")
		if auth != "" {
			req.Header.Set("Authorization", auth)
		}
		response, err := client.Do(req)
		if err != nil {
			return errors.New("delivery transport failed; remote acceptance may be unknown")
		}
		body, readErr := io.ReadAll(io.LimitReader(response.Body, (2<<20)+1))
		response.Body.Close()
		if readErr != nil || len(body) > 2<<20 {
			return errors.New("delivery response unreadable or too large")
		}
		if response.StatusCode >= 200 && response.StatusCode < 300 {
			if output != nil && json.Unmarshal(body, output) != nil {
				return errors.New("invalid delivery response")
			}
			return nil
		}
		// An explicit rate-limit rejection is safe to retry. Ambiguous mutations are not.
		retry := response.StatusCode == 429 || (method == http.MethodGet && response.StatusCode >= 500)
		if !retry || attempt >= retries {
			return fmt.Errorf("delivery HTTP status %d", response.StatusCode)
		}
		wait := time.Duration(1<<attempt) * time.Second
		if value := response.Header.Get("Retry-After"); value != "" {
			if seconds, err := strconv.Atoi(value); err == nil {
				wait = time.Duration(seconds) * time.Second
			} else if date, err := http.ParseTime(value); err == nil {
				wait = time.Until(date)
			}
		}
		if wait < 0 {
			wait = 0
		}
		if wait > 30*time.Second {
			return errors.New("delivery rate limited; retry on a later scan")
		}
		timer := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			timer.Stop()
			return errors.New("delivery cancelled")
		case <-timer.C:
		}
	}
}
