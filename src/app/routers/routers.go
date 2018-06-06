package routers

import (
	// "net/http"
	"app/controllers"

	"github.com/globalsign/mgo"
	"github.com/labstack/echo"
)

type Routers struct {
	Db *mgo.Database
	Ec *echo.Echo
}

func (rts Routers) GetRouter() {
	ctrl := controllers.BlogController{rts.Db}
	rts.Ec.GET("blog/test", ctrl.TestTimeIndex)
	rts.Ec.GET("blog", ctrl.IndexBlog)
	rts.Ec.GET("blog/:id", ctrl.GetBlog)
	rts.Ec.POST("blog", ctrl.StoreBlog)
	rts.Ec.PATCH("blog/:id", ctrl.UpdateBlog)
	rts.Ec.DELETE("blog/:id", ctrl.DeleteBlog)
	rts.Ec.GET("*", func(c echo.Context) error {
		return c.String(404, "404 Page not found")
	})
}
