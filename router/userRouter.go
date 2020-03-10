package router

import (
	handlers "github.com/Alaedeen/goWebProjectTemplate/handlers"
	"github.com/Alaedeen/goWebProjectTemplate/helpers"
	"github.com/gorilla/mux"
)

// UserRouterHandler ...
type UserRouterHandler struct {
	Router  *mux.Router
	Handler handlers.UserHandler
}

// HandleFunctions ...
func (r *UserRouterHandler) HandleFunctions() {
	// Route Handlers / Endpoints
	r.Router.Handle("/api/v1/users", helpers.IsAuthorized(r.Handler.GetUsers)).Methods("GET")
	r.Router.Handle("/api/v1/usersbyname", helpers.IsAuthorized(r.Handler.GetUsersByName)).Methods("GET")
	r.Router.HandleFunc("/api/v1/user", r.Handler.GetUser).Methods("GET")
	r.Router.Handle("/api/v1/userby", helpers.IsAuthorized(r.Handler.GetUserBy)).Methods("GET")
	r.Router.HandleFunc("/api/v1/login", r.Handler.Login).Methods("GET")
	r.Router.HandleFunc("/api/v1/users", r.Handler.CreateUser).Methods("POST")
	r.Router.Handle("/api/v1/users", helpers.IsAuthorized(r.Handler.UpdateUser)).Methods("PUT")
	r.Router.Handle("/api/v1/users", helpers.IsAuthorized(r.Handler.DeleteUser)).Methods("DELETE")
	r.Router.HandleFunc("/api/v1/user/reset_password", r.Handler.ResetPassword).Methods("PUT")
}
