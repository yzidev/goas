//go:build !security

package main

import (
	"encoding/json"
	"net/http"

	"github.com/aizacoders/openapigo/openapi"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	r := openapi.NewRouter()

	r.GET("/users", func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode([]User{{ID: "1", Name: "Alice"}})
	})

	r.POST("/users", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	openapi.Register(r, openapi.Config{Title: "User API", Version: "1.0.0"})

	_ = http.ListenAndServe(":8080", r)
}
