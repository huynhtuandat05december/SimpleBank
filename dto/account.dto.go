package dto

type CreateAccountDTO struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

type GetListAccountDTO struct {
	Page  int32 `form:"page" binding:"required"`
	Limit int32 `form:"limit" binding:"required"`
}
