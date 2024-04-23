package user

type User struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type CreateUser struct {
	Name string `json:"name"`
}
