package models

import (
	"gorm.io/gorm"
	"time"
)

type Penyadap struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	NamaPenyadap string `gorm:"size:100;not null" json:"nama_penyadap"`
	NIK          string `gorm:"size:100;not null;uniqueIndex" json:"nik"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// FIXED: Hapus relasi ke BakuPenyadap karena tidak digunakan
	// BakuPenyadap tidak ada dalam daftar models yang digunakan
}

func (Penyadap) TableName() string {
	return "penyadaps"
}
