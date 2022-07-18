package service

import (
	"context"
	"fmt"

	db "github.com/512/simple_bank/db/sqlc"
	"github.com/512/simple_bank/dto"
)

type TransferService interface {
	CreateTransfer(transferDTO dto.CreateTransferDTO) (db.TransferTxResult, error)
}

type transferService struct {
	store *db.Store
}

func NewTransferService(store *db.Store) TransferService {
	return &transferService{
		store: store,
	}

}

func (service *transferService) CreateTransfer(transferDTO dto.CreateTransferDTO) (db.TransferTxResult, error) {
	arg := db.TransferTxParams{
		FromAccountID: transferDTO.FromAccountID,
		ToAccountID:   transferDTO.ToAccountID,
		Amount:        transferDTO.Amount,
	}
	resultTransfer, err := service.store.TransferTx(context.Background(), arg)
	fmt.Print(resultTransfer)
	return resultTransfer, err
}
