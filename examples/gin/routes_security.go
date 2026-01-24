//go:build gin && security && !typed

package main

import (
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/aizacoders/openapigo/adapters/ginadapter"
	"github.com/aizacoders/openapigo/openapi"
	"github.com/aizacoders/openapigo/openapi/oas"
)

func openAPICfgSecurity() (openapi.Config, *openapi3.SecurityRequirement, *openapi3.SecurityRequirement) {
	cfg := openapi.Config{
		Title:       "User API (Gin + Security)",
		Version:     "1.0.0",
		Description: "An examples API with secured endpoints using Gin and OpenAPIGO",
		Tags: openapi3.Tags{
			{Name: "Secure Users", Description: "Secured endpoints (Bearer / X-API-Key)"},
		},
		SecuritySchemes: map[string]*openapi3.SecuritySchemeRef{
			"bearerAuth": {Value: &openapi3.SecurityScheme{Type: "http", Scheme: "bearer", BearerFormat: "JWT"}},
			"apiKeyAuth": {Value: &openapi3.SecurityScheme{Type: "apiKey", In: "header", Name: "X-API-Key"}},
		},
	}
	bearer := openapi3.NewSecurityRequirement().Authenticate("bearerAuth")
	apiKey := openapi3.NewSecurityRequirement().Authenticate("apiKeyAuth")
	return cfg, &bearer, &apiKey
}

func registerSecureRoutes(r *oas.GinRouter, bearer, apiKey *openapi3.SecurityRequirement) {
	b := oas.NewSpec()
	b.GroupTags("", []string{"Secure Users"}, func(s *oas.SpecBuilder) {
		s.GET("/secure/healthz").Security(bearer).Res(map[string]string{}).OK()
		s.GET("/secure/users").Security(bearer).Res([]SecUser{}).OK()
		s.POST("/secure/users").Security(apiKey).Res(struct{}{}).Created()

		// Upload secure user file: multipart/form-data.
		s.POST("/secure/users/upload").Security(apiKey).MultipartUpload("file", openapi.MultipartField{Name: "note", Type: openapi.ParamString}).Res(map[string]string{}).OK()

		s.GET("/secure/demo-errors").Security(bearer).Res(map[string]string{}).OK()
	})
	r.Spec = b.Spec()

	r.GET("/secure/healthz", handleSecureHealthz)

	secure := r.Group("", ginadapter.WithTags("Secure Users"))
	secure.GET("/secure/users", handleSecureListUsers)
	secure.POST("/secure/users", handleSecureCreateUser)
	secure.POST("/secure/users/upload", handleSecureUploadUserFile)
	secure.GET("/secure/demo-errors", handleSecureDemoErrors)
}
