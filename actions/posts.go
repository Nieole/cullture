package actions

import (
	"culture/models"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gobuffalo/x/responder"
	"github.com/gofrs/uuid"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Post)
// DB Table: Plural (posts)
// Resource: Plural (Posts)
// Path: Plural (/posts)
// View Template Folder: Plural (/templates/posts/)

// PostsResource is the resource for the Post model
type PostsResource struct {
	buffalo.Resource
}

var (
	mu = sync.RWMutex{}
)

// List gets all Posts. This function is mapped to the path
// GET /posts
func (v PostsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	posts := &models.Posts{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	result, err := models.REDIS.Get(fmt.Sprintf("cache:%v:%v", c.Param("updated_at"), c.Param("project_id"))).Result()
	if err != nil {
		mu.Lock()
		// Retrieve all Posts from the DB
		if err := q.Eager("Tags", "User", "Project", "Comments").Scope(ByPage(c.Param("updated_at"))).Scope(ByProject(c.Param("project_id"))).Where("is_delete = ?", false).Order("updated_at desc").All(posts); err != nil {
			return err
		}
		models.REDIS.Set(fmt.Sprintf("cache:%v:%v", c.Param("updated_at"), c.Param("project_id")), posts, time.Second*3)
		mu.Unlock()
	} else {
		err := posts.FromString(result)
		if err != nil {
			return c.Render(http.StatusBadRequest, Fail("解析数据失败 %v", err))
		}
	}
	user, _ := currentUser(c)
	*posts = posts.Fill(user)

	return responder.Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(posts))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(posts))
	}).Respond(c)
}

//ByPage 分页查询posts
func ByPage(updatedAt string) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		if updatedAt != "" {
			q.Where("updated_at < ?", updatedAt)
		}
		return q
	}
}

//ByProject 通过项目id过滤
func ByProject(projectID string) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		if projectID != "" {
			q.Where("project_id = ?", projectID)
		}
		return q
	}
}

//MyList MyList
func MyList(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	posts := &models.Posts{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	user, err := currentUser(c)
	if err != nil {
		return c.Render(http.StatusBadRequest, Fail(err.Error()))
	}

	result, err := models.REDIS.Get(fmt.Sprintf("cache:my:%v:%v:%v", c.Param("updated_at"), c.Param("project_id"), user.ID)).Result()
	if err != nil {
		// Retrieve all Posts from the DB
		if err := q.Eager("Tags", "User", "Project", "Comments").Scope(ByPage(c.Param("updated_at"))).Scope(ByProject(c.Param("project_id"))).Where("user_id = ?", user.ID).Where("is_delete = ?", false).Order("updated_at desc").All(posts); err != nil {
			return err
		}
		models.REDIS.Set(fmt.Sprintf("cache:my:%v:%v:%v", c.Param("updated_at"), c.Param("project_id"), user.ID), posts.String(), time.Second*3)
	} else {
		err := posts.FromString(result)
		if err != nil {
			return c.Render(http.StatusBadRequest, Fail("解析数据失败 %v", err))
		}
	}
	*posts = posts.Fill(user)
	return responder.Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(posts))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(posts))
	}).Respond(c)
}

// Show gets the data for one Post. This function is mapped to
// the path GET /posts/{post_id}
func (v PostsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Post
	post := &models.Post{}

	// To find the Post the parameter post_id is used.
	if err := tx.Eager("Tags").Find(post, c.Param("post_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return responder.Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(post))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(post))
	}).Respond(c)
}

// Create adds a Post to the DB. This function is mapped to the
// path POST /posts
func (v PostsResource) Create(c buffalo.Context) error {
	// Allocate an empty Post
	publish := new(PublishPost)

	// Bind post to the html form elements
	if err := c.Bind(publish); err != nil {
		return err
	}
	e, err := publish.Validate()
	if err != nil {
		return c.Render(http.StatusBadRequest, Fail("验证表单信息失败 %v", err))
	}
	if e.HasAny() {
		return c.Render(http.StatusBadRequest, Fail("校验表单信息失败 %v", e))
	}
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	tags := models.Tags{}
	for _, t := range publish.Tags {
		tag := new(models.Tag)
		err = tx.Find(tag, t)
		if err != nil {
			return c.Render(http.StatusBadRequest, Fail("查询tag失败 %v", err))
		}
		tags = append(tags, *tag)
	}
	project := new(models.Project)
	err = tx.Find(project, publish.Project)
	if err != nil {
		return c.Render(http.StatusBadRequest, Fail("查询project失败 %v", err))
	}
	user, err := currentUser(c)
	if err != nil {
		return c.Render(http.StatusBadRequest, Fail(err.Error()))
	}
	p := &models.Post{
		Project:  project,
		Image:    nulls.NewString(publish.Image),
		Content:  nulls.NewString(publish.Content),
		User:     user,
		Tags:     tags,
		IsDelete: publish.IsDelete,
	}
	e, err = tx.Eager().ValidateAndSave(p)
	if err != nil {
		return c.Render(http.StatusBadRequest, Fail("验证表单信息失败 %v", err))
	}
	if e.HasAny() {
		return c.Render(http.StatusBadRequest, Fail("校验表单信息失败 %v", e))
	}
	if !publish.IsDelete {
		p.Like(user)
	}
	return c.Render(http.StatusCreated, nil)
}

// Update changes a Post in the DB. This function is mapped to
// the path PUT /posts/{post_id}
func (v PostsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Post
	post := &models.Post{}

	if err := tx.Find(post, c.Param("post_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Bind Post to the html form elements
	if err := c.Bind(post); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(post)
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
		return c.Render(http.StatusOK, r.JSON(post))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(post))
	}).Respond(c)
}

//Like Like
func Like(c buffalo.Context) error {
	err := query(c, true, true)
	if err != nil {
		return c.Render(http.StatusBadRequest, Fail("点赞失败 : %v", err.Error()))
	}
	return c.Render(http.StatusCreated, nil)
}

//UnLike UnLike
func UnLike(c buffalo.Context) error {
	err := query(c, true, false)
	if err != nil {
		return c.Render(http.StatusBadRequest, Fail("取消点赞失败 : %v", err.Error()))
	}
	return c.Render(http.StatusCreated, nil)
}

func query(c buffalo.Context, like, append bool) error {
	user, err := currentUser(c)
	if err != nil {
		return c.Render(http.StatusBadRequest, Fail(err.Error()))
	}
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Post
	post := new(models.Post)

	// To find the Post the parameter post_id is used.
	if err := tx.Eager("Tags").Find(post, c.Param("post_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}
	if append {
		if like {
			post.Like(user)
		}
	} else {
		if like {
			post.UnLike(user)
		}
	}
	return nil
}

// Destroy deletes a Post from the DB. This function is mapped
// to the path DELETE /posts/{post_id}
func (v PostsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Post
	post := &models.Post{}

	user, err := currentUser(c)
	if err != nil {
		return c.Render(http.StatusBadRequest, Fail(err.Error()))
	}
	// To find the Post the parameter post_id is used.
	if err := tx.Where("user_id = ?", user.ID).Find(post, c.Param("post_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}
	post.IsDelete = true
	if err := tx.Update(post); err != nil {
		return err
	}

	return responder.Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(post))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(post))
	}).Respond(c)
}

//PublishPost PublishPost
type PublishPost struct {
	Project  uuid.UUID   `json:"project"`
	Tags     []uuid.UUID `json:"tags"`
	Image    string      `json:"image"`
	Content  string      `json:"content"`
	IsDelete bool        `json:"is_delete"`
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *PublishPost) Validate() (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: p.Image, Name: "Image", Message: "发布图片不能为空"},
		&validators.FuncValidator{
			Fn: func() bool {
				return strings.HasPrefix(p.Image, envy.Get("prefix", "https://v2cs-oss.oss-cn-beijing.aliyuncs.com/"))
			},
			Field:   p.Image,
			Name:    "Image",
			Message: "图片格式错误",
		},
		&validators.FuncValidator{
			Fn: func() bool {
				return len(p.Tags) <= 2
			},
			Field:   "Tags",
			Name:    "Tags",
			Message: "标签不能超过2个",
		},
	), nil
}

func currentUser(c buffalo.Context) (*models.User, error) {
	user, ok := c.Session().Get("current_user").(*models.User)
	if !ok {
		return nil, errors.New("未找到当前用户信息")
	}
	return user, nil
}
