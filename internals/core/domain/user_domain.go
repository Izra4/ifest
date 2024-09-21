package domain

type User struct {
	Username string `json:"username"`
	Password string `json:"-"`
	Name     string `json:"name"`
}
