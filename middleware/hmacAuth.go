package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// For inspiration: https://github.com/tjoudeh/WebApiHMACAuthentication/blob/master/HMACAuthentication.WebApi/Filters/HMACAuthenticationAttribute.cs
// article: http://bitoftech.net/2014/12/15/secure-asp-net-web-api-using-api-key-authentication-hmac-authentication/

// HmacAuthentication checks the Authorization header for proper HMAC values
func HmacAuthentication(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "null" || len(authHeader) < 1 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	secrets := strings.Split(authHeader, ":")

	if len(secrets) < 4 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	realm := secrets[0]
	signature := secrets[1]
	nonce := secrets[2]
	timestring := secrets[3]

	if !verifyRealm(realm) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if !verifySignature(signature, c) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// TODO implement check for nonce
	fmt.Println("Need to verify Nonce:", nonce)

	timeInt, err := strconv.ParseInt(timestring, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if !verifyTimestamp(timeInt) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	return
}

func verifyRealm(realm string) bool {
	allowedRealmsVar := os.Getenv("SQLREST_ALLOWED_REALMS")
	// allowedRealmsVar := "burton-func, testing-func"

	allowedRealms := strings.Split(allowedRealmsVar, ",")

	for _, allowed := range allowedRealms {
		if realm == strings.TrimSpace(allowed) {
			return true
		}
	}
	return false
}

func verifySignature(signature string, c *gin.Context) bool {
	// get the request body and re-set it to use later
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	bodyString := string(bodyBytes)

	if c.Request.Method == "GET" && c.Request.ContentLength < 1 {
		return true
	}

	apiKey := os.Getenv("SQLREST_API_KEY")
	// apiKey := "sqlrestTestKey"

	mac := hmac.New(sha256.New, []byte(apiKey))
	mac.Write([]byte(bodyString))
	expected := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(signature))
}

func verifyTimestamp(timestamp int64) bool {
	currentTime := time.Now().UnixNano() / 1000000 // in milliseconds
	difference := currentTime - timestamp
	return difference < 2000 // deny requests more than 2 seconds ago
}
