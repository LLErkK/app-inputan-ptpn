package config

import (
	"app-inputan-ptpn/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// Nama file database SQLite, akan dibuat otomatis kalau belum ada
	dsn := "produksi .db"

	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate tables
	err = DB.AutoMigrate(
		&models.ProduksibakuDetail{},
		&models.ProduksiBaku{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("SQLite database connected and migrated successfully")
}

func GetDB() *gorm.DB {
	return DB
}
