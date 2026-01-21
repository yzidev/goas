//go:build gin && !typed && !security

package main

import (
	"net/http"

	ginlib "github.com/gin-gonic/gin"

	"github.com/aizacoders/openapigo/adapters/gin"
	"github.com/aizacoders/openapigo/openapi"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	r := gin.New()

	r.GET("/users", func(c *ginlib.Context) {
		gin.JSON(c, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
	})

	r.POST("/users", func(c *ginlib.Context) {
		c.Status(http.StatusCreated)
	})

	gin.Register(r, openapi.Config{Title: "User API", Version: "1.0.0"})
	_ = r.Engine.Run(":8080")
}
