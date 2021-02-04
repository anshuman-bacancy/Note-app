package dbserver

import (
	"database/sql"
	. "helper"
	"html/template"
	"log"

	_ "github.com/lib/pq"
)

// DBInit initialized Database and loads HTML templates
func DBInit() {
	//load templates
	Tpl = template.Must(template.ParseGlob("static/*.html"))

	//database connection
	connection := "postgres://postgres:password@localhost/notedb?sslmode=disable"
	Db, DbError = sql.Open("postgres", connection)
	CheckErr(DbError)
	log.Println("Database connection successful...")
}