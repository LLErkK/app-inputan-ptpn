package models

import "time"

type ProduksiBaku struct {
	ID                     uint `gorm:"primarykey"`
	TahunTanam             uint
	Mandor                 string
	Tanggal                time.Time
	JumlahPabrikBasahLatex float64
	JumlahPabrikBasahLump  float64
	JumlahKebunBasahLatex  float64
	JumlahKebunBasahLump   float64
	JumlahPabrikSheet      float64
	JumlahPabrikBrCr       float64
	k3                     float64
	Selisih                float64
	SelisihPersentase      float64
}

type ProduksibakuDetail struct {
	ID           uint `gorm:"primarykey"`
	NIK          string
	NamaPenyadap string
	BasahLatex   float64
	Sheet        float64
	basahLump    float64
	BrCr         float64
}
