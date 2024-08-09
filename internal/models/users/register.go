package users

type RegisterReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
