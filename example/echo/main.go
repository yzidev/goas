//go:build echo && !typed && !security

package main

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	echolib "github.com/labstack/echo/v4"

	"github.com/aizacoders/openapigo/adapters/echo"
	"github.com/aizacoders/openapigo/openapi"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type UpdateUser struct {
	Name string `json:"name"`
}

type CreateUser struct {
	Name string `json:"name"`
}

func main() {
	r := echo.New()

	users := r.Group("", echo.WithTags("Users"))

	users.GET("/users", func(c echolib.Context) error {
		return echo.JSON(c, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
	})

	users.GET("/search", func(c echolib.Context) error {
		_ = c.QueryParam("q")
		return c.NoContent(http.StatusOK)
	}, echo.WithQueryParams(
		openapi.QueryParam{Name: "q", Type: openapi.ParamString, Required: true, Description: "Search term"},
		openapi.QueryParam{Name: "limit", Type: openapi.ParamInteger, Required: false, Description: "Max results"},
	))

	users.POST("/users", func(c echolib.Context) error {
		var in CreateUser
		if err := c.Bind(&in); err != nil {
			return echo.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
		}
		return c.NoContent(http.StatusCreated)
	}, echo.JSONRoute(CreateUser{}, struct{}{}, http.StatusCreated)...)

	users.GET("/users/:id", func(c echolib.Context) error {
		id := c.Param("id")
		if id == "404" {
			return echo.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
		}
		return echo.JSON(c, http.StatusOK, User{ID: id, Name: "Alice"})
	})

	users.PUT("/users/:id", func(c echolib.Context) error {
		id := c.Param("id")
		var in UpdateUser
		if err := c.Bind(&in); err != nil {
			return echo.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
		}
		if id == "404" {
			return echo.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
		}
		return echo.JSON(c, http.StatusOK, User{ID: id, Name: in.Name})
	}, echo.JSONRoute(UpdateUser{}, User{}, http.StatusOK)...)

	users.PATCH("/users/:id", func(c echolib.Context) error {
		id := c.Param("id")
		var in UpdateUser
		if err := c.Bind(&in); err != nil {
			return echo.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
		}
		if id == "404" {
			return echo.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
		}
		return echo.JSON(c, http.StatusOK, User{ID: id, Name: in.Name})
	}, echo.JSONRoute(UpdateUser{}, User{}, http.StatusOK)...)

	users.DELETE("/users/:id", func(c echolib.Context) error {
		id := c.Param("id")
		if id == "404" {
			return echo.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
		}
		return c.NoContent(http.StatusNoContent)
	}, echo.JSONRoute(nil, struct{}{}, http.StatusNoContent)...)

	echo.Register(r, openapi.Config{
		Title:   "User API",
		Version: "1.0.0",
		Tags: openapi3.Tags{
			{Name: "Users", Description: "User management endpoints"},
		},
	})
	_ = r.Echo.Start(":8080")
}
