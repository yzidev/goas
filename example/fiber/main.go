//go:build fiber && !typed && !security

package main

import (
	"net/http"

	fiberlib "github.com/gofiber/fiber/v2"

	"github.com/aizacoders/openapigo/adapters/fiber"
	"github.com/aizacoders/openapigo/openapi"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	r := fiber.New()

	r.GET("/users", func(c *fiberlib.Ctx) error {
		return fiber.JSON(c, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
	})

	r.POST("/users", func(c *fiberlib.Ctx) error {
		return c.SendStatus(http.StatusCreated)
	})

	fiber.Register(r, openapi.Config{Title: "User API", Version: "1.0.0"})
	_ = r.App.Listen(":8080")
}
