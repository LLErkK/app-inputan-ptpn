package models

import (
	"time"

	"gorm.io/gorm"
)

// Upload represents the file upload data structure
type Upload struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Tanggal   time.Time      `gorm:"not null;index" json:"tanggal"`
	FileName  string         `gorm:"type:varchar(255);not null" json:"fileName"`
	FilePath  string         `gorm:"type:varchar(500);not null" json:"filePath"`
	FileSize  int64          `gorm:"not null" json:"fileSize"`
	MimeType  string         `gorm:"type:varchar(100)" json:"mimeType"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName specifies the table name for Upload model
func (Upload) TableName() string {
	return "uploads"
}

// UploadResponse represents the response structure
type UploadResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	ID       uint   `json:"id,omitempty"`
	Tanggal  string `json:"tanggal,omitempty"`
	FileName string `json:"fileName,omitempty"`
	FileSize int64  `json:"fileSize,omitempty"`
	FilePath string `json:"filePath,omitempty"`
}
