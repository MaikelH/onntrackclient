package main

import (
	"context"
	"fmt"
	"github.com/MaikelH/onntrackclient"
	"log/slog"
	"os"
	"time"
)

func main() {
	// Create a logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Get API key from environment variable
	apiKey := os.Getenv("JIMI_API_KEY")
	if apiKey == "" {
		logger.Error("JIMI_API_KEY environment variable not set")
		os.Exit(1)
	}

	// Create a new client with options
	client, err := onntrackclient.NewClient(onntrackclient.WithBaseURL("https://api.jimi-platform.com"), onntrackclient.WithLogger(logger))
	if err != nil {
		logger.Error("Failed to create client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Example: List devices
	devices, resp, err := client.Devices.List(ctx, nil)
	if err != nil {
		logger.Error("Failed to list devices",
			slog.String("error", err.Error()),
			slog.Int("status", resp.StatusCode),
		)
		os.Exit(1)
	}

	// Print devices
	fmt.Println("Devices:")
	for _, device := range devices {
		fmt.Printf("- %s (%s): %s\n", device.Name, device.ID, device.Status)
	}
}
