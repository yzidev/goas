//go:build echo && !typed && security

package main

import (
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	echolib "github.com/labstack/echo/v4"

	"github.com/aizacoders/openapigo/adapters/echo"
	"github.com/aizacoders/openapigo/openapi"
)

type SecUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	r := echo.New()

	cfg := openapi.Config{
		Title:   "User API (Echo + Security)",
		Version: "1.0.0",
		SecuritySchemes: map[string]*openapi3.SecuritySchemeRef{
			"bearerAuth": {Value: &openapi3.SecurityScheme{Type: "http", Scheme: "bearer", BearerFormat: "JWT"}},
			"apiKeyAuth": {Value: &openapi3.SecurityScheme{Type: "apiKey", In: "header", Name: "X-API-Key"}},
		},
	}

	bearer := openapi3.NewSecurityRequirement().Authenticate("bearerAuth")
	apiKey := openapi3.NewSecurityRequirement().Authenticate("apiKeyAuth")

	// Basic handler + Bearer security documented in OpenAPI
	r.GET("/secure/users", func(c echolib.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return c.NoContent(http.StatusUnauthorized)
		}
		return echo.JSON(c, http.StatusOK, []SecUser{{ID: "1", Name: "Alice"}})
	}, func(meta *openapi.RouteMeta) { echo.WithSecurity(&bearer)(meta) })

	// Basic handler + API key security documented in OpenAPI
	r.POST("/secure/users", func(c echolib.Context) error {
		if c.Request().Header.Get("X-API-Key") == "" {
			return c.NoContent(http.StatusUnauthorized)
		}
		return c.NoContent(http.StatusCreated)
	}, func(meta *openapi.RouteMeta) { echo.WithSecurity(&apiKey)(meta) })

	echo.Register(r, cfg)
	_ = r.Echo.Start(":8080")
}
