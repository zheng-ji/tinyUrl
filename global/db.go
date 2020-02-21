package global

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var GDB *sql.DB

func InitDBConnection() {
	// mysql 连接
	GDB, _ = sql.Open(GlobalConfig.MySQL.Name, GlobalConfig.MySQL.Connection)
	GDB.SetMaxIdleConns(GlobalConfig.MySQL.MaxIdleConnections)
}
