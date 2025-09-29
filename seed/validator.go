package seed

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"fmt"
)

// ValidateSeederData memvalidasi bahwa semua data master tersedia sebelum seeding transaksi
func ValidateSeederData() bool {
	fmt.Println("\n🔍 VALIDASI DATA MASTER SEBELUM SEEDING")
	fmt.Println("================================================")

	allValid := true

	// 1. Cek Mandor ID 1
	var mandor models.BakuMandor
	if err := config.DB.First(&mandor, 1).Error; err != nil {
		fmt.Printf("❌ Mandor ID 1 tidak ditemukan: %v\n", err)
		fmt.Println("   Solusi: Pastikan SeedMandor() menghasilkan minimal 1 mandor")
		allValid = false
	} else {
		fmt.Printf("✓ Mandor ID 1 ditemukan: %s (Afdeling: %s, Tipe: %s, Tahun: %d)\n",
			mandor.Mandor, mandor.Afdeling, mandor.Tipe, mandor.TahunTanam)
	}

	// 2. Cek Penyadap ID 1
	var penyadap models.Penyadap
	if err := config.DB.First(&penyadap, 1).Error; err != nil {
		fmt.Printf("❌ Penyadap ID 1 tidak ditemukan: %v\n", err)
		fmt.Println("   Solusi: Pastikan SeedPenyadap() menghasilkan minimal 1 penyadap")
		allValid = false
	} else {
		fmt.Printf("✓ Penyadap ID 1 ditemukan: %s (NIK: %s)\n", penyadap.NamaPenyadap, penyadap.NIK)
	}

	// 3. Cek jumlah total mandor
	var mandorCount int64
	config.DB.Model(&models.BakuMandor{}).Count(&mandorCount)
	fmt.Printf("ℹ️  Total Mandor di database: %d\n", mandorCount)
	if mandorCount == 0 {
		fmt.Println("❌ Tidak ada mandor di database!")
		allValid = false
	}

	// 4. Cek jumlah total penyadap
	var penyadapCount int64
	config.DB.Model(&models.Penyadap{}).Count(&penyadapCount)
	fmt.Printf("ℹ️  Total Penyadap di database: %d\n", penyadapCount)
	if penyadapCount == 0 {
		fmt.Println("❌ Tidak ada penyadap di database!")
		allValid = false
	}

	// 5. Cek apakah tipe BAKU valid
	if !models.IsValidTipeProduksi(models.TipeBaku) {
		fmt.Printf("❌ Tipe 'BAKU' tidak valid!\n")
		fmt.Println("   Solusi: Cek konstanta TipeBaku di models/baku.go")
		allValid = false
	} else {
		fmt.Printf("✓ Tipe produksi 'BAKU' valid\n")
	}

	// 6. List semua tipe produksi yang tersedia
	fmt.Println("\nℹ️  Tipe produksi yang tersedia:")
	for _, tipe := range models.GetAllTipeProduksi() {
		fmt.Printf("   - %s\n", tipe)
	}

	fmt.Println("================================================")
	if allValid {
		fmt.Println("✅ SEMUA VALIDASI LOLOS - Seeding dapat dilanjutkan")
	} else {
		fmt.Println("❌ VALIDASI GAGAL - Perbaiki masalah di atas sebelum seeding")
	}
	fmt.Println()

	return allValid
}

// DebugBakuPenyadapData menampilkan data BakuPenyadap untuk debugging
func DebugBakuPenyadapData() {
	fmt.Println("\n🔍 DEBUG: ISI TABEL BAKU_PENYADAPS")
	fmt.Println("================================================")

	var records []models.BakuPenyadap
	if err := config.DB.Limit(10).Preload("Mandor").Preload("Penyadap").Find(&records).Error; err != nil {
		fmt.Printf("❌ Error mengambil data: %v\n", err)
		return
	}

	if len(records) == 0 {
		fmt.Println("⚠️  Tabel baku_penyadaps kosong!")
		return
	}

	fmt.Printf("Menampilkan %d record pertama:\n\n", len(records))
	for i, r := range records {
		fmt.Printf("%d. ID=%d, Mandor=%s, Penyadap=%s, Tanggal=%s, Tipe=%s\n",
			i+1, r.ID, r.Mandor.Mandor, r.Penyadap.NamaPenyadap,
			r.Tanggal.Format("2006-01-02"), r.Tipe)
	}

	// Total count
	var total int64
	config.DB.Model(&models.BakuPenyadap{}).Count(&total)
	fmt.Printf("\nTotal BakuPenyadap: %d records\n", total)
	fmt.Println("================================================\n")
}

// CheckDuplicateEntries cek apakah ada duplikasi
func CheckDuplicateEntries() {
	fmt.Println("\n🔍 CEK DUPLIKASI DATA")
	fmt.Println("================================================")

	// Cek duplikasi berdasarkan tanggal + mandor + penyadap
	var duplicates []struct {
		IdBakuMandor uint
		IdPenyadap   uint
		Tanggal      string
		Count        int64
	}

	config.DB.Model(&models.BakuPenyadap{}).
		Select("id_baku_mandor, id_penyadap, DATE(tanggal) as tanggal, COUNT(*) as count").
		Group("id_baku_mandor, id_penyadap, DATE(tanggal)").
		Having("count > 1").
		Scan(&duplicates)

	if len(duplicates) == 0 {
		fmt.Println("✓ Tidak ada duplikasi data")
	} else {
		fmt.Printf("⚠️  Ditemukan %d duplikasi:\n", len(duplicates))
		for _, d := range duplicates {
			fmt.Printf("   - Mandor=%d, Penyadap=%d, Tanggal=%s (muncul %d kali)\n",
				d.IdBakuMandor, d.IdPenyadap, d.Tanggal, d.Count)
		}
	}
	fmt.Println("================================================\n")
}
