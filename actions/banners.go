package actions

import (
	"culture/cache"
	"culture/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/x/responder"
)

//BannersResource BannersResource
type BannersResource struct {
	buffalo.Resource
}

//List Banner List
func (b BannersResource) List(c buffalo.Context) error {
	banners := &models.Banners{}
	if err := cache.Once(fmt.Sprintf("cache:banners"), banners, func() (interface{}, error) {
		// Get the DB connection from the context
		tx, ok := c.Value("tx").(*pop.Connection)
		if !ok {
			return nil, fmt.Errorf("no transaction found")
		}

		if err := tx.Where("is_delete = ?", false).Order("sort asc").All(banners); err != nil {
			return nil, err
		}
		return banners, nil
	}, time.Minute*5); err != nil {
		return c.Render(http.StatusBadRequest, Fail("加载缓存数据失败 %v", err))
	}

	return responder.Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(banners))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(banners))
	}).Respond(c)
}
