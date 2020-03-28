package actions

import (
	"culture/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/x/responder"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Project)
// DB Table: Plural (projects)
// Resource: Plural (Projects)
// Path: Plural (/projects)
// View Template Folder: Plural (/templates/projects/)

// ProjectsResource is the resource for the Project model
type ProjectsResource struct {
	buffalo.Resource
}

// List gets all Projects. This function is mapped to the path
// GET /projects
func (v ProjectsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	projects := &models.Projects{}

	filter := func(projectName string) pop.ScopeFunc {
		return func(q *pop.Query) *pop.Query {
			if projectName != "" {
				q.Where("name like ?", "%"+projectName+"%")
			}
			return q
		}
	}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Projects from the DB
	if err := q.Scope(filter(c.Param("projectName"))).All(projects); err != nil {
		return err
	}

	return responder.Wants("json", func(c buffalo.Context) error {
		return c.Render(200, List(projects, q.Paginator))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(projects))
	}).Respond(c)
}

// Show gets the data for one Project. This function is mapped to
// the path GET /projects/{project_id}
func (v ProjectsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Project
	project := &models.Project{}

	// To find the Project the parameter project_id is used.
	if err := tx.Eager("Organization").Find(project, c.Param("project_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if geo, err := models.REDIS.GeoPos("project_geo", project.ID.String()).Result(); err == nil && len(geo) > 0 {
		project.Latitude = strconv.FormatFloat(geo[0].Latitude, 'f', -1, 64)
		project.Longitude = strconv.FormatFloat(geo[0].Longitude, 'f', -1, 64)
	}

	return responder.Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(project))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(project))
	}).Respond(c)
}

// Create adds a Project to the DB. This function is mapped to the
// path POST /projects
func (v ProjectsResource) Create(c buffalo.Context) error {
	// Allocate an empty Project
	project := &models.Project{}

	// Bind project to the html form elements
	if err := c.Bind(project); err != nil {
		return err
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(project)
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
		return c.Render(http.StatusCreated, r.JSON(project))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.XML(project))
	}).Respond(c)
}

// Update changes a Project in the DB. This function is mapped to
// the path PUT /projects/{project_id}
func (v ProjectsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Project
	project := &models.Project{}

	if err := tx.Find(project, c.Param("project_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Bind Project to the html form elements
	if err := c.Bind(project); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(project)
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
		return c.Render(http.StatusOK, r.JSON(project))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(project))
	}).Respond(c)
}

// Destroy deletes a Project from the DB. This function is mapped
// to the path DELETE /projects/{project_id}
func (v ProjectsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Project
	project := &models.Project{}

	// To find the Project the parameter project_id is used.
	if err := tx.Find(project, c.Param("project_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(project); err != nil {
		return err
	}

	return responder.Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(project))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(project))
	}).Respond(c)
}
