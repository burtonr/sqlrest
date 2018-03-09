package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/BurtonR/sqlrest/database"
	"github.com/gin-gonic/gin"
)

// ProcedureRequest A struct to hold the body of the request with properties "name", "parameters", and "result"
type ProcedureRequest struct {
	Name        string
	Parameters  map[string]*interface{}
	ExecuteOnly bool
}

// ExecuteProcedure HttpHandler to create and execute an 'EXEC' command to the database
func ExecuteProcedure(c *gin.Context) {
	// TODO: There's go to be a better way than parsing the request body bytes...
	// Read the Body content to see if the executeOnly flag was included
	var bodyBytes []byte
	parseExec := false

	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		s := string(bodyBytes[:])
		if strings.Contains(strings.ToLower(s), "executeonly") {
			parseExec = true
		}
	}

	// Restore the io.ReadCloser to its original state
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	var request ProcedureRequest
	parseErr := c.BindJSON(&request)

	if parseErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing the request", "error": parseErr.Error()})
		return
	}

	name := request.Name
	execOnly := true

	if parseExec {
		execOnly = request.ExecuteOnly
	}

	strings.TrimSpace(name)

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "You must supply a procedure name"})
		return
	}

	results, err := database.ExecuteStatement(name, execOnly, request.Parameters)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error returned from database", "error": err.Error()})
		return
	}

	if !execOnly {
		c.JSON(http.StatusOK, results)
		return
	}

	c.Status(http.StatusOK)
}
