package service

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var MysqlClient *sql.DB

func MysqlInit() {
	config := GetConfig()
	//var err error
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", config.Mysql.Username, config.Mysql.Password, config.Mysql.Host, config.Mysql.Port, config.Mysql.Database))
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	MysqlClient = db
}
