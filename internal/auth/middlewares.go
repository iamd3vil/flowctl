package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/cvhariharan/autopilot/internal/models"
	"github.com/labstack/echo/v4"
)

func (h *AuthHandler) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := h.sessMgr.Acquire(nil, c, c)
		if err != nil {
			log.Println(err)
			return h.handleUnauthenticated(c)
		}

		rawIDToken, err := sess.Get("id_token")
		if err != nil || rawIDToken == nil {
			return h.handleUnauthenticated(c)
		}

		// Verify the token
		_, err = h.authconfig.verifier.Verify(context.Background(), rawIDToken.(string))
		if err != nil {
			log.Println(err)
			sess.Delete("id_token")
			sess.Delete("user")
			return h.handleUnauthenticated(c)
		}

		// Set user info in context
		if userInfo, err := sess.Get("user"); err == nil {
			var user models.UserInfo
			userBytes, err := json.Marshal(userInfo)
			if err != nil {
				c.Logger().Error(err)
				return h.handleUnauthenticated(c)
			}

			if err := json.NewDecoder(bytes.NewBuffer(userBytes)).Decode(&user); err != nil {
				c.Logger().Error(err)
				return h.handleUnauthenticated(c)
			}

			u, err := h.store.GetUserByUsername(c.Request().Context(), user.Email)
			if err != nil {
				c.Logger().Error(err)
				return echo.NewHTTPError(http.StatusUnauthorized, "could not authenticate user")
			}
			user.UUID = u.Uuid.String()
			user.ID = u.ID

			c.Set("user", user)
			return next(c)
		}

		log.Println(err)
		return h.handleUnauthenticated(c)
	}
}
