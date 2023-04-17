package entity

type User struct {
	ID       int    `json:"-"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
