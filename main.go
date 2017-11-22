package main

import (
	"github.com/BurtonR/sqlrest/database"
	"github.com/BurtonR/sqlrest/handlers"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	v1 := r.Group("v1")
	{
		v1.POST("/query", handlers.ExecuteQuery)
		v1.POST("/update", handlers.ExecuteUpdate)
		v1.PUT("/insert", handlers.ExecuteInsert)
		v1.DELETE("/delete", handlers.ExecuteDelete)
		v1.POST("/procedure", handlers.ExecuteProcedure)
	}

	return r
}

func main() {
	database.GetConnection()
	r := setupRouter()
	r.Run(":8080")
}
