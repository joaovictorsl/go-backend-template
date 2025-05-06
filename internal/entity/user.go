package entity

type User struct {
	Id         string `json:"id"`
	ProviderId string `json:"provider_id"`
	Email      string `json:"email"`
}
