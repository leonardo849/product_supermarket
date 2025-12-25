package user



type UserRepository interface {
	Create(user *User) error
	FindUserByAuthID(authId string) (*User, error)
	DeleteUserByAuthId(authId string) (error)
	FindUserById(id string) (*User, error)
}