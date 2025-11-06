package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

// CreateUpload handles file upload and date submission
func CreateUpload(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with max memory
	if err := r.ParseMultipartForm(MaxFileSize); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "File terlalu besar atau format tidak valid",
		})
		return
	}
	//Get afdeling from form
	afdeling := r.FormValue("afdeling")
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

	// Ensure upload directory exists
	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal membuat direktori upload",
		})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
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
	defer dst.Close()

	// Copy uploaded file to destination
	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(uploadPath) // Cleanup on error
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal menyimpan file",
		})
		return
	}

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
		os.Remove(uploadPath) // Cleanup file if database save fails
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal menyimpan data ke database: " + err.Error(),
		})
		return
	}
	//convert ke csv
	// convert ke csv - kirim nama file asli
	if err := excelToCSV(uploadPath, "csv", tanggal, afdeling, header.Filename); err != nil {
		log.Fatal(err)
	}

	// Return success response with upload data
	responseData := map[string]interface{}{
		"id":       upload.ID,
		"tanggal":  upload.Tanggal.Format("2006-01-02"),
		"fileName": upload.FileName,
		"fileSize": upload.FileSize,
		"filePath": "/uploads/" + newFileName,
	}

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

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "File berhasil diupload",
		Data:    responseData,
	})
}

// GetAllUploads retrieves all upload records
func GetAllUploads(w http.ResponseWriter, r *http.Request) {
	var uploads []models.Upload

	// Parse query parameters for filtering
	tanggalStr := r.URL.Query().Get("tanggal")

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

	if err := query.Find(&uploads).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data: " + err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data berhasil diambil",
		Data:    uploads,
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

	// Delete file from filesystem
	if err := os.Remove(upload.FilePath); err != nil {
		// Log error but continue with database deletion
		fmt.Printf("Warning: Failed to delete file %s: %v\n", upload.FilePath, err)
	}

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
