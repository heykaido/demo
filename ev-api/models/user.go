package models

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Password   string `json:"-"`
	CustomerID string `json:"customer_id"`
	RoleType   string `json:"role_type"`
}
