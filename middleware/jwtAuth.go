package middleware

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gimtwi/go-library-project/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/juju/ratelimit"
)

func GenerateJWT(user *types.User) (string, error) {
	claims := &jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // token expires in 30 days
	}

	secret := os.Getenv("JWT_SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

}

func CheckPrivilege(repo types.UserRepository, role types.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(os.Getenv("COOKIE_NAME"))

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		token, err := ValidateJWT(tokenStr)

		if err != nil {
			// Handle the JWT parsing error gracefully
			if err == jwt.ErrSignatureInvalid {
				// The token signature is invalid
				c.AbortWithStatus(http.StatusUnauthorized)
			} else if err == jwt.ErrTokenMalformed {
				// The token is malformed
				c.AbortWithStatus(http.StatusUnauthorized)
			} else {
				// Other JWT parsing errors
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			log.Println(err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			idClaim, ok := claims["id"].(float64)
			if !ok || idClaim < 0 {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			id := uint(math.Floor(idClaim))
			user, errUser := repo.GetByID(uint(id))
			if errUser != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			requiredRole := role

			if !types.CheckPrivilege(user.Role, requiredRole) {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			c.Set("user", user)

			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Next()
	}
}

func CheckIfMe(repo types.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		userID, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		tokenStr, err := c.Cookie(os.Getenv("COOKIE_NAME"))

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		token, err := ValidateJWT(tokenStr)

		if err != nil {
			// Handle the JWT parsing error gracefully
			if err == jwt.ErrSignatureInvalid {
				// The token signature is invalid
				c.AbortWithStatus(http.StatusUnauthorized)
			} else if err == jwt.ErrTokenMalformed {
				// The token is malformed
				c.AbortWithStatus(http.StatusUnauthorized)
			} else {
				// Other JWT parsing errors
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			log.Println(err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			idClaim, ok := claims["id"].(float64)
			if !ok || idClaim < 0 {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			id := uint(math.Floor(idClaim))
			user, errUser := repo.GetByID(uint(id))
			if errUser != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			if user.ID != uint(userID) {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			c.Set("user", user)

			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Next()
	}
}

func RateLimitMiddleware() gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithRate(1, 5)

	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}
