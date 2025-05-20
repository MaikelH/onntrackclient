package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	jimi "jimi-platform-client"
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
	client, err := jimi.NewClient(
		apiKey,
		jimi.WithBaseURL("https://api.jimi-platform.com"),
		jimi.WithLogger(logger),
	)
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

	// Example: Get current location of a device
	if len(devices) > 0 {
		deviceID := devices[0].ID
		location, resp, err := client.Tracking.GetCurrentLocation(ctx, deviceID)
		if err != nil {
			logger.Error("Failed to get device location",
				slog.String("error", err.Error()),
				slog.Int("status", resp.StatusCode),
			)
		} else {
			fmt.Printf("\nCurrent location of device %s:\n", deviceID)
			fmt.Printf("- Latitude: %f\n", location.Latitude)
			fmt.Printf("- Longitude: %f\n", location.Longitude)
			fmt.Printf("- Timestamp: %s\n", location.Timestamp.Format(time.RFC3339))
			if location.Address != "" {
				fmt.Printf("- Address: %s\n", location.Address)
			}
		}
	}

	// Example: List alerts
	alerts, resp, err := client.Alerts.List(ctx, nil)
	if err != nil {
		logger.Error("Failed to list alerts",
			slog.String("error", err.Error()),
			slog.Int("status", resp.StatusCode),
		)
	} else {
		fmt.Println("\nRecent alerts:")
		for _, alert := range alerts {
			fmt.Printf("- [%s] %s: %s\n",
				alert.Timestamp.Format(time.RFC3339),
				alert.Type,
				alert.Message,
			)
		}
	}

	// Example: List users
	users, resp, err := client.Users.List(ctx, nil)
	if err != nil {
		logger.Error("Failed to list users",
			slog.String("error", err.Error()),
			slog.Int("status", resp.StatusCode),
		)
	} else {
		fmt.Println("\nUsers:")
		for _, user := range users {
			fmt.Printf("- %s %s (%s): %s\n",
				user.FirstName,
				user.LastName,
				user.Email,
				user.Role,
			)
		}
	}
}
