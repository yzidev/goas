//go:build gin && typed && !security

package main

import (
	"net/http"

	ginlib "github.com/gin-gonic/gin"

	"github.com/aizacoders/openapigo/adapters/gin"
	"github.com/aizacoders/openapigo/openapi"
)

type UserTyped struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateUserTyped struct {
	Name string `json:"name"`
}

func main() {
	r := gin.New()

	// Full auto schema (request+response inferred)
	gin.POSTT[CreateUserTyped, UserTyped](r, "/typed/users", func(c *ginlib.Context, in CreateUserTyped) (UserTyped, int, error) {
		return UserTyped{ID: "1", Name: in.Name}, http.StatusCreated, nil
	})

	// No request body
	gin.GETT[struct{}, []UserTyped](r, "/typed/users", func(c *ginlib.Context, _ struct{}) ([]UserTyped, int, error) {
		return []UserTyped{{ID: "1", Name: "Alice"}}, http.StatusOK, nil
	})

	gin.Register(r, openapi.Config{Title: "User API", Version: "1.0.0"})
	_ = r.Engine.Run(":8080")
}
