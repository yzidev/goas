//go:build !security

package main

import (
	"net/http"

	"github.com/aizacoders/openapigo/adapters/httprouter"
	"github.com/aizacoders/openapigo/openapi"
	"github.com/aizacoders/openapigo/openapi/simple"
	"github.com/getkin/kin-openapi/openapi3"
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
	base := httprouter.New()

	b := simple.NewSpec()
	b.GroupTags("", []string{"Users"}, func(s *simple.SpecBuilder) {
		s.GET("/users").Res([]User{}).OK()
		s.GET("/search").Query(
			openapi.QueryParam{Name: "q", Type: openapi.ParamString, Required: true, Description: "Search term"},
			openapi.QueryParam{Name: "limit", Type: openapi.ParamInteger, Required: false, Description: "Max results"},
		).Res(struct{}{}).OK()
		s.POST("/users").Req(CreateUser{}).Res(struct{}{}).Status(http.StatusCreated)
		// Upload user file: multipart/form-data.
		s.POST("/users/upload").MultipartUpload("file", openapi.MultipartField{Name: "note", Type: openapi.ParamString}).Res(map[string]string{}).OK()
		s.GET("/users/{id}").Res(User{}).OK()
		s.PUT("/users/{id}").Req(UpdateUser{}).Res(User{}).OK()
		s.PATCH("/users/{id}").Req(UpdateUser{}).Res(User{}).OK()
		s.DELETE("/users/{id}").Res(struct{}{}).NoContent()
	})

	spec := b.Spec()

	r := simple.New(base, spec)

	// Clean routes: just HTTP methods + handlers.
	r.GET("/users", func(w http.ResponseWriter, _ *http.Request) {
		openapi.JSON(w, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
	})

	r.GET("/search", func(w http.ResponseWriter, req *http.Request) {
		_, _, _ = openapi.QueryValue[int](req, "limit")
		w.WriteHeader(http.StatusOK)
	})

	r.POST("/users", func(w http.ResponseWriter, req *http.Request) {
		var in CreateUser
		if err := openapi.Bind(req, &in); err != nil || in.Name == "" {
			openapi.JSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return
		}
		w.WriteHeader(http.StatusCreated)
	})

	r.POST("/users/upload", func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseMultipartForm(10 << 20); err != nil {
			openapi.JSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid multipart"})
			return
		}
		f, fh, err := req.FormFile("file")
		if err != nil {
			openapi.JSON(w, http.StatusBadRequest, ErrorResponse{Error: "missing file"})
			return
		}
		_ = f.Close()
		note := req.FormValue("note")
		openapi.JSON(w, http.StatusOK, map[string]string{"filename": fh.Filename, "note": note})
	})

	r.GET("/users/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := openapi.PathValue(req, "id")
		if id == "404" {
			openapi.JSON(w, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return
		}
		openapi.JSON(w, http.StatusOK, User{ID: id, Name: "Alice"})
	})

	r.PUT("/users/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := openapi.PathValue(req, "id")
		var in UpdateUser
		if err := openapi.Bind(req, &in); err != nil {
			openapi.JSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return
		}
		if id == "404" {
			openapi.JSON(w, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return
		}
		openapi.JSON(w, http.StatusOK, User{ID: id, Name: in.Name})
	})

	r.PATCH("/users/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := openapi.PathValue(req, "id")
		var in UpdateUser
		if err := openapi.Bind(req, &in); err != nil {
			openapi.JSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return
		}
		if id == "404" {
			openapi.JSON(w, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return
		}
		openapi.JSON(w, http.StatusOK, User{ID: id, Name: in.Name})
	})

	r.DELETE("/users/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := openapi.PathValue(req, "id")
		if id == "404" {
			openapi.JSON(w, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	httprouter.Register(base, openapi.Config{
		Title:   "User API",
		Version: "1.0.0",
		Tags: openapi3.Tags{
			{Name: "Users", Description: "User management endpoints"},
		},
	})

	_ = http.ListenAndServe(":8080", r)
}
