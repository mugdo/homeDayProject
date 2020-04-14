package backEnd

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)
func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "ajudge"
	
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	checkErr(err)
	
	return db
}