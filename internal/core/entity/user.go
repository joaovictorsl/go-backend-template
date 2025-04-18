package entity

type User struct {
	Id       uint   `json:"id"`
	GoogleId string `json:"google_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
