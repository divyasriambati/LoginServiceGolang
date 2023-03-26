package userdb

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const dbname = "mysql"
const dburl = "root:Divya@123@tcp(localhost:3306)/golangtestdb"

func Connectdb() (*sql.DB, error) {
	db, err := sql.Open(dbname, dburl)
	if err != nil {
		panic(err.Error())
	}

	return db, nil
}
