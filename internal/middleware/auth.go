package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const authContextKey = "auth_context"

type AuthContext struct {
	UID      string
	Username string
	Role     string
	LastSeen int64
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	LastSeen int64  `json:"last_seen"`
	jwt.RegisteredClaims
}

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing_authorization"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid_authorization"})
			return
		}

		tokenStr := strings.TrimSpace(parts[1])
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
			return
		}

		ctx := AuthContext{
			UID:      claims.Subject,
			Username: claims.Username,
			Role:     claims.Role,
			LastSeen: claims.LastSeen,
		}
		c.Set(authContextKey, ctx)
		c.Next()
	}
}

func GetAuthContext(c *gin.Context) (AuthContext, bool) {
	val, ok := c.Get(authContextKey)
	if !ok {
		return AuthContext{}, false
	}
	ctx, ok := val.(AuthContext)
	return ctx, ok
}
