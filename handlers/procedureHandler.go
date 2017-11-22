package handlers

import (
	"github.com/gin-gonic/gin"
)

func ExecuteProcedure(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Handling the PROCEDURE"})
}
