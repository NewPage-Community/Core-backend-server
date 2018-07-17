package core

import (
	"database/sql"
	"fmt"
	"log"

	//mysql
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

//JoinQuery ...
var JoinQuery *sql.Stmt

//StatsQuery ...
var StatsQuery *sql.Stmt

//ConnectDB ...
func ConnectDB() bool {
	dbi, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", GetConfig("sqluser"), GetConfig("sqlpasswd"), GetConfig("sqladdr"), GetConfig("sqlname")))

	err = dbi.Ping()
	if err != nil {
		log.Println("mysql错误:", err)
		return false
	}

	db = dbi

	JoinQuery, _ = db.Prepare("CALL user_join(?, ?, ?, ?, ?, ?, ?)")
	StatsQuery, _ = db.Prepare("CALL user_stats(?, ?, ?, ?, ?, ?, ?)")

	return true
}

//CloseDB ...
func CloseDB() {
	db.Close()
	JoinQuery.Close()
	StatsQuery.Close()
}
