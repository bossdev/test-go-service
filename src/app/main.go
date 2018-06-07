package main

import (

	// "fmt"
	// "net/http"
	"app/db"
	"app/routers"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	var cfEnv db.EnvDB
	errOp := godotenv.Load(".env")
	if errOp != nil {
		log.Println(errOp)
	}
	cfEnv.Host = os.Getenv("DB_HOST")
	cfEnv.Dbname = os.Getenv("DB_DATABASE")
	cfEnv.Port = os.Getenv("DB_PORT")
	rt := routers.Routers{cfEnv.Connect(), e}
	rt.GetRouter()
	e.Start(":8080")
}
