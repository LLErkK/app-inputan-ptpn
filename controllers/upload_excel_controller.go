package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

const (
	MaxFileSize = 10 * 1024 * 1024 // 10MB
	UploadDir   = "./uploads"
)

// ServeUploadPage serves the upload HTML page
func ServeUploadPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/html/upload.html")
}

// CreateUpload handles file upload and date submission with optimizations
func CreateUpload(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
	defer cancel()

	// Parse multipart form with max memory
	if err := r.ParseMultipartForm(MaxFileSize); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "File terlalu besar atau format tidak valid",
		})
		return
	}

	// Get afdeling from form
	afdeling := r.FormValue("afdeling")
	if afdeling == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Afdeling wajib diisi",
		})
		return
	}

	// Get tanggal from form
	tanggalStr := r.FormValue("tanggal")
	if tanggalStr == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Tanggal wajib diisi",
		})
		return
	}

	// Parse tanggal
	tanggal, err := time.Parse("2006-01-02", tanggalStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Format tanggal tidak valid. Gunakan format YYYY-MM-DD",
		})
		return
	}

	// Get uploaded file
	file, header, err := r.FormFile("file")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "File tidak ditemukan atau gagal diupload",
		})
		return
	}
	defer file.Close()

	// Validate file size
	if header.Size > MaxFileSize {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: fmt.Sprintf("Ukuran file terlalu besar (maksimal %dMB)", MaxFileSize/(1024*1024)),
		})
		return
	}

	// Validate file extension
	ext := filepath.Ext(header.Filename)
	if ext != ".xlsx" && ext != ".xls" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Format file tidak didukung. Hanya .xlsx dan .xls yang diizinkan",
		})
		return
	}

	// Ensure upload directory exists
	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal membuat direktori upload",
		})
		return
	}

	// Generate unique filename
	newFileName := fmt.Sprintf("%d_%s%s", time.Now().Unix(), generateRandomString(8), ext)
	uploadPath := filepath.Join(UploadDir, newFileName)

	// Create destination file
	dst, err := os.Create(uploadPath)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal membuat file di server",
		})
		return
	}

	// Copy uploaded file to destination with buffered writing
	buf := make([]byte, 32*1024) // 32KB buffer
	if _, err := io.CopyBuffer(dst, file, buf); err != nil {
		dst.Close()
		os.Remove(uploadPath)
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal menyimpan file",
		})
		return
	}
	dst.Close()

	// Create upload record in database
	upload := models.Upload{
		Tanggal:  tanggal,
		FileName: header.Filename,
		FilePath: uploadPath,
		FileSize: header.Size,
		MimeType: header.Header.Get("Content-Type"),
	}

	// Save to database
	if err := config.DB.Create(&upload).Error; err != nil {
		os.Remove(uploadPath)
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal menyimpan data ke database: " + err.Error(),
		})
		return
	}

	// Process Excel to CSV and database in background with context
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered in background process: %v", r)
			}
		}()

		// Check if context is still valid
		select {
		case <-ctx.Done():
			log.Println("Context cancelled, stopping background process")
			return
		default:
		}

		log.Printf("Starting Excel to CSV conversion for: %s", header.Filename)

		if err := excelToCSV(uploadPath, "csv", tanggal, afdeling, header.Filename); err != nil {
			log.Printf("Error converting Excel to CSV: %v", err)
			return
		}

		log.Println("Excel to CSV conversion completed successfully")

		// Cleanup folders asynchronously
		cleanupFolders()
	}()

	// Return success response immediately
	responseData := map[string]interface{}{
		"id":       upload.ID,
		"tanggal":  upload.Tanggal.Format("2006-01-02"),
		"fileName": upload.FileName,
		"fileSize": upload.FileSize,
		"filePath": "/uploads/" + newFileName,
		"message":  "File sedang diproses di background",
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "File berhasil diupload dan sedang diproses",
		Data:    responseData,
	})
}

// cleanupFolders cleans up upload and csv folders concurrently
func cleanupFolders() {
	var wg sync.WaitGroup

	folders := []string{"uploads", "csv"}

	for _, folder := range folders {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			if err := clearFolder(f); err != nil {
				log.Printf("⚠️  Gagal menghapus isi folder %s: %v\n", f, err)
			} else {
				log.Printf("🗑️  Folder '%s' telah dibersihkan.\n", f)
			}
		}(folder)
	}

	wg.Wait()
}

// GetAllUploads retrieves all upload records with pagination
func GetAllUploads(w http.ResponseWriter, r *http.Request) {
	var uploads []models.Upload

	// Parse query parameters for filtering and pagination
	tanggalStr := r.URL.Query().Get("tanggal")
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	// Default pagination values
	pageNum := 1
	limitNum := 50

	if page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pageNum = p
		}
	}

	if limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			limitNum = l
		}
	}

	offset := (pageNum - 1) * limitNum

	query := config.DB.Order("created_at desc")

	// Filter by tanggal if provided
	if tanggalStr != "" {
		tanggal, err := time.Parse("2006-01-02", tanggalStr)
		if err != nil {
			respondJSON(w, http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "Format tanggal tidak valid",
			})
			return
		}
		query = query.Where("DATE(tanggal) = DATE(?)", tanggal)
	}

	// Get total count
	var total int64
	if err := query.Model(&models.Upload{}).Count(&total).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal menghitung total data: " + err.Error(),
		})
		return
	}

	// Get paginated data
	if err := query.Limit(limitNum).Offset(offset).Find(&uploads).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data: " + err.Error(),
		})
		return
	}

	response := map[string]interface{}{
		"uploads": uploads,
		"pagination": map[string]interface{}{
			"page":       pageNum,
			"limit":      limitNum,
			"total":      total,
			"totalPages": (total + int64(limitNum) - 1) / int64(limitNum),
		},
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data berhasil diambil",
		Data:    response,
	})
}

// GetUploadByID retrieves a single upload record by ID
func GetUploadByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var upload models.Upload

	if err := config.DB.First(&upload, id).Error; err != nil {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Data tidak ditemukan",
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data berhasil ditemukan",
		Data:    upload,
	})
}

// DeleteUpload removes an upload record and its file
func DeleteUpload(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var upload models.Upload

	// Find the upload record
	if err := config.DB.First(&upload, id).Error; err != nil {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Data tidak ditemukan",
		})
		return
	}

	// Delete file from filesystem asynchronously
	go func() {
		if err := os.Remove(upload.FilePath); err != nil {
			log.Printf("Warning: Failed to delete file %s: %v\n", upload.FilePath, err)
		}
	}()

	// Delete from database
	if err := config.DB.Delete(&upload).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal menghapus data dari database",
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data dan file berhasil dihapus",
	})
}

// DownloadFile serves the uploaded file for download
func DownloadFile(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var upload models.Upload

	// Find the upload record
	if err := config.DB.First(&upload, id).Error; err != nil {
		http.Error(w, "File tidak ditemukan", http.StatusNotFound)
		return
	}

	// Check if file exists
	if _, err := os.Stat(upload.FilePath); os.IsNotExist(err) {
		http.Error(w, "File tidak ditemukan di server", http.StatusNotFound)
		return
	}

	// Set headers for file download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", upload.FileName))
	w.Header().Set("Content-Type", upload.MimeType)
	w.Header().Set("Content-Length", strconv.FormatInt(upload.FileSize, 10))

	// Serve the file
	http.ServeFile(w, r, upload.FilePath)
}

// GetUploadsByDateRange retrieves uploads within a date range
func GetUploadsByDateRange(w http.ResponseWriter, r *http.Request) {
	tanggalMulai := r.URL.Query().Get("tanggal_mulai")
	tanggalSelesai := r.URL.Query().Get("tanggal_selesai")

	if tanggalMulai == "" || tanggalSelesai == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Parameter tanggal_mulai dan tanggal_selesai wajib diisi",
		})
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", tanggalMulai)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Format tanggal_mulai tidak valid",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", tanggalSelesai)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Format tanggal_selesai tidak valid",
		})
		return
	}

	// Validate date range
	if startDate.After(endDate) {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Tanggal mulai tidak boleh lebih besar dari tanggal selesai",
		})
		return
	}

	var uploads []models.Upload
	if err := config.DB.
		Where("DATE(tanggal) BETWEEN DATE(?) AND DATE(?)", startDate, endDate).
		Order("tanggal desc, created_at desc").
		Find(&uploads).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data: " + err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: fmt.Sprintf("Data berhasil diambil untuk periode %s s/d %s", tanggalMulai, tanggalSelesai),
		Data:    uploads,
	})
}

// Helper function to generate random string
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}
