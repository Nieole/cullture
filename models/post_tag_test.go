package models

import (
	"github.com/gobuffalo/nulls"
)

func (m *ModelSuite) Test_PostTag() {
	m.DB.Eager().Create(&Post{
		Project: &Project{
			Name:     "aa",
			RemoteID: nulls.NewString("bb"),
		},
		Image:  nulls.NewString("cc"),
		UserID: "dd",
		//Tags: &Tags{
		//	Tag{
		//		Name: "ee",
		//	},
		//},
	})
}
