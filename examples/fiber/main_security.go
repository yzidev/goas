//go:build fiber && !typed && security

package main

import (
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	fiberlib "github.com/gofiber/fiber/v2"

	"github.com/aizacoders/openapigo/adapters/fiberadapter"
	"github.com/aizacoders/openapigo/openapi"
	"github.com/aizacoders/openapigo/openapi/oas"
)

type SecUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	base := fiberadapter.New()

	cfg := openapi.Config{
		Title:   "User API (Fiber + Security)",
		Version: "1.0.0",
		SecuritySchemes: map[string]*openapi3.SecuritySchemeRef{
			"bearerAuth": {Value: &openapi3.SecurityScheme{Type: "http", Scheme: "bearer", BearerFormat: "JWT"}},
			"apiKeyAuth": {Value: &openapi3.SecurityScheme{Type: "apiKey", In: "header", Name: "X-API-Key"}},
		},
	}

	bearer := openapi3.NewSecurityRequirement().Authenticate("bearerAuth")
	apiKey := openapi3.NewSecurityRequirement().Authenticate("apiKeyAuth")

	b := oas.NewSpec()
	b.GroupTags("", []string{"Secure Users"}, func(s *oas.SpecBuilder) {
		s.GET("/secure/users").Security(&bearer).Res([]SecUser{}).OK()
		s.POST("/secure/users").Security(&apiKey).Res(struct{}{}).Created()

		// Upload secure user file: multipart/form-data.
		s.POST("/secure/users/upload").Security(&apiKey).MultipartUpload("file", openapi.MultipartField{Name: "note", Type: openapi.ParamString}).Res(map[string]string{}).OK()

		// Error showcase: helps Swagger UI show error schemas in security mode.
		s.GET("/secure/demo-errors").Security(&bearer).Res(map[string]string{}).OK()
	})

	spec := b.Spec()

	r := oas.NewFiberRouter(base, spec)
	secure := r.Group("", fiberadapter.WithTags("Secure Users"))

	secure.GET("/secure/users", func(c *fiberlib.Ctx) error {
		auth := c.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return c.SendStatus(http.StatusUnauthorized)
		}
		return fiberadapter.JSON(c, http.StatusOK, []SecUser{{ID: "1", Name: "Alice"}})
	})

	secure.POST("/secure/users", func(c *fiberlib.Ctx) error {
		if c.Get("X-API-Key") == "" {
			return c.SendStatus(http.StatusUnauthorized)
		}
		return c.SendStatus(http.StatusCreated)
	})

	secure.POST("/secure/users/upload", func(c *fiberlib.Ctx) error {
		if c.Get("X-API-Key") == "" {
			return fiberadapter.JSON(c, http.StatusUnauthorized, openapi.ErrorResponse{Error: "unauthorized"})
		}
		fh, err := c.FormFile("file")
		if err != nil {
			return fiberadapter.JSON(c, http.StatusBadRequest, openapi.ErrorResponse{Error: "missing file"})
		}
		note := c.FormValue("note")
		return fiberadapter.JSON(c, http.StatusOK, map[string]string{"filename": fh.Filename, "note": note})
	})

	secure.GET("/secure/demo-errors", func(c *fiberlib.Ctx) error {
		if !strings.HasPrefix(c.Get("Authorization"), "Bearer ") {
			return c.SendStatus(http.StatusUnauthorized)
		}
		switch c.Query("code") {
		case "400":
			return fiberadapter.JSON(c, http.StatusBadRequest, openapi.ErrorResponse{Error: "bad request"})
		case "500":
			return fiberadapter.JSON(c, http.StatusInternalServerError, openapi.ErrorResponse{Error: "internal error"})
		case "503":
			return fiberadapter.JSON(c, http.StatusServiceUnavailable, openapi.ErrorResponse{Error: "service unavailable"})
		default:
			return fiberadapter.JSON(c, http.StatusOK, map[string]string{"status": "ok"})
		}
	})

	fiberadapter.Register(base, cfg)
	_ = base.App.Listen(":8080")
}
