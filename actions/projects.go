package actions

import (
	"culture/cache"
	"culture/models"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/nulls"

	"github.com/go-redis/redis/v7"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/x/defaults"
	"github.com/gobuffalo/x/responder"
)

var ak = envy.Get("ak", "KIYHs3gdIoIazwZjqz2BeAEK567zxasg")

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
	list := &ListResponse{}
	if err := cache.Once(fmt.Sprintf("cache:projects:%v:%v:%v", c.Param("projectName"), c.Param("page"), c.Param("per_page")), list, func() (interface{}, error) {
		// Get the DB connection from the context
		tx, ok := c.Value("tx").(*pop.Connection)
		if !ok {
			return nil, fmt.Errorf("no transaction found")
		}
		// Paginate results. Params "page" and "per_page" control pagination.
		// Default values are "page=1" and "per_page=20".
		q := tx.PaginateFromParams(c.Params())
		filter := func(projectName string) pop.ScopeFunc {
			return func(q *pop.Query) *pop.Query {
				if projectName != "" {
					q.Where("name like ?", "%"+projectName+"%")
				}
				return q
			}
		}
		projects := &models.Projects{}

		// Retrieve all Projects from the DB
		if err := q.Scope(filter(c.Param("projectName"))).Order("name asc").All(projects); err != nil {
			return nil, err
		}
		return &ListResponse{
			Data:      projects,
			Paginator: q.Paginator,
		}, nil
	}, time.Hour*6); err != nil {
		return c.Render(http.StatusBadRequest, Fail("加载缓存数据失败 %v", err))
	}

	return responder.Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(list))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(list))
	}).Respond(c)
}

// Show gets the data for one Project. This function is mapped to
// the path GET /projects/{project_id}
func (v ProjectsResource) Show(c buffalo.Context) error {
	project := &models.Project{}
	if err := cache.Once(fmt.Sprintf("cache:project:%v", c.Param("project_id")), project, func() (interface{}, error) {
		// Get the DB connection from the context
		tx, ok := c.Value("tx").(*pop.Connection)
		if !ok {
			return nil, fmt.Errorf("no transaction found")
		}

		// Allocate an empty Project
		p := &models.Project{}

		// To find the Project the parameter project_id is used.
		if err := tx.Eager("Organization").Find(p, c.Param("project_id")); err != nil {
			return nil, c.Error(http.StatusNotFound, err)
		}

		if geo, err := models.REDIS.GeoPos("project_geo", p.ID.String()).Result(); err == nil && len(geo) > 0 && geo[0] != nil {
			p.Latitude = strconv.FormatFloat(defaults.Float64(geo[0].Latitude, 0), 'f', -1, 64)
			p.Longitude = strconv.FormatFloat(defaults.Float64(geo[0].Longitude, 0), 'f', -1, 64)
		} else {
			p.Latitude = "104.086818"
			p.Longitude = "30.683696"
		}
		return p, nil
	}, time.Hour*6); err != nil {
		return c.Render(http.StatusBadRequest, Fail("加载缓存数据失败 %v", err))
	}

	return responder.Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(project))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(project))
	}).Respond(c)
}

// ShowCountHandler gets the data for one Project. This function is mapped to
// the path GET /projects/count/{project_id}
func ShowCountHandler(c buffalo.Context) error {
	project := &models.Project{}
	if err := cache.Once(fmt.Sprintf("cache:project:count:%v", c.Param("project_id")), project, func() (interface{}, error) {
		// Get the DB connection from the context
		tx, ok := c.Value("tx").(*pop.Connection)
		if !ok {
			return nil, fmt.Errorf("no transaction found")
		}

		// Allocate an empty Project
		p := &models.Project{}

		// To find the Project the parameter project_id is used.
		if err := tx.Eager("Organization").Find(p, c.Param("project_id")); err != nil {
			return nil, c.Error(http.StatusNotFound, err)
		}
		posts := &models.Posts{}
		count, err := tx.Where("project_id = ?", p.ID).Where("is_delete = ?", false).Count(posts)
		if err != nil {
			log.Printf("failed to count posts %v", err)
		} else {
			p.PostsCount = count
		}

		if geo, err := models.REDIS.GeoPos("project_geo", p.ID.String()).Result(); err == nil && len(geo) > 0 && geo[0] != nil {
			p.Latitude = strconv.FormatFloat(defaults.Float64(geo[0].Latitude, 0), 'f', -1, 64)
			p.Longitude = strconv.FormatFloat(defaults.Float64(geo[0].Longitude, 0), 'f', -1, 64)
		} else {
			p.Latitude = "104.086818"
			p.Longitude = "30.683696"
		}
		return p, nil
	}, time.Second*5); err != nil {
		return c.Render(http.StatusBadRequest, Fail("加载缓存数据失败 %v", err))
	}

	return responder.Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(project))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(project))
	}).Respond(c)
}

// Create adds a Project to the DB. This function is mapped to the
// path POST /projects
//func (v ProjectsResource) Create(c buffalo.Context) error {
//	// Allocate an empty Project
//	project := &models.Project{}
//
//	// Bind project to the html form elements
//	if err := c.Bind(project); err != nil {
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
//	verrs, err := tx.ValidateAndCreate(project)
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
//	if project.Latitude != "" && project.Longitude != "" {
//		longitude, err := strconv.ParseFloat(project.Longitude, 64)
//		if err != nil {
//			return c.Render(http.StatusBadRequest, Fail("解析经度失败 %v", err))
//		}
//		latitude, err := strconv.ParseFloat(project.Latitude, 64)
//		if err != nil {
//			return c.Render(http.StatusBadRequest, Fail("解析纬度失败 %v", err))
//		}
//		if _, err := models.REDIS.GeoAdd("project_geo", &redis.GeoLocation{
//			Name:      project.ID.String(),
//			Longitude: longitude,
//			Latitude:  latitude,
//		}).Result(); err != nil {
//			return c.Render(http.StatusUnprocessableEntity, Fail("更新地理位置失败 %v", err))
//		}
//	}
//	return responder.Wants("json", func(c buffalo.Context) error {
//		return c.Render(http.StatusCreated, r.JSON(project))
//	}).Wants("xml", func(c buffalo.Context) error {
//		return c.Render(http.StatusCreated, r.XML(project))
//	}).Respond(c)
//}

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
		log.Printf("bind project error : %v", err)
		return err
	}

	if project.Latitude != "" && project.Longitude != "" {
		resp, err := http.Get(fmt.Sprintf("http://api.map.baidu.com/reverse_geocoding/v3/?ak=%s&output=json&coordtype=wgs84ll&location=%s,%s", ak, project.Latitude, project.Longitude))
		if err != nil {
			log.Printf("failed to get baidu geo %v", err)
		} else {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("failed to get baidu geo %v", err)
			}
			country, err := jsonparser.GetString(body, "result", "addressComponent", "country")
			if err != nil {
				log.Printf("failed to get country %v", err)
			}
			province, err := jsonparser.GetString(body, "result", "addressComponent", "province")
			if err != nil {
				log.Printf("failed to get province %v", err)
			}
			city, err := jsonparser.GetString(body, "result", "addressComponent", "city")
			if err != nil {
				log.Printf("failed to get city %v", err)
			}
			district, err := jsonparser.GetString(body, "result", "addressComponent", "district")
			if err != nil {
				log.Printf("failed to get district %v", err)
			}
			project.Country = nulls.NewString(country)
			project.Province = nulls.NewString(province)
			project.City = nulls.NewString(city)
			project.District = nulls.NewString(district)
		}
	}

	verrs, err := tx.ValidateAndUpdate(project)
	if err != nil {
		log.Printf("ValidateAndUpdate project error : %v", err)
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	if project.Latitude != "" && project.Longitude != "" {
		longitude, err := strconv.ParseFloat(project.Longitude, 64)
		if err != nil {
			return c.Render(http.StatusBadRequest, Fail("解析经度失败 %v", err))
		}
		latitude, err := strconv.ParseFloat(project.Latitude, 64)
		if err != nil {
			return c.Render(http.StatusBadRequest, Fail("解析纬度失败 %v", err))
		}
		if _, err := models.REDIS.GeoAdd("project_geo", &redis.GeoLocation{
			Name:      project.ID.String(),
			Longitude: longitude,
			Latitude:  latitude,
		}).Result(); err != nil {
			return c.Render(http.StatusUnprocessableEntity, Fail("更新地理位置失败 %v", err))
		}
	}
	return responder.Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(project))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(project))
	}).Respond(c)
}

// Destroy deletes a Project from the DB. This function is mapped
// to the path DELETE /projects/{project_id}
//func (v ProjectsResource) Destroy(c buffalo.Context) error {
//	// Get the DB connection from the context
//	tx, ok := c.Value("tx").(*pop.Connection)
//	if !ok {
//		return fmt.Errorf("no transaction found")
//	}
//
//	// Allocate an empty Project
//	project := &models.Project{}
//
//	// To find the Project the parameter project_id is used.
//	if err := tx.Find(project, c.Param("project_id")); err != nil {
//		return c.Error(http.StatusNotFound, err)
//	}
//
//	if err := tx.Destroy(project); err != nil {
//		return err
//	}
//
//	return responder.Wants("json", func(c buffalo.Context) error {
//		return c.Render(http.StatusOK, r.JSON(project))
//	}).Wants("xml", func(c buffalo.Context) error {
//		return c.Render(http.StatusOK, r.XML(project))
//	}).Respond(c)
//}
