package balance

type Balance struct {
	Current   float32 `json:"current"`
	Withdrawn int     `json:"withdrawn"`
}
