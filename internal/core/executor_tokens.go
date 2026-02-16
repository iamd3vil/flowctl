package core

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

const ExecutorTokenPrefix = "fctl_"

// GenerateSigningKey generates a 32-byte random signing key for HMAC-SHA256.
func GenerateSigningKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate signing key: %w", err)
	}
	return key, nil
}

// GenerateExecutorToken creates a stateless HMAC-SHA256 token for an executor.
// Format: fctl_<executor_name>.<base64url(HMAC-SHA256(executor_name, signing_key))>
func GenerateExecutorToken(executorName string, signingKey []byte) (string, error) {
	if executorName == "" {
		return "", fmt.Errorf("executor name cannot be empty")
	}
	if len(signingKey) == 0 {
		return "", fmt.Errorf("signing key cannot be empty")
	}

	mac := hmac.New(sha256.New, signingKey)
	mac.Write([]byte(executorName))
	sig := mac.Sum(nil)

	encoded := base64.RawURLEncoding.EncodeToString(sig)
	return fmt.Sprintf("%s%s.%s", ExecutorTokenPrefix, executorName, encoded), nil
}

// ValidateExecutorToken validates an executor token and returns the executor name.
func ValidateExecutorToken(token string, signingKey []byte) (string, error) {
	if !strings.HasPrefix(token, ExecutorTokenPrefix) {
		return "", fmt.Errorf("invalid token prefix")
	}

	rest := strings.TrimPrefix(token, ExecutorTokenPrefix)

	parts := strings.SplitN(rest, ".", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid token format")
	}

	executorName := parts[0]
	providedSig, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("invalid token signature encoding: %w", err)
	}

	mac := hmac.New(sha256.New, signingKey)
	mac.Write([]byte(executorName))
	expectedSig := mac.Sum(nil)

	if !hmac.Equal(providedSig, expectedSig) {
		return "", fmt.Errorf("invalid token signature")
	}

	return executorName, nil
}
