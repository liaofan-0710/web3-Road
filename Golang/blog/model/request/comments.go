package request

type Comment struct {
	Content string `json:"content"`
	PostID  uint   `json:"post_id"`
}

type CommentQuery struct {
	PostID uint `json:"post_id"`
}
