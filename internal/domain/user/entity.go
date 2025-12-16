package user

import (
	// "regexp"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID
	AuthID string 
	Role Role
	AuthUpdatedAt time.Time
	CreatedAt time.Time 
	UpdatedAt time.Time 
}

func (u *User) CanUserCreateOrEditAProduct() bool {
	return !(u.Role != ROLEMANAGER && u.Role != ROLEWORKER)
} 

func (u *User) UserWasUpdatedAfterToken(iat float64) bool {
	issuedAtUnix := int64(iat) 
	issuedAtTime := time.Unix(issuedAtUnix, 0)
	return issuedAtTime.Before(u.AuthUpdatedAt) 

}

func New(AuthID string, role Role, authUpdatedAt string) (*User, error){
	t, err := time.Parse(time.RFC3339, authUpdatedAt)
	if err != nil {
		return  nil, err
	}
	if !role.isValid() {
		return  nil,  ErrRoleInvalid
	}
	// var mongoIDRegex = regexp.MustCompile("^[a-fA-F0-9]{24}$")

	// if mongoIDRegex.MatchString(AuthID) {
	// 	return nil, ErrItIsNotAMongoID
	// }
	now := time.Now()
	return  &User{
		ID: uuid.New(),
		AuthID: AuthID,
		Role: role,
		CreatedAt: now,
		AuthUpdatedAt: t,
		UpdatedAt: now,
	}, nil
}

