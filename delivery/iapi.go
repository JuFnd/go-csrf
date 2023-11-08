package delivery

import "net/http"

type IApi interface {
	SendResponse(w http.ResponseWriter, response Response)
	Films(w http.ResponseWriter, r *http.Request)
	LogoutSession(w http.ResponseWriter, r *http.Request)
	AuthAccept(w http.ResponseWriter, r *http.Request)
	Signin(w http.ResponseWriter, r *http.Request)
	Signup(w http.ResponseWriter, r *http.Request)
	Film(w http.ResponseWriter, r *http.Request)
	Actor(w http.ResponseWriter, r *http.Request)
	Comment(w http.ResponseWriter, r *http.Request)
	AddComment(w http.ResponseWriter, r *http.Request)
	Profile(w http.ResponseWriter, r *http.Request)
}
