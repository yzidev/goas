//go:build gin && security && !typed

package main

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/aizacoders/openapigo/adapters/gin"
	"github.com/aizacoders/openapigo/openapi"
)

func openAPICfgSecurity() (openapi.Config, *openapi3.SecurityRequirement, *openapi3.SecurityRequirement) {
	cfg := openapi.Config{
		Title:   "User API (Gin + Security)",
		Version: "1.0.0",
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

func registerSecureRoutes(r *gin.Router, bearer, apiKey *openapi3.SecurityRequirement) {
	// No group: keep one example direct on router.
	healthOpts := append(
		[]gin.HandlerOption{gin.WithTags("System"), gin.WithSecurity(bearer)},
		gin.JSONRoute(struct{}{}, map[string]string{}, http.StatusOK)...,
	)
	r.GET("/secure/healthz", handleSecureHealthz, healthOpts...)

	secure := r.Group("", gin.WithTags("Secure Users"))

	secure.GET("/secure/users", handleSecureListUsers, gin.WithSecurity(bearer))

	postOpts := append(
		[]gin.HandlerOption{gin.WithSecurity(apiKey)},
		gin.JSONRoute(nil, struct{}{}, http.StatusCreated)...,
	)
	secure.POST("/secure/users", handleSecureCreateUser, postOpts...)
}
