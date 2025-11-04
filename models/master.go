package models

import "time"

type Master struct {
	ID       uint64    `gorm:"primaryKey;autoIncrement"`
	Tanggal  time.Time `gorm:"type:date;not null"`
	Afdeling string    `gorm:"type:varchar(100);not null"`
	NamaFile string    `gorm:"type:varchar(255);not null"`

	// Relasi ke Produksi dan Rekap
	Produksis []Produksi `gorm:"foreignKey:IdMaster;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Rekaps    []Rekap    `gorm:"foreignKey:IdMaster;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Master) TableName() string {
	return "master"
}
