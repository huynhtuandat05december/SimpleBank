package service

import (
	"context"

	db "github.com/512/simple_bank/db/sqlc"
	"github.com/512/simple_bank/dto"
)

type AccountService interface {
	CreateAccount(account dto.CreateAccountDTO) (db.Account, error)
	GetAccountByID(accountID int64) (db.Account, error)
	GetListAccount(arg dto.GetListAccountDTO) ([]db.Account, error)
}

type accountService struct {
	store *db.Store
}

func NewAccountService(store *db.Store) AccountService {
	return &accountService{
		store: store,
	}

}

func (service *accountService) CreateAccount(account dto.CreateAccountDTO) (db.Account, error) {
	newAccount := db.CreateAccountParams{
		Owner:    account.Owner,
		Currency: account.Currency,
		Balance:  0,
	}
	resultAccount, err := service.store.CreateAccount(context.Background(), newAccount)
	return resultAccount, err
}

func (service *accountService) GetAccountByID(accountID int64) (db.Account, error) {
	resultAccounts, err := service.store.GetAccount(context.Background(), accountID)
	return resultAccounts, err
}

func (service *accountService) GetListAccount(arg dto.GetListAccountDTO) ([]db.Account, error) {
	listAccountParams := db.ListAccountsParams{
		Owner:  arg.Owner,
		Limit:  arg.Limit,
		Offset: (arg.Page - 1) * arg.Limit,
	}
	resultAccounts, err := service.store.ListAccounts(context.Background(), listAccountParams)
	return resultAccounts, err
}
