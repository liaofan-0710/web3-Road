package request

type Post struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type PostUpdate struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type PostDelete struct {
	ID uint `json:"id"`
}
