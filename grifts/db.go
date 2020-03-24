package grifts

import (
	"culture/actions"
	"culture/models"

	"github.com/gobuffalo/nulls"
	"github.com/markbates/grift/grift"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		// Add DB seeding stuff here
		tags := models.Tags{}
		for i := 0; i < 10; i++ {
			tags = append(tags, models.Tag{
				Name: actions.RandString(5),
				Code: nulls.NewString(actions.RandString(10)),
			})
		}
		models.DB.Save(tags)
		projects := models.Projects{}
		for i := 0; i < 10; i++ {
			projects = append(projects, models.Project{
				Name:     actions.RandString(5),
				RemoteID: nulls.NewString(actions.RandString(20)),
			})
		}
		models.DB.Save(projects)
		return nil
	})

})
