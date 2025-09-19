package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetBakuPenyadapToday(w http.ResponseWriter, r *http.Request) {
	// Ambil tanggal hari ini (format YYYY-MM-DD)
	tanggal := time.Now().Format("2006-01-02")

	// Ambil data penyadap hanya untuk tanggal hari ini
	var penyadap []models.BakuPenyadap
	query := config.DB.Preload("Mandor").Preload("Penyadap").
		Where("DATE(tanggal) = ?", tanggal).
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
		Message: "Data penyadap untuk tanggal " + tanggal + " berhasil diambil",
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
func SmartMonitoringSearch(w http.ResponseWriter, r *http.Request) {
	// Parse all possible query parameters
	namaMandor := strings.TrimSpace(r.URL.Query().Get("namaMandor"))
	namaPenyadap := strings.TrimSpace(r.URL.Query().Get("namaPenyadap"))
	tanggalAwal := strings.TrimSpace(r.URL.Query().Get("filterTanggalAwal"))
	tanggalAkhir := strings.TrimSpace(r.URL.Query().Get("filterTanggalAkhir"))
	tipe := strings.TrimSpace(r.URL.Query().Get("filterJenis"))

	// Determine search strategy
	searchStrategy := determineSearchStrategy(namaMandor, namaPenyadap, tanggalAwal, tanggalAkhir, tipe)

	fmt.Printf("DEBUG: Search strategy determined: %s\n", searchStrategy.SearchType)
	fmt.Printf("DEBUG: Will use API: %s\n", searchStrategy.UsedAPI)

	// Execute search based on strategy
	results, err := executeSearch(searchStrategy, namaMandor, namaPenyadap, tanggalAwal, tanggalAkhir, tipe)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal melakukan pencarian: " + err.Error(),
		})
		return
	}

	// Build search info
	searchInfo := MonitoringSearchInfo{
		SearchType:    searchStrategy.SearchType,
		UsedAPI:       searchStrategy.UsedAPI,
		FilterApplied: buildFilterDescription(namaMandor, namaPenyadap, tanggalAwal, tanggalAkhir, tipe),
		TotalRecords:  len(results),
		DateRange:     buildDateRangeDescription(tanggalAwal, tanggalAkhir),
	}

	response := MonitoringSearchResponse{
		Success:    true,
		Message:    fmt.Sprintf("Pencarian berhasil dengan strategi '%s'", searchStrategy.SearchType),
		Data:       results,
		SearchInfo: searchInfo,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

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

func convertBakuPenyadapToMonitoringItems(penyadaps []models.BakuPenyadap) []MonitoringSearchItem {
	var items []MonitoringSearchItem

	for _, p := range penyadaps {
		items = append(items, MonitoringSearchItem{
			ID:           p.ID,
			Tanggal:      p.Tanggal.Format("2006-01-02"),
			TahunTanam:   p.Mandor.TahunTanam,
			Mandor:       p.Mandor.Mandor,
			Afdeling:     p.Mandor.Afdeling,
			NIK:          p.Penyadap.NIK,
			NamaPenyadap: p.Penyadap.NamaPenyadap,
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

// Helper functions for building descriptions
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
		filters = append(filters, "Tanggal: "+tanggalAwal)
	}
	if tipe != "" {
		filters = append(filters, "Tipe: "+tipe)
	}

	if len(filters) == 0 {
		return "Tanpa filter"
	}

	return strings.Join(filters, ", ")
}

func buildDateRangeDescription(tanggalAwal, tanggalAkhir string) string {
	if tanggalAwal != "" && tanggalAkhir != "" {
		return fmt.Sprintf("%s sampai %s", tanggalAwal, tanggalAkhir)
	} else if tanggalAwal != "" {
		return "Dari " + tanggalAwal
	} else if tanggalAkhir != "" {
		return "Sampai " + tanggalAkhir
	}
	return "Hari ini"
}

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
	} else if namaPenyadap != "" && typ == "penyadap" {
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

// ServeMonitoringPage - Updated to serve the monitoring page
func ServeMonitoringPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/html/monitoring.html")
}
