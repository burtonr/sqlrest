package handlers

import (
	"fmt"

	"github.com/BurtonR/sqlrest/database"
	"github.com/gin-gonic/gin"
)

func ExecuteQuery(connection *database.SqlDatabase) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		cmd := "SELECT TOP 5 * FROM [DATABASE].[dbo].[TABLENAME]"

		err := database.Execute(connection.Connection, cmd)
		if err != nil {
			fmt.Println(err)
		}

		c.JSON(200, gin.H{"message": "Handling the query"})
	})
}
