package handlers

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/BurtonR/sqlrest/database"
	"github.com/gin-gonic/gin"
)

// QueryRequest A struct to hold the body of the query request with a property of "query"
type QueryRequest struct {
	Query string
}

// ExecuteQuery HttpHandler to validate and execute a 'SELECT' command to the database
func ExecuteQuery(c *gin.Context) {
	var params QueryRequest
	c.BindJSON(&params)
	cmd := params.Query
	strings.TrimSpace(cmd)

	if cmd == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "You must supply a query"})
		return
	}

	rex := regexp.MustCompile("(?i)select")
	selects := rex.FindAllString(cmd, -1)

	if len(selects) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Query must contain at least 1 'SELECT' statement for 'Query' operation"})
		return
	}

	data, err := database.ExecuteQuery(cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error returned from database", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}
