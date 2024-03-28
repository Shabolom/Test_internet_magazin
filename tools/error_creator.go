package tools

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func CreateError(code int, err error, c *gin.Context) {
	c.JSON(code, gin.H{
		"code":  code,
		"error": err.Error(),
	})
	_ = c.AbortWithError(code, errors.New(err.Error()))
}
