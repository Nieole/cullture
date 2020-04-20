package actions

import (
	"culture/models"
	"fmt"
	"net/http"

	"github.com/gobuffalo/nulls"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/x/responder"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (User)
// DB Table: Plural (users)
// Resource: Plural (Users)
// Path: Plural (/users)
// View Template Folder: Plural (/templates/users/)

// UsersResource is the resource for the User model
type UsersResource struct {
	buffalo.Resource
}

// List gets all Users. This function is mapped to the path
// GET /users
//func (v UsersResource) List(c buffalo.Context) error {
//	// Get the DB connection from the context
//	tx, ok := c.Value("tx").(*pop.Connection)
//	if !ok {
//		return fmt.Errorf("no transaction found")
//	}
//
//	users := &models.Users{}
//
//	// Paginate results. Params "page" and "per_page" control pagination.
//	// Default values are "page=1" and "per_page=20".
//	q := tx.PaginateFromParams(c.Params())
//
//	// Retrieve all Users from the DB
//	if err := q.All(users); err != nil {
//		return err
//	}
//
//	return responder.Wants("json", func(c buffalo.Context) error {
//		return c.Render(200, r.JSON(users))
//	}).Wants("xml", func(c buffalo.Context) error {
//		return c.Render(200, r.XML(users))
//	}).Respond(c)
//}

// Show gets the data for one User. This function is mapped to
// the path GET /users/{user_id}
//func (v UsersResource) Show(c buffalo.Context) error {
//	// Get the DB connection from the context
//	tx, ok := c.Value("tx").(*pop.Connection)
//	if !ok {
//		return fmt.Errorf("no transaction found")
//	}
//
//	// Allocate an empty User
//	user := &models.User{}
//
//	// To find the User the parameter user_id is used.
//	if err := tx.Find(user, c.Param("user_id")); err != nil {
//		return c.Error(http.StatusNotFound, err)
//	}
//
//	return responder.Wants("json", func(c buffalo.Context) error {
//		return c.Render(200, r.JSON(user))
//	}).Wants("xml", func(c buffalo.Context) error {
//		return c.Render(200, r.XML(user))
//	}).Respond(c)
//}

// Create adds a User to the DB. This function is mapped to the
// path POST /users
//func (v UsersResource) Create(c buffalo.Context) error {
//	// Allocate an empty User
//	user := &models.User{}
//
//	// Bind user to the html form elements
//	if err := c.Bind(user); err != nil {
//		return err
//	}
//
//	// Get the DB connection from the context
//	tx, ok := c.Value("tx").(*pop.Connection)
//	if !ok {
//		return fmt.Errorf("no transaction found")
//	}
//
//	// Validate the data from the html form
//	verrs, err := tx.ValidateAndCreate(user)
//	if err != nil {
//		return err
//	}
//
//	if verrs.HasAny() {
//		return responder.Wants("json", func(c buffalo.Context) error {
//			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
//		}).Wants("xml", func(c buffalo.Context) error {
//			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
//		}).Respond(c)
//	}
//
//	return responder.Wants("json", func(c buffalo.Context) error {
//		return c.Render(http.StatusCreated, r.JSON(user))
//	}).Wants("xml", func(c buffalo.Context) error {
//		return c.Render(http.StatusCreated, r.XML(user))
//	}).Respond(c)
//}

// Update changes a User in the DB. This function is mapped to
// the path PUT /users/{user_id}
func (v UsersResource) Update(c buffalo.Context) error {
	u, ok := c.Session().Get("current_user").(*models.User)
	if !ok {
		return c.Render(http.StatusUnauthorized, Fail("用户未登录不允许操作"))
	}
	if c.Param("user_id") != u.ID.String() {
		return c.Render(http.StatusUnauthorized, Fail("只能修改本账号信息"))
	}
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty User
	user := &models.User{}

	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}
	oldLoginName := user.LoginName

	// Bind User to the html form elements
	if err := c.Bind(user); err != nil {
		return err
	}
	if user.LoginName != oldLoginName {
		u := &models.User{}
		exist, err := tx.Where("login_name = ?", user.LoginName).Exists(u)
		if err != nil {
			return c.Render(http.StatusBadRequest, Fail("检测用户名是否可用失败 %v", err))
		}
		if exist {
			return c.Render(http.StatusUnprocessableEntity, Fail("用户名重复"))
		}
	}

	user.IsActive = true

	verrs, err := tx.ValidateAndUpdate(user)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	user.Password = nulls.String{}
	user.PasswordConfirmation = nulls.String{}
	c.Session().Set("current_user", user)
	c.Session().Save()
	return responder.Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(user))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(user))
	}).Respond(c)
}

// Destroy deletes a User from the DB. This function is mapped
// to the path DELETE /users/{user_id}
//func (v UsersResource) Destroy(c buffalo.Context) error {
//	// Get the DB connection from the context
//	tx, ok := c.Value("tx").(*pop.Connection)
//	if !ok {
//		return fmt.Errorf("no transaction found")
//	}
//
//	// Allocate an empty User
//	user := &models.User{}
//
//	// To find the User the parameter user_id is used.
//	if err := tx.Find(user, c.Param("user_id")); err != nil {
//		return c.Error(http.StatusNotFound, err)
//	}
//
//	if err := tx.Destroy(user); err != nil {
//		return err
//	}
//
//	return responder.Wants("json", func(c buffalo.Context) error {
//		return c.Render(http.StatusOK, r.JSON(user))
//	}).Wants("xml", func(c buffalo.Context) error {
//		return c.Render(http.StatusOK, r.XML(user))
//	}).Respond(c)
//}
