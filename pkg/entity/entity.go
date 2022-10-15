package entity

type UserResponseFromDatabse struct {
	ID     int    `db:"id"`
	Name   string `db:"name"`
	Key    string `db:"key"`
	Secret string `db:"secret"`
}

type BookResponseFromDatabase struct {
	ID        int    `db:"id"`
	Isbn      string `db:"isbn"`
	Title     string `db:"title"`
	Author    string `db:"author"`
	Published int    `db:"published"`
	Pages     int    `db:"pages"`
	Status    int    `db:"status"`
}
