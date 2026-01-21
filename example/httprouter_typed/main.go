//go:build !security

package main

import (
	"net/http"

	"github.com/aizacoders/openapigo/adapters/httprouter"
	"github.com/aizacoders/openapigo/openapi"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateUser struct {
	Name string `json:"name"`
}

func main() {
	r := httprouter.New()

	// Full auto schema: request schema inferred from CreateUser, response schema from User.
	httprouter.POSTT[CreateUser, User](r, "/users", func(w http.ResponseWriter, req *http.Request, in CreateUser) (User, int, error) {
		_ = req
		return User{ID: "1", Name: in.Name}, http.StatusCreated, nil
	})

	// No request body: use struct{}
	httprouter.GETT[struct{}, []User](r, "/users", func(w http.ResponseWriter, req *http.Request, _ struct{}) ([]User, int, error) {
		_ = w
		_ = req
		return []User{{ID: "1", Name: "Alice"}}, http.StatusOK, nil
	})

	openapi.Register(r, openapi.Config{Title: "User API", Version: "1.0.0"})
	_ = http.ListenAndServe(":8080", r)
}
