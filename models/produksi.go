package models

import "time"

type Produksi struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	Tanggal      time.Time `gorm:"type:date;not null;index"`
	TipeProduksi string    `gorm:"type:varchar(100);not null"`
	TahunTanam   string    `gorm:"type:varchar(10);not null"`
	Mandor       string    `gorm:"type:varchar(100);not null"`
	NIK          string    `gorm:"type:varchar(50);not null"`
	NamaPenyadap string    `gorm:"type:varchar(100);not null"`
	BasahLatek   float64   `gorm:"not null;default:0"`
	Sheet        float64   `gorm:"not null;default:0"`
	BasahLump    float64   `gorm:"not null;default:0"`
	BrCr         float64   `gorm:"not null;default:0"`
	Afdeling     string    `gorm:"type:text;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (Produksi) TableName() string {
	return "produksi"
}
