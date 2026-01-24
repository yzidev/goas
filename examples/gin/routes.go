//go:build gin && !typed && !security

package main

import (
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/aizacoders/openapigo/adapters/ginadapter"
	"github.com/aizacoders/openapigo/openapi"
	"github.com/aizacoders/openapigo/openapi/oas"
)

// registerRoutes wires the endpoints in a readable and grouped way.
// (Non-typed, non-security variant.)

func registerSystemRoutes(r *oas.GinRouter) {
	r.GET("/healthz", handleHealthz)
}

func registerUserRoutes(r *oas.GinRouter) {
	users := r.Group("", ginadapter.WithTags("Users"))

	users.GET("/users", handleListUsers)
	users.GET("/search", handleSearchUsers)
	users.POST("/users", handleCreateUser)
	users.POST("/users/upload", handleUploadUserFile)
	users.GET("/users/demo-errors", handleDemoErrors)
	users.GET("/users/:id", handleGetUser)
	users.PUT("/users/:id", handlePutUser)
	users.PATCH("/users/:id", handlePatchUser)
	users.DELETE("/users/:id", handleDeleteUser)
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

func openapiSpec() oas.Spec {
	b := oas.NewSpec()
	b.GroupTags("", []string{"System"}, func(s *oas.SpecBuilder) {
		s.GET("/healthz").Res(map[string]string{}).OK()
	})
	b.GroupTags("", []string{"Users"}, func(s *oas.SpecBuilder) {
		s.GET("/users").Res([]User{}).OK()
		s.GET("/search").Query(
			openapi.QueryParam{Name: "q", Type: openapi.ParamString, Required: true, Description: "Search term"},
			openapi.QueryParam{Name: "limit", Type: openapi.ParamInteger, Required: false, Description: "Max results"},
		).Res(struct{}{}).OK()

		// Create user: normal endpoint.
		s.POST("/users").Req(CreateUser{}).Res(struct{}{}).Created()

		// Upload user file: multipart/form-data.
		s.POST("/users/upload").MultipartUpload("file", openapi.MultipartField{Name: "note", Type: openapi.ParamString}).Res(map[string]string{}).OK()

		// Dedicated error showcase endpoint (doesn't depend on security mode).
		s.GET("/users/demo-errors").Headers(
			openapi.HeaderParam{Name: "X-Demo-Fail", Type: openapi.ParamString, Required: false, Description: "Set to 400/401/500/503 to simulate an error"},
		).Res(map[string]string{}).OK()

		s.GET("/users/:id").Res(User{}).OK()
		s.PUT("/users/:id").Req(UpdateUser{}).Res(User{}).OK()
		s.PATCH("/users/:id").Req(UpdateUser{}).Res(User{}).OK()
		s.DELETE("/users/:id").Res(struct{}{}).NoContent()
	})
	return b.Spec()
}
