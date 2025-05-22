// Package jimi provides a client for the Jimi tracking dashboard REST API.
package onntrackclient

import (
	"context"
	"net/http"
)

// AuthService handles communication with the authentication related
// methods of the Jimi API.
type AuthService service

// LoginRequest represents a request to login to the Jimi platform.
type LoginRequest struct {
	Account   string `json:"account"`
	Password  string `json:"password"`
	Language  string `json:"language"`
	ValidCode string `json:"validCode"`
	NodeID    string `json:"nodeId"`
}

// LoginResponse represents the response from a login request.
type LoginResponse struct {
	OK   bool `json:"ok"`
	Data struct {
		UpgradeTips bool   `json:"upgradeTips"`
		Token       string `json:"token"`
	} `json:"data"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Login authenticates a user with the Jimi platform.
//
// Endpoint: homepage/login
func (s *AuthService) Login(ctx context.Context, loginReq *LoginRequest) (*LoginResponse, *http.Response, error) {
	u := "homepage/login"

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, loginReq)
	if err != nil {
		return nil, nil, err
	}

	// Set the must header to true as specified in the curl command
	req.Header.Set("must", "true")

	loginResp := new(LoginResponse)
	resp, err := s.client.Do(req, loginResp)
	if err != nil {
		return nil, resp, err
	}

	// If login is successful, update the client's API key with the token
	if loginResp.OK && loginResp.Data.Token != "" {
		s.client.APIKey = loginResp.Data.Token
	}

	return loginResp, resp, nil
}
