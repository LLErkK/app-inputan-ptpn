package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	User    string `json:"user,omitempty"`
}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims struct (optional typed claims)
type MyClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}
	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Invalid JSON format",
		})
		return
	}

	// Validate input
	if loginReq.Username == "" || loginReq.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Username dan password harus diisi",
		})
		return
	}

	// Find user in database
	var user models.User
	result := config.DB.Where("username = ?", loginReq.Username).First(&user)

	if result.Error != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Username atau password salah",
		})
		return
	}

	// Verify password (bcrypt)
	if !config.ComparePassword(user.Password, loginReq.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Username atau password salah",
		})
		return
	}

	// Generate JWT token
	expireTime := time.Now().Add(24 * time.Hour)
	claims := MyClaims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := config.JWTSecret
	tokenString, err := token.SignedString(secret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Gagal membuat token",
		})
		return
	}

	// Update last login time (non-blocking)
	go func() {
		config.DB.Model(&user).Update("last_login", time.Now())
	}()

	// Set cookie (HttpOnly)
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  expireTime,
		HttpOnly: true,
		Path:     "/",
	}
	// Set Secure true only if HTTPS (use env to force)
	if os.Getenv("APP_ENV") == "production" {
		cookie.Secure = true
		cookie.SameSite = http.SameSiteStrictMode
	}
	http.SetCookie(w, cookie)

	// Respond
	json.NewEncoder(w).Encode(LoginResponse{
		Success: true,
		Message: "Login berhasil",
		Token:   tokenString,
		User:    user.Username,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})
	// FIXED: Redirect ke halaman login
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// helper untuk mengekstrak token dari header/cookie
func extractTokenFromRequest(r *http.Request) string {
	// 1) Authorization header: Bearer <token>
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			return parts[1]
		}
	}
	// 2) Cookie
	if cookie, err := r.Cookie("auth_token"); err == nil && cookie.Value != "" {
		return cookie.Value
	}
	return ""
}

// AuthMiddleware memvalidasi JWT; jika valid -> panggil next, jika tidak -> redirect ke login
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractTokenFromRequest(r)
		if tokenString == "" {
			// FIXED: Redirect ke /login bukan /
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// parse & validate
		claims := &MyClaims{}
		parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Ensure expected method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return config.JWTSecret, nil
		})
		if err != nil || !parsedToken.Valid {
			// FIXED: token invalid or expired, redirect ke /login
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Attach username to context for handlers that need it
		ctx := context.WithValue(r.Context(), "username", claims.Username)
		next(w, r.WithContext(ctx))
	}
}

// ServeLoginPage - menampilkan halaman login
func ServeLoginPage(w http.ResponseWriter, r *http.Request) {
	// Check if user is already logged in
	tokenString := extractTokenFromRequest(r)
	if tokenString != "" {
		// try to parse
		claims := &MyClaims{}
		if t, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return config.JWTSecret, nil
		}); err == nil && t.Valid {
			// FIXED: Redirect ke /dashboard jika sudah login
			http.Redirect(w, r, "/dashboard", http.StatusFound)
			return
		}
	}

	// Serve login HTML file
	http.ServeFile(w, r, "templates/html/login.html")
}
