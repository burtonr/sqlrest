package handlers

import (
	"fmt"
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
func ExecuteQuery(connection *database.SQLDatabase) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		var params QueryRequest
		c.BindJSON(&params)
		cmd := params.Query
		strings.TrimSpace(cmd)

		if cmd == "" {
			c.JSON(400, gin.H{"message": "You must supply a query"})
			return
		}

		rex := regexp.MustCompile("(?i)select")
		selects := rex.FindAllString(cmd, -1)

		if len(selects) < 1 {
			c.JSON(400, gin.H{"message": "Query must contain at least 1 'SELECT' statement for 'Query' operation"})
			return
		}

		data, err := database.Execute(connection.Connection, cmd)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{"message": "Error returned from database", "error": err})
			return
		}

		c.JSON(200, gin.H{"data": data})
	})
}
