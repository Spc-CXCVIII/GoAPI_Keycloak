package models

type UserDataInToken struct {
	ID   int
	Name string
}

type UserIDToken struct {
	ID string
}

type UserIDRequest struct {
	ID int `json:"userid"`
}
