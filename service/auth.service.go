package service

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/512/simple_bank/db/sqlc"
	"github.com/512/simple_bank/dto"
	"github.com/512/simple_bank/ultis"
)

type AuthService interface {
	CreateUser(createUserDTO dto.CreateUserDTO) (db.User, error)
	VerifyAccount(loginDTO dto.LoginDTO) (db.User, error)
}

type authService struct {
	store *db.Store
}

func NewAuthService(store *db.Store) AuthService {
	return &authService{
		store: store,
	}
}

func (service *authService) CreateUser(createUserDTO dto.CreateUserDTO) (db.User, error) {
	hashPassword, err := ultis.HashPassword(createUserDTO.Password)
	if err != nil {
		return db.User{}, err
	}
	arg := db.CreateUserParams{
		Username:       createUserDTO.Username,
		HashedPassword: hashPassword,
		FullName:       createUserDTO.FullName,
		Email:          createUserDTO.Email,
	}
	resultUser, err := service.store.CreateUser(context.Background(), arg)
	return resultUser, err
}

func (service *authService) VerifyAccount(loginDTO dto.LoginDTO) (db.User, error) {
	user, err := service.store.GetUser(context.Background(), loginDTO.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.User{}, errors.New("User not found")
		}
		return db.User{}, err
	}
	err = ultis.CheckPassword(loginDTO.Password, user.HashedPassword)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}
