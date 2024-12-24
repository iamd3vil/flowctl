package models

type Group struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GroupWithUsers struct {
	Group
	Users []User
}
