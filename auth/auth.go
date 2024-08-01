package auth

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/Sph3ricalPeter/go-auth/config"
	"github.com/Sph3ricalPeter/go-auth/storage"
	"github.com/golang-jwt/jwt/v5"
)

type JwtAuth struct {
	storage storage.Storage
}

func NewJwtAuth(storage storage.Storage) *JwtAuth {
	return &JwtAuth{
		storage: storage,
	}
}

func (auth *JwtAuth) HandleLogin(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := auth.storage.VerifyUser(user, pass)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("invalid username or password"))
		return
	}

	accTokenExpTime := time.Now().Add(1 * time.Minute).Unix()
	accToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user,
		"role":     "admin",
		"exp":      accTokenExpTime,
	})
	accTokenStr, err := accToken.SignedString([]byte(config.ACC_SECRET))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	refTokenExpTime := time.Now().Add(3 * time.Minute).Unix()
	refToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user,
		"exp":      refTokenExpTime,
	})
	refTokenStr, err := refToken.SignedString([]byte(config.REF_SECRET))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	auth.storage.RegisterRefreshToken(refTokenStr)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accTokenStr,
		"refresh_token": refTokenStr,
	})
}

func (auth *JwtAuth) HandleLogout(w http.ResponseWriter, r *http.Request) {
	refTokenStr := r.FormValue("refresh_token")
	if refTokenStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	auth.storage.DeleteRefreshToken(refTokenStr)
}

func (auth *JwtAuth) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	refTokenStr := r.FormValue("refresh_token")
	if refTokenStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !auth.storage.IsRefreshTokenValid(refTokenStr) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	refToken, err := auth.validateToken(refTokenStr, config.REF_SECRET)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, ok := refToken.Claims.(jwt.MapClaims)
	if !ok {
		slog.Error("Refresh token claims are invalid")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	username, ok := claims["username"].(string)
	if !ok {
		slog.Error("Refresh token claims are invalid")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	accTokenExpTime := time.Now().Add(1 * time.Minute).Unix()
	accToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"role":     "admin",
		"exp":      accTokenExpTime,
	})
	accTokenStr, err := accToken.SignedString([]byte(config.ACC_SECRET))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accTokenStr,
	})
}

func (auth *JwtAuth) JwtAuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		_, err := auth.validateToken(tokenStr, config.ACC_SECRET)
		if err != nil {
			slog.Error(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (auth *JwtAuth) validateToken(token string, secret string) (*jwt.Token, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !t.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	return t, nil
}
