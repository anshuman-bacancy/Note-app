package routes

import (
	"fmt"
	. "helper"
	"net/http"
	"strings"
	. "templates"
	"time"

	guuid "github.com/google/uuid"
)

// HandleUser handles "/"
func HandleUser(res http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		Tpl.ExecuteTemplate(res, "user.html", nil)
	}
	if req.Method == "POST" {
		UserEmail, _ = req.FormValue("email"), req.FormValue("password")
		userDash := UserDashboard{UserEmail: UserEmail, AllNotes: GetAllNotes(UserEmail)}
		Tpl.ExecuteTemplate(res, "note-dash.html", userDash)
	}
}

// SaveNote handles "/save-note"
func SaveNote(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		//author := req.FormValue("author")
		title := strings.Title(req.FormValue("title"))
		content := req.FormValue("content")

		dtFormat := "01-02-2006 15:04:05"
		curr := time.Now()
		noteDate := curr.Format(dtFormat)

		noteID := guuid.New()

		note := Note{Id: noteID.String(), Author: UserEmail, Title: title, Content: content, CreatedAt: noteDate}
		Save(note)

		userDash := UserDashboard{UserEmail: UserEmail, AllNotes: GetAllNotes(UserEmail)}
		Tpl.ExecuteTemplate(res, "note-dash.html", userDash)
	}
}

// DeleteNote handles "/delete-note"
func DeleteNote(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	delStmt := `DELETE FROM note where id = $1;`
	_, err := Db.Exec(delStmt, id)
	CheckErr(err)

	userDash := UserDashboard{UserEmail: UserEmail, AllNotes: GetAllNotes(UserEmail)}
	Tpl.ExecuteTemplate(res, "note-dash.html", userDash)
}

// UpdateNote handles "/update-note"
func UpdateNote(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	if req.Method == "GET" {
		allNotes := GetAllNotes(UserEmail)
		var noteToUpdate Note

		//find the note to update
		for _, note := range allNotes {
			if note.Id == id {
				noteToUpdate = note
				break
			}
		}
		Tpl.ExecuteTemplate(res, "update-note.html", noteToUpdate)
	}

	if req.Method == "POST" {
		idn := req.FormValue("id")
		title := req.FormValue("title")
		content := req.FormValue("content")

		fmt.Println("new content ..")
		fmt.Println(id, UserEmail, title, content)

		updateStmt := `UPDATE note SET author = $1, title = $2, content = $3 WHERE id = $4;`
		_, updateErr := Db.Exec(updateStmt, UserEmail, title, content, idn)
		CheckErr(updateErr)

		userDash := UserDashboard{UserEmail: UserEmail, AllNotes: GetAllNotes(UserEmail)}
		Tpl.ExecuteTemplate(res, "note-dash.html", userDash)
	}
}