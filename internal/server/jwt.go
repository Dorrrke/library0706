package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Dorrrke/library0706/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	accessTokenDuration  = 15 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *LibraryApi) JWTAuthMiddleware() gin.HandlerFunc {
	log := logger.Get()
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 8 {
			log.Debug().Msg("no authorization header")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}

		tokenStr := authHeader[len("Bearer "):]
		claims, err := parseToken(tokenStr)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse token")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		ctx.Set("user_id", claims.UserID)
		ctx.Next()
	}
}

func createToken(userID string, duration time.Duration) (string, string, error) {
	tokenID := uuid.New().String()
	claims := Claims{
		UserID: userID,

		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", "", err
	}
	return secretToken, tokenID, nil
}

func parseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil

}
