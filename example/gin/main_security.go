//go:build gin && !typed && security

package main

import (
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	ginlib "github.com/gin-gonic/gin"

	"github.com/aizacoders/openapigo/adapters/gin"
	"github.com/aizacoders/openapigo/openapi"
)

type SecUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	r := gin.New()

	cfg := openapi.Config{
		Title:   "User API (Gin + Security)",
		Version: "1.0.0",
		SecuritySchemes: map[string]*openapi3.SecuritySchemeRef{
			"bearerAuth": {Value: &openapi3.SecurityScheme{Type: "http", Scheme: "bearer", BearerFormat: "JWT"}},
			"apiKeyAuth": {Value: &openapi3.SecurityScheme{Type: "apiKey", In: "header", Name: "X-API-Key"}},
		},
	}

	bearer := openapi3.NewSecurityRequirement().Authenticate("bearerAuth")
	apiKey := openapi3.NewSecurityRequirement().Authenticate("apiKeyAuth")

	r.GET("/secure/users", func(c *ginlib.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.Status(http.StatusUnauthorized)
			return
		}
		gin.JSON(c, http.StatusOK, []SecUser{{ID: "1", Name: "Alice"}})
	}, gin.WithSecurity(&bearer))

	r.POST("/secure/users", func(c *ginlib.Context) {
		if c.GetHeader("X-API-Key") == "" {
			c.Status(http.StatusUnauthorized)
			return
		}
		c.Status(http.StatusCreated)
	}, gin.WithSecurity(&apiKey))

	gin.Register(r, cfg)
	_ = r.Engine.Run(":8080")
}
