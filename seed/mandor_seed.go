package seed

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
)

func SeedMandor() {

	mandors := []models.BakuMandor{
		//baku
		{TahunTanam: 1998, Mandor: "SUKIYATNO", NIK: "9006569", Afdeling: "SETRO", Tipe: "BAKU"},
		{TahunTanam: 2013, Mandor: "SUKIYATNO", NIK: "9006569", Afdeling: "SETRO", Tipe: "BAKU"},
		{TahunTanam: 2018, Mandor: "SUKIYATNO", NIK: "9006569", Afdeling: "SETRO", Tipe: "BAKU"},
		{TahunTanam: 2004, Mandor: "KAYIN", NIK: "9006286", Afdeling: "SETRO", Tipe: "BAKU"},
		{TahunTanam: 2005, Mandor: "KAYIN", NIK: "9006286", Afdeling: "SETRO", Tipe: "BAKU"},
		{TahunTanam: 2008, Mandor: "MULYANTO", NIK: "9006468", Afdeling: "SETRO", Tipe: "BAKU"},
		{TahunTanam: 2009, Mandor: "DWI JATI NUGROHO", NIK: "9006467", Afdeling: "SETRO", Tipe: "BAKU"},
		{TahunTanam: 2010, Mandor: "SUPRIYADI", NIK: "9006291", Afdeling: "SETRO", Tipe: "BAKU"},
		{TahunTanam: 2011, Mandor: "WIDAYAT", NIK: "9006594", Afdeling: "SETRO", Tipe: "BAKU"},
		{TahunTanam: 2012, Mandor: "HARTONO", NIK: "9006406", Afdeling: "SETRO", Tipe: "BAKU"},
		{TahunTanam: 2005, Mandor: "RUJIANTORO", NIK: "9006777", Afdeling: "SETRO", Tipe: "BAKU"},
		{TahunTanam: 2006, Mandor: "RUJIANTORO", NIK: "9006777", Afdeling: "SETRO", Tipe: "BAKU"},
		{TahunTanam: 2006, Mandor: "JONI ANWAR", NIK: "9006585", Afdeling: "SETRO", Tipe: "BAKU"},
		//baku borong
		{TahunTanam: 1998, Mandor: "SUKIYATNO", NIK: "9006569", Afdeling: "SETRO", Tipe: "BAKU_BORONG"},
		{TahunTanam: 2013, Mandor: "SUKIYATNO", NIK: "9006569", Afdeling: "SETRO", Tipe: "BAKU_BORONG"},
		{TahunTanam: 2018, Mandor: "SUKIYATNO", NIK: "9006569", Afdeling: "SETRO", Tipe: "BAKU_BORONG"},
		{TahunTanam: 2004, Mandor: "KAYIN", NIK: "9006286", Afdeling: "SETRO", Tipe: "BAKU_BORONG"},
		{TahunTanam: 2005, Mandor: "KAYIN", NIK: "9006286", Afdeling: "SETRO", Tipe: "BAKU_BORONG"},
		{TahunTanam: 2008, Mandor: "MULYANTO", NIK: "9006468", Afdeling: "SETRO", Tipe: "BAKU_BORONG"},
		{TahunTanam: 2009, Mandor: "DWI JATI NUGROHO", NIK: "9006467", Afdeling: "SETRO", Tipe: "BAKU_BORONG"},
		{TahunTanam: 2010, Mandor: "SUPRIYADI", NIK: "9006291", Afdeling: "SETRO", Tipe: "BAKU_BORONG"},
		{TahunTanam: 2011, Mandor: "WIDAYAT", NIK: "9006594", Afdeling: "SETRO", Tipe: "BAKU_BORONG"},
		{TahunTanam: 2012, Mandor: "HARTONO", NIK: "9006406", Afdeling: "SETRO", Tipe: "BAKU_BORONG"},
		{TahunTanam: 2005, Mandor: "RUJIANTORO", NIK: "9006777", Afdeling: "SETRO", Tipe: "BAKU_BORONG"},
		{TahunTanam: 2006, Mandor: "RUJIANTORO", NIK: "9006777", Afdeling: "SETRO", Tipe: "BAKU_BORONG"},
		{TahunTanam: 2006, Mandor: "JONI ANWAR", NIK: "9006585", Afdeling: "SETRO", Tipe: "BAKU_BORONG"},
		//Borong_minggu
		{TahunTanam: 1998, Mandor: "SUKIYATNO", NIK: "9006569", Afdeling: "SETRO", Tipe: "BAKU_MINGGU"},
		{TahunTanam: 2013, Mandor: "SUKIYATNO", NIK: "9006569", Afdeling: "SETRO", Tipe: "BAKU_MINGGU"},
		{TahunTanam: 2018, Mandor: "SUKIYATNO", NIK: "9006569", Afdeling: "SETRO", Tipe: "BAKU_MINGGU"},
		{TahunTanam: 2004, Mandor: "KAYIN", NIK: "9006286", Afdeling: "SETRO", Tipe: "BAKU_MINGGU"},
		{TahunTanam: 2005, Mandor: "KAYIN", NIK: "9006286", Afdeling: "SETRO", Tipe: "BAKU_MINGGU"},
		{TahunTanam: 2008, Mandor: "MULYANTO", NIK: "9006468", Afdeling: "SETRO", Tipe: "BAKU_MINGGU"},
		{TahunTanam: 2009, Mandor: "DWI JATI NUGROHO", NIK: "9006467", Afdeling: "SETRO", Tipe: "BAKU_MINGGU"},
		{TahunTanam: 2010, Mandor: "SUPRIYADI", NIK: "9006291", Afdeling: "SETRO", Tipe: "BAKU_MINGGU"},
		{TahunTanam: 2011, Mandor: "WIDAYAT", NIK: "9006594", Afdeling: "SETRO", Tipe: "BAKU_MINGGU"},
		{TahunTanam: 2012, Mandor: "HARTONO", NIK: "9006406", Afdeling: "SETRO", Tipe: "BAKU_MINGGU"},
		{TahunTanam: 2005, Mandor: "RUJIANTORO", NIK: "9006777", Afdeling: "SETRO", Tipe: "BAKU_MINGGU"},
		{TahunTanam: 2006, Mandor: "RUJIANTORO", NIK: "9006777", Afdeling: "SETRO", Tipe: "BAKU_MINGGU"},
		{TahunTanam: 2006, Mandor: "JONI ANWAR", NIK: "9006585", Afdeling: "SETRO", Tipe: "BAKU_MINGGU"},
		//borong internal
		{TahunTanam: 1998, Mandor: "SUKIYATNO", NIK: "9006569", Afdeling: "SETRO", Tipe: "BAKU_INTERNAL"},
		{TahunTanam: 2013, Mandor: "SUKIYATNO", NIK: "9006569", Afdeling: "SETRO", Tipe: "BAKU_INTERNAL"},
		{TahunTanam: 2018, Mandor: "SUKIYATNO", NIK: "9006569", Afdeling: "SETRO", Tipe: "BAKU_INTERNAL"},
		{TahunTanam: 2004, Mandor: "KAYIN", NIK: "9006286", Afdeling: "SETRO", Tipe: "BAKU_INTERNAL"},
		{TahunTanam: 2005, Mandor: "KAYIN", NIK: "9006286", Afdeling: "SETRO", Tipe: "BAKU_INTERNAL"},
		{TahunTanam: 2008, Mandor: "MULYANTO", NIK: "9006468", Afdeling: "SETRO", Tipe: "BAKU_INTERNAL"},
		{TahunTanam: 2009, Mandor: "DWI JATI NUGROHO", NIK: "9006467", Afdeling: "SETRO", Tipe: "BAKU_INTERNAL"},
		{TahunTanam: 2010, Mandor: "SUPRIYADI", NIK: "9006291", Afdeling: "SETRO", Tipe: "BAKU_INTERNAL"},
		{TahunTanam: 2011, Mandor: "WIDAYAT", NIK: "9006594", Afdeling: "SETRO", Tipe: "BAKU_INTERNAL"},
		{TahunTanam: 2012, Mandor: "HARTONO", NIK: "9006406", Afdeling: "SETRO", Tipe: "BAKU_INTERNAL"},
		{TahunTanam: 2005, Mandor: "RUJIANTORO", NIK: "9006777", Afdeling: "SETRO", Tipe: "BAKU_INTERNAL"},
		{TahunTanam: 2006, Mandor: "RUJIANTORO", NIK: "9006777", Afdeling: "SETRO", Tipe: "BAKU_INTERNAL"},
		{TahunTanam: 2006, Mandor: "JONI ANWAR", NIK: "9006585", Afdeling: "SETRO", Tipe: "BAKU_INTERNAL"},
		//borong_eksternal
		{TahunTanam: 1998, Mandor: "SUKIYATNO", NIK: "9006569", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2013, Mandor: "SUKIYATNO", NIK: "9006569", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2018, Mandor: "SUKIYATNO", NIK: "9006569", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2004, Mandor: "KAYIN", NIK: "9006286", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2005, Mandor: "KAYIN", NIK: "9006286", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2008, Mandor: "MULYANTO", NIK: "9006468", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2009, Mandor: "DWI JATI NUGROHO", NIK: "9006467", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2010, Mandor: "SUPRIYADI", NIK: "9006291", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2011, Mandor: "WIDAYAT", NIK: "9006594", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2012, Mandor: "HARTONO", NIK: "9006406", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2002, Mandor: "AHMAD ARIF", NIK: "9006406", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2003, Mandor: "AHMAD ARIF", NIK: "9006406", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2013, Mandor: "AHMAD ARIF", NIK: "9006406", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2005, Mandor: "RUJIANTORO", NIK: "9006777", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2006, Mandor: "RUJIANTORO", NIK: "9006777", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2006, Mandor: "JONI ANWAR", NIK: "9006585", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		//tetes lanjut
		{TahunTanam: 1998, Mandor: "SUKIYATNO", NIK: "9006569", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2013, Mandor: "SUKIYATNO", NIK: "9006569", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2018, Mandor: "SUKIYATNO", NIK: "9006569", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2004, Mandor: "MULYANTO", NIK: "9006468", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2005, Mandor: "MULYANTO", NIK: "9006468", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2008, Mandor: "TRI SUSANTO", NIK: "9006574", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2009, Mandor: "DWI JATI NUGROHO", NIK: "9006467", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2010, Mandor: "SUPRIYADI", NIK: "9006291", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2011, Mandor: "WIDAYAT", NIK: "9006594", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2012, Mandor: "HARTONO", NIK: "9006406", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2005, Mandor: "JONI ANWAR", NIK: "9006585", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2006, Mandor: "AHMAD ARIF", NIK: "9006406", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
		{TahunTanam: 2006, Mandor: "AHMAD ARIF", NIK: "9006406", Afdeling: "SETRO", Tipe: "BAKU_EKSTERNAL"},
	}

	for _, mandor := range mandors {
		var existingMandor models.BakuMandor
		config.DB.Where("tahun_tanam = ? AND nik = ? AND afdeling = ? AND tipe = ?",
			mandor.TahunTanam, mandor.NIK, mandor.Afdeling, mandor.Tipe).FirstOrCreate(&existingMandor, mandor)
	}
}
