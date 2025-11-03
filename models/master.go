package models

import "time"

type Master struct {
	ID       uint64    `gorm:"primaryKey;autoIncrement"`
	Tanggal  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	Afdeling string    `gorm:"type:text;not null"`
	NamaFile string    `gorm:"type:text;not null"`

	// Relasi
	Rekaps    []Rekap    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:IdMaster"`
	Produksis []Produksi `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:IdMaster"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Master) TableName() string {
	return "master"
}
