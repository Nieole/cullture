package actions

import (
	"culture/models"
	"fmt"
	"github.com/gobuffalo/pop"
	"golang.org/x/crypto/bcrypt"
	"net/http"

	"github.com/gobuffalo/buffalo"
)

//LoginHandler LoginHandler
func LoginHandler(c buffalo.Context) error {
	// Allocate an empty User
	user := &models.User{}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	if err := tx.Select("id,password_hash").Where("name = ?", c.Param("username")).First(user); err != nil {
		return c.Error(http.StatusNotFound, err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(c.Param("password"))); err != nil {
		return c.Render(http.StatusUnauthorized, r.JSON(map[string]string{"message": "failed login"}))
	}
	c.Session().Set("current_user", user)
	c.Session().Save()
	return c.Render(http.StatusCreated, nil)
}

//SignOutHandler SignOutHandler
func SignOutHandler(c buffalo.Context) error {
	c.Session().Clear()
	c.Session().Save()
	return c.Render(http.StatusAccepted, nil)
}
