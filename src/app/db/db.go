package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type EnvDB struct{
	host string
	username string
	passoword string
	dbname string
}

func Db() *gorm.DB {
	var cfEnv EnvDB
	cfEnv.host = "localhost"
	cfEnv.username = "root"
	cfEnv.passoword = ""
	cfEnv.dbname = "service"
	db , err := gorm.Open("mysql", cfEnv.username+":"+cfEnv.passoword+"@tcp(localhost:3306)/"+cfEnv.dbname)
	if err!=nil {
		fmt.Println(err)
	}
	// defer db.Close()
	return db
}