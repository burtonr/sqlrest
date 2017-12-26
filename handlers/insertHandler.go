package handlers

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/BurtonR/sqlrest/database"
	"github.com/gin-gonic/gin"
)

// InsertRequest A struct to hold the body of the query request with a property of "insert"
type InsertRequest struct {
	Insert string
}

// ExecuteInsert HttpHandler to validate and execute an 'INSERT' command to the database
func ExecuteInsert(c *gin.Context) {
	var params InsertRequest
	c.BindJSON(&params)
	cmd := params.Insert
	strings.TrimSpace(cmd)

	if cmd == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "You must supply an INSERT statement"})
		return
	}

	urex := regexp.MustCompile("(?i)insert into")
	inserts := urex.FindAllString(cmd, -1)

	if len(inserts) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Insert must contain at least 1 'INSERT INTO' statement for 'Insert' operation"})
		return
	}

	err := database.ExecuteWithTransaction(cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error returned from database", "error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
