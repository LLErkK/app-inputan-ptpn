// models/baku.go - Updated version

package models

import (
	"time"

	"gorm.io/gorm"
)

// Enum untuk tipe produksi
type TipeProduksi string

const (
	TipeBaku         TipeProduksi = "BAKU"
	TipeBakuBorong   TipeProduksi = "BAKU_BORONG"
	TipeBorgExternal TipeProduksi = "BORONG_EXTERNAL"
	TipeBorgInternal TipeProduksi = "BORONG_INTERNAL"
	TipeTetesLanjut  TipeProduksi = "TETES_LANJUT"
	TipeBorgMinggu   TipeProduksi = "BORONG_MINGGU"
)

// GetAllTipeProduksi returns all available production types
func GetAllTipeProduksi() []TipeProduksi {
	return []TipeProduksi{
		TipeBaku,
		TipeBakuBorong,
		TipeBorgExternal,
		TipeBorgInternal,
		TipeTetesLanjut,
		TipeBorgMinggu,
	}
}

// IsValidTipeProduksi checks if the given type is valid
func IsValidTipeProduksi(tipe TipeProduksi) bool {
	validTypes := GetAllTipeProduksi()
	for _, validType := range validTypes {
		if tipe == validType {
			return true
		}
	}
	return false
}

// UPDATED: BakuMandor now includes Tipe field
type BakuMandor struct {
	ID         uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	TahunTanam uint         `gorm:"not null" json:"tahun_tanam"`
	Mandor     string       `gorm:"size:100;not null" json:"mandor"`
	Afdeling   string       `gorm:"size:100;not null" json:"afdeling"`
	Tipe       TipeProduksi `gorm:"type:text; not null; default:'BAKU'; index" json:"tipe"` // NEW FIELD
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

// UPDATED: BakuPenyadap still has Tipe but will be auto-set from mandor
type BakuPenyadap struct {
	ID           uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	IdBakuMandor uint64       `gorm:"not null;index" json:"idBakuMandor"`
	IdPenyadap   uint64       `gorm:"not null;index" json:"idPenyadap"`
	Tanggal      time.Time    `gorm:"not null;index" json:"tanggal"`
	Tipe         TipeProduksi `gorm:"type:text; not null; default:'BAKU'; index" json:"tipe"` // Auto-set from mandor

	BasahLatex float64 `gorm:"default:0" json:"basahLatex"`
	Sheet      float64 `gorm:"default:0" json:"sheet"`
	BasahLump  float64 `gorm:"default:0" json:"basahLump"`
	BrCr       float64 `gorm:"default:0" json:"brCr"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Mandor   BakuMandor `gorm:"foreignKey:IdBakuMandor;references:ID" json:"mandor"`
	Penyadap Penyadap   `gorm:"foreignKey:IdPenyadap;references:ID" json:"penyadap"`
}

type BakuDetail struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	Tanggal  time.Time    `gorm:"not null;index" json:"tanggal"`
	Mandor   string       `gorm:"size:100;not null;index" json:"mandor"` // Nama mandor
	Afdeling string       `gorm:"size:100;not null" json:"afdeling"`     // Afdeling
	Tipe     TipeProduksi `gorm:"type:text; not null; default:'BAKU'; index" json:"tipe"`

	JumlahPabrikBasahLatek      float64 `gorm:"default:0" json:"jumlah_pabrik_basah_latek"`
	JumlahKebunBasahLatek       float64 `gorm:"default:0" json:"jumlah_kebun_basah_latek"`
	SelisihBasahLatek           float64 `gorm:"default:0" json:"selisih_basah_latek"`
	PersentaseSelisihBasahLatek float64 `gorm:"default:0" json:"persentase_selisih_basah_latek"`

	JumlahSheet float64 `gorm:"default:0" json:"jumlah_sheet"`
	K3Sheet     float64 `gorm:"default:0" json:"k3_sheet"`

	JumlahPabrikBasahLump      float64 `gorm:"default:0" json:"jumlah_pabrik_basah_lump"`
	JumlahKebunBasahLump       float64 `gorm:"default:0" json:"jumlah_kebun_basah_lump"`
	SelisihBasahLump           float64 `gorm:"default:0" json:"selisih_basah_lump"`
	PersentaseSelisihBasahLump float64 `gorm:"default:0" json:"persentase_selisih_basah_lump"`

	JumlahBrCr float64 `gorm:"default:0" json:"jumlah_br_cr"`
	K3BrCr     float64 `gorm:"default:0" json:"k3_br_cr"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// TableName untuk memastikan nama tabel yang benar
func (BakuDetail) TableName() string {
	return "baku_details"
}

// BeforeCreate hook untuk GORM
func (bd *BakuDetail) BeforeCreate(tx *gorm.DB) error {
	// Memastikan tanggal disimpan tanpa timestamp jam
	bd.Tanggal = bd.Tanggal.Truncate(24 * time.Hour)
	// Set default tipe jika kosong
	if bd.Tipe == "" {
		bd.Tipe = TipeBaku
	}
	return nil
}

// BeforeUpdate hook untuk GORM
func (bd *BakuDetail) BeforeUpdate(tx *gorm.DB) error {
	// Memastikan tanggal disimpan tanpa timestamp jam
	bd.Tanggal = bd.Tanggal.Truncate(24 * time.Hour)
	return nil
}

// UPDATED: BeforeCreate hook for BakuMandor to set default tipe
func (bm *BakuMandor) BeforeCreate(tx *gorm.DB) error {
	// Set default tipe jika kosong
	if bm.Tipe == "" {
		bm.Tipe = TipeBaku
	}
	return nil
}

// UPDATED: BeforeCreate hook untuk BakuPenyadap - tipe will be auto-set from mandor
func (bp *BakuPenyadap) BeforeCreate(tx *gorm.DB) error {
	// Tipe akan di-set otomatis dari mandor di controller, tidak perlu set default di sini
	return nil
}
