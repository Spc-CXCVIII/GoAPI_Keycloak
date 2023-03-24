package models

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Register struct {
	Title     string `json:"title"`
	IDNumber  string `json:"idnum"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
	Phone     string `json:"phone"`
	TH        bool   `json:"th"`
}
