package postgres

import (
	"errors"

	domainUser "github.com/leonardo849/product_supermarket/internal/domain/user"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func toDomainUser(model *UserModel) *domainUser.User {
	return  &domainUser.User{
		ID: model.ID,
		AuthID: model.AuthId,
		Role: domainUser.Role(model.Role),
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return  &UserRepository{
		db: db,
	}
}

func (u *UserRepository) Create(user *domainUser.User) error {
	model := UserModel{
		ID: user.ID,
		AuthId: user.AuthID,
		Role: string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return  u.db.Create(&model).Error
}

func (u *UserRepository) FindUserByAuthID(authId string) (*domainUser.User, error) {
	var user UserModel
	err := u.db.First(&user, "auth_id = ?", authId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainUser.ErrUserNotFound
		}
		return nil, err
	}
	return toDomainUser(&user), nil
}