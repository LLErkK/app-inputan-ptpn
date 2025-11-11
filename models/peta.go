package models

type Peta struct {
	ID          uint    `gorm:"primary_key"`
	Blok        string  `gorm:"type:varchar(255);"`
	Code        string  `gorm:"type:varchar(255);not null"`
	Afdeling    string  `gorm:"type:varchar(100);not null"`
	Luas        float32 `gorm:"not null;default:0"`
	JumlahPohon int64   `gorm:"not null;default:0"`
	JenisKebun  string  `gorm:"type:varchar(255);"`
	TahunTanam  string  `gorm:"type:varchar(255);"`
	Kloon       string  `gorm:"type:varchar(255);"`
}
