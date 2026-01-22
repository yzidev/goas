//go:build !security

package main

import (
	"net/http"

	"github.com/aizacoders/openapigo/openapi"
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
	r := openapi.NewRouter()

	users := r.Group("", openapi.WithTags("Users"))

	users.GET("/users", func(w http.ResponseWriter, _ *http.Request) {
		openapi.JSON(w, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
	}, openapi.JSONRoute(nil, []User{}, http.StatusOK)...)

	users.GET("/search", func(w http.ResponseWriter, req *http.Request) {
		_, _, _ = openapi.QueryValue[int](req, "limit")
		w.WriteHeader(http.StatusOK)
	},
		append(
			[]openapi.HandlerOption{openapi.WithQueryParams(
				openapi.QueryParam{Name: "q", Type: openapi.ParamString, Required: true, Description: "Search term"},
				openapi.QueryParam{Name: "limit", Type: openapi.ParamInteger, Required: false, Description: "Max results"},
			)},
			openapi.JSONRoute(nil, struct{}{}, http.StatusOK)...,
		)...,
	)

	users.POST("/users", func(w http.ResponseWriter, req *http.Request) {
		var in CreateUser
		if err := openapi.Bind(req, &in); err != nil || in.Name == "" {
			openapi.JSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return
		}
		w.WriteHeader(http.StatusCreated)
	}, openapi.JSONRoute(CreateUser{}, struct{}{}, http.StatusCreated)...)

	users.GET("/users/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := openapi.PathValue(req, "id")
		if id == "404" {
			openapi.JSON(w, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return
		}
		openapi.JSON(w, http.StatusOK, User{ID: id, Name: "Alice"})
	}, openapi.JSONRoute(nil, User{}, http.StatusOK)...)

	users.PUT("/users/{id}", func(w http.ResponseWriter, req *http.Request) {
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
	}, openapi.JSONRoute(UpdateUser{}, User{}, http.StatusOK)...)

	users.PATCH("/users/{id}", func(w http.ResponseWriter, req *http.Request) {
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
	}, openapi.JSONRoute(UpdateUser{}, User{}, http.StatusOK)...)

	users.DELETE("/users/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := openapi.PathValue(req, "id")
		if id == "404" {
			openapi.JSON(w, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}, openapi.JSONRoute(nil, struct{}{}, http.StatusNoContent)...)

	openapi.Register(r, openapi.Config{
		Title:   "User API",
		Version: "1.0.0",
		Tags: openapi3.Tags{
			{Name: "Users", Description: "User management endpoints"},
		},
	})

	_ = http.ListenAndServe(":8080", r)
}
