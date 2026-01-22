//go:build gin && !typed && !security

package main

import (
	"net/http"

	ginlib "github.com/gin-gonic/gin"
)

func handleListUsers(c *ginlib.Context) {
	c.JSON(http.StatusOK, []User{{ID: "1", Name: "Alice"}})
}

func handleSearchUsers(c *ginlib.Context) {
	_ = c.Query("q")
	c.Status(http.StatusOK)
}

func handleCreateUser(c *ginlib.Context) {
	var in CreateUser
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
		return
	}
	c.Status(http.StatusCreated)
}

func handleGetUser(c *ginlib.Context) {
	id := c.Param("id")
	if id == "404" {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "user not found"})
		return
	}
	c.JSON(http.StatusOK, User{ID: id, Name: "Alice"})
}

func handlePutUser(c *ginlib.Context) {
	id := c.Param("id")
	var in UpdateUser
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
		return
	}
	if id == "404" {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "user not found"})
		return
	}
	c.JSON(http.StatusOK, User{ID: id, Name: in.Name})
}

func handlePatchUser(c *ginlib.Context) {
	handlePutUser(c)
}

func handleDeleteUser(c *ginlib.Context) {
	id := c.Param("id")
	if id == "404" {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "user not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
