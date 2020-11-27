package Golang

import "testing"
import "github.com/gin-gonic/gin"

func TestGin(t testing.T) {
	engine := gin.Default()
	engine.Group("v1")
	engine.Use()
}
