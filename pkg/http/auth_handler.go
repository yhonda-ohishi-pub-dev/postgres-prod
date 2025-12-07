package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
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
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(
	googleClient *auth.GoogleOAuthClient,
	lineClient *auth.LineOAuthClient,
	jwtService *auth.JWTService,
	appUserRepo *repository.AppUserRepository,
	oauthRepo *repository.OAuthAccountRepository,
) *AuthHandler {
	return &AuthHandler{
		googleClient: googleClient,
		lineClient:   lineClient,
		jwtService:   jwtService,
		appUserRepo:  appUserRepo,
		oauthRepo:    oauthRepo,
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
		http.Error(w, "missing code parameter", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Exchange code for tokens
	tokenResp, err := h.googleClient.ExchangeCode(ctx, code)
	if err != nil {
		log.Printf("Failed to exchange Google code: %v", err)
		http.Error(w, "failed to exchange code", http.StatusInternalServerError)
		return
	}

	// Get user info
	userInfo, err := h.googleClient.GetUserInfo(ctx, tokenResp.AccessToken)
	if err != nil {
		log.Printf("Failed to get Google user info: %v", err)
		http.Error(w, "failed to get user info", http.StatusInternalServerError)
		return
	}

	// Find or create user and generate JWT
	authResp, err := h.processOAuthLogin(ctx, "google", userInfo.ID, &userInfo.Email, userInfo.Name, &userInfo.Picture, tokenResp.AccessToken, tokenResp.RefreshToken, tokenResp.ExpiresIn)
	if err != nil {
		log.Printf("Failed to process Google login: %v", err)
		http.Error(w, "failed to process login", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authResp)
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
		http.Error(w, "missing code parameter", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Exchange code for tokens
	tokenResp, err := h.lineClient.ExchangeCode(ctx, code)
	if err != nil {
		log.Printf("Failed to exchange LINE code: %v", err)
		http.Error(w, "failed to exchange code", http.StatusInternalServerError)
		return
	}

	// Get user info
	userInfo, err := h.lineClient.GetUserInfo(ctx, tokenResp.AccessToken)
	if err != nil {
		log.Printf("Failed to get LINE user info: %v", err)
		http.Error(w, "failed to get user info", http.StatusInternalServerError)
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
		http.Error(w, "failed to process login", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authResp)
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
