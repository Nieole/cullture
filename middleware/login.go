package middleware

import (
	"github.com/gobuffalo/buffalo"
	"math/rand"
	"time"
)

//LoginMiddleware LoginMiddleware
func LoginMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if user, ok := c.Session().Get("current_user_phone").(string); !ok && user == "" {
			//return c.Error(http.StatusUnauthorized, errors.New("Unauthorized"))
			c.Session().Set("current_user_phone", RandString(20))
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
