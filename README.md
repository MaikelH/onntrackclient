# Onntrack Platform Client

A Go client library for the Onntrack tracking dashboard REST API.

## Installation

```bash
go get github.com/MaikelH/onntrackclient
```

## Usage

### Creating a client with an API key

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/MaikelH/onntrackclient"
)

func main() {
    // Create a new client with an API key
    client, err := jimi.NewClient("your-api-key")
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Use the client to make API calls
    ctx := context.Background()
    // Example API call
    // result, err := client.SomeEndpoint(ctx, params)

    // Handle results
    // ...
}
```

### Authentication with username and password

```go
package main

import (
    "context"
    "fmt"
    "log"

    jimi "github.com/MaikelH/onntrackclient"
)

func main() {
    // Create a new client without an API key
    client, err := jimi.NewClient("")
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Create a login request
    loginReq := &jimi.LoginRequest{
        Account:   "your-email@example.com",
        Password:  "your-password",
        Language:  "en",
        ValidCode: "",
        NodeID:    "",
    }

    // Login to get an authentication token
    ctx := context.Background()
    loginResp, _, err := client.Auth.Login(ctx, loginReq)
    if err != nil {
        log.Fatalf("Failed to login: %v", err)
    }

    if !loginResp.Success {
        log.Fatalf("Login failed: %s", loginResp.Message)
    }

    fmt.Println("Login successful!")

    // The client is now authenticated and can be used to make API calls
    // Example API call
    // result, err := client.SomeEndpoint(ctx, params)
}
```

## Features

- Simple, idiomatic Go API
- Authentication with username/password or API key
- Context support for cancellation and timeouts
- Configurable HTTP client
- Comprehensive error handling
- Structured logging with slog

## Documentation

For detailed documentation, see the [GoDoc](https://pkg.go.dev/github.com/MaikelH/onntrackclient).

## License

MIT
