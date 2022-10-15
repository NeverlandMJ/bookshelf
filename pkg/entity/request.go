package entity

type UserSignUpRequest struct {
	Name   string `json:"name"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

type CreatBookRequest struct {
	Isbn string `json:"isbn"`
}

type EditBook struct {
	Status int `json:"status"`
}

type Info struct {
	NumberOfPages int `json:"number_of_pages"`
	Authors       []struct {
		Key string `json:"key"`
	} `json:"authors"`
	Title       string `json:"title"`
	PublishDate string `json:"publish_date"`
}

type Author struct {
	PersonalName string `json:"personal_name,omitempty"`
}
