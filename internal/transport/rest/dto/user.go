package dto

type UserSignUpRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Pwd   string `json:"pwd"`
}

type UserSignUpResponse struct {
	Token string `json:"token"`
}

type UserSignInRequest struct {
	Email string `json:"email"`
	Pwd   string `json:"pwd"`
}

type UserSignInResponse struct {
	Token string `json:"token"`
}
