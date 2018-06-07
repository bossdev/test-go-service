package db

import (
	"fmt"

	"github.com/globalsign/mgo"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type EnvDB struct {
	Host      string
	Username  string
	Passoword string
	Dbname    string
	Port      string
}

func (b *EnvDB) Connect() *mgo.Database {
	session, err := mgo.Dial(b.Host + ":" + b.Port)
	if err != nil {
		fmt.Println(err)
	}
	Db := session.DB(b.Dbname)
	return Db
}
