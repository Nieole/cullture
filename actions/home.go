package actions

import (
	"culture/models"
	"fmt"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	users := &models.Users{}
	tx.All(users)
	post := &models.Post{}
	tx.Find(post, "9adb39a7-6078-4f65-a8ae-64fa44de2870")
	comment := &models.Comment{
		Content: nulls.NewString("aaa"),
		User:    &((*users)[0]),
		Post:    post,
	}
	tx.Save(comment)
	for _, user := range *users {
		c := &models.Comment{
			Content: nulls.NewString("aaa"),
			User:    &user,
			Post:    post,
			Comment: comment,
		}
		tx.Save(c)
		comment = c
	}
	//m1 := &models.Comment{
	//	Content: nulls.NewString("aaa"),
	//	User:    user,
	//	Post:    post,
	//	Comment: &models.Comment{
	//		Content: nulls.NewString("aaa"),
	//		User:    user,
	//		Post:    post,
	//		Comment: &models.Comment{
	//			Content: nulls.NewString("aaa"),
	//			User:    user,
	//			Post:    post,
	//			//	Comment: &models.Comment{
	//			//		Content: nulls.NewString("aaa"),
	//			//		User:    user,
	//			//		Post:    post,
	//			//		Comment: &models.Comment{
	//			//			Content: nulls.NewString("aaa"),
	//			//			User:    user,
	//			//			Post:    post,
	//			//			Comment: &models.Comment{
	//			//				Content: nulls.NewString("aaa"),
	//			//				User:    user,
	//			//				Post:    post,
	//			//				Comment: &models.Comment{
	//			//					Content: nulls.NewString("aaa"),
	//			//					User:    user,
	//			//					Post:    post,
	//			//					Comment: &models.Comment{
	//			//						Content: nulls.NewString("aaa"),
	//			//						User:    user,
	//			//						Post:    post,
	//			//						Comment: &models.Comment{
	//			//							Content: nulls.NewString("aaa"),
	//			//							User:    user,
	//			//							Post:    post,
	//			//							Comment: &models.Comment{
	//			//								Content: nulls.NewString("aaa"),
	//			//								User:    user,
	//			//								Post:    post,
	//			//								Comment: &models.Comment{
	//			//									Content: nulls.NewString("aaa"),
	//			//									User:    user,
	//			//									Post:    post,
	//			//									Comment: &models.Comment{
	//			//										Content: nulls.NewString("aaa"),
	//			//										User:    user,
	//			//										Post:    post,
	//			//										Comment: &models.Comment{
	//			//											Content: nulls.NewString("aaa"),
	//			//											User:    user,
	//			//											Post:    post,
	//			//										},
	//			//									},
	//			//								},
	//			//							},
	//			//						},
	//			//					},
	//			//				},
	//			//			},
	//			//		},
	//			//	},
	//		},
	//	},
	//}
	//if err := tx.Eager().Save(m1); err != nil {
	//	fmt.Println(err)
	//}
	return c.Render(200, r.JSON(map[string]string{"message": "Welcome to Buffalo!"}))
}
