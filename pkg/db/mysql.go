package db

import (
	"database/sql"
	"fmt"
	"golang_im/config"

	_ "github.com/go-sql-driver/mysql"
)

var DBCli *sql.DB

func Mysql_init() {
	var err error
	fmt.Println("config.LogicConf.MySQL", config.DBConf.MySQL)
	DBCli, err = sql.Open("mysql", config.DBConf.MySQL)
	if err != nil {
		panic(err)
	}
}
