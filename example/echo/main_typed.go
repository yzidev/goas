//go:build echo && typed && !security

package main

import (
	"net/http"

	echolib "github.com/labstack/echo/v4"

	"github.com/aizacoders/openapigo/adapters/echo"
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
	r := echo.New()

	echo.POSTT[CreateUserTyped, UserTyped](r, "/typed/users", func(c echolib.Context, in CreateUserTyped) (UserTyped, int, error) {
		return UserTyped{ID: "1", Name: in.Name}, http.StatusCreated, nil
	})

	echo.GETT[struct{}, []UserTyped](r, "/typed/users", func(c echolib.Context, _ struct{}) ([]UserTyped, int, error) {
		return []UserTyped{{ID: "1", Name: "Alice"}}, http.StatusOK, nil
	})

	echo.Register(r, openapi.Config{Title: "User API", Version: "1.0.0"})
	_ = r.Echo.Start(":8080")
}
