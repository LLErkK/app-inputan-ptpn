package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetBakuPenyadapToday(w http.ResponseWriter, r *http.Request) {
	// Ambil tanggal hari ini (format YYYY-MM-DD)
	today := time.Now().Format("2006-01-02")
	tanggal, err := time.Parse("2006-01-02", today)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal memproses tanggal hari ini: " + err.Error(),
		})
		return
	}

	// Ambil data penyadap hanya untuk tanggal hari ini
	var penyadap []models.BakuPenyadap
	query := config.DB.Preload("Mandor").Preload("Penyadap").
		Where("DATE(tanggal) = DATE(?)", tanggal).
		Order("created_at desc")

	if err := query.Find(&penyadap).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data penyadap: " + err.Error(),
		})
		return
	}

	// Kirim response dengan struktur sama seperti GetAllBakuPenyadap
	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data penyadap untuk tanggal " + today + " berhasil diambil",
		Data:    penyadap,
	})
}

type MonitoringSearchRequest struct {
	NamaMandor         string `json:"namaMandor"`
	NamaPenyadap       string `json:"namaPenyadap"`
	FilterTanggalAwal  string `json:"filterTanggalAwal"`
	FilterTanggalAkhir string `json:"filterTanggalAkhir"`
	FilterJenis        string `json:"filterJenis"`
}

// MonitoringSearchResponse represents the unified search response
type MonitoringSearchResponse struct {
	Success    bool                   `json:"success"`
	Message    string                 `json:"message"`
	Data       []MonitoringSearchItem `json:"data"`
	SearchInfo MonitoringSearchInfo   `json:"searchInfo"`
}

// MonitoringSearchItem represents a unified data item for monitoring display
type MonitoringSearchItem struct {
	ID           uint    `json:"id"`
	Tanggal      string  `json:"tanggal"`
	TahunTanam   uint    `json:"tahunTanam"`
	Mandor       string  `json:"mandor"`
	Afdeling     string  `json:"afdeling"`
	NIK          string  `json:"nik"`
	NamaPenyadap string  `json:"namaPenyadap"`
	Periode      string  `json:"periode"`
	Tipe         string  `json:"tipe"`
	BasahLatex   float64 `json:"basahLatex"`
	Sheet        float64 `json:"sheet"`
	BasahLump    float64 `json:"basahLump"`
	BrCr         float64 `json:"brCr"`
	Source       string  `json:"source"` // "mandor_summary", "penyadap_detail", "baku_penyadap"
}

// MonitoringSearchInfo provides search context information
type MonitoringSearchInfo struct {
	SearchType    string `json:"searchType"`    // "mandor", "penyadap", "both", "general"
	UsedAPI       string `json:"usedAPI"`       // API endpoint that was used
	FilterApplied string `json:"filterApplied"` // Description of filters applied
	TotalRecords  int    `json:"totalRecords"`
	DateRange     string `json:"dateRange"`
}

// SmartMonitoringSearch - Main function that determines the best search strategy

// SearchStrategy defines which API to use and how
type SearchStrategy struct {
	SearchType string // "mandor_only", "penyadap_only", "both_specific", "general_search", "date_filter"
	UsedAPI    string // The API endpoint to use
	Method     string // "single", "range", "detail", "all"
}

// determineSearchStrategy analyzes input and decides the best search approach
func determineSearchStrategy(mandor, penyadap, tanggalAwal, tanggalAkhir, tipe string) SearchStrategy {

	// Priority 1: Both mandor and penyadap specified - use global search
	if mandor != "" && penyadap != "" {
		if tanggalAwal != "" && tanggalAkhir != "" {
			return SearchStrategy{
				SearchType: "both_with_date_range",
				UsedAPI:    "/api/search/all",
				Method:     "range",
			}
		}
		return SearchStrategy{
			SearchType: "both_specific",
			UsedAPI:    "/api/search/all",
			Method:     "single",
		}
	}

	// Priority 2: Only mandor specified
	if mandor != "" && penyadap == "" {
		if tanggalAwal != "" && tanggalAkhir != "" {
			return SearchStrategy{
				SearchType: "mandor_with_date_range",
				UsedAPI:    "/api/search/mandor/range",
				Method:     "range",
			}
		} else if tanggalAwal != "" {
			return SearchStrategy{
				SearchType: "mandor_with_single_date",
				UsedAPI:    "/api/search/mandor",
				Method:     "single",
			}
		}
		return SearchStrategy{
			SearchType: "mandor_only",
			UsedAPI:    "/api/search/mandor/detail",
			Method:     "detail",
		}
	}

	// Priority 3: Only penyadap specified
	if penyadap != "" && mandor == "" {
		if tanggalAwal != "" && tanggalAkhir != "" {
			return SearchStrategy{
				SearchType: "penyadap_with_date_range",
				UsedAPI:    "/api/search/penyadap/range",
				Method:     "range",
			}
		} else if tanggalAwal != "" {
			return SearchStrategy{
				SearchType: "penyadap_with_single_date",
				UsedAPI:    "/api/search/penyadap",
				Method:     "single",
			}
		}
		return SearchStrategy{
			SearchType: "penyadap_only",
			UsedAPI:    "/api/search/penyadap",
			Method:     "all",
		}
	}

	// Priority 4: Only date filters (no name search)
	if (tanggalAwal != "" || tanggalAkhir != "") && mandor == "" && penyadap == "" {
		if tanggalAwal != "" && tanggalAkhir != "" {
			return SearchStrategy{
				SearchType: "date_range_only",
				UsedAPI:    "/api/reporting/penyadap/range",
				Method:     "range",
			}
		}
		return SearchStrategy{
			SearchType: "single_date_only",
			UsedAPI:    "/api/reporting/penyadap",
			Method:     "date",
		}
	}

	// Priority 5: Only tipe filter
	if tipe != "" && mandor == "" && penyadap == "" && tanggalAwal == "" {
		return SearchStrategy{
			SearchType: "tipe_filter_only",
			UsedAPI:    "/api/reporting/penyadap",
			Method:     "tipe",
		}
	}

	// Default: Get all recent data
	return SearchStrategy{
		SearchType: "default_all",
		UsedAPI:    "/api/baku/rekap/today",
		Method:     "today",
	}
}

// executeSearch performs the actual search based on the determined strategy
func executeSearch(strategy SearchStrategy, mandor, penyadap, tanggalAwal, tanggalAkhir, tipe string) ([]MonitoringSearchItem, error) {
	switch strategy.SearchType {

	case "both_specific", "both_with_date_range":
		return executeGlobalSearch(mandor, penyadap, tanggalAwal, tanggalAkhir, tipe)

	case "mandor_only", "mandor_with_single_date", "mandor_with_date_range":
		return executeMandorSearch(mandor, tanggalAwal, tanggalAkhir, tipe, strategy.Method)

	case "penyadap_only", "penyadap_with_single_date", "penyadap_with_date_range":
		return executePenyadapSearch(penyadap, tanggalAwal, tanggalAkhir, tipe, strategy.Method)

	case "date_range_only", "single_date_only":
		return executeDateOnlySearch(tanggalAwal, tanggalAkhir, tipe, strategy.Method)

	case "tipe_filter_only":
		return executeTipeOnlySearch(tipe)

	case "default_all":
		return executeDefaultSearch()

	default:
		return executeDefaultSearch()
	}
}

// executeGlobalSearch uses the global search API
func executeGlobalSearch(mandor, penyadap, tanggalAwal, tanggalAkhir, tipe string) ([]MonitoringSearchItem, error) {
	params := url.Values{}

	if mandor != "" {
		params.Add("namaMandor", mandor) // <-- penting
		params.Add("type", "mandor")
	}
	if penyadap != "" {
		params.Add("namaPenyadap", penyadap) // <-- penting
		params.Add("type", "penyadap")
	}

	if mandor != "" && penyadap != "" {
		params.Set("type", "both") // <-- tandai kalau kombinasi
	}

	if tanggalAwal != "" && tanggalAkhir != "" {
		params.Add("tanggal_mulai", tanggalAwal)
		params.Add("tanggal_selesai", tanggalAkhir)
	} else if tanggalAwal != "" {
		params.Add("tanggal", tanggalAwal)
	}

	if tipe != "" {
		params.Add("tipe", tipe)
	}

	return callSearchAllAPI(params)
}

// executeMandorSearch performs mandor-specific search
func executeMandorSearch(mandor, tanggalAwal, tanggalAkhir, tipe, method string) ([]MonitoringSearchItem, error) {
	var summaries []MandorSummary
	var err error

	switch method {
	case "range":
		summaries, err = searchMandorSummariesWithDateRange(mandor, tanggalAwal, tanggalAkhir, tipe)
	case "single":
		summaries, err = searchMandorSummaries(mandor, tanggalAwal, tipe)
	case "detail":
		summaries, err = searchMandorWithDetails(mandor, tanggalAwal, tipe)
	default:
		summaries, err = searchMandorSummaries(mandor, "", tipe)
	}

	if err != nil {
		return nil, err
	}

	return convertMandorSummariesToMonitoringItems(summaries), nil
}

// executePenyadapSearch performs penyadap-specific search
func executePenyadapSearch(penyadap, tanggalAwal, tanggalAkhir, tipe, method string) ([]MonitoringSearchItem, error) {
	var details []PenyadapDetail
	var err error

	switch method {
	case "range":
		details, err = searchPenyadapDetailsWithDateRange(penyadap, tanggalAwal, tanggalAkhir, tipe)
	case "single":
		details, err = searchPenyadapDetails(penyadap, tanggalAwal, tipe)
	default:
		details, err = searchPenyadapDetails(penyadap, "", tipe)
	}

	if err != nil {
		return nil, err
	}

	return convertPenyadapDetailsToMonitoringItems(details), nil
}

// executeDateOnlySearch performs date-only filtering
func executeDateOnlySearch(tanggalAwal, tanggalAkhir, tipe, method string) ([]MonitoringSearchItem, error) {
	var details []PenyadapDetail
	var err error

	if method == "range" {
		details, err = getPenyadapDetailsByDateRange(tanggalAwal, tanggalAkhir, tipe)
	} else {
		details, err = getPenyadapDetails(tanggalAwal, tipe)
	}

	if err != nil {
		return nil, err
	}

	return convertPenyadapDetailsToMonitoringItems(details), nil
}

// executeTipeOnlySearch performs tipe-only filtering
func executeTipeOnlySearch(tipe string) ([]MonitoringSearchItem, error) {
	details, err := getPenyadapDetailsByTipeHarian(tipe)
	if err != nil {
		return nil, err
	}
	return convertPenyadapDetailsToMonitoringItems(details), nil
}

// executeDefaultSearch returns today's data
func executeDefaultSearch() ([]MonitoringSearchItem, error) {
	var penyadap []models.BakuPenyadap
	tanggal := time.Now().Format("2006-01-02")

	query := config.DB.Preload("Mandor").Preload("Penyadap").
		Where("DATE(tanggal) = ?", tanggal).
		Order("created_at desc")

	if err := query.Find(&penyadap).Error; err != nil {
		return nil, err
	}

	return convertBakuPenyadapToMonitoringItems(penyadap), nil
}

// Helper functions for data conversion

func convertMandorSummariesToMonitoringItems(summaries []MandorSummary) []MonitoringSearchItem {
	var items []MonitoringSearchItem

	for _, summary := range summaries {
		// For each mandor summary, create items for each penyadap detail
		for _, detail := range summary.DetailPenyadap {
			items = append(items, MonitoringSearchItem{
				ID:           detail.ID,
				Tanggal:      detail.Tanggal, //seharusnya dari detail
				TahunTanam:   summary.TahunTanam,
				Mandor:       summary.Mandor,
				Afdeling:     summary.Afdeling,
				NIK:          detail.NIK,
				NamaPenyadap: detail.NamaPenyadap,
				Tipe:         string(detail.Tipe),
				BasahLatex:   detail.TotalBasahLatex,
				Sheet:        detail.TotalSheet,
				BasahLump:    detail.TotalBasahLump,
				BrCr:         detail.TotalBrCr,
				Source:       "mandor_summary",
			})
		}

		// If no detail penyadap, create a summary item
		if len(summary.DetailPenyadap) == 0 {
			items = append(items, MonitoringSearchItem{
				ID:         summary.ID,
				Mandor:     summary.Mandor,
				Afdeling:   summary.Afdeling,
				Periode:    string(summary.Tipe),
				Tipe:       string(summary.Tipe),
				BasahLatex: summary.TotalBasahLatex,
				Sheet:      summary.TotalSheet,
				BasahLump:  summary.TotalBasahLump,
				BrCr:       summary.TotalBrCr,
				Source:     "mandor_summary",
			})
		}
	}

	return items
}

func convertPenyadapDetailsToMonitoringItems(details []PenyadapDetail) []MonitoringSearchItem {
	var items []MonitoringSearchItem

	for _, detail := range details {
		items = append(items, MonitoringSearchItem{
			ID:           detail.ID,
			Tanggal:      detail.Tanggal,
			TahunTanam:   detail.TahunTanam,
			Mandor:       detail.Mandor,
			Afdeling:     detail.Afdeling,
			NIK:          detail.NIK,
			NamaPenyadap: detail.NamaPenyadap,
			Periode:      string(detail.Tipe),
			Tipe:         string(detail.Tipe),
			BasahLatex:   detail.TotalBasahLatex,
			Sheet:        detail.TotalSheet,
			BasahLump:    detail.TotalBasahLump,
			BrCr:         detail.TotalBrCr,
			Source:       "penyadap_detail",
		})
	}

	return items
}

// Helper functions for building descriptions

// callSearchAllAPI: implemented to query BakuPenyadap table with filters
func callSearchAllAPI(params url.Values) ([]MonitoringSearchItem, error) {
	var results []models.BakuPenyadap

	query := config.DB.Preload("Mandor").Preload("Penyadap")

	namaMandor := strings.TrimSpace(params.Get("namaMandor"))
	namaPenyadap := strings.TrimSpace(params.Get("namaPenyadap"))
	typ := strings.TrimSpace(params.Get("type"))

	// CASE: both mandor + penyadap
	if namaMandor != "" && namaPenyadap != "" {
		query = query.
			Joins("JOIN baku_mandors bm ON bm.id = baku_penyadaps.id_baku_mandor").
			Joins("JOIN penyadaps p ON p.id = baku_penyadaps.id_penyadap").
			Where("bm.mandor LIKE ? AND p.nama_penyadap LIKE ?", "%"+namaMandor+"%", "%"+namaPenyadap+"%")
	} else if namaMandor != "" && typ == "mandor" {
		query = query.
			Joins("JOIN baku_mandors bm ON bm.id = baku_penyadaps.id_baku_mandor").
			Where("bm.mandor LIKE ?", "%"+namaMandor+"%")
	} else if namaPenyadap != "" {
		query = query.
			Joins("JOIN penyadaps p ON p.id = baku_penyadaps.id_penyadap").
			Where("p.nama_penyadap LIKE ?", "%"+namaPenyadap+"%")
	}

	// tanggal tunggal atau range
	tanggal := strings.TrimSpace(params.Get("tanggal"))
	tanggalMulai := strings.TrimSpace(params.Get("tanggal_mulai"))
	tanggalSelesai := strings.TrimSpace(params.Get("tanggal_selesai"))

	if tanggal != "" {
		query = query.Where("DATE(baku_penyadaps.tanggal) = ?", tanggal)
	} else if tanggalMulai != "" && tanggalSelesai != "" {
		query = query.Where("DATE(baku_penyadaps.tanggal) BETWEEN ? AND ?", tanggalMulai, tanggalSelesai)
	}

	// tipe produksi
	tipe := strings.TrimSpace(params.Get("tipe"))
	if tipe != "" {
		query = query.Where("baku_penyadaps.tipe = ?", tipe)
	}

	query = query.Order("baku_penyadaps.tanggal desc, baku_penyadaps.created_at desc")

	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	return convertBakuPenyadapToMonitoringItems(results), nil
}
func getPenyadapDetailsByTipeHarian(tipe string) ([]PenyadapDetail, error) {
	query := `
        SELECT 
            p.id,
            p.nama_penyadap,
            p.nik,
            bm.mandor,
            bm.afdeling,
            bm.tahun_tanam,
            bp.tipe,
            DATE(bp.tanggal) as tanggal,
            COALESCE(SUM(bp.basah_latex), 0) as total_basah_latex,
            COALESCE(SUM(bp.sheet), 0) as total_sheet,
            COALESCE(SUM(bp.basah_lump), 0) as total_basah_lump,
            COALESCE(SUM(bp.br_cr), 0) as total_br_cr
        FROM baku_penyadaps bp
        LEFT JOIN penyadaps p ON p.id = bp.id_penyadap
        LEFT JOIN baku_mandors bm ON bp.id_baku_mandor = bm.id
        WHERE bp.deleted_at IS NULL
          AND bp.tipe = ?
        GROUP BY p.id, p.nama_penyadap, p.nik, bm.mandor, bm.afdeling, bm.tahun_tanam, DATE(bp.tanggal), bp.tipe
        ORDER BY DATE(bp.tanggal) DESC
    `

	var details []PenyadapDetail
	if err := config.DB.Raw(query, tipe).Scan(&details).Error; err != nil {
		return nil, err
	}

	return details, nil
}

// executeSmartSearch - Updated main search function
func executeSmartSearch(mandor, penyadap, tanggalAwal, tanggalAkhir, tipe string) ([]MonitoringSearchItem, MonitoringSearchInfo, error) {
	var results []MonitoringSearchItem
	var searchType, usedAPI string

	// Build base query
	query := config.DB.Preload("Mandor").Preload("Penyadap")

	// Determine search strategy and apply filters
	if mandor != "" && penyadap == "" {
		// Case 1: Only mandor specified - get baku penyadap data filtered by mandor
		searchType = "mandor_search"
		usedAPI = "/api/baku/penyadap/by-mandor"

		query = query.Joins("JOIN baku_mandors bm ON bm.id = baku_penyadaps.id_baku_mandor").
			Where("bm.mandor LIKE ?", "%"+mandor+"%")

	} else if penyadap != "" && mandor == "" {
		// Case 2: Only penyadap specified - get baku penyadap data filtered by penyadap
		searchType = "penyadap_search"
		usedAPI = "/api/baku/penyadap/by-penyadap"

		query = query.Joins("JOIN penyadaps p ON p.id = baku_penyadaps.id_penyadap").
			Where("p.nama_penyadap LIKE ?", "%"+penyadap+"%")

	} else if mandor != "" && penyadap != "" {
		// Case 3: Both specified - get baku penyadap data filtered by both
		searchType = "both_search"
		usedAPI = "/api/baku/penyadap/by-both"

		query = query.
			Joins("JOIN baku_mandors bm ON bm.id = baku_penyadaps.id_baku_mandor").
			Joins("JOIN penyadaps p ON p.id = baku_penyadaps.id_penyadap").
			Where("bm.mandor LIKE ? AND p.nama_penyadap LIKE ?", "%"+mandor+"%", "%"+penyadap+"%")

	} else {
		// Case 4: No specific name search - get all data (default to current month)
		searchType = "general_search"
		usedAPI = "/api/baku/penyadap/all"

		// If no date filters provided, default to current month
		if tanggalAwal == "" && tanggalAkhir == "" {
			now := time.Now()
			startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			endOfMonth := startOfMonth.AddDate(0, 1, -1)

			tanggalAwal = startOfMonth.Format("2006-01-02")
			tanggalAkhir = endOfMonth.Format("2006-01-02")
			searchType = "current_month_search"
		}
	}

	// Apply date filters
	if tanggalAwal != "" && tanggalAkhir != "" {
		query = query.Where("DATE(baku_penyadaps.tanggal) BETWEEN ? AND ?", tanggalAwal, tanggalAkhir)
	} else if tanggalAwal != "" {
		query = query.Where("DATE(baku_penyadaps.tanggal) >= ?", tanggalAwal)
	} else if tanggalAkhir != "" {
		query = query.Where("DATE(baku_penyadaps.tanggal) <= ?", tanggalAkhir)
	}

	// Apply tipe filter
	if tipe != "" {
		query = query.Where("baku_penyadaps.tipe = ?", tipe)
	}

	// Order by latest date and creation time
	query = query.Order("baku_penyadaps.tanggal DESC, baku_penyadaps.created_at DESC")

	// Execute query
	var bakuPenyadaps []models.BakuPenyadap
	if err := query.Find(&bakuPenyadaps).Error; err != nil {
		return nil, MonitoringSearchInfo{}, err
	}

	// Convert to monitoring items
	results = convertBakuPenyadapToMonitoringItems(bakuPenyadaps)

	// Build search info
	searchInfo := MonitoringSearchInfo{
		SearchType:    searchType,
		UsedAPI:       usedAPI,
		FilterApplied: buildFilterDescription(mandor, penyadap, tanggalAwal, tanggalAkhir, tipe),
		TotalRecords:  len(results),
		DateRange:     buildDateRangeDescription(tanggalAwal, tanggalAkhir),
	}

	return results, searchInfo, nil
}

// Helper function to build filter description - Updated
func buildFilterDescription(mandor, penyadap, tanggalAwal, tanggalAkhir, tipe string) string {
	var filters []string

	if mandor != "" {
		filters = append(filters, "Mandor: "+mandor)
	}
	if penyadap != "" {
		filters = append(filters, "Penyadap: "+penyadap)
	}
	if tanggalAwal != "" && tanggalAkhir != "" {
		filters = append(filters, fmt.Sprintf("Tanggal: %s s/d %s", tanggalAwal, tanggalAkhir))
	} else if tanggalAwal != "" {
		filters = append(filters, "Tanggal mulai: "+tanggalAwal)
	} else if tanggalAkhir != "" {
		filters = append(filters, "Tanggal sampai: "+tanggalAkhir)
	}
	if tipe != "" {
		filters = append(filters, "Tipe: "+tipe)
	}

	if len(filters) == 0 {
		return "Semua data bulan ini"
	}

	return strings.Join(filters, ", ")
}

// Helper function to build date range description - Updated
func buildDateRangeDescription(tanggalAwal, tanggalAkhir string) string {
	if tanggalAwal != "" && tanggalAkhir != "" {
		return fmt.Sprintf("%s sampai %s", tanggalAwal, tanggalAkhir)
	} else if tanggalAwal != "" {
		return "Dari " + tanggalAwal
	} else if tanggalAkhir != "" {
		return "Sampai " + tanggalAkhir
	}

	// Default to current month if no dates provided
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	return fmt.Sprintf("Bulan ini (%s - %s)",
		startOfMonth.Format("2006-01-02"),
		endOfMonth.Format("2006-01-02"))
}

// Additional helper function to get current month data when no filters applied
func getCurrentMonthBakuPenyadap() ([]models.BakuPenyadap, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	var penyadaps []models.BakuPenyadap
	query := config.DB.Preload("Mandor").Preload("Penyadap").
		Where("DATE(tanggal) BETWEEN ? AND ?",
			startOfMonth.Format("2006-01-02"),
			endOfMonth.Format("2006-01-02")).
		Order("tanggal DESC, created_at DESC")

	if err := query.Find(&penyadaps).Error; err != nil {
		return nil, err
	}

	return penyadaps, nil
}

// SmartMonitoringSearch - Handles all 32 combinations of 5 parameters
func SmartMonitoringSearch(w http.ResponseWriter, r *http.Request) {
	// Parse all 5 query parameters
	namaMandor := strings.TrimSpace(r.URL.Query().Get("namaMandor"))
	namaPenyadap := strings.TrimSpace(r.URL.Query().Get("namaPenyadap"))
	tanggalAwal := strings.TrimSpace(r.URL.Query().Get("tanggalAwal"))
	tanggalAkhir := strings.TrimSpace(r.URL.Query().Get("tanggalAkhir"))
	tipe := strings.TrimSpace(r.URL.Query().Get("tipe"))

	// Execute search with all combinations
	results, searchInfo, err := executeSmartSearchAllCombinations(namaMandor, namaPenyadap, tanggalAwal, tanggalAkhir, tipe)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal melakukan pencarian: " + err.Error(),
		})
		return
	}

	response := MonitoringSearchResponse{
		Success:    true,
		Message:    fmt.Sprintf("Pencarian berhasil dengan kombinasi '%s'", searchInfo.SearchType),
		Data:       results,
		SearchInfo: searchInfo,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SearchCombination represents the combination of parameters
type SearchCombination struct {
	ID          int
	Name        string
	Description string
	HasMandor   bool
	HasPenyadap bool
	HasTglAwal  bool
	HasTglAkhir bool
	HasTipe     bool
}

// executeSmartSearchAllCombinations - Handles all 32 combinations
func executeSmartSearchAllCombinations(mandor, penyadap, tanggalAwal, tanggalAkhir, tipe string) ([]MonitoringSearchItem, MonitoringSearchInfo, error) {
	// Determine which parameters are provided
	hasMandor := mandor != ""
	hasPenyadap := penyadap != ""
	hasTglAwal := tanggalAwal != ""
	hasTglAkhir := tanggalAkhir != ""
	hasTipe := tipe != ""

	// Get combination info
	combination := identifyCombination(hasMandor, hasPenyadap, hasTglAwal, hasTglAkhir, hasTipe)

	fmt.Printf("DEBUG: Combination detected: %s (ID: %d)\n", combination.Name, combination.ID)
	fmt.Printf("DEBUG: Parameters - Mandor: %s, Penyadap: %s, TglAwal: %s, TglAkhir: %s, Tipe: %s\n",
		mandor, penyadap, tanggalAwal, tanggalAkhir, tipe)

	// Build base query
	query := config.DB.Preload("Mandor").Preload("Penyadap")

	// Apply filters based on combination
	query = applyFiltersForCombination(query, mandor, penyadap, tanggalAwal, tanggalAkhir, tipe, combination)

	// Execute query
	var bakuPenyadaps []models.BakuPenyadap
	if err := query.Find(&bakuPenyadaps).Error; err != nil {
		return nil, MonitoringSearchInfo{}, err
	}

	// Convert to monitoring items
	results := convertBakuPenyadapToMonitoringItems(bakuPenyadaps)

	// Build search info
	searchInfo := MonitoringSearchInfo{
		SearchType:    combination.Name,
		UsedAPI:       "/api/smart-search",
		FilterApplied: buildFilterDescriptionAllCombinations(mandor, penyadap, tanggalAwal, tanggalAkhir, tipe),
		TotalRecords:  len(results),
		DateRange:     buildDateRangeDescriptionAllCombinations(tanggalAwal, tanggalAkhir),
	}

	return results, searchInfo, nil
}

// identifyCombination - Identifies which of the 32 combinations is being used
func identifyCombination(hasMandor, hasPenyadap, hasTglAwal, hasTglAkhir, hasTipe bool) SearchCombination {
	combinations := []SearchCombination{
		// ID 0: 00000 - No parameters
		{0, "default", "Semua data bulan ini", false, false, false, false, false},

		// ID 1: 00001 - Only Tipe
		{1, "tipe_only", "Filter berdasarkan tipe saja", false, false, false, false, true},

		// ID 2: 00010 - Only TanggalAkhir
		{2, "tanggal_akhir_only", "Filter sampai tanggal akhir", false, false, false, true, false},

		// ID 3: 00011 - TanggalAkhir + Tipe
		{3, "tanggal_akhir_tipe", "Filter sampai tanggal akhir dengan tipe", false, false, false, true, true},

		// ID 4: 00100 - Only TanggalAwal
		{4, "tanggal_awal_only", "Filter dari tanggal awal", false, false, true, false, false},

		// ID 5: 00101 - TanggalAwal + Tipe
		{5, "tanggal_awal_tipe", "Filter dari tanggal awal dengan tipe", false, false, true, false, true},

		// ID 6: 00110 - TanggalAwal + TanggalAkhir (Range)
		{6, "range_tanggal", "Filter range tanggal", false, false, true, true, false},

		// ID 7: 00111 - TanggalAwal + TanggalAkhir + Tipe
		{7, "range_tanggal_tipe", "Filter range tanggal dengan tipe", false, false, true, true, true},

		// ID 8: 01000 - Only Penyadap
		{8, "penyadap_only", "Filter berdasarkan penyadap saja", false, true, false, false, false},

		// ID 9: 01001 - Penyadap + Tipe
		{9, "penyadap_tipe", "Filter penyadap dengan tipe", false, true, false, false, true},

		// ID 10: 01010 - Penyadap + TanggalAkhir
		{10, "penyadap_tanggal_akhir", "Filter penyadap sampai tanggal akhir", false, true, false, true, false},

		// ID 11: 01011 - Penyadap + TanggalAkhir + Tipe
		{11, "penyadap_tanggal_akhir_tipe", "Filter penyadap sampai tanggal akhir dengan tipe", false, true, false, true, true},

		// ID 12: 01100 - Penyadap + TanggalAwal
		{12, "penyadap_tanggal_awal", "Filter penyadap dari tanggal awal", false, true, true, false, false},

		// ID 13: 01101 - Penyadap + TanggalAwal + Tipe
		{13, "penyadap_tanggal_awal_tipe", "Filter penyadap dari tanggal awal dengan tipe", false, true, true, false, true},

		// ID 14: 01110 - Penyadap + Range Tanggal
		{14, "penyadap_range_tanggal", "Filter penyadap dengan range tanggal", false, true, true, true, false},

		// ID 15: 01111 - Penyadap + Range Tanggal + Tipe
		{15, "penyadap_range_tanggal_tipe", "Filter penyadap dengan range tanggal dan tipe", false, true, true, true, true},

		// ID 16: 10000 - Only Mandor
		{16, "mandor_only", "Filter berdasarkan mandor saja", true, false, false, false, false},

		// ID 17: 10001 - Mandor + Tipe
		{17, "mandor_tipe", "Filter mandor dengan tipe", true, false, false, false, true},

		// ID 18: 10010 - Mandor + TanggalAkhir
		{18, "mandor_tanggal_akhir", "Filter mandor sampai tanggal akhir", true, false, false, true, false},

		// ID 19: 10011 - Mandor + TanggalAkhir + Tipe
		{19, "mandor_tanggal_akhir_tipe", "Filter mandor sampai tanggal akhir dengan tipe", true, false, false, true, true},

		// ID 20: 10100 - Mandor + TanggalAwal
		{20, "mandor_tanggal_awal", "Filter mandor dari tanggal awal", true, false, true, false, false},

		// ID 21: 10101 - Mandor + TanggalAwal + Tipe
		{21, "mandor_tanggal_awal_tipe", "Filter mandor dari tanggal awal dengan tipe", true, false, true, false, true},

		// ID 22: 10110 - Mandor + Range Tanggal
		{22, "mandor_range_tanggal", "Filter mandor dengan range tanggal", true, false, true, true, false},

		// ID 23: 10111 - Mandor + Range Tanggal + Tipe
		{23, "mandor_range_tanggal_tipe", "Filter mandor dengan range tanggal dan tipe", true, false, true, true, true},

		// ID 24: 11000 - Mandor + Penyadap
		{24, "mandor_penyadap", "Filter berdasarkan mandor dan penyadap", true, true, false, false, false},

		// ID 25: 11001 - Mandor + Penyadap + Tipe
		{25, "mandor_penyadap_tipe", "Filter mandor dan penyadap dengan tipe", true, true, false, false, true},

		// ID 26: 11010 - Mandor + Penyadap + TanggalAkhir
		{26, "mandor_penyadap_tanggal_akhir", "Filter mandor dan penyadap sampai tanggal akhir", true, true, false, true, false},

		// ID 27: 11011 - Mandor + Penyadap + TanggalAkhir + Tipe
		{27, "mandor_penyadap_tanggal_akhir_tipe", "Filter mandor dan penyadap sampai tanggal akhir dengan tipe", true, true, false, true, true},

		// ID 28: 11100 - Mandor + Penyadap + TanggalAwal
		{28, "mandor_penyadap_tanggal_awal", "Filter mandor dan penyadap dari tanggal awal", true, true, true, false, false},

		// ID 29: 11101 - Mandor + Penyadap + TanggalAwal + Tipe
		{29, "mandor_penyadap_tanggal_awal_tipe", "Filter mandor dan penyadap dari tanggal awal dengan tipe", true, true, true, false, true},

		// ID 30: 11110 - Mandor + Penyadap + Range Tanggal
		{30, "mandor_penyadap_range_tanggal", "Filter mandor dan penyadap dengan range tanggal", true, true, true, true, false},

		// ID 31: 11111 - All parameters
		{31, "all_parameters", "Filter dengan semua parameter", true, true, true, true, true},
	}

	// Calculate combination ID using binary representation
	combinationID := 0
	if hasMandor {
		combinationID += 16 // 2^4
	}
	if hasPenyadap {
		combinationID += 8 // 2^3
	}
	if hasTglAwal {
		combinationID += 4 // 2^2
	}
	if hasTglAkhir {
		combinationID += 2 // 2^1
	}
	if hasTipe {
		combinationID += 1 // 2^0
	}

	return combinations[combinationID]
}

// applyFiltersForCombination - Applies filters based on the combination
func applyFiltersForCombination(query *gorm.DB, mandor, penyadap, tanggalAwal, tanggalAkhir, tipe string, combination SearchCombination) *gorm.DB {
	// Apply Mandor filter
	if combination.HasMandor && mandor != "" {
		query = query.Joins("JOIN baku_mandors bm ON bm.id = baku_penyadaps.id_baku_mandor").
			Where("bm.mandor LIKE ?", "%"+mandor+"%")
	}

	// Apply Penyadap filter
	if combination.HasPenyadap && penyadap != "" {
		if combination.HasMandor {
			// Already joined baku_mandors, now join penyadaps
			query = query.Joins("JOIN penyadaps p ON p.id = baku_penyadaps.id_penyadap").
				Where("p.nama_penyadap LIKE ?", "%"+penyadap+"%")
		} else {
			// Only join penyadaps
			query = query.Joins("JOIN penyadaps p ON p.id = baku_penyadaps.id_penyadap").
				Where("p.nama_penyadap LIKE ?", "%"+penyadap+"%")
		}
	}

	// Apply Date filters
	if combination.HasTglAwal && combination.HasTglAkhir && tanggalAwal != "" && tanggalAkhir != "" {
		// Range date
		query = query.Where("DATE(baku_penyadaps.tanggal) BETWEEN ? AND ?", tanggalAwal, tanggalAkhir)
	} else if combination.HasTglAwal && tanggalAwal != "" {
		// From start date
		query = query.Where("DATE(baku_penyadaps.tanggal) >= ?", tanggalAwal)
	} else if combination.HasTglAkhir && tanggalAkhir != "" {
		// Until end date
		query = query.Where("DATE(baku_penyadaps.tanggal) <= ?", tanggalAkhir)
	} else if combination.ID == 0 {
		// Default case: current month
		now := time.Now()
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endOfMonth := startOfMonth.AddDate(0, 1, -1)
		query = query.Where("DATE(baku_penyadaps.tanggal) BETWEEN ? AND ?",
			startOfMonth.Format("2006-01-02"),
			endOfMonth.Format("2006-01-02"))
	}

	// Apply Tipe filter
	if combination.HasTipe && tipe != "" {
		query = query.Where("baku_penyadaps.tipe = ?", tipe)
	}

	// Always order by latest
	query = query.Order("baku_penyadaps.tanggal DESC, baku_penyadaps.created_at DESC")

	return query
}

// buildFilterDescriptionAllCombinations - Build comprehensive filter description
func buildFilterDescriptionAllCombinations(mandor, penyadap, tanggalAwal, tanggalAkhir, tipe string) string {
	var filters []string

	if mandor != "" {
		filters = append(filters, fmt.Sprintf("Mandor: %s", mandor))
	}
	if penyadap != "" {
		filters = append(filters, fmt.Sprintf("Penyadap: %s", penyadap))
	}
	if tanggalAwal != "" && tanggalAkhir != "" {
		filters = append(filters, fmt.Sprintf("Periode: %s s/d %s", tanggalAwal, tanggalAkhir))
	} else if tanggalAwal != "" {
		filters = append(filters, fmt.Sprintf("Mulai: %s", tanggalAwal))
	} else if tanggalAkhir != "" {
		filters = append(filters, fmt.Sprintf("Sampai: %s", tanggalAkhir))
	}
	if tipe != "" {
		filters = append(filters, fmt.Sprintf("Tipe: %s", tipe))
	}

	if len(filters) == 0 {
		return "Semua data bulan ini"
	}

	return strings.Join(filters, " | ")
}

// buildDateRangeDescriptionAllCombinations - Build date range description
func buildDateRangeDescriptionAllCombinations(tanggalAwal, tanggalAkhir string) string {
	if tanggalAwal != "" && tanggalAkhir != "" {
		return fmt.Sprintf("%s hingga %s", tanggalAwal, tanggalAkhir)
	} else if tanggalAwal != "" {
		return fmt.Sprintf("Dari %s", tanggalAwal)
	} else if tanggalAkhir != "" {
		return fmt.Sprintf("Sampai %s", tanggalAkhir)
	}

	// Default to current month
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	return fmt.Sprintf("Bulan ini (%s - %s)",
		startOfMonth.Format("02 Jan 2006"),
		endOfMonth.Format("02 Jan 2006"))
}

// convertBakuPenyadapToMonitoringItems - Enhanced conversion with better error handling
func convertBakuPenyadapToMonitoringItems(penyadaps []models.BakuPenyadap) []MonitoringSearchItem {
	var items []MonitoringSearchItem

	for _, p := range penyadaps {
		// Safe extraction of mandor data
		mandorName := ""
		afdeling := ""
		tahunTanam := uint(0)

		if p.Mandor.ID != 0 {
			mandorName = p.Mandor.Mandor
			afdeling = p.Mandor.Afdeling
			tahunTanam = p.Mandor.TahunTanam
		}

		// Safe extraction of penyadap data
		penyadapName := ""
		nik := ""
		if p.Penyadap.ID != 0 {
			penyadapName = p.Penyadap.NamaPenyadap
			nik = p.Penyadap.NIK
		}

		items = append(items, MonitoringSearchItem{
			ID:           p.ID,
			Tanggal:      p.Tanggal.Format("2006-01-02"),
			TahunTanam:   tahunTanam,
			Mandor:       mandorName,
			Afdeling:     afdeling,
			NIK:          nik,
			NamaPenyadap: penyadapName,
			Periode:      string(p.Tipe),
			Tipe:         string(p.Tipe),
			BasahLatex:   p.BasahLatex,
			Sheet:        p.Sheet,
			BasahLump:    p.BasahLump,
			BrCr:         p.BrCr,
			Source:       "baku_penyadap",
		})
	}

	return items
}

// GetAllCombinationsInfo - Helper endpoint to see all possible combinations
func GetAllCombinationsInfo(w http.ResponseWriter, r *http.Request) {
	combinations := make([]map[string]interface{}, 32)

	for i := 0; i < 32; i++ {
		hasMandor := (i & 16) != 0
		hasPenyadap := (i & 8) != 0
		hasTglAwal := (i & 4) != 0
		hasTglAkhir := (i & 2) != 0
		hasTipe := (i & 1) != 0

		combination := identifyCombination(hasMandor, hasPenyadap, hasTglAwal, hasTglAkhir, hasTipe)

		// Build example query
		var params []string
		if hasMandor {
			params = append(params, "namaMandor=John")
		}
		if hasPenyadap {
			params = append(params, "namaPenyadap=Smith")
		}
		if hasTglAwal {
			params = append(params, "tanggalAwal=2024-01-01")
		}
		if hasTglAkhir {
			params = append(params, "tanggalAkhir=2024-01-31")
		}
		if hasTipe {
			params = append(params, "tipe=Harian")
		}

		exampleURL := "/api/smart-search"
		if len(params) > 0 {
			exampleURL += "?" + strings.Join(params, "&")
		}

		combinations[i] = map[string]interface{}{
			"id":          combination.ID,
			"name":        combination.Name,
			"description": combination.Description,
			"binary":      fmt.Sprintf("%05b", i),
			"parameters": map[string]bool{
				"mandor":       hasMandor,
				"penyadap":     hasPenyadap,
				"tanggalAwal":  hasTglAwal,
				"tanggalAkhir": hasTglAkhir,
				"tipe":         hasTipe,
			},
			"exampleURL": exampleURL,
		}
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Semua 32 kombinasi pencarian berhasil diambil",
		Data:    combinations,
	})
}

// ServeMonitoringPage - Updated to serve the monitoring page
func ServeMonitoringPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/html/monitoring.html")
}
