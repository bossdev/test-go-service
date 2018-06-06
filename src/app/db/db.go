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
	fmt.Println("come on")
	fmt.Println(b)
	session, err := mgo.Dial(b.Host + ":" + b.Port)
	if err != nil {
		fmt.Println(err)
	}
	Db := session.DB(b.Dbname)
	return Db
}

// func Db() *gorm.DB {
// var cfEnv EnvDB

// cfEnv.host = os.Getenv("DB_HOST")
// cfEnv.username = os.Getenv("DB_USER")
// cfEnv.passoword = os.Getenv("DB_PASSWORD")
// cfEnv.dbname = os.Getenv("DB_DATABASE")
// cfEnv.port = os.Getenv("DB_PORT")

// // Host := []string{
// // 	cfEnv.host + ":" + cfEnv.port,
// // 	// replica set addrs...
// // }

// // db, err := gorm.Open("mysql", cfEnv.username+":"+cfEnv.passoword+"@tcp("+cfEnv.host+":"+cfEnv.port+")/"+cfEnv.dbname)
// // db, err := mgo.DialWithInfo(&mgo.DialInfo{
// // 	Addrs: Host,
// // 	// 	// Username: Username,
// // 	// 	// Password: Password,
// // 	// 	// Database: Database,
// // 	// 	// DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
// // 	// 	// 	return tls.Dial("tcp", addr.String(), &tls.Config{})
// // 	// 	// },
// // })

// session, err := mgo.Dial(cfEnv.host + ":" + cfEnv.port)
// db := session.DB("service")
// if err != nil {
// 	fmt.Println(err)
// }
// return db
// }
