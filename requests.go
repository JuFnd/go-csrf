package main

type SignupRequest struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigninRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Response struct {
	Status int `json:"status"`
	Body   any `json:"body"`
}
