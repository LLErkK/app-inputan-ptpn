package config

import (
	"app-inputan-ptpn/models"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
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

	// Konfigurasi MySQL dari environment variables atau default values
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = ""
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "produksi_ptpn"
	}

	// MySQL DSN format
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	var err error
	// IMPORTANT: Disable foreign key constraint creation by GORM
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // KEY CHANGE
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Get underlying sql.DB
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("✓ MySQL database connected successfully")
	log.Printf("  Database: %s@%s:%s/%s", dbUser, dbHost, dbPort, dbName)

	// Auto migrate tables
	log.Println("🔄 Migrating database tables...")
	err = DB.AutoMigrate(
		// Independent tables (no foreign keys)
		&models.User{},
		&models.Upload{},
		&models.Master{},
		&models.Peta{},
		&models.Penyadap{},
		&models.Mandor{},
		// Tables with foreign keys
		&models.Produksi{},
		&models.Rekap{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Add foreign key constraints manually (jika diperlukan)
	log.Println("🔄 Adding foreign key constraints...")
	addForeignKeyConstraints()

	// Create default admin user if not exists
	createDefaultUser()

	log.Println("✓ MySQL database migrated successfully")
}

func GetDB() *gorm.DB {
	return DB
}

// HashPassword menggunakan bcrypt
func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
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
// createDefaultUser creates a default admin user
func createDefaultUser() {
	var count int64
	DB.Model(&models.User{}).Count(&count)

	if count == 0 {
		defaultUser := models.User{
			Username:  "admin",
			Password:  HashPassword("admin123"),
			LastLogin: time.Now(), // Add valid datetime value
		}

		result := DB.Create(&defaultUser)
		if result.Error != nil {
			log.Printf("Failed to create default user: %v", result.Error)
		} else {
			log.Println("✓ Default admin user created: username=admin, password=admin123")
		}
	}
}

// addForeignKeyConstraints adds foreign key constraints manually with proper CASCADE
func addForeignKeyConstraints() {
	// Hanya untuk Produksi dan Rekap yang punya FK ke Master
	constraints := []struct {
		table      string
		constraint string
		sql        string
	}{
		{
			table:      "produksis",
			constraint: "fk_produksis_master",
			sql: `ALTER TABLE produksis 
				  ADD CONSTRAINT fk_produksis_master 
				  FOREIGN KEY (id_master) REFERENCES masters(id) 
				  ON DELETE CASCADE ON UPDATE CASCADE`,
		},
		{
			table:      "rekaps",
			constraint: "fk_rekaps_master",
			sql: `ALTER TABLE rekaps 
				  ADD CONSTRAINT fk_rekaps_master 
				  FOREIGN KEY (id_master) REFERENCES masters(id) 
				  ON DELETE CASCADE ON UPDATE CASCADE`,
		},
	}

	for _, c := range constraints {
		// Check if constraint exists
		var count int64
		DB.Raw(`SELECT COUNT(*) FROM information_schema.table_constraints 
				WHERE constraint_schema = DATABASE() 
				AND table_name = ? 
				AND constraint_name = ?`, c.table, c.constraint).Scan(&count)

		if count == 0 {
			if err := DB.Exec(c.sql).Error; err != nil {
				// Log warning but don't fail - mungkin kolom belum ada
				log.Printf("  ⚠️  Could not add constraint %s: %v", c.constraint, err)
			} else {
				log.Printf("  ✓ Added constraint: %s", c.constraint)
			}
		} else {
			log.Printf("  ✓ Constraint exists: %s", c.constraint)
		}
	}
}
