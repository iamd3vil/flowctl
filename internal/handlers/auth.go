package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/cvhariharan/flowctl/internal/core/models"
	"github.com/labstack/echo/v4"
	"github.com/zerodha/simplesessions/v3"
	"golang.org/x/oauth2"
)

const (
	SessionTimeout     = 2 * time.Hour
	RedirectPath       = "/auth/callback"
	LoginPath          = "/login"
	RedirectAfterLogin = "/"
)

type TokenData struct {
	Provider   string
	RawIDToken string
}

// GetEncoded returns the token data in base64 encoding
func (s *TokenData) GetEncoded() (string, error) {
	j, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("error encoding token data: %w", err)
	}

	return base64.URLEncoding.EncodeToString(j), nil
}

// Decode converts and unmarshals encoded token data
func (s *TokenData) Decode(e string) error {
	j, err := base64.URLEncoding.DecodeString(e)
	if err != nil {
		return fmt.Errorf("error decoding token data: %w", err)
	}

	if err := json.Unmarshal(j, &s); err != nil {
		return fmt.Errorf("invalid state: %w", err)
	}
	return nil
}

type StateData struct {
	Provider string
	Nonce    string
}

// GetEncoded returns the state data in base64 encoding
func (s *StateData) GetEncoded() (string, error) {
	j, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("error encoding state: %w", err)
	}

	return base64.URLEncoding.EncodeToString(j), nil
}

// Decode converts and unmarshals encoded state
func (s *StateData) Decode(e string) error {
	j, err := base64.URLEncoding.DecodeString(e)
	if err != nil {
		return fmt.Errorf("error decoding state: %w", err)
	}

	if err := json.Unmarshal(j, &s); err != nil {
		return fmt.Errorf("invalid state: %w", err)
	}
	return nil
}

// isSafeRedirect determines if the redirect URL is safe
// Must start with '/' but not with '//' or '/\'.
func isSafeRedirect(u string) bool {
	return len(u) > 0 && u[0] == '/' && (len(u) == 1 || (u[1] != '/' && u[1] != '\\'))
}

func (h *Handler) initOIDC() error {
	for _, oauthConfig := range h.config.OIDC {
		provider, err := oidc.NewProvider(context.Background(), oauthConfig.Issuer)
		if err != nil {
			return fmt.Errorf("could not initialize new OIDC provider client: %w", err)
		}

		redirectURL, err := url.JoinPath(h.config.App.RootURL, RedirectPath)
		if err != nil {
			return fmt.Errorf("failed to create redirect URL: %w", err)
		}

		if oauthConfig.RedirectURL != "" {
			redirectURL = oauthConfig.RedirectURL
		}

		endpoint := provider.Endpoint()
		if oauthConfig.AuthURL != "" {
			endpoint.AuthURL = oauthConfig.AuthURL
		}
		if oauthConfig.TokenURL != "" {
			endpoint.TokenURL = oauthConfig.TokenURL
		}

		oauth2Config := &oauth2.Config{
			ClientID:     oauthConfig.ClientID,
			ClientSecret: oauthConfig.ClientSecret,
			RedirectURL:  redirectURL,
			Endpoint:     endpoint,
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		}

		verifier := provider.Verifier(&oidc.Config{
			ClientID: oauthConfig.ClientID,
		})

		h.authconfig[oauthConfig.Name] = OIDCAuthConfig{
			provider:     provider,
			verifier:     verifier,
			oauth2Config: oauth2Config,
		}

	}

	return nil
}

func (h *Handler) HandleLoginPage(c echo.Context) error {
	var req AuthReq
	if err := c.Bind(&req); err != nil {
		return wrapError(ErrInvalidInput, "invalid request", err, nil)
	}

	sess, err := h.sessMgr.Acquire(nil, c, c)

	if err == simplesessions.ErrInvalidSession {
		sess, err = h.sessMgr.NewSession(c, c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	if req.Username == "" || req.Password == "" {
		return wrapError(ErrRequiredFieldMissing, "username or password cannot be empty", fmt.Errorf("username or password cannot be empty"), nil)
	}

	user, err := h.co.GetUserByUsernameWithGroups(c.Request().Context(), req.Username)
	if err != nil {
		return wrapError(ErrAuthenticationFailed, "could not authenticate user", err, nil)
	}

	// not using password based login
	if user.LoginType != models.StandardLoginType {
		return wrapError(ErrAuthenticationFailed, "invalid authentication method", fmt.Errorf("invalid authentication method for user: %s", user.ID), nil)
	}

	if err := user.CheckPassword(req.Password); err != nil {
		return wrapError(ErrInvalidCredentials, "invalid credentials", err, nil)
	}

	sess.Set("method", "password")

	var groups []string
	for _, v := range user.Groups {
		groups = append(groups, v.ID)
	}

	sess.Set("user", user.ToUserInfo())

	redirectAfterLogin := RedirectAfterLogin
	if redirectURL := c.QueryParam("redirect_url"); redirectURL != "" && isSafeRedirect(redirectURL) {
		redirectAfterLogin = redirectURL
	}

	c.Response().Header().Set("x-redirect", redirectAfterLogin)
	return c.NoContent(http.StatusOK)
}

func (h *Handler) HandleOIDCLogin(c echo.Context) error {
	provider := c.Param("provider")
	if provider == "" {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("oidc provider cannot be empty"))
	}

	sess, err := h.sessMgr.Acquire(nil, c, c)

	if err == simplesessions.ErrInvalidSession {
		sess, err = h.sessMgr.NewSession(c, c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	nonce, err := generateRandomState()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not generate a login state")
	}

	state := StateData{
		Provider: provider,
		Nonce:    nonce,
	}

	encodedState, err := state.GetEncoded()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	sess.Set("state", encodedState)

	if redirectURL := c.QueryParam("redirect_url"); redirectURL != "" && isSafeRedirect(redirectURL) {
		sess.Set("redirect_url", redirectURL)
	}

	authURL := h.authconfig[provider].oauth2Config.AuthCodeURL(encodedState)
	return c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random state: %w", err)
	}

	state := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
	return state, nil
}

func (h *Handler) HandleAuthCallback(c echo.Context) error {
	sess, err := h.sessMgr.Acquire(nil, c, c)
	if err != nil {
		return wrapError(ErrInvalidInput, "session does not exist", err, nil)
	}

	rawSessionState, err := sess.Get("state")
	if err != nil {
		return wrapError(ErrInvalidInput, "state not found", err, nil)
	}
	var sessionState StateData
	if err := sessionState.Decode(rawSessionState.(string)); err != nil {
		return wrapError(ErrInternalError, "invalid session state", err, nil)
	}

	var callbackState StateData
	if err := callbackState.Decode(c.QueryParam("state")); err != nil {
		return wrapError(ErrInternalError, "invalid callback state", err, nil)
	}

	if sessionState.Nonce != callbackState.Nonce || sessionState.Provider != callbackState.Provider {
		return wrapError(ErrInvalidInput, "invalid state parameter", nil, nil)
	}

	token, err := h.authconfig[sessionState.Provider].oauth2Config.Exchange(context.Background(), c.QueryParam("code"))
	if err != nil {
		return wrapError(ErrOperationFailed, "failed to exchange token", err, nil)
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return wrapError(ErrOperationFailed, "no id_token in token response", nil, nil)
	}

	idToken, err := h.authconfig[sessionState.Provider].verifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		return wrapError(ErrOperationFailed, "failed to verify ID token", err, nil)
	}

	var claims struct {
		Email  string   `json:"email"`
		Name   string   `json:"name"`
		Groups []string `json:"groups"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return wrapError(ErrOperationFailed, "failed to parse claims", err, nil)
	}

	user, err := h.co.GetUserByUsernameWithGroups(c.Request().Context(), claims.Email)
	if err != nil {
		return wrapError(ErrForbidden, "user does not exist in flowctl", err, nil)
	}

	sess.Set("method", "oidc")

	td := TokenData{
		Provider:   sessionState.Provider,
		RawIDToken: rawIDToken,
	}
	tokenData, err := td.GetEncoded()
	if err != nil {
		return wrapError(ErrInternalError, err.Error(), err, nil)
	}

	sess.Set("id_token", tokenData)

	sess.Set("user", user.ToUserInfo())

	redirectAfterLogin := RedirectAfterLogin
	if redirectURL, err := sess.Get("redirect_url"); err == nil && redirectURL != nil {
		if url, ok := redirectURL.(string); ok && isSafeRedirect(url) {
			redirectAfterLogin = url
		}
	}

	return c.Redirect(http.StatusTemporaryRedirect, redirectAfterLogin)
}

func (h *Handler) HandleLogout(c echo.Context) error {
	sess, err := h.sessMgr.Acquire(nil, c, c)
	if err != nil {
		return c.NoContent(http.StatusOK)
	}

	err = sess.Destroy()
	if err != nil {
		return wrapError(ErrInternalError, "could not destroy session", err, nil)
	}

	return c.NoContent(http.StatusOK)
}

// HandleGetCasbinPermissions returns policies for casbin.js
// It doesn't use the subject query param set by the frontend but instead uses the session user
func (h *Handler) HandleGetCasbinPermissions(c echo.Context) error {
	user, err := h.getUserInfo(c)
	if err != nil {
		return wrapError(ErrAuthenticationFailed, "could not get user details", err, nil)
	}

	p, err := h.co.GetPermissionsForUser(user.ID)
	if err != nil {
		return wrapError(ErrResourceNotFound, "could not get permissions for user", err, nil)
	}

	return c.JSON(http.StatusOK, struct {
		Data string `json:"data"`
	}{
		Data: p,
	})
}

func (h *Handler) HandleGetSSOProviders(c echo.Context) error {
	var providers []SSOProvider

	for _, v := range h.config.OIDC {
		label := v.Label
		if label == "" {
			label = fmt.Sprintf("Sign in with %s", v.Name)
		}

		providers = append(providers, SSOProvider{
			ID:    v.Name,
			Label: label,
		})
	}

	return c.JSON(http.StatusOK, providers)
}
