package actions

import (
	"culture/models"
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/x/responder"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Comment)
// DB Table: Plural (comments)
// Resource: Plural (Comments)
// Path: Plural (/comments)
// View Template Folder: Plural (/templates/comments/)

// CommentsResource is the resource for the Comment model
type CommentsResource struct {
	buffalo.Resource
}

// List gets all Comments. This function is mapped to the path
// GET /comments
func (v CommentsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	comments := &models.Comments{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Comments from the DB
	if err := q.Eager("User").Where("is_delete = ?", false).All(comments); err != nil {
		return err
	}

	return responder.Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(comments))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(comments))
	}).Respond(c)
}

// Show gets the data for one Comment. This function is mapped to
// the path GET /comments/{comment_id}
func (v CommentsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Comment
	comment := &models.Comment{}

	// To find the Comment the parameter comment_id is used.
	if err := tx.Eager("User").Where("is_delete = ?", false).Find(comment, c.Param("comment_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return responder.Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(comment))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(comment))
	}).Respond(c)
}

// Create adds a Comment to the DB. This function is mapped to the
// path POST /comments
func (v CommentsResource) Create(c buffalo.Context) error {
	// Allocate an empty Comment
	comment := &models.Comment{}

	// Bind comment to the html form elements
	if err := c.Bind(comment); err != nil {
		return err
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(comment)
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

	return responder.Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.JSON(comment))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.XML(comment))
	}).Respond(c)
}

// Update changes a Comment in the DB. This function is mapped to
// the path PUT /comments/{comment_id}
func (v CommentsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Comment
	comment := &models.Comment{}

	if err := tx.Find(comment, c.Param("comment_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Bind Comment to the html form elements
	if err := c.Bind(comment); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(comment)
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

	return responder.Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(comment))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(comment))
	}).Respond(c)
}

// Destroy deletes a Comment from the DB. This function is mapped
// to the path DELETE /comments/{comment_id}
func (v CommentsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Comment
	comment := &models.Comment{}

	// To find the Comment the parameter comment_id is used.
	if err := tx.Find(comment, c.Param("comment_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}
	comment.IsDelete = true

	if err := tx.Update(comment); err != nil {
		return err
	}

	return responder.Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(comment))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(comment))
	}).Respond(c)
}
