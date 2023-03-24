package models

type UserDataInToken struct {
	ID   int
	Name string
}

type UserID struct {
	ID int
}

type UserIDToken struct {
	ID string
}

type UserIDRequest struct {
	ID int `json:"userid"`
}
