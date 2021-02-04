package helper

import (
	"database/sql"
	"html/template"
	"templates"
)

var Db *sql.DB
var DbError error
var Tpl *template.Template
var UserEmail string


//CheckErr logs error
func CheckErr(err error) { 
	if err != nil {
		panic(err)
	}
}

// Save to database
func Save(note templates.Note) {
	insStmt := `INSERT INTO note (id, author, title, content, createdat) VALUES ($1, $2, $3, $4, $5)`
	_, err := Db.Exec(insStmt, note.Id, note.Author, note.Title, note.Content, note.CreatedAt)
	CheckErr(err)
}

// GetAllNotes returns all the notes from Db
func GetAllNotes(userEmail string) []templates.Note {
	var allNotes []templates.Note
	rows, noteErr := Db.Query(`select * from note WHERE author = $1;`, userEmail)
	CheckErr(noteErr)

	for rows.Next() {
		var note templates.Note
		rows.Scan(&note.Id, &note.Author, &note.Title, &note.Content, &note.CreatedAt)
		allNotes = append(allNotes, note)
	}
	return allNotes
}
