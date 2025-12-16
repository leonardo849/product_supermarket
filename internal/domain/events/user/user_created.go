package user

type UserCreated struct {
	ID string `json:"id"`
	AuthUpdatedAt string `json:"auth_updated_at"`
	Role string `json:"role"`
}