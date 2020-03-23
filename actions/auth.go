package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

//LoginHandler LoginHandler
func LoginHandler(c buffalo.Context) error {
	//// Allocate an empty User
	//user := &models.User{}
	//
	//// Get the DB connection from the context
	//tx, ok := c.Value("tx").(*pop.Connection)
	//if !ok {
	//	return fmt.Errorf("no transaction found")
	//}
	//if err := tx.Select("id,password_hash").Where("name = ?", c.Param("username")).First(user); err != nil {
	//	return c.Error(http.StatusNotFound, err)
	//}
	//if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(c.Param("password"))); err != nil {
	//	return c.Render(http.StatusUnauthorized, r.JSON(map[string]string{"message": "failed login"}))
	//}
	c.Session().Set("current_user_id", "aaa")
	return c.Render(http.StatusOK, r.String("aaa"))
}

//SignOutHandler SignOutHandler
func SignOutHandler(c buffalo.Context) error {
	c.Session().Clear()
	return c.Render(http.StatusAccepted, nil)
}

//RegisterHandler RegisterHandler
func RegisterHandler(c buffalo.Context) error {
	//// Allocate an empty User
	//user := &models.User{}
	//
	//// Bind user to the html form elements
	//if err := c.Bind(user); err != nil {
	//	return err
	//}
	//
	//// Get the DB connection from the context
	//tx, ok := c.Value("tx").(*pop.Connection)
	//if !ok {
	//	return fmt.Errorf("no transaction found")
	//}
	//
	//// Validate the data from the html form
	//verrs, err := tx.ValidateAndCreate(user)
	//if err != nil {
	//	return err
	//}
	//
	//if verrs.HasAny() {
	//	return responder.Wants("json", func(c buffalo.Context) error {
	//		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	//	}).Wants("xml", func(c buffalo.Context) error {
	//		return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
	//	}).Respond(c)
	//}

	c.Session().Set("current_user_id", "aaa")
	//return responder.Wants("json", func(c buffalo.Context) error {
	//	return c.Render(http.StatusCreated, nil)
	//}).Wants("xml", func(c buffalo.Context) error {
	//	return c.Render(http.StatusCreated, nil)
	//}).Respond(c)
	return nil
}
