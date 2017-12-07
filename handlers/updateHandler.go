package handlers

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/BurtonR/sqlrest/database"
	"github.com/gin-gonic/gin"
)

// UpdateRequest A struct to hold the body of the query request with a property of "update"
type UpdateRequest struct {
	Update string
}

// ExecuteUpdate HttpHandler to validate and execute an 'UPDATE' command to the database
func ExecuteUpdate(c *gin.Context) {
	var params UpdateRequest
	c.BindJSON(&params)
	cmd := params.Update
	strings.TrimSpace(cmd)

	if cmd == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "You must supply an UPDATE statement"})
		return
	}

	urex := regexp.MustCompile("(?i)update")
	updates := urex.FindAllString(cmd, -1)

	if len(updates) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Update must contain at least 1 'UPDATE' statement for 'Update' operation"})
		return
	}

	wrex := regexp.MustCompile("(?i)where")
	wheres := wrex.FindAllString(cmd, -1)

	if len(wheres) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Update must contain at least 1 'WHERE' clause for your protection"})
		return
	}

	err := database.ExecuteUpdate(cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error returned from database", "error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
