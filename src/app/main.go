package main

import (
	// "fmt"
	// "net/http"
	// "time"
	"app/db"
	"app/routers"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	rt := routers.Routers{db.Db(), e}
	rt.GetRouter()
	e.Start(":8080")
}
