//go:build fiber && !typed && !security

package main

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	fiberlib "github.com/gofiber/fiber/v2"

	"github.com/aizacoders/openapigo/adapters/fiber"
	"github.com/aizacoders/openapigo/openapi"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UpdateUser struct {
	Name string `json:"name"`
}

type CreateUser struct {
	Name string `json:"name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	r := fiber.New()

	users := r.Group("", fiber.WithTags("Users"))

	users.GET("/users", func(c *fiberlib.Ctx) error {
		return fiber.JSON(c, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
	}, fiber.JSONRoute(nil, []User{}, http.StatusOK)...)

	users.GET("/search", func(c *fiberlib.Ctx) error {
		_ = c.Query("q")
		return c.SendStatus(http.StatusOK)
	},
		fiber.WithQueryParams(
			openapi.QueryParam{Name: "q", Type: openapi.ParamString, Required: true, Description: "Search term"},
			openapi.QueryParam{Name: "limit", Type: openapi.ParamInteger, Required: false, Description: "Max results"},
		),
		fiber.JSONRoute(nil, struct{}{}, http.StatusOK)...,
	)

	users.POST("/users", func(c *fiberlib.Ctx) error {
		var in CreateUser
		if err := fiber.Bind(c, &in); err != nil || in.Name == "" {
			_ = fiber.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return nil
		}
		return c.SendStatus(http.StatusCreated)
	}, fiber.JSONRoute(CreateUser{}, struct{}{}, http.StatusCreated)...)

	users.GET("/users/:id", func(c *fiberlib.Ctx) error {
		id := c.Params("id")
		if id == "404" {
			return fiber.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
		}
		return fiber.JSON(c, http.StatusOK, User{ID: id, Name: "Alice"})
	}, fiber.JSONRoute(nil, User{}, http.StatusOK)...)

	users.PUT("/users/:id", func(c *fiberlib.Ctx) error {
		id := c.Params("id")
		var in UpdateUser
		if err := fiber.Bind(c, &in); err != nil {
			_ = fiber.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return nil
		}
		if id == "404" {
			_ = fiber.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return nil
		}
		return fiber.JSON(c, http.StatusOK, User{ID: id, Name: in.Name})
	}, fiber.JSONRoute(UpdateUser{}, User{}, http.StatusOK)...)

	users.PATCH("/users/:id", func(c *fiberlib.Ctx) error {
		id := c.Params("id")
		var in UpdateUser
		if err := fiber.Bind(c, &in); err != nil {
			_ = fiber.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return nil
		}
		if id == "404" {
			_ = fiber.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return nil
		}
		return fiber.JSON(c, http.StatusOK, User{ID: id, Name: in.Name})
	}, fiber.JSONRoute(UpdateUser{}, User{}, http.StatusOK)...)

	users.DELETE("/users/:id", func(c *fiberlib.Ctx) error {
		id := c.Params("id")
		if id == "404" {
			_ = fiber.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return nil
		}
		return c.SendStatus(http.StatusNoContent)
	}, fiber.JSONRoute(nil, struct{}{}, http.StatusNoContent)...)

	fiber.Register(r, openapi.Config{
		Title:   "User API",
		Version: "1.0.0",
		Tags: openapi3.Tags{
			{Name: "Users", Description: "User management endpoints"},
		},
	})
	_ = r.App.Listen(":8080")
}
