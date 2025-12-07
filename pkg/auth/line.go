package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var (
	ErrLineAuthFailed = errors.New("line authentication failed")
)

// LineUserInfo represents the user info from LINE
type LineUserInfo struct {
	UserID        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	PictureURL    string `json:"pictureUrl,omitempty"`
	StatusMessage string `json:"statusMessage,omitempty"`
}

// LineIDTokenPayload represents the ID token payload from LINE
type LineIDTokenPayload struct {
	Iss     string `json:"iss"`
	Sub     string `json:"sub"` // User ID
	Aud     string `json:"aud"`
	Exp     int64  `json:"exp"`
	Iat     int64  `json:"iat"`
	Name    string `json:"name,omitempty"`
	Picture string `json:"picture,omitempty"`
	Email   string `json:"email,omitempty"`
}

// LineOAuthClient handles LINE OAuth2 operations
type LineOAuthClient struct {
	channelID     string
	channelSecret string
	redirectURI   string
	httpClient    *http.Client
}

// NewLineOAuthClient creates a new LINE OAuth client
func NewLineOAuthClient(channelID, channelSecret, redirectURI string) *LineOAuthClient {
	return &LineOAuthClient{
		channelID:     channelID,
		channelSecret: channelSecret,
		redirectURI:   redirectURI,
		httpClient:    &http.Client{},
	}
}

// LineTokenResponse represents the token response from LINE
type LineTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	IDToken      string `json:"id_token"`
	Scope        string `json:"scope"`
}

// ExchangeCode exchanges an authorization code for tokens
func (c *LineOAuthClient) ExchangeCode(ctx context.Context, code string) (*LineTokenResponse, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", c.channelID)
	data.Set("client_secret", c.channelSecret)
	data.Set("redirect_uri", c.redirectURI)
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.line.me/oauth2/v2.1/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%w: %s", ErrLineAuthFailed, string(body))
	}

	var tokenResp LineTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResp, nil
}

// GetUserInfo retrieves user info using the access token
func (c *LineOAuthClient) GetUserInfo(ctx context.Context, accessToken string) (*LineUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.line.me/v2/profile", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%w: %s", ErrLineAuthFailed, string(body))
	}

	var userInfo LineUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}

// VerifyIDToken verifies the ID token and returns the payload
func (c *LineOAuthClient) VerifyIDToken(ctx context.Context, idToken string) (*LineIDTokenPayload, error) {
	data := url.Values{}
	data.Set("id_token", idToken)
	data.Set("client_id", c.channelID)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.line.me/oauth2/v2.1/verify", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to verify id token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%w: %s", ErrLineAuthFailed, string(body))
	}

	var payload LineIDTokenPayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to decode id token payload: %w", err)
	}

	return &payload, nil
}

// RefreshAccessToken refreshes the access token using a refresh token
func (c *LineOAuthClient) RefreshAccessToken(ctx context.Context, refreshToken string) (*LineTokenResponse, error) {
	data := url.Values{}
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", c.channelID)
	data.Set("client_secret", c.channelSecret)
	data.Set("grant_type", "refresh_token")

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.line.me/oauth2/v2.1/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%w: %s", ErrLineAuthFailed, string(body))
	}

	var tokenResp LineTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResp, nil
}

// GetAuthURL returns the LINE OAuth authorization URL
func (c *LineOAuthClient) GetAuthURL(state string) string {
	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", c.channelID)
	params.Set("redirect_uri", c.redirectURI)
	params.Set("scope", "profile openid email")
	params.Set("state", state)

	return "https://access.line.me/oauth2/v2.1/authorize?" + params.Encode()
}
