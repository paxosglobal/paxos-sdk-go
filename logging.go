package paxos

import (
	"log/slog"
	"os"
	"strings"
)

// WithLogger sets a structured logger for debug output.
func WithLogger(logger *slog.Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}

// initDefaultLogger sets up the default logger from PAXOS_LOG_LEVEL env var.
func (c *Client) initDefaultLogger() {
	if c.logger == nil {
		if strings.EqualFold(os.Getenv("PAXOS_LOG_LEVEL"), "debug") {
			c.logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
		}
	}
}

// logDebug logs a debug message if the logger is configured.
func (c *Client) logDebug(msg string, args ...any) {
	if c.logger != nil {
		c.logger.Debug(msg, args...)
	}
}
