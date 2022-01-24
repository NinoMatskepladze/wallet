package models

type AddBalanceRequest struct {
	Amount int `json:"amount,required" valid:"type(int)"`
}
