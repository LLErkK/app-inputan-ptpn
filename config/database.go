package config

import (
	"app-inputan-ptpn/models"
	"crypto/sha256"
	"encoding/hex"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// Nama file database SQLite, akan dibuat otomatis kalau belum ada
	dsn := "produksi.db"

	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate tables
	err = DB.AutoMigrate(
		&models.User{},
		&models.ProduksiBaku{},
		&models.ProduksibakuDetail{},
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

// HashPassword creates SHA256 hash of password
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
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
