//go:build fiber && !typed && !security

package main

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	fiberlib "github.com/gofiber/fiber/v2"

	"github.com/aizacoders/openapigo/adapters/fiberadapter"
	"github.com/aizacoders/openapigo/openapi"
	"github.com/aizacoders/openapigo/openapi/oas"
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
	base := fiberlib.New()

	b := oas.NewSpec()
	b.GroupTags("", []string{"Users"}, func(s *oas.SpecBuilder) {
		s.GET("/users").Res([]User{}).OK()
		s.GET("/search").Query(
			openapi.QueryParam{Name: "q", Type: openapi.ParamString, Required: true, Description: "Search term"},
			openapi.QueryParam{Name: "limit", Type: openapi.ParamInteger, Required: false, Description: "Max results"},
		).Res(struct{}{}).OK()
		s.POST("/users").Req(CreateUser{}).Res(struct{}{}).Created()

		// Upload user file: multipart/form-data.
		s.POST("/users/upload").MultipartUpload("file", openapi.MultipartField{Name: "note", Type: openapi.ParamString}).Res(map[string]string{}).OK()

		s.GET("/users/:id").Res(User{}).OK()
		s.PUT("/users/:id").Req(UpdateUser{}).Res(User{}).OK()
		s.PATCH("/users/:id").Req(UpdateUser{}).Res(User{}).OK()
		s.DELETE("/users/:id").Res(struct{}{}).NoContent()
	})

	spec := b.Spec()

	// wrap existing fiber App into adapter
	r := fiberadapter.NewFiberAdapters(base)
	sr := oas.NewFiberRouter(r, spec)

	users := sr.Group("", fiberadapter.WithTags("Users"))
	users.GET("/users", func(c *fiberlib.Ctx) error {
		return fiberadapter.JSON(c, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
	})

	users.GET("/search", func(c *fiberlib.Ctx) error {
		_ = c.Query("q")
		return c.SendStatus(http.StatusOK)
	})

	users.POST("/users", func(c *fiberlib.Ctx) error {
		var in CreateUser
		if err := fiberadapter.Bind(c, &in); err != nil || in.Name == "" {
			_ = fiberadapter.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return nil
		}
		return c.SendStatus(http.StatusCreated)
	})

	users.POST("/users/upload", func(c *fiberlib.Ctx) error {
		fh, err := c.FormFile("file")
		if err != nil {
			return fiberadapter.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "missing file"})
		}
		note := c.FormValue("note")
		return fiberadapter.JSON(c, http.StatusOK, map[string]string{"filename": fh.Filename, "note": note})
	})

	users.GET("/users/:id", func(c *fiberlib.Ctx) error {
		id := c.Params("id")
		if id == "404" {
			return fiberadapter.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
		}
		return fiberadapter.JSON(c, http.StatusOK, User{ID: id, Name: "Alice"})
	})

	users.PUT("/users/:id", func(c *fiberlib.Ctx) error {
		id := c.Params("id")
		var in UpdateUser
		if err := fiberadapter.Bind(c, &in); err != nil {
			_ = fiberadapter.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return nil
		}
		if id == "404" {
			_ = fiberadapter.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return nil
		}
		return fiberadapter.JSON(c, http.StatusOK, User{ID: id, Name: in.Name})
	})

	users.PATCH("/users/:id", func(c *fiberlib.Ctx) error {
		id := c.Params("id")
		var in UpdateUser
		if err := fiberadapter.Bind(c, &in); err != nil {
			_ = fiberadapter.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return nil
		}
		if id == "404" {
			_ = fiberadapter.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return nil
		}
		return fiberadapter.JSON(c, http.StatusOK, User{ID: id, Name: in.Name})
	})

	users.DELETE("/users/:id", func(c *fiberlib.Ctx) error {
		id := c.Params("id")
		if id == "404" {
			_ = fiberadapter.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return nil
		}
		return c.SendStatus(http.StatusNoContent)
	})

	fiberadapter.Register(r, openapi.Config{
		Title:   "User API",
		Version: "1.0.0",
		Tags:    openapi3.Tags{{Name: "Users", Description: "User management endpoints"}},
	})
	_ = r.App.Listen(":8080")
}
