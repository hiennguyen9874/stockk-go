package users

import "net/http"

type Handlers interface {
	Create() func(w http.ResponseWriter, r *http.Request)
	Get() func(w http.ResponseWriter, r *http.Request)
	GetMulti() func(w http.ResponseWriter, r *http.Request)
	Delete() func(w http.ResponseWriter, r *http.Request)
	Update() func(w http.ResponseWriter, r *http.Request)
	SignIn() func(w http.ResponseWriter, r *http.Request)
	Me() func(w http.ResponseWriter, r *http.Request)
	UpdateMe() func(w http.ResponseWriter, r *http.Request)
	UpdatePassword() func(w http.ResponseWriter, r *http.Request)
	UpdatePasswordMe() func(w http.ResponseWriter, r *http.Request)
	RefreshToken() func(w http.ResponseWriter, r *http.Request)
	GetPublicKey() func(w http.ResponseWriter, r *http.Request)
	Logout() func(w http.ResponseWriter, r *http.Request)
	LogoutAllToken() func(w http.ResponseWriter, r *http.Request)
	LogoutAllAdmin() func(w http.ResponseWriter, r *http.Request)
	VerifyEmail() func(w http.ResponseWriter, r *http.Request)
	ForgotPassword() func(w http.ResponseWriter, r *http.Request)
	ResetPassword() func(w http.ResponseWriter, r *http.Request)
}
