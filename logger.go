// Package jimi provides a client for the Jimi tracking dashboard REST API.
package onntrackclient

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

// LoggingTransport is an http.RoundTripper that logs requests and responses.
type LoggingTransport struct {
	Transport http.RoundTripper
	Logger    *slog.Logger
}

// RoundTrip implements the http.RoundTripper interface.
func (t *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	// Log the request
	t.Logger.InfoContext(req.Context(), "API request",
		slog.String("method", req.Method),
		slog.String("url", req.URL.String()),
	)

	// Execute the request
	resp, err := t.Transport.RoundTrip(req)

	// Calculate duration
	duration := time.Since(start)

	if err != nil {
		// Log the error
		t.Logger.ErrorContext(req.Context(), "API request failed",
			slog.String("method", req.Method),
			slog.String("url", req.URL.String()),
			slog.Duration("duration", duration),
			slog.String("error", err.Error()),
		)
		return resp, err
	}

	// Log the response
	t.Logger.InfoContext(req.Context(), "API response",
		slog.String("method", req.Method),
		slog.String("url", req.URL.String()),
		slog.Int("status", resp.StatusCode),
		slog.Duration("duration", duration),
	)

	return resp, err
}

// WithLogger returns a ClientOption that sets a logger for the client.
func WithLogger(logger *slog.Logger) ClientOption {
	return func(c *Client) error {
		// Create a logging transport that wraps the existing transport
		transport := c.HTTPClient.Transport
		if transport == nil {
			transport = http.DefaultTransport
		}

		c.HTTPClient.Transport = &LoggingTransport{
			Transport: transport,
			Logger:    logger,
		}

		return nil
	}
}

// SetLogger sets a logger for the client.
func (c *Client) SetLogger(logger *slog.Logger) {
	// Create a logging transport that wraps the existing transport
	transport := c.HTTPClient.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	c.HTTPClient.Transport = &LoggingTransport{
		Transport: transport,
		Logger:    logger,
	}
}

// ContextWithLogger returns a new context with the logger attached.
func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// LoggerFromContext returns the logger from the context.
func LoggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// contextKey is a private type for context keys.
type contextKey int

const (
	// loggerKey is the context key for the logger.
	loggerKey contextKey = iota
)
