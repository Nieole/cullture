package actions

import (
	"culture/models"
	"fmt"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/x/responder"
	"net/http"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (PostTag)
// DB Table: Plural (post_tags)
// Resource: Plural (PostTags)
// Path: Plural (/post_tags)
// View Template Folder: Plural (/templates/post_tags/)

// PostTagsResource is the resource for the PostTag model
type PostTagsResource struct {
	buffalo.Resource
}

// List gets all PostTags. This function is mapped to the path
// GET /post_tags
func (v PostTagsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	postTags := &models.PostTags{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all PostTags from the DB
	if err := q.All(postTags); err != nil {
		return err
	}

	return responder.Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(postTags))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(postTags))
	}).Respond(c)
}

// Show gets the data for one PostTag. This function is mapped to
// the path GET /post_tags/{post_tag_id}
func (v PostTagsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty PostTag
	postTag := &models.PostTag{}

	// To find the PostTag the parameter post_tag_id is used.
	if err := tx.Find(postTag, c.Param("post_tag_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return responder.Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(postTag))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(postTag))
	}).Respond(c)
}

// Create adds a PostTag to the DB. This function is mapped to the
// path POST /post_tags
func (v PostTagsResource) Create(c buffalo.Context) error {
	// Allocate an empty PostTag
	postTag := &models.PostTag{}

	// Bind postTag to the html form elements
	if err := c.Bind(postTag); err != nil {
		return err
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(postTag)
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
		return c.Render(http.StatusCreated, r.JSON(postTag))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.XML(postTag))
	}).Respond(c)
}

// Update changes a PostTag in the DB. This function is mapped to
// the path PUT /post_tags/{post_tag_id}
func (v PostTagsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty PostTag
	postTag := &models.PostTag{}

	if err := tx.Find(postTag, c.Param("post_tag_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Bind PostTag to the html form elements
	if err := c.Bind(postTag); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(postTag)
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
		return c.Render(http.StatusOK, r.JSON(postTag))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(postTag))
	}).Respond(c)
}

// Destroy deletes a PostTag from the DB. This function is mapped
// to the path DELETE /post_tags/{post_tag_id}
func (v PostTagsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty PostTag
	postTag := &models.PostTag{}

	// To find the PostTag the parameter post_tag_id is used.
	if err := tx.Find(postTag, c.Param("post_tag_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(postTag); err != nil {
		return err
	}

	return responder.Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(postTag))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(postTag))
	}).Respond(c)
}
