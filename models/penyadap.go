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

	// Relasi: Penyadap punya banyak BakuPenyadap
	BakuPenyadaps []BakuPenyadap `gorm:"foreignKey:IdPenyadap;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}
