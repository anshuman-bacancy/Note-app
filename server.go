package main

import (
  "fmt"
  "log"
  "time"
  "net/http"
  "database/sql"
  "html/template"
  _ "github.com/lib/pq"
  guuid "github.com/google/uuid"
)

//----------GLOBAL VARIABLES----------
var db *sql.DB
var dbError error
var tpl *template.Template
var userEmail string
//------------------------------------


//----------TEMPLATES----------
type Note struct {
  Id string
  Author string
  Title string
  Content string
  CreatedAt string
}

type UserDashboard struct {
  UserEmail string
  AllNotes []Note
}
//------------------------------------


//----------HELPER FUNCTIONS----------
func checkErr(err error) {
  if err != nil {
    fmt.Println(err)
    //panic(err)
  }
}

func save(note Note) {
  insStmt := `INSERT INTO note (id, author, title, content, createdat) VALUES ($1, $2, $3, $4, $5)`
  _, err := db.Exec(insStmt, note.Id, note.Author, note.Title, note.Content, note.CreatedAt)
  checkErr(err)
}

func getAllNotes() []Note {
  var allNotes []Note
  rows, noteErr := db.Query("select * from note")
  checkErr(noteErr)

  for rows.Next() {
    var note Note
    rows.Scan(&note.Id, &note.Author, &note.Title, &note.Content, &note.CreatedAt)
    allNotes = append(allNotes, note)
  }
  return allNotes
}

func init() {
  //load templates
  tpl = template.Must(template.ParseGlob("static/*.html"))

  //database connection
  connection := "postgres://postgres:password@localhost/notedb?sslmode=disable"
  db, dbError = sql.Open("postgres", connection)
  checkErr(dbError)
  log.Println("Database connection successful...")
}
//------------------------------------


// ----------- ROUTE HANDLERS ----------
func handleUser(res http.ResponseWriter, req *http.Request) {
  if req.Method == "GET" {
    tpl.ExecuteTemplate(res, "user.html", nil)
  }
  if req.Method == "POST" {
    userEmail, _ = req.FormValue("email"), req.FormValue("password")
    //fmt.Println(getAllNotes())
    userDash := UserDashboard{UserEmail: userEmail, AllNotes: getAllNotes()}
    tpl.ExecuteTemplate(res, "note-dash.html", userDash)
  }
}

func saveNote(res http.ResponseWriter, req *http.Request) {
  if req.Method == "POST" {
    //author := req.FormValue("author")
    title := req.FormValue("title")
    content := req.FormValue("content")

    dtFormat := "01-02-2006 15:04:05 Mon"
    curr := time.Now()
    noteDate := curr.Format(dtFormat)

    noteId := guuid.New()

    fmt.Println("author ..", userEmail)
    fmt.Println("title ..", title)
    fmt.Println("content ..", content)

    note := Note{Id: noteId.String(), Author: userEmail, Title: title, Content: content, CreatedAt: noteDate}
    save(note)

    userDash := UserDashboard{UserEmail: userEmail, AllNotes: getAllNotes()}
    tpl.ExecuteTemplate(res, "note-dash.html", userDash)
  }
}

func deleteNote(res http.ResponseWriter, req *http.Request) {
  id := req.URL.Query().Get("id")
  delStmt := `DELETE FROM note where id = $1;`
  _, err := db.Exec(delStmt, id)
  checkErr(err)

  userDash := UserDashboard{UserEmail: userEmail, AllNotes: getAllNotes()}
  tpl.ExecuteTemplate(res, "note-dash.html", userDash)
}

func updateNote(res http.ResponseWriter, req *http.Request) {
  id := req.URL.Query().Get("id")
  fmt.Println(id)

  if req.Method == "GET" {
    allNotes := getAllNotes()
    var noteToUpdate Note

    //find the note to update
    for _, note := range allNotes {
      if note.Id == id {
        noteToUpdate = note
        break
      }
    }
    fmt.Println("note to update ..", noteToUpdate)
    tpl.ExecuteTemplate(res, "update-note.html", noteToUpdate)
  }

  if req.Method == "POST" {
    idn := req.FormValue("id")
    //author := req.FormValue("author")
    title := req.FormValue("title")
    content := req.FormValue("content")

    fmt.Println("new content ..")
    fmt.Println(id, userEmail, title, content)

    updateStmt := `UPDATE note SET author = $1, title = $2, content = $3 WHERE id = $4;`
    _, updateErr := db.Exec(updateStmt, userEmail, title, content, idn)
    checkErr(updateErr)

    userDash := UserDashboard{UserEmail: userEmail, AllNotes: getAllNotes()}
    tpl.ExecuteTemplate(res, "note-dash.html", userDash)
  }
}

func main() {
  defer db.Close()
  log.Println("Server running....")

  http.HandleFunc("/user", handleUser)
  http.HandleFunc("/save-note", saveNote)
  http.HandleFunc("/delete-note", deleteNote)
  http.HandleFunc("/update-note", updateNote)
  //http.HandleFunc("/confirm-update", confirmUpdate)

  http.ListenAndServe(":8080", nil)
}

