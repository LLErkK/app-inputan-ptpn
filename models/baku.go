package models

import (
	"time"

	"gorm.io/gorm"
)

type BakuMandor struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	TahunTanam uint   `gorm:"not null"`
	Mandor     string `gorm:"size:100;not null"`
	Afdeling   string `gorm:"size:100;not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`

	// Relasi: Mandor punya banyak penyadap
	Penyadap []BakuPenyadap `gorm:"foreignKey:IdBakuMandor;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type BakuPenyadap struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	IdBakuMandor uint64    `gorm:"not null;index"` // FK ke mandor
	NIK          string    `gorm:"size:50;uniqueIndex;not null"`
	NamaPenyadap string    `gorm:"size:100;not null"`
	Tanggal      time.Time `gorm:"not null;index"`
	BasahLatex   float64   `gorm:"default:0"`
	Sheet        float64   `gorm:"default:0"`
	BasahLump    float64   `gorm:"default:0"`
	BrCr         float64   `gorm:"default:0"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	// Relasi ke Mandor
	Mandor BakuMandor `gorm:"foreignKey:IdBakuMandor;references:ID"`
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
