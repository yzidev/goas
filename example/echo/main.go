//go:build echo && !typed && !security

package main

import (
	"net/http"

	echolib "github.com/labstack/echo/v4"

	"github.com/aizacoders/openapigo/adapters/echo"
	"github.com/aizacoders/openapigo/openapi"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	r := echo.New()

	r.GET("/users", func(c echolib.Context) error {
		return echo.JSON(c, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
	})

	r.POST("/users", func(c echolib.Context) error {
		return c.NoContent(http.StatusCreated)
	})

	echo.Register(r, openapi.Config{Title: "User API", Version: "1.0.0"})
	_ = r.Echo.Start(":8080")
}
