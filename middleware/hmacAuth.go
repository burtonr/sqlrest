package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// HmacAuthentication checks the Authorization header for proper HMAC values
func HmacAuthentication(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	fmt.Println(authHeader)

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
	timestamp := secrets[3]

	if !verifyRealm(realm) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	fmt.Println("Need to verify Signature:", signature)
	fmt.Println("Need to verify Nonce:", nonce)
	fmt.Println("Need to verify TimeStamp:", timestamp)

	return
}

func verifyRealm(realm string) bool {
	// allowedRealmsVar := os.Getenv("SQLREST_ALLOWED_REALMS")
	allowedRealmsVar := "kgb-functions"

	allowedRealms := strings.Split(allowedRealmsVar, ",")

	for _, allowed := range allowedRealms {
		if realm == allowed {
			return true
		}
	}
	return false
}
