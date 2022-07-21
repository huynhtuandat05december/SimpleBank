package dto

type CreateAccountDTO struct {
	Owner    string `json:"owner,omitempty"`
	Currency string `json:"currency" binding:"required,currency"`
}

type GetListAccountDTO struct {
	Owner string `json:"owner,omitempty"`
	Page  int32  `form:"page" binding:"required"`
	Limit int32  `form:"limit" binding:"required"`
}
