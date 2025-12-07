package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/auth"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// AuthHandler handles HTTP auth endpoints
type AuthHandler struct {
	googleClient *auth.GoogleOAuthClient
	lineClient   *auth.LineOAuthClient
	jwtService   *auth.JWTService
	appUserRepo  *repository.AppUserRepository
	oauthRepo    *repository.OAuthAccountRepository
	frontendURL  string
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(
	googleClient *auth.GoogleOAuthClient,
	lineClient *auth.LineOAuthClient,
	jwtService *auth.JWTService,
	appUserRepo *repository.AppUserRepository,
	oauthRepo *repository.OAuthAccountRepository,
	frontendURL string,
) *AuthHandler {
	return &AuthHandler{
		googleClient: googleClient,
		lineClient:   lineClient,
		jwtService:   jwtService,
		appUserRepo:  appUserRepo,
		oauthRepo:    oauthRepo,
		frontendURL:  frontendURL,
	}
}

// HandleGoogleAuth redirects to Google OAuth
func (h *AuthHandler) HandleGoogleAuth(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state == "" {
		state = "default"
	}
	url := h.googleClient.GetAuthURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

// HandleGoogleCallback handles Google OAuth callback
func (h *AuthHandler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		h.redirectWithError(w, r, "missing code parameter")
		return
	}

	ctx := r.Context()

	// Exchange code for tokens
	tokenResp, err := h.googleClient.ExchangeCode(ctx, code)
	if err != nil {
		log.Printf("Failed to exchange Google code: %v", err)
		h.redirectWithError(w, r, "failed to exchange code")
		return
	}

	// Get user info
	userInfo, err := h.googleClient.GetUserInfo(ctx, tokenResp.AccessToken)
	if err != nil {
		log.Printf("Failed to get Google user info: %v", err)
		h.redirectWithError(w, r, "failed to get user info")
		return
	}

	// Find or create user and generate JWT
	authResp, err := h.processOAuthLogin(ctx, "google", userInfo.ID, &userInfo.Email, userInfo.Name, &userInfo.Picture, tokenResp.AccessToken, tokenResp.RefreshToken, tokenResp.ExpiresIn)
	if err != nil {
		log.Printf("Failed to process Google login: %v", err)
		h.redirectWithError(w, r, "failed to process login")
		return
	}

	h.redirectWithToken(w, r, authResp)
}

// HandleLineAuth redirects to LINE OAuth
func (h *AuthHandler) HandleLineAuth(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state == "" {
		state = "default"
	}
	url := h.lineClient.GetAuthURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

// HandleLineCallback handles LINE OAuth callback
func (h *AuthHandler) HandleLineCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		h.redirectWithError(w, r, "missing code parameter")
		return
	}

	ctx := r.Context()

	// Exchange code for tokens
	tokenResp, err := h.lineClient.ExchangeCode(ctx, code)
	if err != nil {
		log.Printf("Failed to exchange LINE code: %v", err)
		h.redirectWithError(w, r, "failed to exchange code")
		return
	}

	// Get user info
	userInfo, err := h.lineClient.GetUserInfo(ctx, tokenResp.AccessToken)
	if err != nil {
		log.Printf("Failed to get LINE user info: %v", err)
		h.redirectWithError(w, r, "failed to get user info")
		return
	}

	// Try to get email from ID token
	var email *string
	if tokenResp.IDToken != "" {
		if payload, err := h.lineClient.VerifyIDToken(ctx, tokenResp.IDToken); err == nil && payload.Email != "" {
			email = &payload.Email
		}
	}

	var pictureURL *string
	if userInfo.PictureURL != "" {
		pictureURL = &userInfo.PictureURL
	}

	// Find or create user and generate JWT
	authResp, err := h.processOAuthLogin(ctx, "line", userInfo.UserID, email, userInfo.DisplayName, pictureURL, tokenResp.AccessToken, tokenResp.RefreshToken, tokenResp.ExpiresIn)
	if err != nil {
		log.Printf("Failed to process LINE login: %v", err)
		h.redirectWithError(w, r, "failed to process login")
		return
	}

	h.redirectWithToken(w, r, authResp)
}

// AuthResponse is the JSON response for auth endpoints
type AuthResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int64    `json:"expires_in"`
	User         UserInfo `json:"user"`
}

// UserInfo contains user information
type UserInfo struct {
	ID           string  `json:"id"`
	Email        *string `json:"email,omitempty"`
	DisplayName  string  `json:"display_name"`
	AvatarURL    *string `json:"avatar_url,omitempty"`
	IsSuperadmin bool    `json:"is_superadmin"`
}

func (h *AuthHandler) processOAuthLogin(ctx context.Context, provider, providerUserID string, email *string, displayName string, avatarURL *string, accessToken, refreshToken string, expiresIn int) (*AuthResponse, error) {
	tokenExpiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)

	// Try to find existing OAuth account
	oauthAccount, err := h.oauthRepo.GetByProviderAndProviderUserID(ctx, provider, providerUserID)
	var user *repository.AppUser

	if err == nil {
		// Found existing OAuth account, get the user
		user, err = h.appUserRepo.GetByID(ctx, oauthAccount.AppUserID)
		if err != nil {
			return nil, err
		}

		// Update tokens
		_, err = h.oauthRepo.UpdateTokens(ctx, oauthAccount.ID, &accessToken, &refreshToken, &tokenExpiresAt)
		if err != nil {
			return nil, err
		}
	} else {
		// Create new user and OAuth account
		user, err = h.appUserRepo.Create(ctx, email, displayName, avatarURL, false)
		if err != nil {
			return nil, err
		}

		_, err = h.oauthRepo.Create(ctx, user.ID, provider, providerUserID, email, &accessToken, &refreshToken, &tokenExpiresAt)
		if err != nil {
			_ = h.appUserRepo.Delete(ctx, user.ID)
			return nil, err
		}
	}

	// Generate JWT
	var userEmail string
	if user.Email != nil {
		userEmail = *user.Email
	}
	tokenPair, err := h.jwtService.GenerateTokenPair(user.ID, userEmail, user.DisplayName, user.IsSuperadmin)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User: UserInfo{
			ID:           user.ID,
			Email:        user.Email,
			DisplayName:  user.DisplayName,
			AvatarURL:    user.AvatarURL,
			IsSuperadmin: user.IsSuperadmin,
		},
	}, nil
}

// RegisterRoutes registers auth HTTP routes
func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/auth/google", h.HandleGoogleAuth)
	mux.HandleFunc("/auth/google/callback", h.HandleGoogleCallback)
	mux.HandleFunc("/auth/line", h.HandleLineAuth)
	mux.HandleFunc("/auth/line/callback", h.HandleLineCallback)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}

// redirectWithToken redirects to frontend with token in URL parameters
func (h *AuthHandler) redirectWithToken(w http.ResponseWriter, r *http.Request, authResp *AuthResponse) {
	if h.frontendURL == "" {
		// Fallback to JSON response if no frontend URL configured
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(authResp)
		return
	}

	// Build redirect URL with token parameters
	redirectURL, err := url.Parse(h.frontendURL)
	if err != nil {
		log.Printf("Invalid frontend URL: %v", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(authResp)
		return
	}

	// Add token as URL fragment (hash) for security - not sent to server
	// Frontend: https://example.com/auth/callback#access_token=xxx&refresh_token=xxx
	fragment := fmt.Sprintf("access_token=%s&refresh_token=%s&expires_in=%d",
		url.QueryEscape(authResp.AccessToken),
		url.QueryEscape(authResp.RefreshToken),
		authResp.ExpiresIn,
	)
	redirectURL.Fragment = fragment

	http.Redirect(w, r, redirectURL.String(), http.StatusFound)
}

// redirectWithError redirects to frontend with error message
func (h *AuthHandler) redirectWithError(w http.ResponseWriter, r *http.Request, errMsg string) {
	if h.frontendURL == "" {
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	redirectURL, err := url.Parse(h.frontendURL)
	if err != nil {
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	// Add error as URL fragment
	redirectURL.Fragment = fmt.Sprintf("error=%s", url.QueryEscape(errMsg))
	http.Redirect(w, r, redirectURL.String(), http.StatusFound)
}
