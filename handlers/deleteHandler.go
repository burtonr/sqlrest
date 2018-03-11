package handlers

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/BurtonR/sqlrest/database"
	"github.com/gin-gonic/gin"
)

// DeleteRequest A struct to hold the body of the delete request with a property of "delete"
type DeleteRequest struct {
	Delete string
}

// ExecuteDelete HttpHandler to validate and execute a 'DELETE' command to the database
func ExecuteDelete(c *gin.Context) {
	var params DeleteRequest
	c.BindJSON(&params)
	cmd := params.Delete
	strings.TrimSpace(cmd)

	if cmd == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "You must supply a DELETE statement"})
		return
	}

	drex := regexp.MustCompile("(?i)delete into")
	deletes := drex.FindAllString(cmd, -1)

	if len(deletes) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Delete must contain at least 1 'DELETE' statement for 'Delete' operation"})
		return
	}

	wrex := regexp.MustCompile("(?i)where")
	wheres := wrex.FindAllString(cmd, -1)

	if len(wheres) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Delete must contain at least 1 'WHERE' clause for your protection"})
		return
	}

	err := database.ExecuteWithTransaction(cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error returned from database", "error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
