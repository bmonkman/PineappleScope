package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Success renders a template
func Success(c *gin.Context, name string, vars gin.H) {
	vars["test"] = "testing"
	for k, v := range c.Keys {
		vars[k] = v
	}
	c.HTML(http.StatusOK, name, vars)

}
