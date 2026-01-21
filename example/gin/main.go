//go:build gin && !typed && !security

package main

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	ginlib "github.com/gin-gonic/gin"

	"github.com/aizacoders/openapigo/adapters/gin"
	"github.com/aizacoders/openapigo/openapi"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	r := gin.New()

	r.GET("/users", func(c *ginlib.Context) {
		gin.JSON(c, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
	}, gin.WithTags("Users"), gin.WithResponses(
		openapi.ResponseSpec{Status: http.StatusOK, Schema: []User{}, Description: "OK"},
		openapi.ResponseSpec{Status: http.StatusInternalServerError, Schema: ErrorResponse{}, Description: "Internal Server Error"},
	))

	r.POST("/users", func(c *ginlib.Context) {
		c.Status(http.StatusCreated)
	}, gin.WithTags("Users"), gin.WithResponses(
		openapi.ResponseSpec{Status: http.StatusCreated, Schema: struct{}{}, Description: "Created"},
		openapi.ResponseSpec{Status: http.StatusBadRequest, Schema: ErrorResponse{}, Description: "Bad Request"},
		openapi.ResponseSpec{Status: http.StatusInternalServerError, Schema: ErrorResponse{}, Description: "Internal Server Error"},
	))

	gin.Register(r, openapi.Config{
		Title:   "User API",
		Version: "1.0.0",
		Tags: openapi3.Tags{
			{Name: "Users", Description: "User management endpoints"},
		},
	})
	_ = r.Engine.Run(":8080")
}
