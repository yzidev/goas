//go:build gin && security && !typed

package main

import (
	"net/http"
	"strings"

	ginlib "github.com/gin-gonic/gin"
)

func requireBearer(c *ginlib.Context) bool {
	auth := c.GetHeader("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		c.Status(http.StatusUnauthorized)
		return false
	}
	return true
}

func requireAPIKey(c *ginlib.Context) bool {
	if c.GetHeader("X-API-Key") == "" {
		c.Status(http.StatusUnauthorized)
		return false
	}
	return true
}

func handleSecureHealthz(c *ginlib.Context) {
	if !requireBearer(c) {
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func handleSecureListUsers(c *ginlib.Context) {
	if !requireBearer(c) {
		return
	}
	c.JSON(http.StatusOK, []SecUser{{ID: "1", Name: "Alice"}})
}

func handleSecureCreateUser(c *ginlib.Context) {
	if !requireAPIKey(c) {
		return
	}
	c.Status(http.StatusCreated)
}
