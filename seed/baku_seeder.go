package seed

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"time"
)

func SeedBaku() {
	bakus := []models.BakuPenyadap{
		{IdBakuMandor: 0, IdPenyadap: 0, Tanggal: time.Now(), Tipe: "BAKU", TahunTanam: 0, BasahLatex: 0, Sheet: 0, BasahLump: 0, BrCr: 0},
	}
	config.DB.Create(&bakus)
}
