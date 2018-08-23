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

//CheckBanQuery ...
var CheckBanQuery *sql.Stmt

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
	CheckBanQuery, _ = db.Prepare("SELECT `bid`, `bantype`, `sid`, `mid`, `ends`, `adminname`, `reason` FROM `np_bans` WHERE `steamid` = '?' AND `bRemovedBy` = -1 AND (`ends` > ? OR `length` = 0) ORDER BY `created` DESC")

	return true
}

//CloseDB ...
func CloseDB() {
	db.Close()
	JoinQuery.Close()
	StatsQuery.Close()
	CheckBanQuery.Close()
}
