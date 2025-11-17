package models

import "time"

type Master struct {
	ID       uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Tanggal  time.Time `gorm:"type:date;not null" json:"tanggal"`
	Afdeling string    `gorm:"type:varchar(100);not null" json:"afdeling"`
	NamaFile string    `gorm:"type:varchar(255);not null" json:"nama_file"`

	// Relasi ke Produksi dan Rekap - CASCADE sudah benar
	Produksis []Produksi `gorm:"foreignKey:IdMaster;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Rekaps    []Rekap    `gorm:"foreignKey:IdMaster;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Master) TableName() string {
	return "masters" // FIXED: Plural form untuk konsistensi
}
