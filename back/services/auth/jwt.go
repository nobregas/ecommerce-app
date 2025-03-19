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
	UserKey     contextKey = "userId"
	UserRoleKey contextKey = "userRole"
)

func CreateJWT(secret []byte, userId int, userRole string) (string, error) {
	expiration := time.Second * time.Duration(configs.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    strconv.Itoa(userId),
		"expiredAt": time.Now().Add(expiration).Unix(),
		"userRole":  userRole,
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func WithRoleAuth(requiredRole string, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the user role from context
		userRole := getUserRoleFromContext(r.Context())

		// check if the user has the required role
		if userRole != requiredRole {
			permissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}

func WithJwtAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the JWT token from the request header
		tokenString := utils.GetTokenFromRequest(r)

		// validate the token
		token, err := validateToken(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			unauthorized(w)
			return
		}

		if !token.Valid {
			log.Printf("invalid token")
			unauthorized(w)
			return
		}

		// if the token is valid, fetch the userID from the DB (id from token)
		claims := token.Claims.(jwt.MapClaims)
		str := claims["userId"].(string)

		userID, _ := strconv.Atoi(str)

		u, err := store.GetUserByID(userID)
		if err != nil {
			log.Printf("failed to get user: %v", err)
			unauthorized(w)
			return
		}

		// set context "userID" to the user ID

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
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

func unauthorized(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIdFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}

	return userID
}

func getUserRoleFromContext(ctx context.Context) string {
	role, ok := ctx.Value(UserRoleKey).(string)
	if !ok {
		return ""
	}

	return role
}
