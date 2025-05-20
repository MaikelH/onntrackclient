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

	// Create a new client without an API key
	client, err := onntrackclient.NewClient(
		"", // No API key needed for login
		onntrackclient.WithBaseURL("https://platform.onntrack.nl/v3/new/"),
		onntrackclient.WithLogger(logger),
	)
	if err != nil {
		logger.Error("Failed to create client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get credentials from environment variables
	account := os.Getenv("ONNTRACK_ACCOUNT")
	password := os.Getenv("ONNTRACK_PASSWORD")
	if account == "" || password == "" {
		logger.Error("ONNTRACK_ACCOUNT and ONNTRACK_PASSWORD environment variables must be set")
		os.Exit(1)
	}

	// Create login request
	loginReq := &onntrackclient.LoginRequest{
		Account:   account,
		Password:  password,
		Language:  "en",
		ValidCode: "",
		NodeID:    "",
	}

	// Call the Login method
	loginResp, resp, err := client.Auth.Login(ctx, loginReq)
	if err != nil {
		logger.Error("Failed to login",
			slog.String("error", err.Error()),
			slog.Int("status", resp.StatusCode),
		)
		os.Exit(1)
	}

	// Check login response
	if !loginResp.Success {
		logger.Error("Login failed",
			slog.String("message", loginResp.Message),
		)
		os.Exit(1)
	}

	logger.Info("Login successful",
		slog.String("token", loginResp.Token),
		slog.String("message", loginResp.Message),
	)

	// The client's API key is now set to the token from the login response
	// We can now use the client for authenticated requests

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

	// Note: In a real application, you might want to store the token
	// and reuse it for future sessions instead of logging in each time
}
