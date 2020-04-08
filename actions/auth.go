package actions

import (
	"culture/models"
	"fmt"
	"net/http"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"golang.org/x/crypto/bcrypt"

	"github.com/gobuffalo/buffalo"
)

//LoginHandler LoginHandler
func LoginHandler(c buffalo.Context) error {
	// Allocate an empty Post
	login := new(Login)

	// Bind post to the html form elements
	if err := c.Bind(login); err != nil {
		return err
	}
	e, err := login.Validate()
	if err != nil {
		return c.Render(http.StatusBadRequest, Fail("验证表单信息失败 %v", err))
	}
	if e.HasAny() {
		return c.Render(http.StatusBadRequest, Fail("校验表单信息失败 %v", e))
	}

	// Allocate an empty User
	user := &models.User{}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	if err := tx.Where("login_name = ?", login.LoginName).Where("is_active = ?", true).Where("is_delete = ?", false).First(user); err != nil {
		return c.Render(http.StatusNotFound, Fail("未找到用户 : %v", err))
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(login.Password)); err != nil {
		return c.Render(http.StatusUnauthorized, r.JSON(map[string]string{"message": "failed login"}))
	}
	c.Session().Set("current_user", user)
	c.Session().Save()
	return c.Render(http.StatusCreated, nil)
}

//Login Login
type Login struct {
	LoginName string `json:"login_name"`
	Password  string `json:"password"`
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (l *Login) Validate() (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: l.LoginName, Name: "LoginName", Message: "用户名不能为空"},
		&validators.StringIsPresent{Field: l.Password, Name: "Password", Message: "密码不能为空"},
	), nil
}

//SignOutHandler SignOutHandler
func SignOutHandler(c buffalo.Context) error {
	c.Session().Clear()
	c.Session().Save()
	return c.Render(http.StatusAccepted, nil)
}
