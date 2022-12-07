package types

type UserLogin struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}