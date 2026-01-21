//go:build security

package main

import (
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/aizacoders/openapigo/adapters/httprouter"
	"github.com/aizacoders/openapigo/openapi"
)

// This example shows how to:
// - declare multiple security schemes (bearer JWT + x-api-key)
// - apply security per-route
// - do a tiny header validation (demo only)

type SecUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SecCreateUser struct {
	Name string `json:"name"`
}

func main() {
	r := httprouter.New()

	cfg := openapi.Config{
		Title:   "User API (Security)",
		Version: "1.0.0",
		SecuritySchemes: map[string]*openapi3.SecuritySchemeRef{
			"bearerAuth": {Value: &openapi3.SecurityScheme{Type: "http", Scheme: "bearer", BearerFormat: "JWT"}},
			"apiKeyAuth": {Value: &openapi3.SecurityScheme{Type: "apiKey", In: "header", Name: "X-API-Key"}},
		},
	}

	bearer := openapi3.NewSecurityRequirement().Authenticate("bearerAuth")
	apiKey := openapi3.NewSecurityRequirement().Authenticate("apiKeyAuth")

	// Bearer-protected endpoint
	httprouter.POSTT[SecCreateUser, SecUser](r, "/secure/users", func(w http.ResponseWriter, req *http.Request, in SecCreateUser) (SecUser, int, error) {
		auth := req.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return SecUser{}, http.StatusUnauthorized, nil
		}
		return SecUser{ID: "1", Name: in.Name}, http.StatusCreated, nil
	}, httprouter.WithSecurity(&bearer))

	// API-key-protected endpoint
	httprouter.GETT[struct{}, []SecUser](r, "/secure/users", func(w http.ResponseWriter, req *http.Request, _ struct{}) ([]SecUser, int, error) {
		if req.Header.Get("X-API-Key") == "" {
			return nil, http.StatusUnauthorized, nil
		}
		return []SecUser{{ID: "1", Name: "Alice"}}, http.StatusOK, nil
	}, httprouter.WithSecurity(&apiKey))

	openapi.Register(r, cfg)
	_ = http.ListenAndServe(":8080", r)
}
