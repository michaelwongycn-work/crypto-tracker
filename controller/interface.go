package controller

import "net/http"

type Controller interface {
	Ping(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)

	ShowUserAsset(w http.ResponseWriter, r *http.Request)
	InsertUserAsset(w http.ResponseWriter, r *http.Request)
	DeleteUserAsset(w http.ResponseWriter, r *http.Request)
}
