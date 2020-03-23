package middleware

import (
	"errors"
	"net/http"

	"github.com/gobuffalo/buffalo"
)

//LoginMiddleware LoginMiddleware
func LoginMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if user, ok := c.Session().Get("current_user_id").(string); !ok && user == "" {
			return c.Error(http.StatusUnauthorized, errors.New("Unauthorized"))
		}
		return next(c)
	}
}
