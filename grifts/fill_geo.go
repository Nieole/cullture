package grifts

import (
	"culture/models"
	"log"

	"github.com/go-redis/redis/v7"
	. "github.com/markbates/grift/grift"
)

var _ = Namespace("project", func() {

	Desc("fill_geo", "填充所有项目的地理位置信息")
	Add("fill_geo", func(c *Context) error {
		projects := &models.Projects{}
		if err := models.DB.Select("id").All(projects); err != nil {
			log.Fatalf("quert all project failed : %v", err)
		}
		geos := make([]*redis.GeoLocation, 0, len(*projects))
		for _, project := range *projects {
			geos = append(geos, &redis.GeoLocation{
				Name: project.ID.String(),
			})
		}
		models.REDIS.GeoAdd("project_geo", geos...)
		return nil
	})

})
