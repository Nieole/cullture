package middleware

import (
	"math/rand"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gofrs/uuid"
)

//LoginMiddleware LoginMiddleware
func LoginMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if user, ok := c.Session().Get("current_user_phone").(string); !ok && user == "" {
			//return c.Error(http.StatusUnauthorized, errors.New("Unauthorized"))
			id, _ := uuid.NewV4()
			c.Session().Set("current_user_phone", id.String())
			c.Session().Set("current_user_name", RandString(5))
			c.Session().Save()
		}
		return next(c)
	}
}

var ran *rand.Rand

func init() {
	ran = rand.New(rand.NewSource(time.Now().Unix()))
}

//RandString RandString
func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(ran.Intn(26) + 65)
	}
	return string(bytes)
}
