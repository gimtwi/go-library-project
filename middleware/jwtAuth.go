package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
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

func CheckPrivilege(ur types.UserRepository, role types.UserRole) gin.HandlerFunc {
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

			id, ok := claims["id"].(string)
			if !ok || id == "" {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			user, errUser := ur.GetByID(id)
			if errUser != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			if user == nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
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

func CompareCookiesAndParameter(ur types.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")

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

			id, ok := claims["id"].(string)
			if !ok || id == "" {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			user, errUser := ur.GetByID(id)
			if errUser != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			if user.ID != userID {
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

func GetUserIDFromTheToken(c *gin.Context) string {
	tokenStr, err := c.Cookie(os.Getenv("COOKIE_NAME"))
	if err != nil {
		return ""
	}

	token, err := ValidateJWT(tokenStr)
	if err != nil {
		return ""
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return ""
		}

		id, ok := claims["id"].(string)
		if !ok || id == "" {
			return ""
		}

		return id

	} else {
		return ""
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
