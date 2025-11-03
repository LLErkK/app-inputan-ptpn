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

func excelToCSV(excelFile string, outputFolder string, tanggal time.Time) error {
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

	//masukan ke database
	ConvertCSVAutoBaseWithFilter(tanggal)
	return nil
}
