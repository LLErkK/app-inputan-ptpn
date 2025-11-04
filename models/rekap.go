package models

import "time"

type Rekap struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Tanggal          time.Time `gorm:"type:date;not null;index;column:tanggal" json:"tanggal"`
	TipeProduksi     string    `gorm:"type:text;not null;column:tipe_produksi" json:"tipe_produksi"`
	TahunTanam       string    `gorm:"type:varchar(10);not null;column:tahun_tanam" json:"tahun_tanam"`
	NIK              string    `gorm:"type:varchar(20);not null;index;column:nik" json:"nik"`
	Mandor           string    `gorm:"type:varchar(100);not null;column:mandor" json:"mandor"`
	HKOHariIni       int       `gorm:"default:0;column:hko_hari_ini" json:"hko_hari_ini"`
	HKOSampaiHariIni int       `gorm:"default:0;column:hko_sampai_hari_ini" json:"hko_sampai_hari_ini"`

	HariIniBasahLatekKebun  float64 `gorm:"type:decimal(10,2);default:0;column:hari_ini_basah_latek_kebun" json:"hari_ini_basah_latek_kebun"`
	HariIniBasahLatekPabrik float64 `gorm:"type:decimal(10,2);default:0;column:hari_ini_basah_latek_pabrik" json:"hari_ini_basah_latek_pabrik"`
	HariIniBasahLatekPersen float64 `gorm:"type:decimal(5,2);default:0;column:hari_ini_basah_latek_persen" json:"hari_ini_basah_latek_persen"`
	HariIniBasahLumpKebun   float64 `gorm:"type:decimal(10,2);default:0;column:hari_ini_basah_lump_kebun" json:"hari_ini_basah_lump_kebun"`
	HariIniBasahLumpPabrik  float64 `gorm:"type:decimal(10,2);default:0;column:hari_ini_basah_lump_pabrik" json:"hari_ini_basah_lump_pabrik"`
	HariIniBasahLumpPersen  float64 `gorm:"type:decimal(5,2);default:0;column:hari_ini_basah_lump_persen" json:"hari_ini_basah_lump_persen"`
	HariIniK3Sheet          float64 `gorm:"type:decimal(10,2);default:0;column:hari_ini_k3_sheet" json:"hari_ini_k3_sheet"`
	HariIniKeringSheet      float64 `gorm:"type:decimal(10,2);default:0;column:hari_ini_kering_sheet" json:"hari_ini_kering_sheet"`
	HariIniKeringBrCr       float64 `gorm:"type:decimal(10,2);default:0;column:hari_ini_kering_br_cr" json:"hari_ini_kering_br_cr"`
	HariIniKeringJumlah     float64 `gorm:"type:decimal(10,2);default:0;column:hari_ini_kering_jumlah" json:"hari_ini_kering_jumlah"`

	SampaiHariIniBasahLatekKebun  float64 `gorm:"type:decimal(10,2);default:0;column:sampai_hari_ini_basah_latek_kebun" json:"sampai_hari_ini_basah_latek_kebun"`
	SampaiHariIniBasahLatekPabrik float64 `gorm:"type:decimal(10,2);default:0;column:sampai_hari_ini_basah_latek_pabrik" json:"sampai_hari_ini_basah_latek_pabrik"`
	SampaiHariIniBasahLatekPersen float64 `gorm:"type:decimal(5,2);default:0;column:sampai_hari_ini_basah_latek_persen" json:"sampai_hari_ini_basah_latek_persen"`
	SampaiHariIniBasahLumpKebun   float64 `gorm:"type:decimal(10,2);default:0;column:sampai_hari_ini_basah_lump_kebun" json:"sampai_hari_ini_basah_lump_kebun"`
	SampaiHariIniBasahLumpPabrik  float64 `gorm:"type:decimal(10,2);default:0;column:sampai_hari_ini_basah_lump_pabrik" json:"sampai_hari_ini_basah_lump_pabrik"`
	SampaiHariIniBasahLumpPersen  float64 `gorm:"type:decimal(5,2);default:0;column:sampai_hari_ini_basah_lump_persen" json:"sampai_hari_ini_basah_lump_persen"`
	SampaiHariIniK3Sheet          float64 `gorm:"type:decimal(10,2);default:0;column:sampai_hari_ini_k3_sheet" json:"sampai_hari_ini_k3_sheet"`
	SampaiHariIniKeringSheet      float64 `gorm:"type:decimal(10,2);default:0;column:sampai_hari_ini_kering_sheet" json:"sampai_hari_ini_kering_sheet"`
	SampaiHariIniKeringBrCr       float64 `gorm:"type:decimal(10,2);default:0;column:sampai_hari_ini_kering_br_cr" json:"sampai_hari_ini_kering_br_cr"`
	SampaiHariIniKeringJumlah     float64 `gorm:"type:decimal(10,2);default:0;column:sampai_hari_ini_kering_jumlah" json:"sampai_hari_ini_kering_jumlah"`

	ProduksiPerTaperHariIni       float64 `gorm:"type:decimal(10,2);default:0;column:produksi_per_taper_hari_ini" json:"produksi_per_taper_hari_ini"`
	ProduksiPerTaperSampaiHariIni float64 `gorm:"type:decimal(10,2);default:0;column:produksi_per_taper_sampai_hari_ini" json:"produksi_per_taper_sampai_hari_ini"`

	Afdeling string `gorm:"type:text;not null;index;column:afdeling" json:"afdeling"`

	// Foreign key
	IdMaster uint64 `gorm:"not null;index;column:id_master" json:"id_master"`
	Master   Master `gorm:"foreignKey:IdMaster;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (Rekap) TableName() string {
	return "rekap"
}
