package models

import "time"

type Produksi struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Tanggal      time.Time `gorm:"type:date;not null;index" json:"tanggal"`
	TipeProduksi string    `gorm:"type:varchar(100);not null" json:"tipe_produksi"`
	TahunTanam   string    `gorm:"type:varchar(10);not null" json:"tahun_tanam"`
	Mandor       string    `gorm:"type:varchar(100);not null" json:"mandor"`
	NIK          string    `gorm:"type:varchar(50);not null;index" json:"nik"`
	NamaPenyadap string    `gorm:"type:varchar(100);not null" json:"nama_penyadap"`
	BasahLatek   float64   `gorm:"not null;default:0" json:"basah_latek"`
	Sheet        float64   `gorm:"not null;default:0" json:"sheet"`
	BasahLump    float64   `gorm:"not null;default:0" json:"basah_lump"`
	BrCr         float64   `gorm:"not null;default:0" json:"br_cr"`
	Afdeling     string    `gorm:"type:varchar(100);not null" json:"afdeling"`

	// Foreign key - CASCADE sudah benar
	IdMaster uint64 `gorm:"not null;index" json:"id_master"`
	Master   Master `gorm:"foreignKey:IdMaster;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Produksi) TableName() string {
	return "produksis" // FIXED: Plural form untuk konsistensi
}
