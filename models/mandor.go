package models

type Mandor struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	TahunTanam string `gorm:"not null" json:"tahun_tanam"`
	NIK        string `json:"nik"`
	Nama       string `json:"nama"`
}
