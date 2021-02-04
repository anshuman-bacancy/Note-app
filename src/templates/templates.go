package templates

// Signifies a note
type Note struct {
	Id        string
	Author    string
	Title     string
	Content   string
	CreatedAt string
}

// Show data on User dashboard
type UserDashboard struct {
	UserEmail string
	AllNotes  []Note
}
