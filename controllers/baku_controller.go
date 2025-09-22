package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
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
	ID              uint                `json:"id"`
	TahunTanam      uint                `json:"tahunTanam"`
	Mandor          string              `json:"mandor"`
	Tanggal         time.Time           `json:"tanggal"`
	Afdeling        string              `json:"afdeling"`
	Tipe            models.TipeProduksi `json:"tipe"`
	TotalBasahLatex float64             `json:"totalBasahLatex"`
	TotalSheet      float64             `json:"totalSheet"`
	TotalBasahLump  float64             `json:"totalBasahLump"`
	TotalBrCr       float64             `json:"totalBrCr"`
	JumlahPenyadap  int                 `json:"jumlahPenyadap"`
	DetailPenyadap  []PenyadapDetail    `json:"detailPenyadap,omitempty"`
}

type PenyadapDetail struct {
	ID              uint                `json:"id"`
	NamaPenyadap    string              `json:"namaPenyadap"`
	NIK             string              `json:"nik"`
	TahunTanam      uint                `json:"tahunTanam"`
	Tipe            models.TipeProduksi `json:"tipe"`
	Tanggal         string              `json:"tanggal"`
	TotalBasahLatex float64             `json:"totalBasahLatex"`
	TotalSheet      float64             `json:"totalSheet"`
	TotalBasahLump  float64             `json:"totalBasahLump"`
	TotalBrCr       float64             `json:"totalBrCr"`
	JumlahHariKerja int                 `json:"jumlahHariKerja"`
	Mandor          string              `json:"mandor,omitempty"`
	Afdeling        string              `json:"afdeling,omitempty"`
}

type ReportingResponse struct {
	Success    bool            `json:"success"`
	Message    string          `json:"message"`
	Data       []MandorSummary `json:"data"`
	FilterInfo FilterInfo      `json:"filterInfo"`
}

type FilterInfo struct {
	TanggalMulai   string              `json:"tanggalMulai,omitempty"`
	TanggalSelesai string              `json:"tanggalSelesai,omitempty"`
	Tanggal        string              `json:"tanggal,omitempty"`
	Tipe           models.TipeProduksi `json:"tipe,omitempty"`
	TotalRecord    int                 `json:"totalRecord"`
	Periode        string              `json:"periode"`
	JumlahHari     int                 `json:"jumlahHari,omitempty"`
}

type BakuPageData struct {
	Title        string
	MandorList   []models.BakuMandor
	PenyadapList []models.BakuPenyadap
	TipeList     []models.TipeProduksi // For mandor input form
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

// ======== NEW: Get Available Production Types ========
func GetTipeProduksiList(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Daftar tipe produksi berhasil diambil",
		Data:    models.GetAllTipeProduksi(),
	})
}

// ======== CRUD OPERATIONS ========
// GetBakuPenyadapByDate - Get penyadap data for a specific date with each penyadap's total, not the aggregate.

// GetAllBakuPenyadap - Get all penyadap records with optional tipe filter
func GetAllBakuPenyadap(w http.ResponseWriter, r *http.Request) {
	tipeFilter := r.URL.Query().Get("tipe")

	var penyadap []models.BakuPenyadap
	query := config.DB.Preload("Mandor").Preload("Penyadap").Order("created_at desc")

	// Filter by tipe if provided
	if tipeFilter != "" {
		if !models.IsValidTipeProduksi(models.TipeProduksi(tipeFilter)) {
			respondJSON(w, http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "Tipe produksi tidak valid",
			})
			return
		}
		query = query.Where("tipe = ?", tipeFilter)
	}

	if err := query.Find(&penyadap).Error; err != nil {
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
	if penyadap.TahunTanam == 0 {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Tahun tanam wajib diisi",
		})
		return
	}

	// Auto-set tipe dari mandor
	var mandor models.BakuMandor
	if err := config.DB.First(&mandor, penyadap.IdBakuMandor).Error; err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Mandor dengan ID tersebut tidak ditemukan",
		})
		return
	}

	penyadap.Tipe = mandor.Tipe
	fmt.Printf("DEBUG: Auto-setting tipe '%s' from mandor '%s'\n", penyadap.Tipe, mandor.Mandor)

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
		Message: fmt.Sprintf("Data penyadap berhasil ditambahkan dengan tipe %s (tahun tanam %d, dari profil mandor)", penyadap.Tipe, penyadap.TahunTanam),
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

	// UPDATED: If mandor changed, update tipe from new mandor profile
	if updates.IdBakuMandor != 0 && updates.IdBakuMandor != existing.IdBakuMandor {
		var newMandor models.BakuMandor
		if err := config.DB.First(&newMandor, updates.IdBakuMandor).Error; err != nil {
			respondJSON(w, http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "Mandor baru dengan ID tersebut tidak ditemukan",
			})
			return
		}

		// Auto-set tipe from new mandor
		updates.Tipe = newMandor.Tipe
		fmt.Printf("DEBUG: Mandor changed, updating tipe to '%s' from mandor '%s'\n", updates.Tipe, newMandor.Mandor)
	} else {
		// UPDATED: If mandor didn't change, keep existing tipe from mandor profile
		var currentMandor models.BakuMandor
		if err := config.DB.First(&currentMandor, existing.IdBakuMandor).Error; err == nil {
			updates.Tipe = currentMandor.Tipe
			fmt.Printf("DEBUG: Mandor unchanged, keeping tipe '%s' from mandor '%s'\n", updates.Tipe, currentMandor.Mandor)
		}
	}

	if err := config.DB.Model(&existing).Updates(updates).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal update penyadap: " + err.Error(),
		})
		return
	}

	// Update detail berdasarkan tanggal, mandor, dan tipe
	updateBakuDetail(existing, "update", &oldCopy)

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: fmt.Sprintf("Data penyadap berhasil diperbarui dengan tipe %s", existing.Tipe),
	})
}

// ======== DETAIL OPERATIONS WITH TIPE SUPPORT ========

// GetAllBakuDetail - Get all detail records with optional tipe filter
func GetAllBakuDetail(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("DEBUG: GetAllBakuDetail called - Method: %s, URL: %s\n", r.Method, r.URL.Path)

	tipeFilter := r.URL.Query().Get("tipe")

	var details []models.BakuDetail
	query := config.DB.Order("tanggal desc, mandor asc")

	// Filter by tipe if provided
	if tipeFilter != "" {
		if !models.IsValidTipeProduksi(models.TipeProduksi(tipeFilter)) {
			respondJSON(w, http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "Tipe produksi tidak valid",
			})
			return
		}
		query = query.Where("tipe = ?", tipeFilter)
	}

	if err := query.Find(&details).Error; err != nil {
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

// GetBakuDetailByDate - Get detail by specific date with optional tipe filter
func GetBakuDetailByDate(w http.ResponseWriter, r *http.Request) {
	tanggalStr := mux.Vars(r)["tanggal"]
	tipeFilter := r.URL.Query().Get("tipe")

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
	query := config.DB.Where("DATE(tanggal) = DATE(?)", tanggal).Order("mandor asc")

	// Filter by tipe if provided
	if tipeFilter != "" {
		if !models.IsValidTipeProduksi(models.TipeProduksi(tipeFilter)) {
			respondJSON(w, http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "Tipe produksi tidak valid",
			})
			return
		}
		query = query.Where("tipe = ?", tipeFilter)
	}

	if err := query.Find(&details).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil detail: " + err.Error(),
		})
		return
	}

	if len(details) == 0 {
		message := "Detail untuk tanggal " + tanggalStr
		if tipeFilter != "" {
			message += " dengan tipe " + tipeFilter
		}
		message += " tidak ditemukan"

		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: message,
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Detail berhasil ditemukan",
		Data:    details,
	})
}

// GetBakuDetailByDateAndMandor - Get detail by date, mandor, and optional tipe
func GetBakuDetailByDateAndMandor(w http.ResponseWriter, r *http.Request) {
	tanggalStr := r.URL.Query().Get("tanggal")
	mandor := r.URL.Query().Get("mandor")
	tipeFilter := r.URL.Query().Get("tipe")

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
	query := config.DB.Where("DATE(tanggal) = DATE(?) AND mandor = ?", tanggal, mandor)

	// Filter by tipe if provided
	if tipeFilter != "" {
		if !models.IsValidTipeProduksi(models.TipeProduksi(tipeFilter)) {
			respondJSON(w, http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "Tipe produksi tidak valid",
			})
			return
		}
		query = query.Where("tipe = ?", tipeFilter)
	}

	if err := query.First(&detail).Error; err != nil {
		message := fmt.Sprintf("Detail untuk tanggal %s dan mandor %s", tanggalStr, mandor)
		if tipeFilter != "" {
			message += " dengan tipe " + tipeFilter
		}
		message += " tidak ditemukan"

		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: message,
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Detail berhasil ditemukan",
		Data:    detail,
	})
}

// updateBakuDetail - Helper function to update daily summary with tipe support
// Mencari berdasarkan kombinasi tanggal, mandor, DAN tipe
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
	err := config.DB.Where("DATE(tanggal) = DATE(?) AND mandor = ? AND tipe = ?", targetDate, mandor.Mandor, entry.Tipe).First(&detail).Error

	if err != nil {
		// Jika belum ada detail untuk kombinasi tanggal + mandor + tipe → buat baru
		if action == "create" {
			detail = models.BakuDetail{
				Tanggal:  targetDate,
				Mandor:   mandor.Mandor,
				Afdeling: mandor.Afdeling,
				Tipe:     entry.Tipe,
				// semua field default 0
			}
		} else {
			fmt.Printf("WARNING: Tidak ada BakuDetail untuk tanggal %s mandor %s tipe %s pada action %s\n",
				targetDate.Format("2006-01-02"), mandor.Mandor, entry.Tipe, action)
			return
		}
	}

	// ================== Update nilai berdasarkan action ==================
	switch action {
	case "create":
		fmt.Printf("CREATE: Menambah data untuk %s mandor %s tipe %s\n", targetDate.Format("2006-01-02"), mandor.Mandor, entry.Tipe)
		detail.JumlahPabrikBasahLatek += entry.BasahLatex
		detail.JumlahPabrikBasahLump += entry.BasahLump
		detail.JumlahSheet += entry.Sheet
		detail.JumlahBrCr += entry.BrCr

	case "update":
		if oldEntry != nil {
			fmt.Printf("UPDATE: Mengupdate data untuk %s mandor %s tipe %s\n", targetDate.Format("2006-01-02"), mandor.Mandor, entry.Tipe)

			// Jika tipe berubah, perlu update 2 detail berbeda
			if oldEntry.Tipe != entry.Tipe {
				// Kurangi dari detail lama
				var oldDetail models.BakuDetail
				if err := config.DB.Where("DATE(tanggal) = DATE(?) AND mandor = ? AND tipe = ?", targetDate, mandor.Mandor, oldEntry.Tipe).First(&oldDetail).Error; err == nil {
					oldDetail.JumlahPabrikBasahLatek -= oldEntry.BasahLatex
					oldDetail.JumlahPabrikBasahLump -= oldEntry.BasahLump
					oldDetail.JumlahSheet -= oldEntry.Sheet
					oldDetail.JumlahBrCr -= oldEntry.BrCr

					// Prevent negative values
					if oldDetail.JumlahPabrikBasahLatek < 0 {
						oldDetail.JumlahPabrikBasahLatek = 0
					}
					if oldDetail.JumlahPabrikBasahLump < 0 {
						oldDetail.JumlahPabrikBasahLump = 0
					}
					if oldDetail.JumlahSheet < 0 {
						oldDetail.JumlahSheet = 0
					}
					if oldDetail.JumlahBrCr < 0 {
						oldDetail.JumlahBrCr = 0
					}

					config.DB.Save(&oldDetail)
				}

				// Tambah ke detail baru
				detail.JumlahPabrikBasahLatek += entry.BasahLatex
				detail.JumlahPabrikBasahLump += entry.BasahLump
				detail.JumlahSheet += entry.Sheet
				detail.JumlahBrCr += entry.BrCr
			} else {
				// Tipe sama, hitung selisih
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
		}

	case "delete":
		fmt.Printf("DELETE: Mengurangi data untuk %s mandor %s tipe %s\n", targetDate.Format("2006-01-02"), mandor.Mandor, entry.Tipe)
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
		fmt.Printf("SUCCESS: BakuDetail %s mandor %s tipe %s terupdate\n",
			targetDate.Format("2006-01-02"), mandor.Mandor, detail.Tipe)
		fmt.Printf("  - Pabrik Latex: %.2f | Kebun Latex: %.2f\n", detail.JumlahPabrikBasahLatek, detail.JumlahKebunBasahLatek)
		fmt.Printf("  - Pabrik Lump: %.2f | Kebun Lump: %.2f\n", detail.JumlahPabrikBasahLump, detail.JumlahKebunBasahLump)
	}
}

// RecalculateBakuDetail - Fungsi untuk hitung ulang BakuDetail berdasarkan tanggal, mandor, dan tipe
func RecalculateBakuDetail(tanggal time.Time, mandorID uint, tipe models.TipeProduksi) error {
	targetDate := tanggal.Truncate(24 * time.Hour)

	// Ambil data mandor
	var mandor models.BakuMandor
	if err := config.DB.First(&mandor, mandorID).Error; err != nil {
		return fmt.Errorf("mandor tidak ditemukan: %v", err)
	}

	// Hitung ulang total dari semua BakuPenyadap untuk tanggal, mandor, dan tipe tersebut
	var totals struct {
		TotalBasahLatex float64
		TotalSheet      float64
		TotalBasahLump  float64
		TotalBrCr       float64
	}

	err := config.DB.Model(&models.BakuPenyadap{}).
		Where("DATE(tanggal) = DATE(?) AND id_baku_mandor = ? AND tipe = ?", targetDate, mandorID, tipe).
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
	err = config.DB.Where("DATE(tanggal) = DATE(?) AND mandor = ? AND tipe = ?", targetDate, mandor.Mandor, tipe).First(&detail).Error

	if err != nil {
		// Buat baru
		detail = models.BakuDetail{
			Tanggal:  targetDate,
			Mandor:   mandor.Mandor,
			Afdeling: mandor.Afdeling,
			Tipe:     tipe,
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
		TipeList:     models.GetAllTipeProduksi(), // BARU: Tambahkan list tipe
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

// ======== REPORTING FUNCTIONS WITH TIPE SUPPORT ========

// GetMandorSummaryAll - Get summary of all mandors for all time with optional tipe filter
func GetMandorSummaryAll(w http.ResponseWriter, r *http.Request) {
	tipeFilter := r.URL.Query().Get("tipe")

	summaries, err := getMandorSummaries("", tipeFilter)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data summary: " + err.Error(),
		})
		return
	}

	periode := "Semua waktu"
	if tipeFilter != "" {
		periode += " - Tipe: " + tipeFilter
	}

	response := ReportingResponse{
		Success: true,
		Message: "Data summary mandor berhasil diambil",
		Data:    summaries,
		FilterInfo: FilterInfo{
			TotalRecord: len(summaries),
			Periode:     periode,
			Tipe:        models.TipeProduksi(tipeFilter),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetMandorSummaryByDate - Get mandor summary for specific date with optional tipe filter
func GetMandorSummaryByDate(w http.ResponseWriter, r *http.Request) {
	tanggalStr := mux.Vars(r)["tanggal"]
	tipeFilter := r.URL.Query().Get("tipe")

	// Validasi format tanggal
	if _, err := time.Parse("2006-01-02", tanggalStr); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Format tanggal tidak valid. Gunakan format YYYY-MM-DD",
		})
		return
	}

	summaries, err := getMandorSummaries(tanggalStr, tipeFilter)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data summary: " + err.Error(),
		})
		return
	}

	periode := "Tanggal: " + tanggalStr
	if tipeFilter != "" {
		periode += " - Tipe: " + tipeFilter
	}

	response := ReportingResponse{
		Success: true,
		Message: "Data summary mandor untuk tanggal " + tanggalStr + " berhasil diambil",
		Data:    summaries,
		FilterInfo: FilterInfo{
			Tanggal:     tanggalStr,
			TotalRecord: len(summaries),
			Periode:     periode,
			Tipe:        models.TipeProduksi(tipeFilter),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetPenyadapDetailAll - Get details of all penyadap for all time with optional tipe filter
func GetPenyadapDetailAll(w http.ResponseWriter, r *http.Request) {
	tipeFilter := r.URL.Query().Get("tipe")

	details, err := getPenyadapDetails("", tipeFilter)
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

// GetPenyadapDetailByDate - Get penyadap details for specific date with optional tipe filter
func GetPenyadapDetailByDate(w http.ResponseWriter, r *http.Request) {
	tanggalStr := mux.Vars(r)["tanggal"]
	tipeFilter := r.URL.Query().Get("tipe")

	if _, err := time.Parse("2006-01-02", tanggalStr); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Format tanggal tidak valid. Gunakan format YYYY-MM-DD",
		})
		return
	}

	details, err := getPenyadapDetails(tanggalStr, tipeFilter)
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

// ======== SEARCH FUNCTIONS WITH TIPE SUPPORT ========

// SearchMandorByName - Search mandor by name with optional date and tipe filter
func SearchMandorByName(w http.ResponseWriter, r *http.Request) {
	// Query parameters
	namaMandor := r.URL.Query().Get("nama")
	tanggal := r.URL.Query().Get("tanggal")
	tipeFilter := r.URL.Query().Get("tipe")

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

	// Validasi tipe jika ada
	if tipeFilter != "" && !models.IsValidTipeProduksi(models.TipeProduksi(tipeFilter)) {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Tipe produksi tidak valid",
		})
		return
	}

	summaries, err := searchMandorSummaries(namaMandor, tanggal, tipeFilter)
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
	if tipeFilter != "" {
		periode += " - Tipe: " + tipeFilter
	}

	response := ReportingResponse{
		Success: true,
		Message: "Hasil pencarian mandor '" + namaMandor + "'",
		Data:    summaries,
		FilterInfo: FilterInfo{
			Tanggal:     tanggal,
			TotalRecord: len(summaries),
			Periode:     periode,
			Tipe:        models.TipeProduksi(tipeFilter),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SearchPenyadapByName - Search penyadap by name with optional date and tipe filter
func SearchPenyadapByName(w http.ResponseWriter, r *http.Request) {
	namaPenyadap := r.URL.Query().Get("nama")
	tanggal := r.URL.Query().Get("tanggal")
	tipeFilter := r.URL.Query().Get("tipe")

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

	// Validasi tipe jika ada
	if tipeFilter != "" && !models.IsValidTipeProduksi(models.TipeProduksi(tipeFilter)) {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Tipe produksi tidak valid",
		})
		return
	}

	details, err := searchPenyadapDetails(namaPenyadap, tanggal, tipeFilter)
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
	if tipeFilter != "" {
		periode += " - Tipe: " + tipeFilter
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Hasil pencarian penyadap '" + namaPenyadap + "' - " + periode,
		Data:    details,
	})
}

// GetMandorWithPenyadapDetail - Get mandor details with all penyadap, with optional tipe filter
func GetMandorWithPenyadapDetail(w http.ResponseWriter, r *http.Request) {
	namaMandor := r.URL.Query().Get("nama")
	tanggal := r.URL.Query().Get("tanggal")
	tipeFilter := r.URL.Query().Get("tipe")

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

	if tipeFilter != "" && !models.IsValidTipeProduksi(models.TipeProduksi(tipeFilter)) {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Tipe produksi tidak valid",
		})
		return
	}

	summaries, err := searchMandorWithDetails(namaMandor, tanggal, tipeFilter)
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
	if tipeFilter != "" {
		periode += " - Tipe: " + tipeFilter
	}

	response := ReportingResponse{
		Success: true,
		Message: "Detail mandor '" + namaMandor + "' beserta penyadapnya",
		Data:    summaries,
		FilterInfo: FilterInfo{
			Tanggal:     tanggal,
			TotalRecord: len(summaries),
			Periode:     periode,
			Tipe:        models.TipeProduksi(tipeFilter),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ======== HELPER FUNCTIONS WITH TIPE SUPPORT ========

func getMandorSummaries(tanggal, tipeFilter string) ([]MandorSummary, error) {
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

		if tipeFilter != "" {
			query = query.Where("tipe = ?", tipeFilter)
			summary.Tipe = models.TipeProduksi(tipeFilter)
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

// getPenyadapDetails - Updated with tipe support
func getPenyadapDetails(tanggal, tipeFilter string) ([]PenyadapDetail, error) {
	// Jika ada filter tanggal, ambil dari baku_detail berdasarkan tanggal, mandor, dan tipe
	if tanggal != "" {
		targetDate, err := time.Parse("2006-01-02", tanggal)
		if err != nil {
			return nil, err
		}

		// Ambil semua BakuDetail untuk tanggal tersebut
		var bakuDetails []models.BakuDetail
		query := config.DB.Where("DATE(tanggal) = DATE(?)", targetDate)

		if tipeFilter != "" {
			query = query.Where("tipe = ?", tipeFilter)
		}

		err = query.Find(&bakuDetails).Error
		if err != nil {
			return nil, err
		}

		if len(bakuDetails) == 0 {
			return []PenyadapDetail{}, nil
		}

		var allDetails []PenyadapDetail

		// Untuk setiap mandor yang ada di BakuDetail
		for _, bakuDetail := range bakuDetails {
			// Ambil daftar penyadap yang aktif untuk mandor, tanggal, dan tipe tersebut
			queryPenyadap := `
				SELECT DISTINCT
					p.id,
					p.nama_penyadap,
					p.nik,
					bm.mandor,
					bm.afdeling,
					bp.tipe,
					COUNT(bp.id) as jumlah_hari_kerja
				FROM penyadaps p
				INNER JOIN baku_penyadaps bp ON p.id = bp.id_penyadap
				INNER JOIN baku_mandors bm ON bp.id_baku_mandor = bm.id
				WHERE bp.deleted_at IS NULL 
				AND DATE(bp.tanggal) = DATE(?)
				AND bm.mandor = ?
				AND bp.tipe = ?
				GROUP BY p.id, p.nama_penyadap, p.nik, bm.mandor, bm.afdeling, bp.tipe
				ORDER BY p.nama_penyadap
			`

			var mandorDetails []PenyadapDetail
			err = config.DB.Raw(queryPenyadap, targetDate, bakuDetail.Mandor, bakuDetail.Tipe).Scan(&mandorDetails).Error
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
			bp.tipe,
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
	if tipeFilter != "" {
		query += " AND bp.tipe = ?"
		args = append(args, tipeFilter)
	}

	query += " GROUP BY p.id, p.nama_penyadap, p.nik, bp.tipe ORDER BY p.nama_penyadap"

	var details []PenyadapDetail
	if err := config.DB.Raw(query, args...).Scan(&details).Error; err != nil {
		return nil, err
	}

	return details, nil
}

func searchMandorSummaries(namaMandor, tanggal, tipeFilter string) ([]MandorSummary, error) {
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

		if tipeFilter != "" {
			bakuQuery = bakuQuery.Where("tipe = ?", tipeFilter)
			summary.Tipe = models.TipeProduksi(tipeFilter)
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
func searchPenyadapDetails(namaPenyadap, tanggal, tipeFilter string) ([]PenyadapDetail, error) {
	query := `
		SELECT 
			p.id,
			p.nama_penyadap,
			p.nik,
			bm.mandor,
			bm.afdeling,
			bm.tahun_tanam,                          -- <-- ambil tahun tanam
			bp.tipe,
			DATE(bp.tanggal) as tanggal,
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

	if tipeFilter != "" {
		query += " AND bp.tipe = ?"
		args = append(args, tipeFilter)
	}

	// tambahkan bm.tahun_tanam ke GROUP BY
	query += `
		GROUP BY p.id, p.nama_penyadap, p.nik, bm.mandor, bm.afdeling, bm.tahun_tanam, bp.tipe, bp.tanggal
		ORDER BY bp.tanggal, p.nama_penyadap
	`

	var details []PenyadapDetail
	if err := config.DB.Raw(query, args...).Scan(&details).Error; err != nil {
		return nil, err
	}

	return details, nil
}

func searchMandorWithDetails(namaMandor, tanggal, tipeFilter string) ([]MandorSummary, error) {
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

		if tipeFilter != "" {
			bakuQuery = bakuQuery.Where("tipe = ?", tipeFilter)
			summary.Tipe = models.TipeProduksi(tipeFilter)
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
    	bp.tipe,
    	DATE(bp.tanggal) as tanggal,   -- <== ambil tanggal
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

		if tipeFilter != "" {
			penyadapQuery += " AND bp.tipe = ?"
			penyadapArgs = append(penyadapArgs, tipeFilter)
		}

		penyadapQuery += " GROUP BY p.id, p.nama_penyadap, p.nik, bp.tipe ORDER BY p.nama_penyadap"

		var penyadapDetails []PenyadapDetail
		if err := config.DB.Raw(penyadapQuery, penyadapArgs...).Scan(&penyadapDetails).Error; err != nil {
			return nil, err
		}

		summary.DetailPenyadap = penyadapDetails
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

// Advanced search functions with tipe support
func advancedSearchMandor(nama, tanggal, afdeling, tahunTanam, tipeFilter string) ([]MandorSummary, error) {
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

		if tipeFilter != "" {
			bakuQuery = bakuQuery.Where("tipe = ?", tipeFilter)
			summary.Tipe = models.TipeProduksi(tipeFilter)
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
				bp.tipe,
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

		if tipeFilter != "" {
			penyadapQuery += " AND bp.tipe = ?"
			args = append(args, tipeFilter)
		}

		penyadapQuery += " GROUP BY p.id, p.nama_penyadap, p.nik, bp.tipe ORDER BY p.nama_penyadap"

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

func advancedSearchPenyadap(nama, tanggal, afdeling, tipeFilter string) ([]PenyadapDetail, error) {
	query := `
		SELECT 
			p.id,
			p.nama_penyadap,
			p.nik,
			bm.mandor,
			bm.afdeling,
			bp.tipe,
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

	// Filter tipe
	if tipeFilter != "" {
		query += " AND bp.tipe = ?"
		args = append(args, tipeFilter)
	}

	query += " GROUP BY p.id, p.nama_penyadap, p.nik, bm.mandor, bm.afdeling, bp.tipe ORDER BY p.nama_penyadap"

	var details []PenyadapDetail
	if err := config.DB.Raw(query, args...).Scan(&details).Error; err != nil {
		return nil, err
	}

	return details, nil
}

// SearchAll - Global search with tipe support
func SearchAll(w http.ResponseWriter, r *http.Request) {
	searchType := r.URL.Query().Get("type") // "mandor" atau "penyadap"
	nama := r.URL.Query().Get("nama")
	tanggal := r.URL.Query().Get("tanggal")
	afdeling := r.URL.Query().Get("afdeling")
	tahunTanam := r.URL.Query().Get("tahun")
	tipeFilter := r.URL.Query().Get("tipe")

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

	// Validasi tipe jika ada
	if tipeFilter != "" && !models.IsValidTipeProduksi(models.TipeProduksi(tipeFilter)) {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Tipe produksi tidak valid",
		})
		return
	}

	var result interface{}
	var err error

	switch searchType {
	case "mandor":
		result, err = advancedSearchMandor(nama, tanggal, afdeling, tahunTanam, tipeFilter)
	case "penyadap":
		result, err = advancedSearchPenyadap(nama, tanggal, afdeling, tipeFilter)
	default:
		// Auto detect berdasarkan hasil pencarian
		mandorResult, _ := advancedSearchMandor(nama, tanggal, afdeling, tahunTanam, tipeFilter)
		penyadapResult, _ := advancedSearchPenyadap(nama, tanggal, afdeling, tipeFilter)

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

// GetBakuPenyadapByID - Get single penyadap record by ID
func GetBakuPenyadapByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var penyadap models.BakuPenyadap
	if err := config.DB.First(&penyadap, id).Error; err != nil {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Data penyadap tidak ditemukan",
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data penyadap berhasil ditemukan",
		Data:    penyadap,
	})
}

// DeleteBakuPenyadap - Delete penyadap record by ID
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

	// Update detail harian setelah delete
	updateBakuDetail(penyadap, "delete", nil)

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data penyadap berhasil dihapus",
	})
}
func parseDateRange(tanggalMulai, tanggalSelesai string) (time.Time, time.Time, error) {
	var startDate, endDate time.Time
	var err error

	if tanggalMulai != "" {
		startDate, err = time.Parse("2006-01-02", tanggalMulai)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("format tanggal mulai tidak valid: %v", err)
		}
	}

	if tanggalSelesai != "" {
		endDate, err = time.Parse("2006-01-02", tanggalSelesai)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("format tanggal selesai tidak valid: %v", err)
		}
	}

	// Validasi: tanggal mulai tidak boleh lebih besar dari tanggal selesai
	if !startDate.IsZero() && !endDate.IsZero() && startDate.After(endDate) {
		return time.Time{}, time.Time{}, fmt.Errorf("tanggal mulai tidak boleh lebih besar dari tanggal selesai")
	}

	return startDate, endDate, nil
}

func formatDateRangePeriode(tanggalMulai, tanggalSelesai string) (string, int) {
	if tanggalMulai != "" && tanggalSelesai != "" {
		startDate, _ := time.Parse("2006-01-02", tanggalMulai)
		endDate, _ := time.Parse("2006-01-02", tanggalSelesai)
		days := int(endDate.Sub(startDate).Hours()/24) + 1
		return fmt.Sprintf("Periode: %s s/d %s", tanggalMulai, tanggalSelesai), days
	} else if tanggalMulai != "" {
		return fmt.Sprintf("Dari tanggal: %s", tanggalMulai), -1
	} else if tanggalSelesai != "" {
		return fmt.Sprintf("Sampai tanggal: %s", tanggalSelesai), -1
	}
	return "Semua waktu", -1
}

func GetBakuDetailByDateRange(w http.ResponseWriter, r *http.Request) {
	tanggalMulai := r.URL.Query().Get("tanggal_mulai")
	tanggalSelesai := r.URL.Query().Get("tanggal_selesai")
	tipeFilter := r.URL.Query().Get("tipe")

	if tanggalMulai == "" || tanggalSelesai == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Parameter tanggal_mulai dan tanggal_selesai wajib diisi",
		})
		return
	}

	// Parse and validate date range
	startDate, endDate, err := parseDateRange(tanggalMulai, tanggalSelesai)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	var details []models.BakuDetail
	query := config.DB.
		Where("tanggal >= ? AND tanggal < ?", startDate, endDate.Add(24*time.Hour)).
		Order("tanggal desc, mandor asc")

	// Filter by tipe if provided
	if tipeFilter != "" {
		if !models.IsValidTipeProduksi(models.TipeProduksi(tipeFilter)) {
			respondJSON(w, http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "Tipe produksi tidak valid",
			})
			return
		}
		query = query.Where("tipe = ?", tipeFilter)
	}

	if err := query.Find(&details).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil detail: " + err.Error(),
		})
		return
	}

	periode, _ := formatDateRangePeriode(tanggalMulai, tanggalSelesai)
	if tipeFilter != "" {
		periode += " - Tipe: " + tipeFilter
	}

	if len(details) == 0 {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Detail untuk " + periode + " tidak ditemukan",
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Detail berhasil ditemukan untuk " + periode,
		Data:    details,
	})

}
func GetMandorSummaryByDateRange(w http.ResponseWriter, r *http.Request) {
	tanggalMulai := r.URL.Query().Get("tanggal_mulai")
	tanggalSelesai := r.URL.Query().Get("tanggal_selesai")
	tipeFilter := r.URL.Query().Get("tipe")

	if tanggalMulai == "" || tanggalSelesai == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Parameter tanggal_mulai dan tanggal_selesai wajib diisi",
		})
		return
	}

	// Validate date range
	_, _, err := parseDateRange(tanggalMulai, tanggalSelesai)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	summaries, err := getMandorSummariesByDateRange(tanggalMulai, tanggalSelesai, tipeFilter)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data summary: " + err.Error(),
		})
		return
	}

	periode, jumlahHari := formatDateRangePeriode(tanggalMulai, tanggalSelesai)
	if tipeFilter != "" {
		periode += " - Tipe: " + tipeFilter
	}

	response := ReportingResponse{
		Success: true,
		Message: "Data summary mandor untuk " + periode + " berhasil diambil",
		Data:    summaries,
		FilterInfo: FilterInfo{
			TanggalMulai:   tanggalMulai,
			TanggalSelesai: tanggalSelesai,
			TotalRecord:    len(summaries),
			Periode:        periode,
			JumlahHari:     jumlahHari,
			Tipe:           models.TipeProduksi(tipeFilter),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getMandorSummariesByDateRange(tanggalMulai, tanggalSelesai, tipeFilter string) ([]MandorSummary, error) {
	startDate, endDate, _ := parseDateRange(tanggalMulai, tanggalSelesai)

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

		// Query for totals from all penyadap for this mandor in date range
		query := config.DB.Model(&models.BakuPenyadap{}).Where("id_baku_mandor = ?", mandor.ID)

		// Apply date range filter
		query = buildDateRangeQuery(query, startDate, endDate)

		if tipeFilter != "" {
			query = query.Where("tipe = ?", tipeFilter)
			summary.Tipe = models.TipeProduksi(tipeFilter)
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

func buildDateRangeQuery(query *gorm.DB, startDate, endDate time.Time) *gorm.DB {
	if !startDate.IsZero() && !endDate.IsZero() {
		return query.Where("DATE(tanggal) BETWEEN DATE(?) AND DATE(?)", startDate, endDate)
	} else if !startDate.IsZero() {
		return query.Where("DATE(tanggal) >= DATE(?)", startDate)
	} else if !endDate.IsZero() {
		return query.Where("DATE(tanggal) <= DATE(?)", endDate)
	}
	return query
}

func GetPenyadapDetailByDateRange(w http.ResponseWriter, r *http.Request) {
	tanggalMulai := r.URL.Query().Get("tanggal_mulai")
	tanggalSelesai := r.URL.Query().Get("tanggal_selesai")
	tipeFilter := r.URL.Query().Get("tipe")

	if tanggalMulai == "" || tanggalSelesai == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Parameter tanggal_mulai dan tanggal_selesai wajib diisi",
		})
		return
	}

	// Validate date range
	_, _, err := parseDateRange(tanggalMulai, tanggalSelesai)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	details, err := getPenyadapDetailsByDateRange(tanggalMulai, tanggalSelesai, tipeFilter)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil detail penyadap: " + err.Error(),
		})
		return
	}

	periode, _ := formatDateRangePeriode(tanggalMulai, tanggalSelesai)
	if tipeFilter != "" {
		periode += " - Tipe: " + tipeFilter
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Detail penyadap untuk " + periode + " berhasil diambil",
		Data:    details,
	})
}

func SearchMandorWithDateRange(w http.ResponseWriter, r *http.Request) {
	namaMandor := r.URL.Query().Get("nama")
	tanggalMulai := r.URL.Query().Get("tanggal_mulai")
	tanggalSelesai := r.URL.Query().Get("tanggal_selesai")
	tipeFilter := r.URL.Query().Get("tipe")

	if namaMandor == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Parameter 'nama' wajib diisi",
		})
		return
	}

	// Validate date range if provided
	if tanggalMulai != "" || tanggalSelesai != "" {
		_, _, err := parseDateRange(tanggalMulai, tanggalSelesai)
		if err != nil {
			respondJSON(w, http.StatusBadRequest, APIResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}
	}

	// Validate tipe if provided
	if tipeFilter != "" && !models.IsValidTipeProduksi(models.TipeProduksi(tipeFilter)) {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Tipe produksi tidak valid",
		})
		return
	}

	summaries, err := searchMandorSummariesWithDateRange(namaMandor, tanggalMulai, tanggalSelesai, tipeFilter)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mencari data mandor: " + err.Error(),
		})
		return
	}

	periode, jumlahHari := formatDateRangePeriode(tanggalMulai, tanggalSelesai)
	if tipeFilter != "" {
		periode += " - Tipe: " + tipeFilter
	}

	response := ReportingResponse{
		Success: true,
		Message: "Hasil pencarian mandor '" + namaMandor + "'",
		Data:    summaries,
		FilterInfo: FilterInfo{
			TanggalMulai:   tanggalMulai,
			TanggalSelesai: tanggalSelesai,
			TotalRecord:    len(summaries),
			Periode:        periode,
			JumlahHari:     jumlahHari,
			Tipe:           models.TipeProduksi(tipeFilter),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func SearchPenyadapWithDateRange(w http.ResponseWriter, r *http.Request) {
	namaPenyadap := r.URL.Query().Get("nama")
	tanggalMulai := r.URL.Query().Get("tanggal_mulai")
	tanggalSelesai := r.URL.Query().Get("tanggal_selesai")
	tipeFilter := r.URL.Query().Get("tipe")

	if namaPenyadap == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Parameter 'nama' wajib diisi",
		})
		return
	}

	// Validate date range if provided
	if tanggalMulai != "" || tanggalSelesai != "" {
		_, _, err := parseDateRange(tanggalMulai, tanggalSelesai)
		if err != nil {
			respondJSON(w, http.StatusBadRequest, APIResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}
	}

	// Validate tipe if provided
	if tipeFilter != "" && !models.IsValidTipeProduksi(models.TipeProduksi(tipeFilter)) {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Tipe produksi tidak valid",
		})
		return
	}

	details, err := searchPenyadapDetailsWithDateRange(namaPenyadap, tanggalMulai, tanggalSelesai, tipeFilter)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mencari detail penyadap: " + err.Error(),
		})
		return
	}

	periode, _ := formatDateRangePeriode(tanggalMulai, tanggalSelesai)
	if tipeFilter != "" {
		periode += " - Tipe: " + tipeFilter
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Hasil pencarian penyadap '" + namaPenyadap + "' - " + periode,
		Data:    details,
	})
}
func getPenyadapDetailsByDateRange(tanggalMulai, tanggalSelesai string, tipeFilter string) ([]PenyadapDetail, error) {
	startDate, endDate, err := parseDateRange(tanggalMulai, tanggalSelesai)
	if err != nil {
		return nil, err
	}

	var details []PenyadapDetail
	query := config.DB.Table("baku_penyadaps").Select(`
        penyadaps.id,
        penyadaps.nama_penyadap,
        penyadaps.nik,
        baku_mandors.mandor,
        baku_mandors.afdeling,
        baku_penyadaps.tipe,
        COALESCE(SUM(baku_penyadaps.basah_latex), 0) as total_basah_latex,
        COALESCE(SUM(baku_penyadaps.sheet), 0) as total_sheet,
        COALESCE(SUM(baku_penyadaps.basah_lump), 0) as total_basah_lump,
        COALESCE(SUM(baku_penyadaps.br_cr), 0) as total_br_cr,
        COUNT(baku_penyadaps.id) as jumlah_hari_kerja
    `).
		Joins("JOIN penyadaps ON penyadaps.id = baku_penyadaps.id_penyadap").
		Joins("JOIN baku_mandors ON baku_penyadaps.id_baku_mandor = baku_mandors.id").
		Where("DATE(baku_penyadaps.tanggal) BETWEEN DATE(?) AND DATE(?)", startDate, endDate)

	// Apply tipe filter if provided
	if tipeFilter != "" {
		query = query.Where("baku_penyadaps.tipe = ?", tipeFilter)
	}

	// Execute query
	if err := query.Group("penyadaps.id, penyadaps.nama_penyadap, penyadaps.nik, baku_mandors.mandor, baku_mandors.afdeling, baku_penyadaps.tipe").
		Scan(&details).Error; err != nil {
		return nil, err
	}

	return details, nil
}

func searchMandorSummariesWithDateRange(namaMandor, tanggalMulai, tanggalSelesai, tipeFilter string) ([]MandorSummary, error) {
	startDate, endDate, err := parseDateRange(tanggalMulai, tanggalSelesai)
	if err != nil {
		return nil, err
	}

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

		// Query to get totals from all penyadap for this mandor in the date range
		query := config.DB.Model(&models.BakuPenyadap{}).
			Where("id_baku_mandor = ?", mandor.ID)

		// Apply date range filter
		query = buildDateRangeQuery(query, startDate, endDate)

		if tipeFilter != "" {
			query = query.Where("tipe = ?", tipeFilter)
			summary.Tipe = models.TipeProduksi(tipeFilter)
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
func searchPenyadapDetailsWithDateRange(namaPenyadap, tanggalMulai, tanggalSelesai, tipeFilter string) ([]PenyadapDetail, error) {
	startDate, endDate, err := parseDateRange(tanggalMulai, tanggalSelesai)
	if err != nil {
		return nil, err
	}

	var details []PenyadapDetail
	query := config.DB.Table("baku_penyadaps").Select(`
        penyadaps.id,
        penyadaps.nama_penyadap,
        penyadaps.nik,
        baku_mandors.mandor,
        baku_mandors.afdeling,
        baku_penyadaps.tipe,
        COALESCE(SUM(baku_penyadaps.basah_latex), 0) as total_basah_latex,
        COALESCE(SUM(baku_penyadaps.sheet), 0) as total_sheet,
        COALESCE(SUM(baku_penyadaps.basah_lump), 0) as total_basah_lump,
        COALESCE(SUM(baku_penyadaps.br_cr), 0) as total_br_cr,
        COUNT(baku_penyadaps.id) as jumlah_hari_kerja
    `).
		Joins("JOIN penyadaps ON penyadaps.id = baku_penyadaps.id_penyadap").
		Joins("JOIN baku_mandors ON baku_penyadaps.id_baku_mandor = baku_mandors.id").
		Where("penyadaps.nama_penyadap LIKE ?", "%"+namaPenyadap+"%").
		Where("DATE(baku_penyadaps.tanggal) BETWEEN DATE(?) AND DATE(?)", startDate, endDate)

	if tipeFilter != "" {
		query = query.Where("baku_penyadaps.tipe = ?", tipeFilter)
	}

	if err := query.Group("penyadaps.id, penyadaps.nama_penyadap, penyadaps.nik, baku_mandors.mandor, baku_mandors.afdeling, baku_penyadaps.tipe").
		Scan(&details).Error; err != nil {
		return nil, err
	}

	return details, nil
}
