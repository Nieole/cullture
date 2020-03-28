package actions

import (
	"culture/models"
	"net/http"
	"strconv"

	"github.com/gobuffalo/buffalo"
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
	for _, id := range ids {
		geo, err := models.REDIS.GeoPos("project_geo", id).Result()
		if err != nil {
			return c.Render(http.StatusBadRequest, Fail("查询地理位置错误 : %v", err))
		}
		if len(geo) > 0 {
			geos = append(geos, GeoResult{
				Longitude: strconv.FormatFloat(geo[0].Longitude, 'f', -1, 64),
				Latitude:  strconv.FormatFloat(geo[0].Latitude, 'f', -1, 64),
				ID:        id,
			})
		} else {
			geos = append(geos, GeoResult{
				Longitude: "0",
				Latitude:  "0",
				ID:        id,
			})
		}
	}
	return c.Render(http.StatusOK, r.JSON(geos))
}

//GeoResult GeoResult
type GeoResult struct {
	ID        string `json:"id"`
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}
