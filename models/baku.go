package models

import (
	"time"
)

type ProduksiBaku struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	TahunTanam   int       `json:"tahun_tanam"`
	Mandor       string    `json:"mandor"`
	NIK          string    `json:"nik"`
	NamaPenyadap string    `json:"nama_penyadap"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
type DetailProduksi struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ProduksiBakuID uint      `json:"produksi_baku_id"`
	JumlahPabrik   float64   `json:"jumlah_pabrik"`
	K3             float64   `json:"k3"`
	JumlahKebun    float64   `json:"jumlah_kebun"`
	Selisih        float64   `json:"selisih"`
	Persentase     float64   `json:"persentase"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
type Ringkasan struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Tahun           int       `json:"tahun"`
	TotalPabrik     float64   `json:"total_pabrik"`
	TotalK3         float64   `json:"total_k3"`
	TotalKebun      float64   `json:"total_kebun"`
	TotalSelisih    float64   `json:"total_selisih"`
	TotalPersentase float64   `json:"total_persentase"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
