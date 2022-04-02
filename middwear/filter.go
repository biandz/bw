package middwear

import (
	"github.com/gin-gonic/gin"
)

func Filter() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Next()
	}
}
