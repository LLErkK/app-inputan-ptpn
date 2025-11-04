package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"fmt"
)

// UpdatePenyadapMandor akan memanggil dua fungsi pembaruan data
func UpdatePenyadapMandor(idMaster uint64) {
	updatePenyadap(idMaster)
	updateMandor(idMaster)
}

// ------------------------------------------
// Update Penyadap berdasarkan data Produksi
// ------------------------------------------
func updatePenyadap(idMaster uint64) {
	db := config.GetDB()

	var produksis []models.Produksi
	if err := db.Where("id_master = ?", idMaster).Find(&produksis).Error; err != nil {
		fmt.Println("Gagal mengambil data produksi:", err)
		return
	}

	// Gunakan map untuk menghindari duplikasi berdasarkan NIK
	uniquePenyadap := make(map[string]models.Penyadap)

	for _, p := range produksis {
		if _, exists := uniquePenyadap[p.NIK]; !exists {
			uniquePenyadap[p.NIK] = models.Penyadap{
				NIK:          p.NIK,
				NamaPenyadap: p.NamaPenyadap,
			}
		}
	}

	// Simpan data unik ke tabel penyadap
	for _, penyadap := range uniquePenyadap {
		var existing models.Penyadap
		// Cek apakah sudah ada
		err := db.Where("nik = ?", penyadap.NIK).First(&existing).Error
		if err != nil {
			if err := db.Create(&penyadap).Error; err != nil {
				fmt.Println("Gagal menambahkan penyadap:", err)
			}
		}
	}
}

// ------------------------------------------
// Update Mandor berdasarkan data Rekap
// ------------------------------------------
func updateMandor(idMaster uint64) {
	db := config.GetDB()

	var rekaps []models.Rekap
	if err := db.Where("id_master = ?", idMaster).Find(&rekaps).Error; err != nil {
		fmt.Println("Gagal mengambil data rekap:", err)
		return
	}

	uniqueMandor := make(map[string]models.Mandor)

	for _, r := range rekaps {
		if _, exists := uniqueMandor[r.NIK]; !exists {
			uniqueMandor[r.NIK] = models.Mandor{
				NIK:        r.NIK,
				Nama:       r.Mandor,
				TahunTanam: r.TahunTanam,
			}
		}
	}

	for _, mandor := range uniqueMandor {
		var existing models.Mandor
		err := db.Where("nik = ? AND tahun_tanam", mandor.NIK, mandor.TahunTanam).First(&existing).Error
		if err != nil {
			if err := db.Create(&mandor).Error; err != nil {
				fmt.Println("Gagal menambahkan mandor:", err)
			}
		}
	}
}
