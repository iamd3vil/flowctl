package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cvhariharan/flowctl/internal/core"
	"github.com/cvhariharan/flowctl/internal/core/models"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Check for executor API key first
		executorName, err := h.authenticateExecutor(c)
		if err != nil {
			return wrapError(ErrAuthenticationFailed, "invalid executor token", err, nil)
		}
		if executorName != "" {
			c.Set("executor_name", executorName)
			c.Set("is_executor", true)
			return next(c)
		}

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
			td, err := sess.Get("id_token")
			if err != nil {
				return wrapError(ErrAuthenticationFailed, "could not get id token", err, nil)
			}
			var tokenData TokenData
			if err := tokenData.Decode(td.(string)); err != nil {
				return wrapError(ErrAuthenticationFailed, "invalid token data", err, nil)
			}

			_, err = h.authconfig[tokenData.Provider].verifier.Verify(context.Background(), tokenData.RawIDToken)
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

// authenticateExecutor validates the executor API key from the Authorization header,
// resolves the user from X-User-UUID, and sets the user in the context.
// Returns the executor name if valid, or empty string if not an executor request.
func (h *Handler) authenticateExecutor(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer "+core.ExecutorTokenPrefix) {
		return "", nil
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	executorName, err := core.ValidateExecutorToken(token, h.executorSigningKey)
	if err != nil {
		return "", err
	}

	// Resolve user from X-User-UUID
	if userUUID := c.Request().Header.Get("X-User-UUID"); userUUID != "" {
		userWithGroups, err := h.co.GetUserWithUUIDWithGroups(c.Request().Context(), userUUID)
		if err == nil {
			c.Set("user", userWithGroups.ToUserInfo())
		}
	}

	return executorName, nil
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

// AuthorizeNamespaceAction checks if a user is allowed to perform an action on the given resource.
// For flow/execution resources, it resolves the flow prefix to build a domain-scoped check.
func (h *Handler) AuthorizeNamespaceAction(resource models.Resource, action models.RBACAction) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// skip RBAC for executors
			if isExecutor, _ := c.Get("is_executor").(bool); isExecutor {
				return next(c)
			}

			user, err := h.getUserInfo(c)
			if err != nil {
				return wrapError(ErrAuthenticationFailed, "could not get user details", err, nil)
			}

			namespaceID, ok := c.Get("namespace").(string)
			if !ok {
				return wrapError(ErrRequiredFieldMissing, "could not get namespace", nil, nil)
			}

			// Build domain: default to namespace-level
			domain := core.NamespaceDomain(namespaceID)

			// For flow/execution resources, resolve the flow prefix for domain-scoped checks
			if resource == models.ResourceFlow || resource == models.ResourceExecution {
				if flowID := c.Param("flowID"); flowID != "" {
					if f, err := h.co.GetFlowByID(flowID, namespaceID); err == nil {
						domain = core.FlowDomain(namespaceID, f.Meta.Prefix)
					}
				} else if flowSlug := c.Param("flow"); flowSlug != "" {
					if f, err := h.co.GetFlowByID(flowSlug, namespaceID); err == nil {
						domain = core.FlowDomain(namespaceID, f.Meta.Prefix)
					}
				} else if group := c.Param("group"); group != "" {
					domain = core.FlowDomain(namespaceID, group)
				}
			}

			allowed, err := h.co.CheckPermission(c.Request().Context(), user.ID, domain, resource, action)
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

// AuthorizeNamespaceAdmins checks if a user is an admin in at least one namespace
// This is used for global resources that namespace admins should be able to access
func (h *Handler) AuthorizeNamespaceAdmins() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, err := h.getUserInfo(c)
			if err != nil {
				return wrapError(ErrAuthenticationFailed, "could not get user details", err, nil)
			}

			// Get all namespaces the user has access to
			namespaces, err := h.co.GetUserNamespaces(c.Request().Context(), user.ID)
			if err != nil {
				return wrapError(ErrOperationFailed, "could not get user namespaces", err, nil)
			}

			// Check if user is admin in any namespace
			for _, ns := range namespaces {
				if ns.Role == models.NamespaceRoleAdmin {
					return next(c)
				}
			}

			return wrapError(ErrForbidden, "insufficient permissions", nil, nil)
		}
	}
}

// AuthorizeAction checks if a user is allowed to perform an action on the resource irrespective of the namespace
func (h *Handler) AuthorizeAction(resource models.Resource, action models.RBACAction) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, err := h.getUserInfo(c)
			if err != nil {
				return wrapError(ErrAuthenticationFailed, "could not get user details", err, nil)
			}

			allowed, err := h.co.CheckPermission(c.Request().Context(), user.ID, "/*", resource, action)
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

// NamespaceMiddleware resolves the namespace name to ID and checks user access
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

		// skip permission checks for executors
		if isExecutor, _ := c.Get("is_executor").(bool); isExecutor {
			c.Set("namespace", ns.ID)
			return next(c)
		}

		user, err := h.getUserInfo(c)
		if err != nil {
			return wrapError(ErrAuthenticationFailed, "could not get user details", err, nil)
		}

		// Namespace gate: all namespace members have namespace:view via role:user base policies
		domain := core.NamespaceDomain(ns.ID)
		hasAccess, err := h.co.CheckPermission(c.Request().Context(), user.ID, domain, models.ResourceNamespace, models.RBACActionView)
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
	// Check context first (set by Authenticate for both executor and session requests)
	if user, ok := c.Get("user").(models.UserInfo); ok {
		return user, nil
	}

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
