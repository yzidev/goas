//go:build gin && !typed && !security

package main

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/aizacoders/openapigo/adapters/gin"
	"github.com/aizacoders/openapigo/openapi"
)

// registerRoutes wires the endpoints in a readable and grouped way.
// (Non-typed, non-security variant.)

func registerSystemRoutes(r *gin.Router) {
	// No group: show simplest usage.
	opts := append([]gin.HandlerOption{gin.WithTags("System")}, gin.JSONRoute(struct{}{}, map[string]string{}, http.StatusOK)...)
	r.GET("/healthz", handleHealthz, opts...)
}

func registerUserRoutes(r *gin.Router) {
	users := r.Group("", gin.WithTags("Users"))

	users.GET("/users", handleListUsers)

	users.GET("/search",
		handleSearchUsers,
		gin.WithQueryParams(
			openapi.QueryParam{Name: "q", Type: openapi.ParamString, Required: true, Description: "Search term"},
			openapi.QueryParam{Name: "limit", Type: openapi.ParamInteger, Required: false, Description: "Max results"},
		),
		gin.WithResponses(openapi.ResponseSpec{Status: http.StatusOK, Schema: struct{}{}, Description: "OK"}),
	)

	// Schemas declared via JSONRoute; handler stays plain gin.HandlerFunc.
	users.POST("/users", handleCreateUser, gin.JSONRoute(CreateUser{}, struct{}{}, http.StatusCreated)...)
	users.GET("/users/:id", handleGetUser)
	users.PUT("/users/:id", handlePutUser, gin.JSONRoute(UpdateUser{}, User{}, http.StatusOK)...)
	users.PATCH("/users/:id", handlePatchUser, gin.JSONRoute(UpdateUser{}, User{}, http.StatusOK)...)
	users.DELETE("/users/:id", handleDeleteUser, gin.JSONRoute(nil, struct{}{}, http.StatusNoContent)...)
}

func openAPICfg() openapi.Config {
	return openapi.Config{
		Title:   "User API",
		Version: "1.0.0",
		Tags: openapi3.Tags{
			{Name: "Users", Description: "User management endpoints"},
		},
	}
}
