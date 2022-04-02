package middwear

import (
	"github.com/gin-gonic/gin"
)

func Intercept() func(c *gin.Context) {
	return func(c *gin.Context) {
		//注意:c.Next()和c.Abort()
		c.Next()
	}
}
