package actions

import (
	"culture/models"
	"fmt"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"log"
	"net/http"
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
			id, _ := uuid.NewV4()
			user := &models.User{
				Name:      name,
				LoginName: id.String(),
			}
			err := tx.Save(user)
			if err != nil {
				log.Println(errors.WithStack(err))
				return c.Render(http.StatusBadRequest, Fail("保存用户失败 : %v", err))
			}
			c.Session().Set("current_user", user)
			c.Session().Save()
			// 清洗
			posts := &models.Posts{}
			if err := tx.Where("user_phone = ?", phone).All(posts); err != nil {
				log.Printf("failed to select : %v", err)
				return c.Render(http.StatusBadRequest, Fail("查询posts失败 : %v", err))
			}
			posts.ChangeLike(user, tx, phone)
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
