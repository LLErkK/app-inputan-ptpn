package controllers

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"
)

// excelToCSV converts Excel file to CSV with optimizations
func excelToCSV(excelFile string, outputFolder string, tanggal time.Time, afdeling string, originalFileName string) error {
	// Jika outputFolder kosong, gunakan folder yang sama dengan file Excel
	if outputFolder == "" {
		outputFolder = filepath.Dir(excelFile)
		if outputFolder == "" {
			outputFolder = "."
		}
	}

	// Buat folder output jika belum ada
	if err := os.MkdirAll(outputFolder, 0755); err != nil {
		return fmt.Errorf("gagal membuat folder output: %v", err)
	}

	// Buka file Excel dengan options untuk performa
	f, err := excelize.OpenFile(excelFile, excelize.Options{
		UnzipSizeLimit: 100 * 1024 * 1024, // 100MB limit
	})
	if err != nil {
		return fmt.Errorf("gagal membuka file Excel: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Warning: gagal menutup file Excel: %v", err)
		}
	}()

	sheets := f.GetSheetList()

	fmt.Printf("Mengkonversi file: %s\n", excelFile)
	fmt.Printf("Total sheets: %d\n\n", len(sheets))

	// Process sheets concurrently with worker pool
	type sheetJob struct {
		name string
		idx  int
	}

	jobs := make(chan sheetJob, len(sheets))
	results := make(chan error, len(sheets))

	// Worker pool dengan 4 workers
	numWorkers := 4
	var wg sync.WaitGroup

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				if err := processSheet(f, job.name, outputFolder); err != nil {
					results <- fmt.Errorf("sheet %s: %v", job.name, err)
				} else {
					results <- nil
				}
			}
		}()
	}

	// Send jobs
	for idx, sheetName := range sheets {
		jobs <- sheetJob{name: sheetName, idx: idx}
	}
	close(jobs)

	// Wait for all workers
	wg.Wait()
	close(results)

	// Check for errors
	var errors []string
	for err := range results {
		if err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		log.Printf("Beberapa sheet gagal diproses: %v", errors)
	}

	fmt.Printf("\nSelesai! %d file CSV telah dibuat di: %s\n", len(sheets), outputFolder)

	// Process database operations
	tanggalInt := tanggal.Day()
	fmt.Printf("Memproses membuat table master dengan nama file: %s\n", originalFileName)

	idMaster, err := CreateMaster(tanggal, afdeling, originalFileName)
	if err != nil {
		return fmt.Errorf("gagal membuat master: %v", err)
	}

	fmt.Println("\nMemproses CSV ke database...")

	// Run database operations concurrently
	type dbResult struct {
		name   string
		saved  int
		failed int
		errs   []string
		err    error
	}

	dbResults := make(chan dbResult, 2)

	// Function 1: ConvertCSVAutoBaseWithFilter
	go func() {
		saved, failed, errs, err := ConvertCSVAutoBaseWithFilter(tanggal, afdeling, idMaster)
		dbResults <- dbResult{
			name:   "ConvertCSVAutoBaseWithFilter",
			saved:  saved,
			failed: failed,
			errs:   errs,
			err:    err,
		}
	}()

	// Function 2: ConvertCSVTanggalFormat
	go func() {
		saved, failed, errs, err := ConvertCSVTanggalFormat(tanggalInt, afdeling, idMaster)
		dbResults <- dbResult{
			name:   "ConvertCSVTanggalFormat",
			saved:  saved,
			failed: failed,
			errs:   errs,
			err:    err,
		}
	}()

	// Collect results
	successCount := 0
	for i := 0; i < 2; i++ {
		result := <-dbResults
		if result.err != nil {
			fmt.Printf("✗ %s gagal: %v\n", result.name, result.err)
		} else {
			fmt.Printf("✓ %s: %d berhasil, %d gagal\n", result.name, result.saved, result.failed)
			successCount++
			if len(result.errs) > 0 {
				fmt.Println("  Detail error:")
				for _, e := range result.errs {
					fmt.Printf("   - %s\n", e)
				}
			}
		}
	}
	close(dbResults)

	// Evaluate results
	if successCount == 2 {
		fmt.Println("\n✅ Semua proses berhasil dilakukan!")
		// Update table mandor dan penyadap
		go UpdatePenyadapMandor(idMaster) // Run async
	} else {
		fmt.Println("\n⚠️  Beberapa proses gagal, periksa log di atas.")
	}

	return nil
}

// processSheet processes a single Excel sheet to CSV
func processSheet(f *excelize.File, sheetName string, outputFolder string) error {
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("gagal membaca sheet: %v", err)
	}

	// Buat nama file CSV sama persis dengan nama sheet
	cleanName := strings.TrimSpace(sheetName)
	csvFilename := fmt.Sprintf("%s.csv", cleanName)
	csvPath := filepath.Join(outputFolder, csvFilename)

	csvFile, err := os.Create(csvPath)
	if err != nil {
		return fmt.Errorf("gagal membuat file CSV: %v", err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	writer.Comma = ','

	// Write rows with error handling
	for rowIdx, row := range rows {
		if err := writer.Write(row); err != nil {
			log.Printf("Warning: gagal menulis baris %d ke CSV %s: %v\n", rowIdx, csvFilename, err)
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return fmt.Errorf("error flushing CSV writer: %v", err)
	}

	fmt.Printf("✓ Sheet '%s' → %s (Rows: %d)\n", sheetName, csvFilename, len(rows))
	return nil
}

// clearFolder clears all files in a folder with improved error handling
func clearFolder(folder string) error {
	// Check if folder exists
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		return nil // Folder doesn't exist, nothing to clear
	}

	files, err := os.ReadDir(folder)
	if err != nil {
		return fmt.Errorf("gagal membaca folder %s: %w", folder, err)
	}

	var wg sync.WaitGroup
	errors := make(chan error, len(files))

	for _, file := range files {
		wg.Add(1)
		go func(f os.DirEntry) {
			defer wg.Done()
			path := filepath.Join(folder, f.Name())

			if f.IsDir() {
				// hapus folder secara rekursif
				if err := os.RemoveAll(path); err != nil {
					errors <- fmt.Errorf("gagal menghapus subfolder %s: %v", path, err)
				}
			} else {
				// hapus file biasa
				if err := os.Remove(path); err != nil {
					errors <- fmt.Errorf("gagal menghapus file %s: %v", path, err)
				}
			}
		}(file)
	}

	wg.Wait()
	close(errors)

	// Collect errors
	var errorList []string
	for err := range errors {
		errorList = append(errorList, err.Error())
	}

	if len(errorList) > 0 {
		return fmt.Errorf("beberapa file gagal dihapus: %v", errorList)
	}

	return nil
}
