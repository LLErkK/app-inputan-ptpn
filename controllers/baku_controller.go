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

	// Update detail harian berdasarkan tanggal dan mandor
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

	// Simpan copy untuk update detail
	oldCopy := existing

	if err := config.DB.Model(&existing).Updates(updates).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal update penyadap: " + err.Error(),
		})
		return
	}

	// Update detail berdasarkan tanggal dan mandor
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

	// Update detail berdasarkan tanggal dan mandor
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
	if err := config.DB.Order("tanggal desc, mandor asc").Find(&details).Error; err != nil {
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

	var details []models.BakuDetail

	// Query untuk mendapatkan semua detail untuk tanggal tersebut (semua mandor)
	if err := config.DB.Where("DATE(tanggal) = DATE(?)", tanggal).Order("mandor asc").Find(&details).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil detail: " + err.Error(),
		})
		return
	}

	if len(details) == 0 {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Detail untuk tanggal " + tanggalStr + " tidak ditemukan",
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Detail berhasil ditemukan",
		Data:    details,
	})
}

// GetBakuDetailByDateAndMandor - Get detail by date and mandor
func GetBakuDetailByDateAndMandor(w http.ResponseWriter, r *http.Request) {
	tanggalStr := r.URL.Query().Get("tanggal")
	mandor := r.URL.Query().Get("mandor")

	if tanggalStr == "" || mandor == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Parameter tanggal dan mandor wajib diisi",
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

	var detail models.BakuDetail
	if err := config.DB.Where("DATE(tanggal) = DATE(?) AND mandor = ?", tanggal, mandor).First(&detail).Error; err != nil {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: fmt.Sprintf("Detail untuk tanggal %s dan mandor %s tidak ditemukan", tanggalStr, mandor),
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
// Mencari berdasarkan kombinasi tanggal DAN mandor
func updateBakuDetail(entry models.BakuPenyadap, action string, oldEntry *models.BakuPenyadap) {
	// Ambil tanggal tanpa jam untuk konsistensi
	targetDate := entry.Tanggal.Truncate(24 * time.Hour)

	// Ambil data mandor
	var mandor models.BakuMandor
	if err := config.DB.First(&mandor, entry.IdBakuMandor).Error; err != nil {
		fmt.Printf("ERROR: Tidak dapat menemukan mandor dengan ID %d: %v\n", entry.IdBakuMandor, err)
		return
	}

	var detail models.BakuDetail
	err := config.DB.Where("DATE(tanggal) = DATE(?) AND mandor = ?", targetDate, mandor.Mandor).First(&detail).Error

	if err != nil {
		// Jika belum ada detail untuk kombinasi ini → buat baru
		if action == "create" {
			detail = models.BakuDetail{
				Tanggal:  targetDate,
				Mandor:   mandor.Mandor,
				Afdeling: mandor.Afdeling,
				// semua field default 0
			}
		} else {
			fmt.Printf("WARNING: Tidak ada BakuDetail untuk tanggal %s mandor %s pada action %s\n",
				targetDate.Format("2006-01-02"), mandor.Mandor, action)
			return
		}
	}

	// ================== Update nilai berdasarkan action ==================
	switch action {
	case "create":
		fmt.Printf("CREATE: Menambah data untuk %s mandor %s\n", targetDate.Format("2006-01-02"), mandor.Mandor)
		detail.JumlahPabrikBasahLatek += entry.BasahLatex
		detail.JumlahPabrikBasahLump += entry.BasahLump
		detail.JumlahSheet += entry.Sheet
		detail.JumlahBrCr += entry.BrCr

	case "update":
		if oldEntry != nil {
			fmt.Printf("UPDATE: Mengupdate data untuk %s mandor %s\n", targetDate.Format("2006-01-02"), mandor.Mandor)

			// hitung selisih antara nilai baru dan lama
			deltaBasahLatex := entry.BasahLatex - oldEntry.BasahLatex
			deltaBasahLump := entry.BasahLump - oldEntry.BasahLump
			deltaSheet := entry.Sheet - oldEntry.Sheet
			deltaBrCr := entry.BrCr - oldEntry.BrCr

			detail.JumlahPabrikBasahLatek += deltaBasahLatex
			detail.JumlahPabrikBasahLump += deltaBasahLump
			detail.JumlahSheet += deltaSheet
			detail.JumlahBrCr += deltaBrCr

			fmt.Printf("  Delta - BasahLatex: %.2f, BasahLump: %.2f, Sheet: %.2f, BrCr: %.2f\n",
				deltaBasahLatex, deltaBasahLump, deltaSheet, deltaBrCr)
		}

	case "delete":
		fmt.Printf("DELETE: Mengurangi data untuk %s mandor %s\n", targetDate.Format("2006-01-02"), mandor.Mandor)
		detail.JumlahPabrikBasahLatek -= entry.BasahLatex
		detail.JumlahPabrikBasahLump -= entry.BasahLump
		detail.JumlahSheet -= entry.Sheet
		detail.JumlahBrCr -= entry.BrCr

		// jaga jangan sampai negatif
		if detail.JumlahPabrikBasahLatek < 0 {
			detail.JumlahPabrikBasahLatek = 0
		}
		if detail.JumlahPabrikBasahLump < 0 {
			detail.JumlahPabrikBasahLump = 0
		}
		if detail.JumlahSheet < 0 {
			detail.JumlahSheet = 0
		}
		if detail.JumlahBrCr < 0 {
			detail.JumlahBrCr = 0
		}
	}

	// ================== Hitung Selisih & Persentase ==================
	detail.SelisihBasahLatek = detail.JumlahPabrikBasahLatek - detail.JumlahKebunBasahLatek
	if detail.JumlahKebunBasahLatek > 0 {
		detail.PersentaseSelisihBasahLatek = (detail.SelisihBasahLatek / detail.JumlahKebunBasahLatek) * 100
	} else {
		detail.PersentaseSelisihBasahLatek = 0
	}

	detail.SelisihBasahLump = detail.JumlahPabrikBasahLump - detail.JumlahKebunBasahLump
	if detail.JumlahKebunBasahLump > 0 {
		detail.PersentaseSelisihBasahLump = (detail.SelisihBasahLump / detail.JumlahKebunBasahLump) * 100
	} else {
		detail.PersentaseSelisihBasahLump = 0
	}

	// ================== Simpan ke DB ==================
	if err := config.DB.Save(&detail).Error; err != nil {
		fmt.Printf("ERROR: Gagal menyimpan BakuDetail: %v\n", err)
	} else {
		fmt.Printf("SUCCESS: BakuDetail %s mandor %s terupdate\n",
			targetDate.Format("2006-01-02"), mandor.Mandor)
		fmt.Printf("  - Pabrik Latex: %.2f | Kebun Latex: %.2f\n", detail.JumlahPabrikBasahLatek, detail.JumlahKebunBasahLatek)
		fmt.Printf("  - Pabrik Lump: %.2f | Kebun Lump: %.2f\n", detail.JumlahPabrikBasahLump, detail.JumlahKebunBasahLump)
	}
}

// RecalculateBakuDetail - Fungsi untuk hitung ulang BakuDetail berdasarkan tanggal dan mandor
func RecalculateBakuDetail(tanggal time.Time, mandorID uint) error {
	targetDate := tanggal.Truncate(24 * time.Hour)

	// Ambil data mandor
	var mandor models.BakuMandor
	if err := config.DB.First(&mandor, mandorID).Error; err != nil {
		return fmt.Errorf("mandor tidak ditemukan: %v", err)
	}

	// Hitung ulang total dari semua BakuPenyadap untuk tanggal dan mandor tersebut
	var totals struct {
		TotalBasahLatex float64
		TotalSheet      float64
		TotalBasahLump  float64
		TotalBrCr       float64
	}

	err := config.DB.Model(&models.BakuPenyadap{}).
		Where("DATE(tanggal) = DATE(?) AND id_baku_mandor = ?", targetDate, mandorID).
		Select(`
			COALESCE(SUM(basah_latex), 0) as total_basah_latex,
			COALESCE(SUM(sheet), 0) as total_sheet,
			COALESCE(SUM(basah_lump), 0) as total_basah_lump,
			COALESCE(SUM(br_cr), 0) as total_br_cr
		`).
		Scan(&totals).Error

	if err != nil {
		return err
	}

	// Update atau create BakuDetail
	var detail models.BakuDetail
	err = config.DB.Where("DATE(tanggal) = DATE(?) AND mandor = ?", targetDate, mandor.Mandor).First(&detail).Error

	if err != nil {
		// Buat baru
		detail = models.BakuDetail{
			Tanggal:  targetDate,
			Mandor:   mandor.Mandor,
			Afdeling: mandor.Afdeling,
		}
	}

	// Set nilai yang dihitung ulang
	detail.JumlahKebunBasahLatek = totals.TotalBasahLatex
	detail.JumlahSheet = totals.TotalSheet
	detail.JumlahKebunBasahLump = totals.TotalBasahLump
	detail.JumlahBrCr = totals.TotalBrCr

	// Hitung ulang selisih
	detail.SelisihBasahLatek = detail.JumlahPabrikBasahLatek - detail.JumlahKebunBasahLatek
	if detail.JumlahKebunBasahLatek > 0 {
		detail.PersentaseSelisihBasahLatek = (detail.SelisihBasahLatek / detail.JumlahKebunBasahLatek) * 100
	} else {
		detail.PersentaseSelisihBasahLatek = 0
	}

	detail.SelisihBasahLump = detail.JumlahPabrikBasahLump - detail.JumlahKebunBasahLump
	if detail.JumlahKebunBasahLump > 0 {
		detail.PersentaseSelisihBasahLump = (detail.SelisihBasahLump / detail.JumlahKebunBasahLump) * 100
	} else {
		detail.PersentaseSelisihBasahLump = 0
	}

	return config.DB.Save(&detail).Error
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
			query = query.Where("DATE(tanggal) = DATE(?)", tanggal)
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

// getPenyadapDetails - DIPERBAIKI: Menggunakan BakuDetail yang sudah akurat
func getPenyadapDetails(tanggal string) ([]PenyadapDetail, error) {
	// Jika ada filter tanggal, ambil dari baku_detail berdasarkan tanggal dan mandor
	if tanggal != "" {
		targetDate, err := time.Parse("2006-01-02", tanggal)
		if err != nil {
			return nil, err
		}

		// Ambil semua BakuDetail untuk tanggal tersebut (semua mandor)
		var bakuDetails []models.BakuDetail
		err = config.DB.Where("DATE(tanggal) = DATE(?)", targetDate).Find(&bakuDetails).Error
		if err != nil {
			return nil, err
		}

		if len(bakuDetails) == 0 {
			// Tidak ada data untuk tanggal ini
			return []PenyadapDetail{}, nil
		}

		var allDetails []PenyadapDetail

		// Untuk setiap mandor yang ada di BakuDetail
		for _, bakuDetail := range bakuDetails {
			// Ambil daftar penyadap yang aktif untuk mandor dan tanggal tersebut
			query := `
				SELECT DISTINCT
					p.id,
					p.nama_penyadap,
					p.nik,
					bm.mandor,
					bm.afdeling,
					COUNT(bp.id) as jumlah_hari_kerja
				FROM penyadaps p
				INNER JOIN baku_penyadaps bp ON p.id = bp.id_penyadap
				INNER JOIN baku_mandors bm ON bp.id_baku_mandor = bm.id
				WHERE bp.deleted_at IS NULL 
				AND DATE(bp.tanggal) = DATE(?)
				AND bm.mandor = ?
				GROUP BY p.id, p.nama_penyadap, p.nik, bm.mandor, bm.afdeling
				ORDER BY p.nama_penyadap
			`

			var mandorDetails []PenyadapDetail
			err = config.DB.Raw(query, targetDate, bakuDetail.Mandor).Scan(&mandorDetails).Error
			if err != nil {
				return nil, err
			}

			// Distribusikan total dari BakuDetail ke semua penyadap mandor tersebut
			totalPenyadap := len(mandorDetails)
			if totalPenyadap > 0 {
				for i := range mandorDetails {
					// Bagi rata sesuai jumlah penyadap per mandor
					mandorDetails[i].TotalBasahLatex = bakuDetail.JumlahKebunBasahLatek / float64(totalPenyadap)
					mandorDetails[i].TotalSheet = bakuDetail.JumlahSheet / float64(totalPenyadap)
					mandorDetails[i].TotalBasahLump = bakuDetail.JumlahKebunBasahLump / float64(totalPenyadap)
					mandorDetails[i].TotalBrCr = bakuDetail.JumlahBrCr / float64(totalPenyadap)
					mandorDetails[i].JumlahHariKerja = 1 // Untuk tanggal spesifik, hari kerja = 1
				}
			}

			allDetails = append(allDetails, mandorDetails...)
		}

		return allDetails, nil
	}

	// Untuk semua waktu (tanpa filter tanggal), menggunakan agregasi dari baku_penyadaps
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
		GROUP BY p.id, p.nama_penyadap, p.nik 
		ORDER BY p.nama_penyadap
	`

	var details []PenyadapDetail
	if err := config.DB.Raw(query).Scan(&details).Error; err != nil {
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
			bakuQuery = bakuQuery.Where("DATE(tanggal) = DATE(?)", tanggal)
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
		query += " AND DATE(bp.tanggal) = DATE(?)"
		args = append(args, tanggal)
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
			bakuQuery = bakuQuery.Where("DATE(tanggal) = DATE(?)", tanggal)
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
			penyadapQuery += " AND DATE(bp.tanggal) = DATE(?)"
			penyadapArgs = append(penyadapArgs, tanggal)
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

		// Hitung total penyadap + produksi
		bakuQuery := config.DB.Model(&models.BakuPenyadap{}).
			Where("id_baku_mandor = ?", mandor.ID)

		if tanggal != "" {
			bakuQuery = bakuQuery.Where("DATE(tanggal) = DATE(?)", tanggal)
		}

		var result struct {
			TotalBasahLatex float64
			TotalSheet      float64
			TotalBasahLump  float64
			TotalBrCr       float64
			JumlahPenyadap  int64
		}

		if err := bakuQuery.Select(`
			COALESCE(SUM(basah_latex), 0) as total_basah_latex,
			COALESCE(SUM(sheet), 0) as total_sheet,
			COALESCE(SUM(basah_lump), 0) as total_basah_lump,
			COALESCE(SUM(br_cr), 0) as total_br_cr,
			COUNT(DISTINCT id_penyadap) as jumlah_penyadap
		`).Scan(&result).Error; err != nil {
			return nil, err
		}

		summary.TotalBasahLatex = result.TotalBasahLatex
		summary.TotalSheet = result.TotalSheet
		summary.TotalBasahLump = result.TotalBasahLump
		summary.TotalBrCr = result.TotalBrCr
		summary.JumlahPenyadap = int(result.JumlahPenyadap)

		// Ambil detail per penyadap
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

		args := []interface{}{mandor.ID}
		if tanggal != "" {
			penyadapQuery += " AND DATE(bp.tanggal) = DATE(?)"
			args = append(args, tanggal)
		}

		penyadapQuery += " GROUP BY p.id, p.nama_penyadap, p.nik ORDER BY p.nama_penyadap"

		var penyadapDetails []PenyadapDetail
		if err := config.DB.Raw(penyadapQuery, args...).Scan(&penyadapDetails).Error; err != nil {
			return nil, err
		}

		// Tambahkan mandor & afdeling ke detail
		for i := range penyadapDetails {
			penyadapDetails[i].Mandor = mandor.Mandor
			penyadapDetails[i].Afdeling = mandor.Afdeling
		}

		summary.DetailPenyadap = penyadapDetails
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
		WHERE bp.deleted_at IS NULL 
		AND p.nama_penyadap LIKE ? 
	`

	args := []interface{}{"%" + nama + "%"}

	// Filter tanggal
	if tanggal != "" {
		query += " AND DATE(bp.tanggal) = DATE(?)"
		args = append(args, tanggal)
	}

	// Filter afdeling
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
