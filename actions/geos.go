package actions

import (
	"culture/cache"
	"culture/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
)

//GeosResource GeosResource
type GeosResource struct {
	buffalo.Resource
}

//List default implementation.
func (v GeosResource) List(c buffalo.Context) error {
	ids, err := models.REDIS.ZRange("project_geo", 0, -1).Result()
	if err != nil {
		return c.Render(http.StatusBadRequest, Fail("查询地理位置集合错误 : %v", err))
	}
	geos := make([]GeoResult, 0, len(ids))
	if err := cache.Once("cache:project_geo", &geos, func() (interface{}, error) {
		//Get the DB connection from the context
		tx, ok := c.Value("tx").(*pop.Connection)
		if !ok {
			return nil, fmt.Errorf("no transaction found")
		}
		var types []Types
		if err := tx.RawQuery("select id,type from projects").All(&types); err != nil {
			return nil, c.Render(http.StatusBadRequest, Fail("查询项目信息错误 : %v", err))
		}
		typeMap := map[string]string{}
		for _, t := range types {
			typeMap[t.ID] = t.Type
		}
		for _, id := range ids {
			geo, err := models.REDIS.GeoPos("project_geo", id).Result()
			if err != nil {
				return nil, c.Render(http.StatusBadRequest, Fail("查询地理位置错误 : %v", err))
			}
			var longitude string
			var latitude string
			if len(geo) > 0 {
				longitude = strconv.FormatFloat(geo[0].Longitude, 'f', -1, 64)
				latitude = strconv.FormatFloat(geo[0].Latitude, 'f', -1, 64)
			} else {
				longitude = "0"
				latitude = "0"
			}
			geos = append(geos, GeoResult{
				Longitude: longitude,
				Latitude:  latitude,
				ID:        id,
				Type:      typeMap[id],
			})
		}
		return geos, nil
	}, time.Hour*3); err != nil {
		return c.Render(http.StatusBadRequest, Fail("加载缓存数据失败 %v", err))
	}
	return c.Render(http.StatusOK, r.JSON(geos))
}

//GeoResult GeoResult
type GeoResult struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}

// Types Types
type Types struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
