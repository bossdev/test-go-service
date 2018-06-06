package main

import (
	"fmt"
	// "fmt"
	// "net/http"
	"app/db"
	"app/routers"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
)

func main() {
	var ts time.Time
	createdFormat := "2006-01-02 15:04:05"
	now := time.Now()
	loc, _ := time.LoadLocation("Asia/Bangkok")
	currentDate := now.In(loc).Format(createdFormat)
	t, _ := time.Parse(createdFormat, currentDate)
	ts = t
	fmt.Println(ts)

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
