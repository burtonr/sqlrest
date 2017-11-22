package handlers

import (
	"github.com/gin-gonic/gin"
)

func ExecuteDelete(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Handling the DELETE"})
}
