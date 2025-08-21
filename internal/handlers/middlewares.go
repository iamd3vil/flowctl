package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/cvhariharan/flowctl/internal/core/models"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := h.sessMgr.Acquire(nil, c, c)
		if err != nil {
			return wrapError(ErrAuthenticationFailed, "could not get user session", err, nil)
		}

		user, err := sess.Get("user")
		if err != nil {
			return wrapError(ErrAuthenticationFailed, "could not get user details", err, nil)
		}

		method, err := sess.String(sess.Get("method"))
		if err != nil {
			return wrapError(ErrAuthenticationFailed, "could not get login method", err, nil)
		}

		// if using oidc, validate the token to check if they have not expired
		if method == "oidc" {
			rawIDToken, err := sess.Get("id_token")
			if err != nil || rawIDToken == nil {
				return wrapError(ErrAuthenticationFailed, "could not get id token", err, nil)
			}

			_, err = h.authconfig.verifier.Verify(context.Background(), rawIDToken.(string))
			if err != nil {
				sess.Delete("method")
				sess.Delete("id_token")
				sess.Delete("user")
				return wrapError(ErrAuthenticationFailed, "could not verify id token", err, nil)
			}
		}

		var userInfo models.UserInfo
		userBytes, err := json.Marshal(user)
		if err != nil {
			return wrapError(ErrAuthenticationFailed, "could not get user details", err, nil)
		}

		if err := json.NewDecoder(bytes.NewBuffer(userBytes)).Decode(&userInfo); err != nil {
			c.Logger().Error(err)
			return wrapError(ErrAuthenticationFailed, "could not get user details", err, nil)
		}
		c.Set("user", userInfo)

		return next(c)
	}
}

func (h *Handler) AuthorizeForRole(expectedRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userInfo, err := h.getUserInfo(c)
			if err != nil {
				return wrapError(ErrAuthenticationFailed, "could not get user details", err, nil)
			}

			if userInfo.Role == expectedRole {
				return next(c)
			}

			return wrapError(ErrUnauthorized, "unauthorized", nil, nil)
		}
	}
}

func (h *Handler) AuthorizeNamespaceAction(resource models.Resource, action models.RBACAction) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, err := h.getUserInfo(c)
			if err != nil {
				return wrapError(ErrAuthenticationFailed, "could not get user details", err, nil)
			}

			namespaceID, ok := c.Get("namespace").(string)
			if !ok {
				return wrapError(ErrRequiredFieldMissing, "could not get namespace", nil, nil)
			}

			allowed, err := h.co.CheckPermission(c.Request().Context(), user.ID, namespaceID, resource, action)
			if err != nil {
				return wrapError(ErrOperationFailed, "could not check permissions", err, nil)
			}

			if !allowed {
				return wrapError(ErrForbidden, "insufficient permissions", nil, nil)
			}

			return next(c)
		}
	}
}

// Updated NamespaceMiddleware for simpler access check
func (h *Handler) NamespaceMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		namespace := c.Param("namespace")
		if namespace == "" {
			return wrapError(ErrRequiredFieldMissing, "namespace cannot be empty", nil, nil)
		}

		ns, err := h.co.GetNamespaceByName(c.Request().Context(), namespace)
		if err != nil {
			return wrapError(ErrResourceNotFound, "could not find namespace", err, nil)
		}

		user, err := h.getUserInfo(c)
		if err != nil {
			return wrapError(ErrAuthenticationFailed, "could not get user details", err, nil)
		}

		// Basic access check - user must have at least view permission
		hasAccess, err := h.co.CheckPermission(c.Request().Context(), user.ID, ns.ID, models.ResourceFlow, models.RBACActionView)
		if err != nil {
			return wrapError(ErrOperationFailed, "could not check namespace access", err, nil)
		}

		if !hasAccess {
			return wrapError(ErrForbidden, "user does not have access to this namespace", nil, nil)
		}

		c.Set("namespace", ns.ID)
		return next(c)
	}
}

func (h *Handler) getUserInfo(c echo.Context) (models.UserInfo, error) {
	sess, err := h.sessMgr.Acquire(nil, c, c)
	if err != nil {
		return models.UserInfo{}, err
	}

	user, err := sess.Get("user")
	if err != nil {
		return models.UserInfo{}, err
	}

	if user == nil {
		err := fmt.Errorf("user session is empty")
		return models.UserInfo{}, err
	}

	var userInfo models.UserInfo
	userBytes, err := json.Marshal(user)
	if err != nil {
		return models.UserInfo{}, err
	}

	if err := json.NewDecoder(bytes.NewBuffer(userBytes)).Decode(&userInfo); err != nil {
		return models.UserInfo{}, err
	}

	return userInfo, nil
}
