package controllers

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func excelToCSV(excelFile string, outputFolder string, tanggal time.Time, afdeling string) error {
	// Jika outputFolder kosong, gunakan folder yang sama dengan file Excel
	if outputFolder == "" {
		outputFolder = filepath.Dir(excelFile)
		if outputFolder == "" {
			outputFolder = "."
		}
	}

	// Buat folder output jika belum ada
	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		err := os.MkdirAll(outputFolder, os.ModePerm)
		if err != nil {
			return fmt.Errorf("gagal membuat folder output: %v", err)
		}
	}

	// Buka file Excel
	f, err := excelize.OpenFile(excelFile)
	if err != nil {
		return fmt.Errorf("gagal membuka file Excel: %v", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()

	fmt.Printf("Mengkonversi file: %s\n", excelFile)
	fmt.Printf("Total sheets: %d\n\n", len(sheets))

	for _, sheetName := range sheets {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			log.Printf("Gagal membaca sheet %s: %v\n", sheetName, err)
			continue
		}

		// Buat nama file CSV sama persis dengan nama sheet
		cleanName := strings.TrimSpace(sheetName)
		csvFilename := fmt.Sprintf("%s.csv", cleanName)
		csvPath := filepath.Join(outputFolder, csvFilename)

		csvFile, err := os.Create(csvPath)
		if err != nil {
			log.Printf("Gagal membuat file CSV %s: %v\n", csvPath, err)
			continue
		}

		writer := csv.NewWriter(csvFile)
		for _, row := range rows {
			if err := writer.Write(row); err != nil {
				log.Printf("Gagal menulis baris ke CSV %s: %v\n", csvFilename, err)
			}
		}
		writer.Flush()
		csvFile.Close()

		fmt.Printf("✓ Sheet '%s' → %s\n", sheetName, csvFilename)
		fmt.Printf("  Rows: %d\n", len(rows))
	}

	fmt.Printf("\nSelesai! %d file CSV telah dibuat di: %s\n", len(sheets), outputFolder)

	// Ambil tanggal (hari) sebagai int
	tanggalInt := tanggal.Day()
	fmt.Printf("memproses membuat table master")
	idMaster, err := CreateMaster(tanggal, afdeling, excelFile)

	fmt.Println("\nMemproses CSV ke database...")

	// --- Jalankan fungsi pertama ---
	saved1, failed1, errs1, err1 := ConvertCSVAutoBaseWithFilter(tanggal, afdeling, idMaster)
	if err1 != nil {
		fmt.Printf("✗ ConvertCSVAutoBaseWithFilter gagal: %v\n", err1)
	} else {
		fmt.Printf("✓ ConvertCSVAutoBaseWithFilter: %d berhasil, %d gagal\n", saved1, failed1)
		if len(errs1) > 0 {
			fmt.Println("  Detail error:")
			for _, e := range errs1 {
				fmt.Printf("   - %s\n", e)
			}
		}
	}

	// --- Jalankan fungsi kedua ---
	saved2, failed2, errs2, err2 := ConvertCSVTanggalFormat(tanggalInt, afdeling, idMaster)
	if err2 != nil {
		fmt.Printf("✗ ConvertCSVTanggalFormat gagal: %v\n", err2)
	} else {
		fmt.Printf("✓ ConvertCSVTanggalFormat: %d berhasil, %d gagal\n", saved2, failed2)
		if len(errs2) > 0 {
			fmt.Println("  Detail error:")
			for _, e := range errs2 {
				fmt.Printf("   - %s\n", e)
			}
		}
	}

	// --- Evaluasi hasil ---
	if err1 == nil && err2 == nil {
		fmt.Println("\n✅ Semua proses berhasil dilakukan!")

		//update table mandor dan penyadap
		UpdatePenyadapMandor(idMaster)

		// Jika kedua fungsi berhasil → hapus isi folder uploads & csv
		if err := clearFolder("uploads"); err != nil {
			fmt.Printf("⚠️  Gagal menghapus isi folder uploads: %v\n", err)
		} else {
			fmt.Println("🗑️  Folder 'uploads' telah dibersihkan.")
		}

		if err := clearFolder("csv"); err != nil {
			fmt.Printf("⚠️  Gagal menghapus isi folder csv: %v\n", err)
		} else {
			fmt.Println("🗑️  Folder 'csv' telah dibersihkan.")
		}

	} else {
		fmt.Println("\n⚠️  Beberapa proses gagal, periksa log di atas.")
	}

	return nil
}

// Fungsi bantu untuk menghapus semua file di dalam folder tertentu
func clearFolder(folder string) error {
	files, err := os.ReadDir(folder)
	if err != nil {
		return fmt.Errorf("gagal membaca folder %s: %w", folder, err)
	}

	for _, file := range files {
		path := filepath.Join(folder, file.Name())
		if file.IsDir() {
			// hapus folder secara rekursif
			if err := os.RemoveAll(path); err != nil {
				log.Printf("Gagal menghapus subfolder %s: %v", path, err)
			}
		} else {
			// hapus file biasa
			if err := os.Remove(path); err != nil {
				log.Printf("Gagal menghapus file %s: %v", path, err)
			}
		}
	}

	return nil
}
