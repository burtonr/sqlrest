package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/BurtonR/sqlrest/database"
	"github.com/BurtonR/sqlrest/handlers"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.GET("/connect", func(c *gin.Context) {
		connected, err := database.GetConnection()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to connect"})
			return
		}
		if connected {
			c.JSON(http.StatusOK, gin.H{"message": "connected"})
		}
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
	connectToDb()

	ticker := time.NewTicker(2 * time.Minute)
	quit := make(chan struct{})
	go pinger(ticker, quit)

	r := setupRouter()
	r.Run(":5050")
}

func pinger(ticker *time.Ticker, quit chan struct{}) {
	for {
		select {
		case <-ticker.C:
			connectToDb()
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func connectToDb() {
	maxRetries := 2
	var conn *sql.DB
	for i := 0; i < maxRetries; i++ {
		connected, err := database.GetConnection()
		if err != nil {
			fmt.Printf("Unable to connect. Attempt %d of %d", i+1, maxRetries)
			fmt.Println()
		}
		if connected {
			return
		}

		time.Sleep(500 * time.Millisecond)
	}

	if conn == nil {
		fmt.Printf("No database connection after %d attempts", maxRetries)
		fmt.Println()
	}

	return
}
