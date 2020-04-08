package actions

import (
	"culture/models"
	"encoding/gob"
	"encoding/json"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/prometheus/common/log"

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
			SessionStore: sessionStore(),
			SessionName:  "_culture_session",
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
		app.POST("/login", LoginHandler)
		app.DELETE("/signout", SignOutHandler)

		app.Resource("/tags", TagsResource{})
		app.Resource("/projects", ProjectsResource{})
		//app.Resource("/post_tags", PostTagsResource{})
		app.Resource("/organizations", OrganizationsResource{})
		app.Resource("/geos", GeosResource{})
		app.Resource("/comments", CommentsResource{})
		app.Resource("/users", UsersResource{}).Use(CheckLoginMiddleware)

		auth := app.Group("/")
		mw := LoginMiddleware
		auth.Use(mw)
		auth.GET("/posts/my", MyList)
		auth.GET("/user/info", func(context buffalo.Context) error {
			user, ok := context.Session().Get("current_user").(*models.User)
			if !ok {
				return context.Render(http.StatusBadRequest, Fail("获取用户信息失败"))
			}
			return context.Render(http.StatusOK, r.JSON(user))
		})
		auth.POST("/like/{post_id}", Like)
		auth.DELETE("/like/{post_id}", UnLike)
		pr := PostsResource{}
		p := auth.Resource("/posts", pr)
		p.Middleware.Skip(mw, pr.List, pr.Show)
	}

	return app
}

func init() {
	gob.Register(&models.User{})
}

func sessionStore() *sessions.CookieStore {
	secret := envy.Get("SESSION_SECRET", "")

	if secret == "" && (ENV == "development" || ENV == "test") {
		secret = "buffalo-secret"
	}

	// In production a SESSION_SECRET must be set!
	if secret == "" {
		log.Warn("Unless you set SESSION_SECRET env variable, your session storage is not protected!")
	}

	cookieStore := sessions.NewCookieStore([]byte(secret))

	//Cookie secure attributes, see: https://www.owasp.org/index.php/Testing_for_cookies_attributes_(OTG-SESS-002)
	cookieStore.Options.HttpOnly = true
	//if ENV == "production" {
	//	cookieStore.Options.Secure = true
	//}
	return cookieStore
}
