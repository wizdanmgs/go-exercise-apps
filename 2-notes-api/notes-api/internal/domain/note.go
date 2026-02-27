package domain

// Note is the core business entity.
// It contains no framework or infrastructure dependency.
type Note struct {
	ID      string
	Title   string
	Content string
}
