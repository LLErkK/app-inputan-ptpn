package config

import (
	"app-inputan-ptpn/models"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// JWT secret (baca dari env JWT_SECRET jika tersedia)
var JWTSecret []byte

func InitDB() {
	// init JWT secret
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "please-change-this-secret-in-production"
	}
	JWTSecret = []byte(secret)

	// Nama file database SQLite
	dsn := "produksi.db"

	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate tables
	err = DB.AutoMigrate(
		&models.User{},
		&models.BakuPenyadap{},
		&models.BakuMandor{},
		&models.BakuDetail{},
		&models.Upload{},
		&models.Master{},
		&models.Rekap{},
		&models.Produksi{},
		&models.Penyadap{},
		&models.Mandor{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Create default admin user if not exists
	createDefaultUser()

	log.Println("SQLite database connected and migrated successfully")
}

func GetDB() *gorm.DB {
	return DB
}

// HashPassword menggunakan bcrypt
func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// fallback (jarang terjadi)
		log.Printf("bcrypt generate error: %v", err)
		return password
	}
	return string(hash)
}

// ComparePassword membandingkan bcrypt hash dengan password plain
func ComparePassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// createDefaultUser creates a default admin user
func createDefaultUser() {
	var count int64
	DB.Model(&models.User{}).Count(&count)

	if count == 0 {
		defaultUser := models.User{
			Username: "admin",
			Password: HashPassword("admin123"),
		}

		result := DB.Create(&defaultUser)
		if result.Error != nil {
			log.Printf("Failed to create default user: %v", result.Error)
		} else {
			log.Println("Default admin user created: username=admin, password=admin123")
		}
	}
}
