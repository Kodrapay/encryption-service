package dto

type TokenizeCardRequest struct {
	Pan       string `json:"pan"`
	ExpMonth  int    `json:"exp_month"`
	ExpYear   int    `json:"exp_year"`
	Cvv       string `json:"cvv"`
	Reference string `json:"reference"`
}

type TokenResponse struct {
	TokenID string `json:"token_id"`
	Last4   string `json:"last4"`
	Brand   string `json:"brand"`
}
