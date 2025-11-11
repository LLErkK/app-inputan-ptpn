package seed

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"fmt"
	"log"
)

// SeedPetaData - Function untuk seed data peta
// Cara pakai: panggil function ini dari main.go atau buat file terpisah untuk run seeder
func SeedPetaData() {
	db := config.GetDB()

	var count int64
	if err := db.Model(&models.Peta{}).Count(&count).Error; err != nil {
		log.Fatal("Error saat mengecek data peta:", err)
	}

	if count > 0 {
		fmt.Printf("⚠️  Table peta sudah memiliki %d data. Seeding dibatalkan.\n", count)
		fmt.Println("Hapus data terlebih dahulu jika ingin menjalankan seeder ulang.")
		return
	}

	fmt.Println("✓ Table peta kosong, melanjutkan seeding...\n")

	// Data peta yang akan diinput
	petas := []models.Peta{
		// Afdeling Gebugan (1-53)
		{Code: "EMPLSEMEN", Afdeling: "Gebugan"},
		{Code: "FM-IE10-03-AR0020", Afdeling: "Gebugan"},
		{Code: "FM-IE10-03-AR0019", Afdeling: "Gebugan", Blok: "Sekendil", TahunTanam: "1999", Luas: 9.83, JumlahPohon: 11800, JenisKebun: "Kopi Arabika"},
		{Code: "FM-IE10-62-RO0016", Afdeling: "Gebugan", Blok: "Wates", TahunTanam: "1999", Luas: 20, JumlahPohon: 16084, JenisKebun: "Kopi Robusta"},
		{Code: "FM-IE10-03-AR0021", Afdeling: "Gebugan", Blok: "Semangun", TahunTanam: "1999", Luas: 12.17, JumlahPohon: 42483, JenisKebun: "Kopi Arabika"},
		{Code: "KM-IE10-14-GEB001", Afdeling: "Gebugan", Blok: "Sekandri", TahunTanam: "2008", Luas: 61.45, JumlahPohon: 23860, JenisKebun: "Karet"},
		{Code: "FM-IE10-75-RO0013", Afdeling: "Gebugan", Blok: "Lempuyang", TahunTanam: "1971", Luas: 18, JumlahPohon: 21403, JenisKebun: "Kopi Robusta"},
		{Code: "FM-IE10-14-RO0014", Afdeling: "Gebugan", Blok: "Sedandang A"},
		{Code: "FM-IE10-16-RO0015", Afdeling: "Gebugan", Blok: "Sedandang B", TahunTanam: "2012", JumlahPohon: 5482, JenisKebun: "Kopi Robusta"},
		{Code: "FI-IE10-19-AR0023", Afdeling: "Gebugan"},
		{Code: "FM-IE10-16-RO0017", Afdeling: "Gebugan", Blok: "Sebanteng", TahunTanam: "2012", Luas: 8.75, JumlahPohon: 8119},
		{Code: "FM-IE10-21-AR024A", Afdeling: "Gebugan", Blok: "WARUDOYONG", TahunTanam: "2017", Luas: 10},
		{Code: "YP-IE10-19-GESR19", Afdeling: "Gebugan", Blok: "Sejati Suren", Luas: 15.25},
		{Code: "R1-IE10-17-GBPL25", Afdeling: "Gebugan", Blok: "Lemahbang", Luas: 10.2},
		{Code: "R1-IE10-17-GBPL24", Afdeling: "Gebugan", Blok: "Kandri", TahunTanam: "2017", Luas: 22.6, JenisKebun: "Pala"},
		{Code: "FM-IE10-14-RO0006", Afdeling: "Gebugan", Blok: "Sepayung A", TahunTanam: "2010", Luas: 11, JenisKebun: "Kopi Robusta"},
		{Code: "FM-IE10-93-RO0007", Afdeling: "Gebugan", Blok: "Sepayung B", TahunTanam: "1989", Luas: 6, JenisKebun: "Kopi Robusta"},
		{Code: "R1-IE10-80-GBPL13", Afdeling: "Gebugan", Blok: "Wagiri A", TahunTanam: "1968", Luas: 8, JenisKebun: "Pala"},
		{Code: "R1-IE10-80-GBPL06", Afdeling: "Gebugan", Blok: "Sepayug D", TahunTanam: "1915", Luas: 13, JenisKebun: "Pala"},
		{Code: "EKS jabon 3 a", Afdeling: "Gebugan"},
		{Code: "R1-IE10-17-GBPL27", Afdeling: "Gebugan", Blok: "Gerbetung", TahunTanam: "2017", Luas: 16, JenisKebun: "Pala"},
		{Code: "R1-IE10-80-GBPL05", Afdeling: "Gebugan", Blok: "Sepayung C", Luas: 7},
		{Code: "R1-IE10-80-GBPL03", Afdeling: "Gebugan", Blok: "Gogik", TahunTanam: "1915", Luas: 10, JenisKebun: "Pala"},
		{Code: "R1-IE10-80-GBPL16", Afdeling: "Gebugan", Blok: "Gogik A", TahunTanam: "1969", Luas: 10, JenisKebun: "Pala"},
		{Code: "FM-IE10-33-RO0008", Afdeling: "Gebugan", Blok: "Wagiri", TahunTanam: "1929", Luas: 3, JenisKebun: "Kopi Robusta"},
		{Code: "R1-IE10-80-GBPL18", Afdeling: "Gebugan", Blok: "Wagiri B", TahunTanam: "1970", Luas: 2, JenisKebun: "Pala"},
		{Code: "R1-IE10-80-GBPL14", Afdeling: "Gebugan", Blok: "Wagiri C", TahunTanam: "1968", Luas: 5, JenisKebun: "Pala"},
		{Code: "FM-IE10-93-RO0001", Afdeling: "Gebugan", Blok: "Masiran", TahunTanam: "1989", Luas: 2, JenisKebun: "Kopi Robusta"},
		{Code: "FM-IE10-58-RO0002", Afdeling: "Gebugan", Blok: "Tegalrejo", TahunTanam: "1954", Luas: 1, JenisKebun: "Kopi Robusta"},
		{Code: "FM-IE10-15-RO0009", Afdeling: "Gebugan", Blok: "Sebulus A", TahunTanam: "2011", Luas: 4, JenisKebun: "Kopi Robusta"},
		{Code: "R1-IE10-80-GBPL02", Afdeling: "Gebugan", Blok: "Segenting B", TahunTanam: "1913", Luas: 7, JenisKebun: "Pala"},
		{Code: "R1-IE10-80-GBPL23", Afdeling: "Gebugan", Blok: "Sebulus A", TahunTanam: "1974", Luas: 2.25, JenisKebun: "Pala"},
		{Code: "R1-IE10-80-GBPL05-2", Afdeling: "Gebugan", Blok: "Sepayung C", TahunTanam: "1915", Luas: 7, JenisKebun: "Pala"},
		{Code: "geb warudoyong", Afdeling: "Gebugan"},
		{Code: "FM-IE10-62-RO0016-2", Afdeling: "Gebugan", Blok: "Seproamng", Luas: 18},
		{Code: "FM-IE10-62-RO0011", Afdeling: "Gebugan", Blok: "Kemploko B", TahunTanam: "1958", Luas: 7, JenisKebun: "Kopi Robusta"},
		{Code: "FM-IE10-15-RO0010", Afdeling: "Gebugan", Blok: "Kemloko A", TahunTanam: "2011", Luas: 4},
		{Code: "FM-IE10-93-RO0012", Afdeling: "Gebugan", Blok: "Kemploko C", TahunTanam: "1989", Luas: 1, JenisKebun: "Kopi Robusta"},
		{Code: "FM-IE10-14-RO0004", Afdeling: "Gebugan", Blok: "Segadung A", TahunTanam: "2010", Luas: 9, JenisKebun: "Kopi Robusta"},
		{Code: "FM-IE10-15-RO0025", Afdeling: "Gebugan", Blok: "Warurejo A", TahunTanam: "2011", Luas: 3},
		{Code: "FM-IE10-59-RO0003", Afdeling: "Gebugan", Blok: "Segadung B", TahunTanam: "1955", Luas: 7, JenisKebun: "Kopi Robusta"},
		{Code: "FM-IE10-14-RO0005", Afdeling: "Gebugan", Blok: "Segenting", TahunTanam: "2010", Luas: 5, JenisKebun: "Kopi Robusta"},
		{Code: "R1-IE10-80-GBPL19", Afdeling: "Gebugan", Blok: "Segenting A", TahunTanam: "1971", Luas: 11.5, JenisKebun: "Pala"},
		{Code: "R1-IE10-80-GBPL04", Afdeling: "Gebugan", Blok: "Gintungan B", TahunTanam: "1915", Luas: 18, JenisKebun: "Pala"},
		{Code: "R1-IE10-80-GBPL07", Afdeling: "Gebugan", Blok: "Warurejo B", TahunTanam: "1915", Luas: 7, JenisKebun: "Pala"},
		{Code: "R1-IE10-80-GBPL08", Afdeling: "Gebugan", Blok: "Warurejo C", TahunTanam: "1915", Luas: 8, JenisKebun: "Pala"},
		{Code: "R1-IE10-80-GBPL09", Afdeling: "Gebugan", Blok: "Senanas", TahunTanam: "1915", Luas: 23, JenisKebun: "Pala"},
		{Code: "R1-IE10-80-GBPL01", Afdeling: "Gebugan", Blok: "Tegalrejo B", TahunTanam: "1913", Luas: 6, JenisKebun: "Pala"},
		{Code: "R1-IE10-80-GBPL11", Afdeling: "Gebugan", Blok: "Sebulus B", TahunTanam: "1966", Luas: 5, JenisKebun: "Pala"},
		{Code: "R1-IE10-80-GBPL10", Afdeling: "Gebugan", Blok: "Sebulus C", TahunTanam: "1929", Luas: 13, JenisKebun: "Pala"},
		{Code: "R1-IE10-80-GBPL20", Afdeling: "Gebugan", Blok: "Tegalrejo C", TahunTanam: "1971", Luas: 0.5, JenisKebun: "Pala"},
		{Code: "R1-IE10-80-GBPL22", Afdeling: "Gebugan", Blok: "Warudoyong B", TahunTanam: "1972", Luas: 1, JenisKebun: "Pala"},
		{Code: "R1-IE10-80-GBPL15", Afdeling: "Gebugan", Blok: "Segadung A", TahunTanam: "1998", Luas: 1.5},

		// Afdeling Setro (54-75)
		{Code: "KM-IE10-10-STR001", Afdeling: "Setro", Blok: "Siwalan", TahunTanam: "2004", Luas: 29.90, JumlahPohon: 16508, JenisKebun: "Karet", Kloon: "Polycloon"},
		{Code: "KM-IE10-11-STR002", Afdeling: "Setro", Blok: "Manggihan", TahunTanam: "2005", Luas: 18.00, JumlahPohon: 5677, JenisKebun: "Karet", Kloon: "Polycloon"},
		{Code: "YP-IE10-18-STMP01", Afdeling: "Setro", Blok: "Kambangan", TahunTanam: "2018", Luas: 18.00, JumlahPohon: 278, JenisKebun: "TDP Miopsis", Kloon: "Miopsis"},
		{Code: "KI-IE10-19-STR017", Afdeling: "Setro", Blok: "Jenggleng", TahunTanam: "2018", Luas: 18.00, JumlahPohon: 9547, JenisKebun: "Karet", Kloon: "GT 1,IRR 118"},
		{Code: "KM-IE10-19-STR016", Afdeling: "Setro", Blok: "Genurit", TahunTanam: "2013", Luas: 20.64, JumlahPohon: 9594, JenisKebun: "Karet", Kloon: "IRR 118"},
		{Code: "KL-IE10-98-STR14A", Afdeling: "Setro", Blok: "Watututup", TahunTanam: "1998", Luas: 53.02, JumlahPohon: 17586, JenisKebun: "Karet", Kloon: "BPM 1,RRIC 110,PB 235,RRIM 712,PR 300,RRIM 600,CAMPURAN"},
		{Code: "KI-IE10-20-STR008", Afdeling: "Setro", Blok: "Mendiro", TahunTanam: "2020", Luas: 38.66, JumlahPohon: 20731, JenisKebun: "Karet", Kloon: "GT 1,IRR 118"},
		{Code: "KM-IE10-18-STR006", Afdeling: "Setro", Blok: "Rempong", TahunTanam: "2012", Luas: 65.66, JumlahPohon: 33707, JenisKebun: "Karet", Kloon: "PB260,IRR 118"},
		{Code: "KM-IE10-16-STR007", Afdeling: "Setro", Blok: "Ngaglik", TahunTanam: "2010", Luas: 39.45, JumlahPohon: 23129, JenisKebun: "Karet", Kloon: "BPM 1"},
		{Code: "KM-IE10-15-STR005", Afdeling: "Setro", Blok: "Bulu", TahunTanam: "2009", Luas: 29.18, JumlahPohon: 19662, JenisKebun: "Karet", Kloon: "BPM 1,BPM 24"},
		{Code: "KM-IE10-11-STR013", Afdeling: "Setro", Blok: "Setro", TahunTanam: "2005", Luas: 37.60, JumlahPohon: 22516, JenisKebun: "Karet", Kloon: "Polycloon"},
		{Code: "KM-IE10-12-STR012", Afdeling: "Setro", Blok: "Klesem", TahunTanam: "2006", Luas: 72.82, JumlahPohon: 41284, JenisKebun: "Karet", Kloon: "RRIC 110,BPM 1"},
		{Code: "KM-IE10-19-STR011", Afdeling: "Setro", Blok: "Tempel", TahunTanam: "2013", Luas: 4.87, JumlahPohon: 2308, JenisKebun: "Karet", Kloon: "PB 260"},
		{Code: "KM-IE10-08-STR009", Afdeling: "Setro", Blok: "Gondoriyo", TahunTanam: "2002", Luas: 28.74, JumlahPohon: 14479, JenisKebun: "Karet", Kloon: "RRIC 110,BPM 1,BPM 24"},
		{Code: "KM-IE10-14-STR004", Afdeling: "Setro", Blok: "Watugajah", TahunTanam: "2008", Luas: 42.76, JumlahPohon: 26966, JenisKebun: "Karet", Kloon: "BPM1,BPM 24 ,PB 260"},
		{Code: "KM-IE10-09-STR010", Afdeling: "Setro", Blok: "Jimbaran", TahunTanam: "2003", Luas: 30.83, JumlahPohon: 15682, JenisKebun: "Karet", Kloon: "BPM 1,BPM 24"},
		{Code: "KM-IE10-17-STR003", Afdeling: "Setro", Blok: "Kalikopeng", TahunTanam: "2011", Luas: 60.17, JumlahPohon: 33109, JenisKebun: "Karet", Kloon: "PB 260, BPM 1"},
		{Code: "KN-IE10-ENTR-ST01", Afdeling: "Setro", Blok: "Entrys Bulu", Luas: 2.82, JenisKebun: "Entrys"},
		{Code: "EMPLASEMENT", Afdeling: "Setro"},
		{Code: "KANTOR AFD SETRO", Afdeling: "Setro"},
		{Code: "PABRIK RSS NGOBO", Afdeling: "Setro"},
		{Code: "EMPLASEMEN GONDORIYO AFD SETRO", Afdeling: "Setro"},

		// Afdeling Jatirunggo (76-99)
		{Code: "KM-IE10-16-JR0012", Afdeling: "Jatirunggo", Blok: "Wonorejo", TahunTanam: "2010", Luas: 104.10, JumlahPohon: 53695, JenisKebun: "Karet", Kloon: "BPM 1"},
		{Code: "RY-IE10-17-JRAK02", Afdeling: "Jatirunggo", Blok: "Akasia Intercrop 1 Jatirunggo 2017", TahunTanam: "-12", JumlahPohon: 2723, JenisKebun: "TDP Akasia Intercrop", Kloon: "Akasia"},
		{Code: "RY-IE10-15-JRSG17", Afdeling: "Jatirunggo", Blok: "Sengon Intercrop 2 Jatirunggo 2015", TahunTanam: "-7,5", JenisKebun: "TDP Sengon Intercrop (terjual)", Kloon: "Sengon"},
		{Code: "TANAH PENGHIJAUAN BLOK BLIMBING EX 1997", Afdeling: "Jatirunggo"},
		{Code: "AGRO DAN IMPLASEMENT", Afdeling: "Jatirunggo"},
		{Code: "EXS OKUPASI KALISALAK 41,91 Ha", Afdeling: "Jatirunggo"},
		{Code: "TANAH OKUPASI TEGALREJO", Afdeling: "Jatirunggo"},
		{Code: "R2-IE10-19-JTSI01", Afdeling: "Jatirunggo", Blok: "Serai Intercrop Jati Runggu", TahunTanam: "2019", Luas: 5, JenisKebun: "Serai Wangi"},
		{Code: "R2-IE10-19-JTSI01-2", Afdeling: "Jatirunggo"},
		{Code: "R2-IE10-19-JTSI01-3", Afdeling: "Jatirunggo"},
		{Code: "KM-IE10-18-JR0013", Afdeling: "Jatirunggo", Blok: "Barat", TahunTanam: "2012", Luas: 55.24, JumlahPohon: 24533, JenisKebun: "Karet", Kloon: "PB 260"},
		{Code: "KM-IE10-13-JR0011", Afdeling: "Jatirunggo", Blok: "Bubak", TahunTanam: "2007", Luas: 55.00, JumlahPohon: 23095, JenisKebun: "Karet", Kloon: "Polykloon"},
		{Code: "KM-IE10-15-JR0010", Afdeling: "Jatirunggo", Blok: "Rejosari", TahunTanam: "2009", Luas: 49.97, JumlahPohon: 15233, JenisKebun: "Karet", Kloon: "BPM 1"},
		{Code: "KL-IE10-19-JR0007", Afdeling: "Jatirunggo", Blok: "Geyongan", TahunTanam: "1999", Luas: 15.64, JenisKebun: "Karet"},
		{Code: "KM-IE10-05-JR0008", Afdeling: "Jatirunggo", Blok: "Kali Salak", TahunTanam: "ex 1999", Luas: 20.65},
		{Code: "KM-IE10-14-JR0009", Afdeling: "Jatirunggo", Blok: "Jatikurung", TahunTanam: "2008", Luas: 51.50, JumlahPohon: 19514, JenisKebun: "Karet", Kloon: "BPM 1"},
		{Code: "KM-IE10-17-JR0005", Afdeling: "Jatirunggo", Blok: "Tugusari", TahunTanam: "2011", Luas: 13.00, JumlahPohon: 8202, JenisKebun: "Karet", Kloon: "IRR 118"},
		{Code: "KM-IE10-17-JR0003", Afdeling: "Jatirunggo", Blok: "Sajen", TahunTanam: "2011", Luas: 33.11, JumlahPohon: 13384, JenisKebun: "Karet", Kloon: "BPM 1"},
		{Code: "KM-IE10-17-JR0014", Afdeling: "Jatirunggo", Blok: "Soko", TahunTanam: "2011", Luas: 24.00, JumlahPohon: 14179, JenisKebun: "Karet", Kloon: "BPM 1"},
		{Code: "CI-IE10-17-JR0001", Afdeling: "Jatirunggo", Blok: "Kalikaseh", TahunTanam: "2017", Luas: 55.00, JumlahPohon: 2533, JenisKebun: "Kebun Koleksi Kakao", Kloon: "DR 1"},
		{Code: "CI-IE10-17-JR0002", Afdeling: "Jatirunggo", Blok: "SAMINAN", TahunTanam: "2023", Luas: 29.86, JumlahPohon: 16575, JenisKebun: "TTI KARET", Kloon: "GT 1"},
		{Code: "YP-IE10-18-JRSG20", Afdeling: "Jatirunggo", Blok: "Tegal Rejo", TahunTanam: "2018", Luas: 18.00, JumlahPohon: 3369, JenisKebun: "TDP Sengon", Kloon: "Sengon"},
		{Code: "YP-IE10-18-JRSG19", Afdeling: "Jatirunggo", Blok: "Mendoh", TahunTanam: "2018", Luas: 30.10, JumlahPohon: 6088, JenisKebun: "TDP Sengon", Kloon: "Sengon"},
		{Code: "JR URUGAN TOL", Afdeling: "Jatirunggo"},

		// Afdeling Klepu (100-115)
		{Code: "KM-IE10-19-KLP001", Afdeling: "Klepu", Blok: "Ngobo", TahunTanam: "2013", Luas: 10.98, JumlahPohon: 4562, JenisKebun: "Karet", Kloon: "IRR 118"},
		{Code: "KM-IE10-06-KLP002", Afdeling: "Klepu", Blok: "Watu Tumpeng", TahunTanam: "ex 2000", Luas: 26.2, JenisKebun: "Karet"},
		{Code: "YP-IE10-17-KLSG05", Afdeling: "Klepu", Blok: "Randualas", TahunTanam: "2017", Luas: 29.87, JumlahPohon: 18901, JenisKebun: "DP Sengon (Sidah Terjual)", Kloon: "Sengon"},
		{Code: "KM-IE10-17-KLP004", Afdeling: "Klepu", Blok: "Tunon", TahunTanam: "2011", Luas: 41.76, JumlahPohon: 25531, JenisKebun: "Karet", Kloon: "BPM 1"},
		{Code: "KM-IE10-14-KLP007", Afdeling: "Klepu", Blok: "Jangkang", TahunTanam: "2008", Luas: 24.68, JumlahPohon: 13644, JenisKebun: "Karet", Kloon: "PB 260"},
		{Code: "KM-IE10-13-KLP008", Afdeling: "Klepu", Blok: "Bodean", TahunTanam: "2007", Luas: 26.75, JumlahPohon: 12137, JenisKebun: "Karet", Kloon: "BPM 1,PB 260"},
		{Code: "KM-IE10-12-KLP003", Afdeling: "Klepu", Blok: "Pluang", TahunTanam: "2006", Luas: 49.29, JumlahPohon: 22757, JenisKebun: "Karet", Kloon: "BPM 1,PB 260"},
		{Code: "KM-IE10-05-KLP009", Afdeling: "Klepu", Blok: "Kaliulo", TahunTanam: "2022", Luas: 31.74, JumlahPohon: 17545, JenisKebun: "Karet", Kloon: "Polykloon"},
		{Code: "KM-IE10-19-KLP010", Afdeling: "Klepu", Blok: "Watu Kursi", TahunTanam: "2013", Luas: 30.03, JumlahPohon: 17634, JenisKebun: "Karet", Kloon: "IRR 118"},
		{Code: "YP-IE10-17-KLSG06", Afdeling: "Klepu", Blok: "Watu Bangkong", TahunTanam: "2017", Luas: 13.77, JenisKebun: "TDP Sengon (Sudah Terjual)", Kloon: "Sengon"},
		{Code: "YP-IE10-17-KLSG11", Afdeling: "Klepu", Blok: "Gondang", TahunTanam: "2017", Luas: 27.35, JumlahPohon: 19573, JenisKebun: "TDP Sengon (Sudah Terjual)", Kloon: "Sengon"},
		{Code: "EMPLASEMENT KLEPU", Afdeling: "Klepu"},
		{Code: "KANTOR AFD KLEPU", Afdeling: "Klepu"},
		{Code: "ENTRYS KLEPU", Afdeling: "Klepu"},
		{Code: "EMPELASEMEN GAJIHAN BARAT JATIRUNGGO", Afdeling: "Klepu"},
		{Code: "EMPLASMENT PLUWANG AFD. KLEPU", Afdeling: "Klepu"},
	}

	// Insert data ke database
	fmt.Println("Mulai seeding data peta...")

	for i, peta := range petas {
		if err := db.Create(&peta).Error; err != nil {
			log.Printf("Error insert data ke-%d (Code: %s): %v", i+1, peta.Code, err)
		} else {
			fmt.Printf("✓ Berhasil insert data ke-%d: %s - %s\n", i+1, peta.Code, peta.Afdeling)
		}
	}

	fmt.Println("\nSeeding selesai! Total data:", len(petas))
}
