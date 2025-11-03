package models

import "time"

type Rekap struct {
	ID               uint      `gorm:"primaryKey;autoIncrement"`
	Tanggal          time.Time `gorm:"type:date;not null;index"`
	TipeProduksi     string    `gorm:"type:text;not null"`
	TahunTanam       string    `gorm:"type:varchar(10);not null"`
	NIK              string    `gorm:"type:varchar(20);not null;index"`
	Mandor           string    `gorm:"type:varchar(100);not null"`
	HKOHariIni       int       `gorm:"default:0"`
	HKOSampaiHariIni int       `gorm:"default:0"`

	// Produksi Hari Ini
	HariIniBasahLatekKebun  float64 `gorm:"type:decimal(10,2);default:0"`
	HariIniBasahLatekPabrik float64 `gorm:"type:decimal(10,2);default:0"`
	HariIniBasahLatekPersen float64 `gorm:"type:decimal(5,2);default:0"`
	HariIniBasahLumpKebun   float64 `gorm:"type:decimal(10,2);default:0"`
	HariIniBasahLumpPabrik  float64 `gorm:"type:decimal(10,2);default:0"`
	HariIniBasahLumpPersen  float64 `gorm:"type:decimal(5,2);default:0"`
	HariIniK3Sheet          float64 `gorm:"type:decimal(10,2);default:0"`
	HariIniKeringSheet      float64 `gorm:"type:decimal(10,2);default:0"`
	HariIniKeringBrCr       float64 `gorm:"type:decimal(10,2);default:0"`
	HariIniKeringJumlah     float64 `gorm:"type:decimal(10,2);default:0"`

	// Produksi Sampai Hari Ini (Kumulatif)
	SampaiHariIniBasahLatekKebun  float64 `gorm:"type:decimal(10,2);default:0"`
	SampaiHariIniBasahLatekPabrik float64 `gorm:"type:decimal(10,2);default:0"`
	SampaiHariIniBasahLatekPersen float64 `gorm:"type:decimal(5,2);default:0"`
	SampaiHariIniBasahLumpKebun   float64 `gorm:"type:decimal(10,2);default:0"`
	SampaiHariIniBasahLumpPabrik  float64 `gorm:"type:decimal(10,2);default:0"`
	SampaiHariIniBasahLumpPersen  float64 `gorm:"type:decimal(5,2);default:0"`
	SampaiHariIniK3Sheet          float64 `gorm:"type:decimal(10,2);default:0"`
	SampaiHariIniKeringSheet      float64 `gorm:"type:decimal(10,2);default:0"`
	SampaiHariIniKeringBrCr       float64 `gorm:"type:decimal(10,2);default:0"`
	SampaiHariIniKeringJumlah     float64 `gorm:"type:decimal(10,2);default:0"`

	// Produktivitas
	ProduksiPerTaperHariIni       float64 `gorm:"type:decimal(10,2);default:0"`
	ProduksiPerTaperSampaiHariIni float64 `gorm:"type:decimal(10,2);default:0"`

	afdeling string `gorm:"type:text;not null;index"`

	// Timestamps
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName overrides the table name
func (Rekap) TableName() string {
	return "rekap"
}
