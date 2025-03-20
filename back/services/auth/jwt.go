package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	configs "github.com/nobregas/ecommerce-mobile-back/config"
	"github.com/nobregas/ecommerce-mobile-back/types"
	"github.com/nobregas/ecommerce-mobile-back/utils"
)

type contextKey string

const (
	userKey     contextKey = "userId"
	userRoleKey contextKey = "userRole"
)

func CreateJWT(secret []byte, userId int, userRole types.UserRole) (string, error) {
	expiration := time.Second * time.Duration(configs.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   strconv.Itoa(userId),
		"exp":      time.Now().Add(expiration).Unix(),
		"userRole": string(userRole),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func WithAdminAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userRole := GetUserRoleFromContext(r.Context())

		if userRole != types.RoleAdmin {
			utils.WriteError(w, http.StatusForbidden, fmt.Errorf("access denied"))
			return
		}

		handlerFunc(w, r)
	}
}

func WithJwtAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := utils.GetTokenFromRequest(r)
		if tokenString == "" {
			unauthorized(w)
			return
		}

		token, err := validateToken(tokenString)
		if err != nil || !token.Valid {
			log.Printf("invalid token: %v", err)
			unauthorized(w)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			unauthorized(w)
			return
		}

		userID, err := parseUserID(claims)
		if err != nil {
			log.Printf("invalid user ID: %v", err)
			unauthorized(w)
			return
		}

		user, err := store.GetUserByID(userID)
		if err != nil {
			log.Printf("user not found: %v", err)
			unauthorized(w)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, user.ID)
		ctx = context.WithValue(ctx, userRoleKey, user.Role)

		handlerFunc(w, r.WithContext(ctx))
	}
}

func validateToken(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(configs.Envs.JWTSecret), nil
	})
}

func parseUserID(claims jwt.MapClaims) (int, error) {
	userIDStr, ok := claims["userId"].(string)
	if !ok {
		return 0, fmt.Errorf("invalid user ID type")
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID format")
	}

	return userID, nil
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, _ := ctx.Value(userKey).(int)
	return userID
}

func GetUserRoleFromContext(ctx context.Context) types.UserRole {
	role, _ := ctx.Value(userRoleKey).(types.UserRole)
	return role
}

func unauthorized(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid or missing authentication token"))
}
