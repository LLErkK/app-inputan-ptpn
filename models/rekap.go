package models

import "time"

type Rekap struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Tanggal          time.Time `gorm:"type:date;not null;index" json:"tanggal"`
	TipeProduksi     string    `gorm:"type:text;not null" json:"tipe_produksi"`
	TahunTanam       string    `gorm:"type:varchar(10);not null" json:"tahun_tanam"`
	NIK              string    `gorm:"type:varchar(20);not null;index" json:"nik"`
	Mandor           string    `gorm:"type:varchar(100);not null" json:"mandor"`
	HKOHariIni       int       `gorm:"default:0" json:"hko_hari_ini"`
	HKOSampaiHariIni int       `gorm:"default:0" json:"hko_sampai_hari_ini"`

	HariIniBasahLatekKebun  float64 `gorm:"type:decimal(10,2);default:0" json:"hari_ini_basah_latek_kebun"`
	HariIniBasahLatekPabrik float64 `gorm:"type:decimal(10,2);default:0" json:"hari_ini_basah_latek_pabrik"`
	HariIniBasahLatekPersen float64 `gorm:"type:decimal(5,2);default:0" json:"hari_ini_basah_latek_persen"`
	HariIniBasahLumpKebun   float64 `gorm:"type:decimal(10,2);default:0" json:"hari_ini_basah_lump_kebun"`
	HariIniBasahLumpPabrik  float64 `gorm:"type:decimal(10,2);default:0" json:"hari_ini_basah_lump_pabrik"`
	HariIniBasahLumpPersen  float64 `gorm:"type:decimal(5,2);default:0" json:"hari_ini_basah_lump_persen"`
	HariIniK3Sheet          float64 `gorm:"type:decimal(10,2);default:0" json:"hari_ini_k3_sheet"`
	HariIniKeringSheet      float64 `gorm:"type:decimal(10,2);default:0" json:"hari_ini_kering_sheet"`
	HariIniKeringBrCr       float64 `gorm:"type:decimal(10,2);default:0" json:"hari_ini_kering_br_cr"`
	HariIniKeringJumlah     float64 `gorm:"type:decimal(10,2);default:0" json:"hari_ini_kering_jumlah"`

	SampaiHariIniBasahLatekKebun  float64 `gorm:"type:decimal(10,2);default:0" json:"sampai_hari_ini_basah_latek_kebun"`
	SampaiHariIniBasahLatekPabrik float64 `gorm:"type:decimal(10,2);default:0" json:"sampai_hari_ini_basah_latek_pabrik"`
	SampaiHariIniBasahLatekPersen float64 `gorm:"type:decimal(5,2);default:0" json:"sampai_hari_ini_basah_latek_persen"`
	SampaiHariIniBasahLumpKebun   float64 `gorm:"type:decimal(10,2);default:0" json:"sampai_hari_ini_basah_lump_kebun"`
	SampaiHariIniBasahLumpPabrik  float64 `gorm:"type:decimal(10,2);default:0" json:"sampai_hari_ini_basah_lump_pabrik"`
	SampaiHariIniBasahLumpPersen  float64 `gorm:"type:decimal(5,2);default:0" json:"sampai_hari_ini_basah_lump_persen"`
	SampaiHariIniK3Sheet          float64 `gorm:"type:decimal(10,2);default:0" json:"sampai_hari_ini_k3_sheet"`
	SampaiHariIniKeringSheet      float64 `gorm:"type:decimal(10,2);default:0" json:"sampai_hari_ini_kering_sheet"`
	SampaiHariIniKeringBrCr       float64 `gorm:"type:decimal(10,2);default:0" json:"sampai_hari_ini_kering_br_cr"`
	SampaiHariIniKeringJumlah     float64 `gorm:"type:decimal(10,2);default:0" json:"sampai_hari_ini_kering_jumlah"`

	ProduksiPerTaperHariIni       float64 `gorm:"type:decimal(10,2);default:0" json:"produksi_per_taper_hari_ini"`
	ProduksiPerTaperSampaiHariIni float64 `gorm:"type:decimal(10,2);default:0" json:"produksi_per_taper_sampai_hari_ini"`

	TotalProduksi float64 `gorm:"type:decimal(10,2);default:0" json:"total_produksi"`

	Afdeling string `gorm:"type:varchar(100);not null;index" json:"afdeling"`

	// Foreign key - CASCADE sudah benar
	IdMaster uint64 `gorm:"not null;index" json:"id_master"`
	Master   Master `gorm:"foreignKey:IdMaster;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Rekap) TableName() string {
	return "rekaps" // FIXED: Plural form untuk konsistensi
}
