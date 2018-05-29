package routers

import (
	// "net/http"
	"app/controllers"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type Routers struct {
	Db *gorm.DB
	Ec *echo.Echo
}

func (rts Routers) GetRouter() {
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World with echo222!")
	// })
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
