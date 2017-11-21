package handlers

import "fmt"
import "github.com/gin-gonic/gin"

func ExecuteProcedure(c *gin.Context) {
	fmt.Println("Executing a procedure")
	return
}
