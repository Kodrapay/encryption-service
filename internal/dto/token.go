package dto

type TokenizeCardRequest struct {
	Pan       string `json:"pan"`
	ExpMonth  int    `json:"exp_month"`
	ExpYear   int    `json:"exp_year"`
	Cvv       int    `json:"cvv"`
	Reference int    `json:"reference"`
}

type TokenResponse struct {
	TokenID int `json:"token_id"`
	Last4   string `json:"last4"`
	Brand   string `json:"brand"`
}
