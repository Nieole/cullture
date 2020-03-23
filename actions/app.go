package actions

import (
	"culture/middleware"
	"culture/models"
	"encoding/json"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	paramlogger "github.com/gobuffalo/mw-paramlogger"

	"github.com/gobuffalo/buffalo-pop/pop/popmw"
	contenttype "github.com/gobuffalo/mw-contenttype"
	"github.com/rs/cors"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env: ENV,
			PreWares: []buffalo.PreWare{
				cors.Default().Handler,
			},
			SessionName: "_culture_session",
		})
		app.ErrorHandlers[http.StatusUnauthorized] = func(status int, err error, c buffalo.Context) error {
			res := c.Response()
			res.Header().Set("Content-Type", "application/json")
			res.WriteHeader(http.StatusUnauthorized)
			bytes, err := json.Marshal(map[string]string{"message": http.StatusText(http.StatusUnauthorized)})
			if err != nil {
				return err
			}
			res.Write(bytes)
			return nil
		}

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Set the request content type to JSON
		app.Use(contenttype.Set("application/json"))

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))

		app.GET("/", HomeHandler)
		app.POST("/login/{phone}", LoginHandler)

		app.Resource("/posts", PostsResource{})
		app.Resource("/tags", TagsResource{})
		app.Resource("/projects", ProjectsResource{})
		//app.Resource("/post_tags", PostTagsResource{})

		auth := app.Group("/auth")
		auth.Use(middleware.LoginMiddleware)
		auth.DELETE("/signout", SignOutHandler)
	}

	return app
}
