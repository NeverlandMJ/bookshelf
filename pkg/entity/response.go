package entity

type Book struct {
	ID        int    `json:"id,omitempty"`
	Isbn      string `json:"isbn,omitempty"`
	Title     string `json:"title,omitempty"`
	Author    string `json:"author,omitempty"`
	Published int    `json:"published,omitempty"`
	Pages     int    `json:"pages,omitempty"`
}
type ResponseBook struct {
	Book   Book `json:"book"`
	Status int  `json:"status"`
}

type ResponseUser struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

type Response struct {
	Data    interface{} `json:"data"`
	IsOk    bool        `json:"isOk"`
	Message string      `json:"message"`
}
