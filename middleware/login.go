package middleware

import (
	"culture/models"
	"fmt"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"log"
	"math/rand"
	"time"

	"github.com/gobuffalo/buffalo"
)

//LoginMiddleware LoginMiddleware
func LoginMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if phone, ok := c.Session().GetOnce("current_user_phone").(string); ok {
			// Get the DB connection from the context
			tx, ok := c.Value("tx").(*pop.Connection)
			if !ok {
				return fmt.Errorf("no transaction found")
			}
			name, ok := c.Session().GetOnce("current_user_name").(string)
			if !ok || name == "" {
				name = RandString(5)
			}
			user := &models.User{
				Name:      name,
				LoginName: phone,
			}
			err := tx.Save(user)
			if err != nil {
				log.Println(errors.WithStack(err))
				return fmt.Errorf("保存用户失败 : %v", err)
			}
			c.Session().Set("current_user", user)
			c.Session().Save()
			// TODO 清洗
			posts := &models.Posts{}
			if err := tx.Where("user_phone = ?", phone).All(posts); err != nil {
				log.Printf("failed to select : %v", err)
				return fmt.Errorf("查询posts失败 : %v", err)
			}
			posts.ChangeLike(user, tx)
			return next(c)
		}
		if _, ok := c.Session().Get("current_user").(*models.User); !ok {
			// Get the DB connection from the context
			tx, ok := c.Value("tx").(*pop.Connection)
			if !ok {
				return fmt.Errorf("no transaction found")
			}
			login, _ := uuid.NewV4()
			user := &models.User{
				Name:      RandString(5),
				LoginName: login.String(),
			}
			tx.Save(user)
			c.Session().Set("current_user", user)
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
