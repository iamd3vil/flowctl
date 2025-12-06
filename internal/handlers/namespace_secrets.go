package handlers

import (
	"fmt"
	"net/http"

	"github.com/cvhariharan/flowctl/internal/core/models"
	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleCreateNamespaceSecret(c echo.Context) error {
	namespace, ok := c.Get("namespace").(string)
	if !ok {
		return wrapError(ErrRequiredFieldMissing, "could not get namespace", nil, nil)
	}

	var req NamespaceSecretReq
	if err := c.Bind(&req); err != nil {
		return wrapError(ErrInvalidInput, "could not decode request", err, nil)
	}

	if err := h.validate.Struct(req); err != nil {
		return wrapError(ErrValidationFailed, fmt.Sprintf("request validation failed: %s", formatValidationErrors(err)), err, nil)
	}

	secret := models.NamespaceSecret{
		Key:         req.Key,
		Value:       req.Value,
		Description: req.Description,
	}

	created, err := h.co.CreateNamespaceSecret(c.Request().Context(), secret, namespace)
	if err != nil {
		return wrapError(ErrOperationFailed, "could not create namespace secret", err, nil)
	}

	return c.JSON(http.StatusCreated, coreNamespaceSecretToNamespaceSecretResp(created))
}

func (h *Handler) HandleGetNamespaceSecret(c echo.Context) error {
	namespace, ok := c.Get("namespace").(string)
	if !ok {
		return wrapError(ErrRequiredFieldMissing, "could not get namespace", nil, nil)
	}

	var req NamespaceSecretGetReq
	if err := c.Bind(&req); err != nil {
		return wrapError(ErrInvalidInput, "could not decode request", err, nil)
	}

	if err := h.validate.Struct(req); err != nil {
		return wrapError(ErrValidationFailed, fmt.Sprintf("request validation failed: %s", formatValidationErrors(err)), err, nil)
	}

	secret, err := h.co.GetNamespaceSecretByID(c.Request().Context(), req.SecretID, namespace)
	if err != nil {
		return wrapError(ErrResourceNotFound, "secret not found", err, nil)
	}

	return c.JSON(http.StatusOK, coreNamespaceSecretToNamespaceSecretResp(secret))
}

func (h *Handler) HandleListNamespaceSecrets(c echo.Context) error {
	namespace, ok := c.Get("namespace").(string)
	if !ok {
		return wrapError(ErrRequiredFieldMissing, "could not get namespace", nil, nil)
	}

	secrets, err := h.co.ListNamespaceSecrets(c.Request().Context(), namespace)
	if err != nil {
		return wrapError(ErrOperationFailed, "could not list namespace secrets", err, nil)
	}

	resp := make([]NamespaceSecretResp, len(secrets))
	for i, secret := range secrets {
		resp[i] = coreNamespaceSecretToNamespaceSecretResp(secret)
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) HandleUpdateNamespaceSecret(c echo.Context) error {
	namespace, ok := c.Get("namespace").(string)
	if !ok {
		return wrapError(ErrRequiredFieldMissing, "could not get namespace", nil, nil)
	}

	var req NamespaceSecretUpdateReq
	if err := c.Bind(&req); err != nil {
		return wrapError(ErrInvalidInput, "could not decode request", err, nil)
	}

	if err := h.validate.Struct(req); err != nil {
		return wrapError(ErrValidationFailed, fmt.Sprintf("request validation failed: %s", formatValidationErrors(err)), err, nil)
	}

	secret := models.NamespaceSecret{
		Value:       req.Value,
		Description: req.Description,
	}

	updated, err := h.co.UpdateNamespaceSecret(c.Request().Context(), req.SecretID, secret, namespace)
	if err != nil {
		return wrapError(ErrOperationFailed, "could not update namespace secret", err, nil)
	}

	return c.JSON(http.StatusOK, coreNamespaceSecretToNamespaceSecretResp(updated))
}

func (h *Handler) HandleDeleteNamespaceSecret(c echo.Context) error {
	namespace, ok := c.Get("namespace").(string)
	if !ok {
		return wrapError(ErrRequiredFieldMissing, "could not get namespace", nil, nil)
	}

	var req NamespaceSecretGetReq
	if err := c.Bind(&req); err != nil {
		return wrapError(ErrInvalidInput, "could not decode request", err, nil)
	}

	if err := h.validate.Struct(req); err != nil {
		return wrapError(ErrValidationFailed, fmt.Sprintf("request validation failed: %s", formatValidationErrors(err)), err, nil)
	}

	err := h.co.DeleteNamespaceSecret(c.Request().Context(), req.SecretID, namespace)
	if err != nil {
		return wrapError(ErrOperationFailed, "could not delete namespace secret", err, nil)
	}

	return c.NoContent(http.StatusOK)
}
