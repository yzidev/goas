//go:build fiber && typed && !security

package main

import (
	"net/http"

	fiberlib "github.com/gofiber/fiber/v2"

	"github.com/aizacoders/openapigo/adapters/fiber"
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
	r := fiber.New()

	fiber.POSTT[CreateUserTyped, UserTyped](r, "/typed/users", func(c *fiberlib.Ctx, in CreateUserTyped) (UserTyped, int, error) {
		return UserTyped{ID: "1", Name: in.Name}, http.StatusCreated, nil
	})

	fiber.GETT[struct{}, []UserTyped](r, "/typed/users", func(c *fiberlib.Ctx, _ struct{}) ([]UserTyped, int, error) {
		return []UserTyped{{ID: "1", Name: "Alice"}}, http.StatusOK, nil
	})

	fiber.Register(r, openapi.Config{Title: "User API", Version: "1.0.0"})
	_ = r.App.Listen(":8080")
}
