// models/baku.go - Fixed version

package models

import (
	"time"

	"gorm.io/gorm"
)

// Enum untuk tipe produksi - FIXED: Consistent naming
type TipeProduksi string

const (
	TipeBaku         TipeProduksi = "BAKU"
	TipeBakuBorong   TipeProduksi = "BAKU_BORONG"
	TipeBorgExternal TipeProduksi = "BAKU_EKSTERNAL" // FIXED: Consistent with seed
	TipeBorgInternal TipeProduksi = "BAKU_INTERNAL"  // FIXED: Consistent with seed
	TipeTetesLanjut  TipeProduksi = "TETES_LANJUT"
	TipeBorgMinggu   TipeProduksi = "BAKU_MINGGU" // FIXED: Consistent with seed
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

type BakuMandor struct {
	ID         uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	TahunTanam uint         `gorm:"not null" json:"tahun_tanam"`
	NIK        string       `gorm:"not null" json:"nik"`
	Mandor     string       `gorm:"size:100;not null" json:"mandor"`
	Afdeling   string       `gorm:"size:100;not null" json:"afdeling"`
	Tipe       TipeProduksi `gorm:"type:text; not null; default:'BAKU'; index" json:"tipe"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`

	// Relasi: Mandor punya banyak BakuPenyadap
	BakuPenyadaps []BakuPenyadap `gorm:"foreignKey:IdBakuMandor;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// FIXED: Change foreign key type to uint (not uint64)
type BakuPenyadap struct {
	ID           uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	IdBakuMandor uint         `gorm:"not null;index" json:"idBakuMandor"` // FIXED: uint instead of uint64
	IdPenyadap   uint         `gorm:"not null;index" json:"idPenyadap"`   // FIXED: uint instead of uint64
	Tanggal      time.Time    `gorm:"not null;index" json:"tanggal"`
	Tipe         TipeProduksi `gorm:"type:text; not null; default:'BAKU'; index" json:"tipe"`
	TahunTanam   uint         `gorm:"" json:"tahun_tanam"`

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

	Tanggal      time.Time    `gorm:"not null;index" json:"tanggal"`
	IdBakuMandor uint         `gorm:"not null;index" json:"idBakuMandor"`
	Mandor       string       `gorm:"size:100;not null;index" json:"mandor"` // Nama mandor
	Afdeling     string       `gorm:"size:100;not null" json:"afdeling"`     // Afdeling
	TahunTanam   uint         `gorm:"" json:"tahun_tanam"`
	Tipe         TipeProduksi `gorm:"type:text; not null; default:'BAKU'; index" json:"tipe"`

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

func (bm *BakuMandor) GetTipeAsString() string {
	return string(bm.Tipe)
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

// BeforeCreate hook for BakuMandor to set default tipe
func (bm *BakuMandor) BeforeCreate(tx *gorm.DB) error {
	// Set default tipe jika kosong
	if bm.Tipe == "" {
		bm.Tipe = TipeBaku
	}
	return nil
}

// BeforeCreate hook untuk BakuPenyadap - set tanggal hari ini dan auto-set tipe from mandor
func (bp *BakuPenyadap) BeforeCreate(tx *gorm.DB) error {
	// Set tanggal ke hari ini jika kosong
	if bp.Tanggal.IsZero() {
		bp.Tanggal = time.Now().Truncate(24 * time.Hour)
	}

	// Auto-set tipe dari mandor jika kosong
	if bp.Tipe == "" && bp.IdBakuMandor > 0 {
		var mandor BakuMandor
		if err := tx.First(&mandor, bp.IdBakuMandor).Error; err == nil {
			bp.Tipe = mandor.Tipe
			bp.TahunTanam = mandor.TahunTanam
		}
	}

	return nil
}
