package main

import (
	"log"

	"github.com/BurtonR/sqlrest/database"
	"github.com/BurtonR/sqlrest/handlers"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

func setupRouter(db *database.SQLDatabase) *gin.Engine {
	r := gin.Default()

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	v1 := r.Group("v1")
	{
		v1.POST("/query", handlers.ExecuteQuery(db))
		v1.POST("/update", handlers.ExecuteUpdate)
		v1.PUT("/insert", handlers.ExecuteInsert)
		v1.DELETE("/delete", handlers.ExecuteDelete)
		v1.POST("/procedure", handlers.ExecuteProcedure)
	}

	return r
}

func main() {
	conn, err := database.GetConnection()

	if err != nil {
		log.Panic(err)
	}

	db := &database.SQLDatabase{Connection: conn}

	r := setupRouter(db)
	r.Run(":8080")
}
