//go:build echo && !typed && !security

package main

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	echolib "github.com/labstack/echo/v4"

	"github.com/aizacoders/openapigo/adapters/echo"
	"github.com/aizacoders/openapigo/openapi"
	"github.com/aizacoders/openapigo/openapi/simple"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type UpdateUser struct {
	Name string `json:"name"`
}

type CreateUser struct {
	Name string `json:"name"`
}

func main() {
	base := echolib.New()

	b := simple.NewSpec()
	b.GroupTags("", []string{"Users"}, func(s *simple.SpecBuilder) {
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

	// wrap the existing echo instance into the adapter
	r := echo.NewFromEcho(base)

	sr := simple.NewEcho(r, spec)
	users := sr.Group("", echo.WithTags("Users"))

	users.GET("/users", func(c echolib.Context) error {
		return echo.JSON(c, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
	})

	users.GET("/search", func(c echolib.Context) error {
		_ = c.QueryParam("q")
		return c.NoContent(http.StatusOK)
	})

	users.POST("/users", func(c echolib.Context) error {
		var in CreateUser
		if err := echo.Bind(c, &in); err != nil || in.Name == "" {
			_ = echo.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return nil
		}
		return c.NoContent(http.StatusCreated)
	})

	users.POST("/users/upload", func(c echolib.Context) error {
		f, err := c.FormFile("file")
		if err != nil {
			_ = echo.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "missing file"})
			return nil
		}
		note := c.FormValue("note")
		return echo.JSON(c, http.StatusOK, map[string]string{"filename": f.Filename, "note": note})
	})

	users.GET("/users/:id", func(c echolib.Context) error {
		id := c.Param("id")
		if id == "404" {
			return echo.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
		}
		return echo.JSON(c, http.StatusOK, User{ID: id, Name: "Alice"})
	})

	users.PUT("/users/:id", func(c echolib.Context) error {
		id := c.Param("id")
		var in UpdateUser
		if err := echo.Bind(c, &in); err != nil {
			_ = echo.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return nil
		}
		if id == "404" {
			_ = echo.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return nil
		}
		return echo.JSON(c, http.StatusOK, User{ID: id, Name: in.Name})
	})

	users.PATCH("/users/:id", func(c echolib.Context) error {
		id := c.Param("id")
		var in UpdateUser
		if err := echo.Bind(c, &in); err != nil {
			_ = echo.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return nil
		}
		if id == "404" {
			_ = echo.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return nil
		}
		return echo.JSON(c, http.StatusOK, User{ID: id, Name: in.Name})
	})

	users.DELETE("/users/:id", func(c echolib.Context) error {
		id := c.Param("id")
		if id == "404" {
			_ = echo.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return nil
		}
		return c.NoContent(http.StatusNoContent)
	})

	echo.Register(r, openapi.Config{
		Title:   "User API",
		Version: "1.0.0",
		Tags:    openapi3.Tags{{Name: "Users", Description: "User management endpoints"}},
	})
	_ = r.Echo.Start(":8080")
}
