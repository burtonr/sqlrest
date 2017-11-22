package handlers

import (
	"database/sql"
	"fmt"

	"github.com/BurtonR/sqlrest/database"
	"github.com/gin-gonic/gin"
)

func ExecuteQuery(c *gin.Context) {
	var connection *sql.DB
	var err error

	connection = database.Connection

	if connection == nil {
		connection, err = database.GetConnection()

		if err != nil {
			fmt.Println(err)
		}
	}

	cmd := "SELECT TOP 5 * FROM [DATABASE].[dbo].[TABLENAME]"

	err = database.Execute(connection, cmd)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(200, gin.H{"message": "Handling the query"})
}
