package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ChangeUsernameRequest struct {
	OldUsername string `json:"oldUsername"`
	NewUsername string `json:"newUsername"`
	Password    string `json:"password"`
}

type ChangePasswordRequest struct {
	Username    string `json:"username"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type AccountResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

func ServeAccountManagementPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/html/manajemenAkun.html")
}

func ChangeUsername(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(AccountResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var req ChangeUsernameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AccountResponse{
			Success: false,
			Message: "Invalid JSON format",
		})
		return
	}

	// Validate input
	if req.OldUsername == "" || req.NewUsername == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AccountResponse{
			Success: false,
			Message: "Semua field harus diisi",
		})
		return
	}

	// Find user with old username
	var user models.User
	result := config.DB.Where("username = ?", req.OldUsername).First(&user)
	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(AccountResponse{
			Success: false,
			Message: "Username lama tidak ditemukan",
		})
		return
	}

	// Verify password
	if !config.ComparePassword(user.Password, req.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(AccountResponse{
			Success: false,
			Message: "Password salah",
		})
		return
	}

	// Check if new username already exists
	var existingUser models.User
	if err := config.DB.Where("username = ?", req.NewUsername).First(&existingUser).Error; err == nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(AccountResponse{
			Success: false,
			Message: "Username baru sudah digunakan",
		})
		return
	}

	// Update username
	if err := config.DB.Model(&user).Update("username", req.NewUsername).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AccountResponse{
			Success: false,
			Message: "Gagal mengubah username",
		})
		return
	}

	// Generate new JWT token with new username
	expireTime := time.Now().Add(24 * time.Hour)
	claims := MyClaims{
		Username: req.NewUsername,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   req.NewUsername,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JWTSecret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AccountResponse{
			Success: false,
			Message: "Gagal membuat token baru",
		})
		return
	}

	// Update cookie with new token
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  expireTime,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)

	json.NewEncoder(w).Encode(AccountResponse{
		Success: true,
		Message: "Username berhasil diubah",
		Token:   tokenString,
	})
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(AccountResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AccountResponse{
			Success: false,
			Message: "Invalid JSON format",
		})
		return
	}

	// Validate input
	if req.Username == "" || req.OldPassword == "" || req.NewPassword == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AccountResponse{
			Success: false,
			Message: "Semua field harus diisi",
		})
		return
	}

	// Validate new password strength (optional)
	if len(req.NewPassword) < 6 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AccountResponse{
			Success: false,
			Message: "Password baru minimal 6 karakter",
		})
		return
	}

	// Find user
	var user models.User
	result := config.DB.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(AccountResponse{
			Success: false,
			Message: "Username tidak ditemukan",
		})
		return
	}

	// Verify old password
	if !config.ComparePassword(user.Password, req.OldPassword) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(AccountResponse{
			Success: false,
			Message: "Password lama salah",
		})
		return
	}

	// Hash new password
	hashedPassword := config.HashPassword(req.NewPassword)

	// Update password
	if err := config.DB.Model(&user).Update("password", hashedPassword).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AccountResponse{
			Success: false,
			Message: "Gagal mengubah password",
		})
		return
	}

	json.NewEncoder(w).Encode(AccountResponse{
		Success: true,
		Message: "Password berhasil diubah",
	})
}
