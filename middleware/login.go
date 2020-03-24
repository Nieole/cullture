package middleware

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gofrs/uuid"
)

//LoginMiddleware LoginMiddleware
func LoginMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if user, ok := c.Session().Get("current_user_phone").(string); !ok && user == "" {
			//return c.Error(http.StatusUnauthorized, errors.New("Unauthorized"))
			id, _ := uuid.NewV4()
			c.Session().Set("current_user_phone", id)
		}
		return next(c)
	}
}
