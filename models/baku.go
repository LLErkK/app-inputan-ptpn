package models

import (
	"time"

	"gorm.io/gorm"
)

type BakuMandor struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	TahunTanam uint   `gorm:"not null"`
	Mandor     string `gorm:"size:100;not null" json:"mandor"`
	Afdeling   string `gorm:"size:100;not null" json:"afdeling"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`

	// Relasi: Mandor punya banyak BakuPenyadap
	BakuPenyadaps []BakuPenyadap `gorm:"foreignKey:IdBakuMandor;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

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

type BakuPenyadap struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	IdBakuMandor uint64    `gorm:"not null;index" json:"idBakuMandor"` // FK ke Mandor
	IdPenyadap   uint64    `gorm:"not null;index" json:"idPenyadap"`   // FK ke master Penyadap
	Tanggal      time.Time `gorm:"not null;index" json:"tanggal"`

	BasahLatex float64 `gorm:"default:0" json:"basahLatex"`
	Sheet      float64 `gorm:"default:0" json:"sheet"`
	BasahLump  float64 `gorm:"default:0" json:"basahLump"`
	BrCr       float64 `gorm:"default:0" json:"brCr"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relasi ke Mandor
	Mandor BakuMandor `gorm:"foreignKey:IdBakuMandor;references:ID" json:"mandor"`
	// Relasi ke Penyadap
	Penyadap Penyadap `gorm:"foreignKey:IdPenyadap;references:ID" json:"penyadap"`
}

type BakuDetail struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	Tanggal time.Time `gorm:"not null;index"`

	JumlahPabrikBasahLatek      float64 `gorm:"default:0"`
	JumlahKebunBasahLatek       float64 `gorm:"default:0"`
	SelisihBasahLatek           float64 `gorm:"default:0"`
	PersentaseSelisihBasahLatek float64 `gorm:"default:0"`

	JumlahSheet float64 `gorm:"default:0"`
	K3Sheet     float64 `gorm:"default:0"`

	JumlahPabrikBasahLump      float64 `gorm:"default:0"`
	JumlahKebunBasahLump       float64 `gorm:"default:0"`
	SelisihBasahLump           float64 `gorm:"default:0"`
	PersentaseSelisihBasahLump float64 `gorm:"default:0"`

	JumlahBrCr float64 `gorm:"default:0"`
	K3BrCr     float64 `gorm:"default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
