package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// ======== RESPONSE STRUCTS ========
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type MandorSummary struct {
	ID              uint             `json:"id"`
	TahunTanam      uint             `json:"tahunTanam"`
	Mandor          string           `json:"mandor"`
	Afdeling        string           `json:"afdeling"`
	TotalBasahLatex float64          `json:"totalBasahLatex"`
	TotalSheet      float64          `json:"totalSheet"`
	TotalBasahLump  float64          `json:"totalBasahLump"`
	TotalBrCr       float64          `json:"totalBrCr"`
	JumlahPenyadap  int              `json:"jumlahPenyadap"`
	DetailPenyadap  []PenyadapDetail `json:"detailPenyadap,omitempty"`
}

type PenyadapDetail struct {
	ID              uint    `json:"id"`
	NamaPenyadap    string  `json:"namaPenyadap"`
	NIK             string  `json:"nik"`
	TotalBasahLatex float64 `json:"totalBasahLatex"`
	TotalSheet      float64 `json:"totalSheet"`
	TotalBasahLump  float64 `json:"totalBasahLump"`
	TotalBrCr       float64 `json:"totalBrCr"`
	JumlahHariKerja int     `json:"jumlahHariKerja"`
	Mandor          string  `json:"mandor,omitempty"`
	Afdeling        string  `json:"afdeling,omitempty"`
}

type ReportingResponse struct {
	Success    bool            `json:"success"`
	Message    string          `json:"message"`
	Data       []MandorSummary `json:"data"`
	FilterInfo FilterInfo      `json:"filterInfo"`
}

type FilterInfo struct {
	Tanggal     string `json:"tanggal,omitempty"`
	TotalRecord int    `json:"totalRecord"`
	Periode     string `json:"periode"`
}

type BakuPageData struct {
	Title        string
	MandorList   []models.BakuMandor
	PenyadapList []models.BakuPenyadap
}

// ======== UTILITY FUNCTIONS ========
func respondJSON(w http.ResponseWriter, status int, payload APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// Template functions
var templateFuncs = template.FuncMap{
	"add": func(a, b int) int { return a + b },
}

// ======== CRUD OPERATIONS ========

// GetAllBakuPenyadap - Get all penyadap records
func GetAllBakuPenyadap(w http.ResponseWriter, r *http.Request) {
	var penyadap []models.BakuPenyadap

	if err := config.DB.Preload("Mandor").Preload("Penyadap").Order("created_at desc").Find(&penyadap).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data: " + err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data berhasil diambil",
		Data:    penyadap,
	})
}

// GetAllBakuMandor - Get all mandor records
func GetAllBakuMandor(w http.ResponseWriter, r *http.Request) {
	var mandor []models.BakuMandor
	if err := config.DB.Order("created_at desc").Find(&mandor).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data: " + err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data berhasil diambil",
		Data:    mandor,
	})
}

// GetBakuPenyadapByID - Get penyadap by ID
func GetBakuPenyadapByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var penyadap models.BakuPenyadap
	if err := config.DB.Preload("Mandor").Preload("Penyadap").First(&penyadap, id).Error; err != nil {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Data tidak ditemukan",
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data berhasil ditemukan",
		Data:    penyadap,
	})
}

// CreateBakuPenyadap - Create new penyadap record
func CreateBakuPenyadap(w http.ResponseWriter, r *http.Request) {
	var penyadap models.BakuPenyadap
	if err := json.NewDecoder(r.Body).Decode(&penyadap); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Format JSON tidak valid: " + err.Error(),
		})
		return
	}

	// Validasi
	if penyadap.IdBakuMandor == 0 || penyadap.IdPenyadap == 0 {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "ID mandor dan ID penyadap wajib diisi",
		})
		return
	}
	if penyadap.Tanggal.IsZero() {
		penyadap.Tanggal = time.Now()
	}

	// Simpan data penyadap
	if err := config.DB.Create(&penyadap).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal menyimpan data penyadap: " + err.Error(),
		})
		return
	}

	// Update detail harian
	updateBakuDetail(penyadap, "create", nil)

	respondJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Message: "Data penyadap berhasil ditambahkan",
		Data:    penyadap,
	})
}

// UpdateBakuPenyadap - Update penyadap record
func UpdateBakuPenyadap(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var existing models.BakuPenyadap
	if err := config.DB.First(&existing, id).Error; err != nil {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Data penyadap tidak ditemukan",
		})
		return
	}

	var updates models.BakuPenyadap
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Format JSON tidak valid: " + err.Error(),
		})
		return
	}

	// Simpan selisih untuk update detail
	oldCopy := existing

	if err := config.DB.Model(&existing).Updates(updates).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal update penyadap: " + err.Error(),
		})
		return
	}

	// Update detail
	updateBakuDetail(existing, "update", &oldCopy)

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data penyadap berhasil diperbarui",
	})
}

// DeleteBakuPenyadap - Delete penyadap record
func DeleteBakuPenyadap(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var penyadap models.BakuPenyadap
	if err := config.DB.First(&penyadap, id).Error; err != nil {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Data penyadap tidak ditemukan",
		})
		return
	}

	if err := config.DB.Delete(&penyadap).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal menghapus data penyadap: " + err.Error(),
		})
		return
	}

	updateBakuDetail(penyadap, "delete", nil)

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data penyadap berhasil dihapus",
	})
}

// ======== DETAIL OPERATIONS ========

// GetAllBakuDetail - Get all detail records
func GetAllBakuDetail(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("DEBUG: GetAllBakuDetail called - Method: %s, URL: %s\n", r.Method, r.URL.Path)

	var details []models.BakuDetail
	if err := config.DB.Order("tanggal desc").Find(&details).Error; err != nil {
		fmt.Printf("DEBUG: Database error: %v\n", err)
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil detail baku: " + err.Error(),
		})
		return
	}

	fmt.Printf("DEBUG: Found %d records\n", len(details))

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data detail berhasil diambil",
		Data:    details,
	})
}

// GetBakuDetailByDate - Get detail by specific date
func GetBakuDetailByDate(w http.ResponseWriter, r *http.Request) {
	tanggalStr := mux.Vars(r)["tanggal"]

	// Validasi format tanggal
	if tanggalStr == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Parameter tanggal wajib diisi",
		})
		return
	}

	// Parse tanggal dari string (format: YYYY-MM-DD)
	tanggal, err := time.Parse("2006-01-02", tanggalStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Format tanggal tidak valid. Gunakan format YYYY-MM-DD (contoh: 2024-09-10)",
		})
		return
	}

	var detail models.BakuDetail

	// Query dengan range tanggal untuk menangani timestamp
	startOfDay := tanggal
	endOfDay := tanggal.Add(24 * time.Hour)

	if err := config.DB.Where("tanggal >= ? AND tanggal < ?", startOfDay, endOfDay).First(&detail).Error; err != nil {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Detail untuk tanggal " + tanggalStr + " tidak ditemukan",
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Detail berhasil ditemukan",
		Data:    detail,
	})
}

// updateBakuDetail - Helper function to update daily summary
func updateBakuDetail(entry models.BakuPenyadap, action string, oldEntry *models.BakuPenyadap) {
	var detail models.BakuDetail
	err := config.DB.Where("tanggal = ?", entry.Tanggal).First(&detail).Error
	if err != nil {
		// kalau belum ada & action create → buat baru
		if action == "create" {
			detail = models.BakuDetail{Tanggal: entry.Tanggal}
		} else {
			return
		}
	}

	switch action {
	case "create":
		detail.JumlahKebunBasahLatek += entry.BasahLatex
		detail.JumlahSheet += entry.Sheet
		detail.JumlahKebunBasahLump += entry.BasahLump
		detail.JumlahBrCr += entry.BrCr
	case "update":
		if oldEntry != nil {
			detail.JumlahKebunBasahLatek += entry.BasahLatex - oldEntry.BasahLatex
			detail.JumlahSheet += entry.Sheet - oldEntry.Sheet
			detail.JumlahKebunBasahLump += entry.BasahLump - oldEntry.BasahLump
			detail.JumlahBrCr += entry.BrCr - oldEntry.BrCr
		}
	case "delete":
		detail.JumlahKebunBasahLatek -= entry.BasahLatex
		detail.JumlahSheet -= entry.Sheet
		detail.JumlahKebunBasahLump -= entry.BasahLump
		detail.JumlahBrCr -= entry.BrCr
	}

	config.DB.Save(&detail)
}

// ======== PAGE RENDERING ========

// ServeBakuPage - Serve the main baku page
func ServeBakuPage(w http.ResponseWriter, r *http.Request) {
	var mandor []models.BakuMandor
	var penyadap []models.BakuPenyadap

	if err := config.DB.Order("created_at desc").Find(&mandor).Error; err != nil {
		http.Error(w, "Gagal mengambil data mandor: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := config.DB.Preload("Mandor").Preload("Penyadap").Order("created_at desc").Find(&penyadap).Error; err != nil {
		http.Error(w, "Gagal mengambil data penyadap: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := BakuPageData{
		Title:        "Data Mandor & Penyadap",
		MandorList:   mandor,
		PenyadapList: penyadap,
	}

	tmpl, err := template.New("baku.html").Funcs(templateFuncs).ParseFiles("templates/html/baku.html")
	if err != nil {
		http.Error(w, "Gagal parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Gagal render template: "+err.Error(), http.StatusInternalServerError)
	}
}

// ======== REPORTING FUNCTIONS ========

// GetMandorSummaryAll - Get summary of all mandors for all time
func GetMandorSummaryAll(w http.ResponseWriter, r *http.Request) {
	summaries, err := getMandorSummaries("")
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data summary: " + err.Error(),
		})
		return
	}

	response := ReportingResponse{
		Success: true,
		Message: "Data summary mandor berhasil diambil",
		Data:    summaries,
		FilterInfo: FilterInfo{
			TotalRecord: len(summaries),
			Periode:     "Semua waktu",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetMandorSummaryByDate - Get mandor summary for specific date
func GetMandorSummaryByDate(w http.ResponseWriter, r *http.Request) {
	tanggalStr := mux.Vars(r)["tanggal"]

	// Validasi format tanggal
	if _, err := time.Parse("2006-01-02", tanggalStr); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Format tanggal tidak valid. Gunakan format YYYY-MM-DD",
		})
		return
	}

	summaries, err := getMandorSummaries(tanggalStr)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data summary: " + err.Error(),
		})
		return
	}

	response := ReportingResponse{
		Success: true,
		Message: "Data summary mandor untuk tanggal " + tanggalStr + " berhasil diambil",
		Data:    summaries,
		FilterInfo: FilterInfo{
			Tanggal:     tanggalStr,
			TotalRecord: len(summaries),
			Periode:     "Tanggal: " + tanggalStr,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetPenyadapDetailAll - Get details of all penyadap for all time
func GetPenyadapDetailAll(w http.ResponseWriter, r *http.Request) {
	details, err := getPenyadapDetails("")
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil detail penyadap: " + err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Detail semua penyadap berhasil diambil",
		Data:    details,
	})
}

// GetPenyadapDetailByDate - Get penyadap details for specific date
func GetPenyadapDetailByDate(w http.ResponseWriter, r *http.Request) {
	tanggalStr := mux.Vars(r)["tanggal"]

	if _, err := time.Parse("2006-01-02", tanggalStr); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Format tanggal tidak valid. Gunakan format YYYY-MM-DD",
		})
		return
	}

	details, err := getPenyadapDetails(tanggalStr)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil detail penyadap: " + err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Detail penyadap untuk tanggal " + tanggalStr + " berhasil diambil",
		Data:    details,
	})
}

// ======== SEARCH FUNCTIONS ========

// SearchMandorByName - Search mandor by name with optional date filter
func SearchMandorByName(w http.ResponseWriter, r *http.Request) {
	// Query parameters
	namaMandor := r.URL.Query().Get("nama")
	tanggal := r.URL.Query().Get("tanggal")

	if namaMandor == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Parameter 'nama' wajib diisi",
		})
		return
	}

	// Validasi tanggal jika ada
	if tanggal != "" {
		if _, err := time.Parse("2006-01-02", tanggal); err != nil {
			respondJSON(w, http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "Format tanggal tidak valid. Gunakan format YYYY-MM-DD",
			})
			return
		}
	}

	summaries, err := searchMandorSummaries(namaMandor, tanggal)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mencari data mandor: " + err.Error(),
		})
		return
	}

	periode := "Semua waktu"
	if tanggal != "" {
		periode = "Tanggal: " + tanggal
	}

	response := ReportingResponse{
		Success: true,
		Message: "Hasil pencarian mandor '" + namaMandor + "'",
		Data:    summaries,
		FilterInfo: FilterInfo{
			Tanggal:     tanggal,
			TotalRecord: len(summaries),
			Periode:     periode,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SearchPenyadapByName - Search penyadap by name with optional date filter
func SearchPenyadapByName(w http.ResponseWriter, r *http.Request) {
	namaPenyadap := r.URL.Query().Get("nama")
	tanggal := r.URL.Query().Get("tanggal")

	if namaPenyadap == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Parameter 'nama' wajib diisi",
		})
		return
	}

	// Validasi tanggal jika ada
	if tanggal != "" {
		if _, err := time.Parse("2006-01-02", tanggal); err != nil {
			respondJSON(w, http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "Format tanggal tidak valid. Gunakan format YYYY-MM-DD",
			})
			return
		}
	}

	details, err := searchPenyadapDetails(namaPenyadap, tanggal)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mencari detail penyadap: " + err.Error(),
		})
		return
	}

	periode := "Semua waktu"
	if tanggal != "" {
		periode = "Tanggal: " + tanggal
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Hasil pencarian penyadap '" + namaPenyadap + "' - " + periode,
		Data:    details,
	})
}

// GetMandorWithPenyadapDetail - Get mandor details with all penyadap
func GetMandorWithPenyadapDetail(w http.ResponseWriter, r *http.Request) {
	namaMandor := r.URL.Query().Get("nama")
	tanggal := r.URL.Query().Get("tanggal")

	if namaMandor == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Parameter 'nama' wajib diisi",
		})
		return
	}

	if tanggal != "" {
		if _, err := time.Parse("2006-01-02", tanggal); err != nil {
			respondJSON(w, http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "Format tanggal tidak valid. Gunakan format YYYY-MM-DD",
			})
			return
		}
	}

	summaries, err := searchMandorWithDetails(namaMandor, tanggal)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil detail mandor: " + err.Error(),
		})
		return
	}

	periode := "Semua waktu"
	if tanggal != "" {
		periode = "Tanggal: " + tanggal
	}

	response := ReportingResponse{
		Success: true,
		Message: "Detail mandor '" + namaMandor + "' beserta penyadapnya",
		Data:    summaries,
		FilterInfo: FilterInfo{
			Tanggal:     tanggal,
			TotalRecord: len(summaries),
			Periode:     periode,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SearchAll - Global search with various filters
func SearchAll(w http.ResponseWriter, r *http.Request) {
	searchType := r.URL.Query().Get("type") // "mandor" atau "penyadap"
	nama := r.URL.Query().Get("nama")
	tanggal := r.URL.Query().Get("tanggal")
	afdeling := r.URL.Query().Get("afdeling")
	tahunTanam := r.URL.Query().Get("tahun")

	if nama == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Parameter 'nama' wajib diisi",
		})
		return
	}

	// Validasi tanggal jika ada
	if tanggal != "" {
		if _, err := time.Parse("2006-01-02", tanggal); err != nil {
			respondJSON(w, http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "Format tanggal tidak valid. Gunakan format YYYY-MM-DD",
			})
			return
		}
	}

	var result interface{}
	var err error

	switch searchType {
	case "mandor":
		result, err = advancedSearchMandor(nama, tanggal, afdeling, tahunTanam)
	case "penyadap":
		result, err = advancedSearchPenyadap(nama, tanggal, afdeling)
	default:
		// Auto detect berdasarkan hasil pencarian
		mandorResult, _ := advancedSearchMandor(nama, tanggal, afdeling, tahunTanam)
		penyadapResult, _ := advancedSearchPenyadap(nama, tanggal, afdeling)

		result = map[string]interface{}{
			"mandor":   mandorResult,
			"penyadap": penyadapResult,
		}
	}

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal melakukan pencarian: " + err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Pencarian berhasil",
		Data:    result,
	})
}

// ======== HELPER FUNCTIONS ========

func getMandorSummaries(tanggal string) ([]MandorSummary, error) {
	var mandors []models.BakuMandor
	if err := config.DB.Find(&mandors).Error; err != nil {
		return nil, err
	}

	var summaries []MandorSummary

	for _, mandor := range mandors {
		summary := MandorSummary{
			ID:         mandor.ID,
			TahunTanam: mandor.TahunTanam,
			Mandor:     mandor.Mandor,
			Afdeling:   mandor.Afdeling,
		}

		// Query untuk mendapat total dari semua penyadap mandor ini
		query := config.DB.Model(&models.BakuPenyadap{}).
			Where("id_baku_mandor = ?", mandor.ID)

		if tanggal != "" {
			query = query.Where("tanggal LIKE ?", tanggal+"%")
		}

		var results []struct {
			TotalBasahLatex float64
			TotalSheet      float64
			TotalBasahLump  float64
			TotalBrCr       float64
			JumlahPenyadap  int64
		}

		err := query.Select(`
			COALESCE(SUM(basah_latex), 0) as total_basah_latex,
			COALESCE(SUM(sheet), 0) as total_sheet,
			COALESCE(SUM(basah_lump), 0) as total_basah_lump,
			COALESCE(SUM(br_cr), 0) as total_br_cr,
			COUNT(DISTINCT id_penyadap) as jumlah_penyadap
		`).Find(&results).Error

		if err != nil {
			return nil, err
		}

		if len(results) > 0 {
			summary.TotalBasahLatex = results[0].TotalBasahLatex
			summary.TotalSheet = results[0].TotalSheet
			summary.TotalBasahLump = results[0].TotalBasahLump
			summary.TotalBrCr = results[0].TotalBrCr
			summary.JumlahPenyadap = int(results[0].JumlahPenyadap)
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func getPenyadapDetails(tanggal string) ([]PenyadapDetail, error) {
	query := `
		SELECT 
			p.id,
			p.nama_penyadap,
			p.nik,
			
			COALESCE(SUM(bp.basah_latex), 0) as total_basah_latex,
			COALESCE(SUM(bp.sheet), 0) as total_sheet,
			COALESCE(SUM(bp.basah_lump), 0) as total_basah_lump,
			COALESCE(SUM(bp.br_cr), 0) as total_br_cr,
			COUNT(bp.id) as jumlah_hari_kerja
		FROM penyadaps p
		LEFT JOIN baku_penyadaps bp ON p.id = bp.id_penyadap
		WHERE bp.deleted_at IS NULL
	`

	args := []interface{}{}
	if tanggal != "" {
		query += " AND bp.tanggal LIKE ?"
		args = append(args, tanggal+"%")
	}

	query += " GROUP BY p.id, p.nama_penyadap, p.nik ORDER BY p.nama_penyadap"

	var details []PenyadapDetail
	if err := config.DB.Raw(query, args...).Scan(&details).Error; err != nil {
		return nil, err
	}

	return details, nil
}

func searchMandorSummaries(namaMandor, tanggal string) ([]MandorSummary, error) {
	var mandors []models.BakuMandor
	query := config.DB.Where("mandor LIKE ?", "%"+namaMandor+"%")

	if err := query.Find(&mandors).Error; err != nil {
		return nil, err
	}

	var summaries []MandorSummary

	for _, mandor := range mandors {
		summary := MandorSummary{
			ID:         mandor.ID,
			TahunTanam: mandor.TahunTanam,
			Mandor:     mandor.Mandor,
			Afdeling:   mandor.Afdeling,
		}

		// Query untuk mendapat total dari semua penyadap mandor ini
		bakuQuery := config.DB.Model(&models.BakuPenyadap{}).
			Where("id_baku_mandor = ?", mandor.ID)

		if tanggal != "" {
			bakuQuery = bakuQuery.Where("tanggal LIKE ?", tanggal+"%")
		}

		var results []struct {
			TotalBasahLatex float64
			TotalSheet      float64
			TotalBasahLump  float64
			TotalBrCr       float64
			JumlahPenyadap  int64
		}

		err := bakuQuery.Select(`
			COALESCE(SUM(basah_latex), 0) as total_basah_latex,
			COALESCE(SUM(sheet), 0) as total_sheet,
			COALESCE(SUM(basah_lump), 0) as total_basah_lump,
			COALESCE(SUM(br_cr), 0) as total_br_cr,
			COUNT(DISTINCT id_penyadap) as jumlah_penyadap
		`).Find(&results).Error

		if err != nil {
			return nil, err
		}

		if len(results) > 0 {
			summary.TotalBasahLatex = results[0].TotalBasahLatex
			summary.TotalSheet = results[0].TotalSheet
			summary.TotalBasahLump = results[0].TotalBasahLump
			summary.TotalBrCr = results[0].TotalBrCr
			summary.JumlahPenyadap = int(results[0].JumlahPenyadap)
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func searchPenyadapDetails(namaPenyadap, tanggal string) ([]PenyadapDetail, error) {
	query := `
		SELECT 
			p.id,
			p.nama_penyadap,
			p.nik,
			bm.mandor,
			bm.afdeling,
			COALESCE(SUM(bp.basah_latex), 0) as total_basah_latex,
			COALESCE(SUM(bp.sheet), 0) as total_sheet,
			COALESCE(SUM(bp.basah_lump), 0) as total_basah_lump,
			COALESCE(SUM(bp.br_cr), 0) as total_br_cr,
			COUNT(bp.id) as jumlah_hari_kerja
		FROM penyadaps p
		LEFT JOIN baku_penyadaps bp ON p.id = bp.id_penyadap
		LEFT JOIN baku_mandors bm ON bp.id_baku_mandor = bm.id
		WHERE bp.deleted_at IS NULL AND p.nama_penyadap LIKE ?
	`

	args := []interface{}{"%" + namaPenyadap + "%"}
	if tanggal != "" {
		query += " AND bp.tanggal LIKE ?"
		args = append(args, tanggal+"%")
	}

	query += " GROUP BY p.id, p.nama_penyadap, p.nik, bm.mandor, bm.afdeling ORDER BY p.nama_penyadap"

	var details []PenyadapDetail
	if err := config.DB.Raw(query, args...).Scan(&details).Error; err != nil {
		return nil, err
	}

	return details, nil
}

func searchMandorWithDetails(namaMandor, tanggal string) ([]MandorSummary, error) {
	var mandors []models.BakuMandor
	query := config.DB.Where("mandor LIKE ?", "%"+namaMandor+"%")

	if err := query.Find(&mandors).Error; err != nil {
		return nil, err
	}

	var summaries []MandorSummary

	for _, mandor := range mandors {
		summary := MandorSummary{
			ID:         mandor.ID,
			TahunTanam: mandor.TahunTanam,
			Mandor:     mandor.Mandor,
			Afdeling:   mandor.Afdeling,
		}

		// Query untuk mendapat total dari semua penyadap mandor ini
		bakuQuery := config.DB.Model(&models.BakuPenyadap{}).
			Where("id_baku_mandor = ?", mandor.ID)

		if tanggal != "" {
			bakuQuery = bakuQuery.Where("tanggal LIKE ?", tanggal+"%")
		}

		var results []struct {
			TotalBasahLatex float64
			TotalSheet      float64
			TotalBasahLump  float64
			TotalBrCr       float64
			JumlahPenyadap  int64
		}

		err := bakuQuery.Select(`
			COALESCE(SUM(basah_latex), 0) as total_basah_latex,
			COALESCE(SUM(sheet), 0) as total_sheet,
			COALESCE(SUM(basah_lump), 0) as total_basah_lump,
			COALESCE(SUM(br_cr), 0) as total_br_cr,
			COUNT(DISTINCT id_penyadap) as jumlah_penyadap
		`).Find(&results).Error

		if err != nil {
			return nil, err
		}

		if len(results) > 0 {
			summary.TotalBasahLatex = results[0].TotalBasahLatex
			summary.TotalSheet = results[0].TotalSheet
			summary.TotalBasahLump = results[0].TotalBasahLump
			summary.TotalBrCr = results[0].TotalBrCr
			summary.JumlahPenyadap = int(results[0].JumlahPenyadap)
		}

		// Mendapatkan detail setiap penyadap untuk mandor ini
		penyadapQuery := `
			SELECT 
				p.id,
				p.nama_penyadap,
				p.nik,
				COALESCE(SUM(bp.basah_latex), 0) as total_basah_latex,
				COALESCE(SUM(bp.sheet), 0) as total_sheet,
				COALESCE(SUM(bp.basah_lump), 0) as total_basah_lump,
				COALESCE(SUM(bp.br_cr), 0) as total_br_cr,
				COUNT(bp.id) as jumlah_hari_kerja
			FROM penyadaps p
			LEFT JOIN baku_penyadaps bp ON p.id = bp.id_penyadap
			WHERE bp.deleted_at IS NULL AND bp.id_baku_mandor = ?
		`

		penyadapArgs := []interface{}{mandor.ID}
		if tanggal != "" {
			penyadapQuery += " AND bp.tanggal LIKE ?"
			penyadapArgs = append(penyadapArgs, tanggal+"%")
		}

		penyadapQuery += " GROUP BY p.id, p.nama_penyadap, p.nik ORDER BY p.nama_penyadap"

		var penyadapDetails []PenyadapDetail
		if err := config.DB.Raw(penyadapQuery, penyadapArgs...).Scan(&penyadapDetails).Error; err != nil {
			return nil, err
		}

		summary.DetailPenyadap = penyadapDetails
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func advancedSearchMandor(nama, tanggal, afdeling, tahunTanam string) ([]MandorSummary, error) {
	query := config.DB.Where("mandor LIKE ?", "%"+nama+"%")

	if afdeling != "" {
		query = query.Where("afdeling LIKE ?", "%"+afdeling+"%")
	}

	if tahunTanam != "" {
		query = query.Where("tahun_tanam = ?", tahunTanam)
	}

	var mandors []models.BakuMandor
	if err := query.Find(&mandors).Error; err != nil {
		return nil, err
	}

	var summaries []MandorSummary
	for _, mandor := range mandors {
		summary := MandorSummary{
			ID:         mandor.ID,
			TahunTanam: mandor.TahunTanam,
			Mandor:     mandor.Mandor,
			Afdeling:   mandor.Afdeling,
		}

		bakuQuery := config.DB.Model(&models.BakuPenyadap{}).
			Where("id_baku_mandor = ?", mandor.ID)

		if tanggal != "" {
			bakuQuery = bakuQuery.Where("tanggal LIKE ?", tanggal+"%")
		}

		var results []struct {
			TotalBasahLatex float64
			TotalSheet      float64
			TotalBasahLump  float64
			TotalBrCr       float64
			JumlahPenyadap  int64
		}

		err := bakuQuery.Select(`
			COALESCE(SUM(basah_latex), 0) as total_basah_latex,
			COALESCE(SUM(sheet), 0) as total_sheet,
			COALESCE(SUM(basah_lump), 0) as total_basah_lump,
			COALESCE(SUM(br_cr), 0) as total_br_cr,
			COUNT(DISTINCT id_penyadap) as jumlah_penyadap
		`).Find(&results).Error

		if err != nil {
			return nil, err
		}

		if len(results) > 0 {
			summary.TotalBasahLatex = results[0].TotalBasahLatex
			summary.TotalSheet = results[0].TotalSheet
			summary.TotalBasahLump = results[0].TotalBasahLump
			summary.TotalBrCr = results[0].TotalBrCr
			summary.JumlahPenyadap = int(results[0].JumlahPenyadap)
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func advancedSearchPenyadap(nama, tanggal, afdeling string) ([]PenyadapDetail, error) {
	query := `
		SELECT 
			p.id,
			p.nama_penyadap,
			p.nik,
			bm.mandor,
			bm.afdeling,
			COALESCE(SUM(bp.basah_latex), 0) as total_basah_latex,
			COALESCE(SUM(bp.sheet), 0) as total_sheet,
			COALESCE(SUM(bp.basah_lump), 0) as total_basah_lump,
			COALESCE(SUM(bp.br_cr), 0) as total_br_cr,
			COUNT(bp.id) as jumlah_hari_kerja
		FROM penyadaps p
		LEFT JOIN baku_penyadaps bp ON p.id = bp.id_penyadap
		LEFT JOIN baku_mandors bm ON bp.id_baku_mandor = bm.id
		WHERE bp.deleted_at IS NULL AND p.nama_penyadap LIKE ?
	`

	args := []interface{}{"%" + nama + "%"}

	if tanggal != "" {
		query += " AND bp.tanggal LIKE ?"
		args = append(args, tanggal+"%")
	}

	if afdeling != "" {
		query += " AND bm.afdeling LIKE ?"
		args = append(args, "%"+afdeling+"%")
	}

	query += " GROUP BY p.id, p.nama_penyadap, p.nik, bm.mandor, bm.afdeling ORDER BY p.nama_penyadap"

	var details []PenyadapDetail
	if err := config.DB.Raw(query, args...).Scan(&details).Error; err != nil {
		return nil, err
	}

	return details, nil
}
