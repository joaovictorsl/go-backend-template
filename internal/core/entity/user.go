package entity

type User struct {
	Id         uint   `json:"id"`
	ProviderId string `json:"provider_id"`
	Email      string `json:"email"`
}
