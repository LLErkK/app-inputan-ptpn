package seed

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"fmt"
	"log"
	"time"
)

// SeedUsers creates default admin user if users table is empty
func SeedUsers() {
	fmt.Println("\n📝 Seeding users...")

	var count int64
	config.DB.Model(&models.User{}).Count(&count)

	if count > 0 {
		fmt.Printf("⏭️  Users table already has %d record(s), skipping seed\n", count)
		return
	}

	// Create default admin user with valid last_login time
	adminUser := models.User{
		Username:  "admin",
		Password:  config.HashPassword("admin123"),
		LastLogin: time.Now(),
	}

	result := config.DB.Create(&adminUser)
	if result.Error != nil {
		log.Printf("❌ Failed to create admin user: %v", result.Error)
		return
	}

	fmt.Println("✅ Default admin user created successfully")
	fmt.Println("   Username: admin")
	fmt.Println("   Password: admin123")
	fmt.Println("   ⚠️  Please change the password after first login!")
}
