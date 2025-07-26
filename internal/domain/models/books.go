package models

type Book struct {
	BookID      string `json:"book_id,omitempty"`
	Author      string `json:"author"`
	Lable       string `json:"lable"`
	Description string `json:"description"`
	Genre       string `json:"genre"`
	WritedAt    string `json:"writed_at"`
	Count       int    `json:"count,omitempty"`
}
