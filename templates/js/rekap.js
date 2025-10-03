// Fungsi untuk memformat angka dengan satu digit di belakang koma
function formatDecimal(value) {
  // Memastikan nilai yang diberikan adalah angka
  if (isNaN(value)) return value;

  // Memformat angka ke satu desimal dan menghapus ".0" jika tidak ada angka di belakang koma
  const formattedValue = parseFloat(value).toFixed(1);

  // Jika angka berakhiran ".0", hapus ".0"
  if (formattedValue.endsWith(".0")) {
    return formattedValue.slice(0, -2); // Menghapus ".0"
  }

  return formattedValue;
}

// Fungsi untuk mengambil data dari API
async function fetchData() {
  try {
    console.log("=== START LOADING DATA ===");

    // Mengambil data dari endpoint
    const response = await fetch('/rekap/today');
    
    // Menangani respon jika status tidak OK
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.Message || 'Gagal memuat data');
    }

    // Mengambil data JSON dari API
    const data = await response.json();
    
    // Menampilkan parsed data untuk debug
    console.log("Parsed JSON data:", data);
    
    // Mengecek apakah data berhasil diambil dan memeriksa apakah data.data ada
    if (data.success && Array.isArray(data.data)) {
      console.log("Data type: " + typeof data.data);
      console.log("Is array? " + Array.isArray(data.data));
      console.log("Data length: " + data.data.length);
      populateTable(data.data);
    } else {
      document.getElementById('rekap-tbody').innerHTML = `<tr><td colspan="14" style="text-align: center;">${data.message}</td></tr>`;
    }
  } catch (error) {
    // Menampilkan error jika terjadi kesalahan
    console.error('Error:', error);
    document.getElementById('rekap-tbody').innerHTML = `<tr><td colspan="14" style="text-align: center;">${error.message}</td></tr>`;
  }
}

// Fungsi untuk mengelompokkan data berdasarkan tipe
function groupByTipe(details) {
  return details.reduce((result, detail) => {
    // Menambahkan data ke array berdasarkan tipe
    if (!result[detail.tipe]) {
      result[detail.tipe] = [];
    }
    result[detail.tipe].push(detail);
    return result;
  }, {});
}

// Fungsi untuk membuat tabel dan menampilkan data berdasarkan tipe
function populateTable(details) {
  // Mengambil elemen container untuk tabel-tabel
  const container = document.getElementById('rekap-tables-container');
  
  // Menghapus data lama sebelum menambah yang baru
  container.innerHTML = '';

  // Mengelompokkan data berdasarkan tipe
  const groupedData = groupByTipe(details);

  // Inisialisasi total untuk semua tipe
  let totalHko = 0;
  let totalJumlahPabrikBasahLatek = 0;
  let totalJumlahKebunBasahLatek = 0;
  let totalSelisihBasahLatek = 0;
  let totalJumlahSheet = 0;
  let totalK3Sheet = 0;
  let totalJumlahPabrikBasahLump = 0;
  let totalJumlahKebunBasahLump = 0;
  let totalSelisihBasahLump = 0;
  let totalPersentaseSelisihBasahLump = 0;
  let totalJumlahBrCr = 0;

  // Loop untuk setiap tipe dan membuat tabel
  Object.keys(groupedData).forEach(tipe => {
    // Membuat elemen untuk menampilkan nama tipe
    const tipeHeader = document.createElement('h2');
    tipeHeader.textContent = `Tipe: ${tipe}`;
    
    // Membuat elemen tabel untuk setiap tipe
    const tableContainer = document.createElement('div');
    const table = document.createElement('table');
    const tbody = document.createElement('tbody');
    const thead = document.createElement('thead');
    
    // Menambahkan nama tipe di atas tabel
    tableContainer.classList.add('rekap-table-container');
    table.classList.add('rekap-table');
    table.appendChild(thead);
    table.appendChild(tbody);
    tableContainer.appendChild(tipeHeader);  // Menambahkan nama tipe
    tableContainer.appendChild(table);
    container.appendChild(tableContainer);

    // Menambahkan header tabel
    thead.innerHTML = `
      <tr>
        <th rowspan="4">TAHUN<br>TANAM</th>
        <th rowspan="4">NIK</th>
        <th rowspan="4">MANDOR</th>
        <th rowspan="2">HKO</th>
        <th colspan="10">PRODUKSI HARI INI</th>
      </tr>
      <tr>
        <th colspan="6">BASAH</th>
        <th rowspan="3">KKK<br>SHEET</th>
        <th colspan="3">KERING</th>
      </tr>
      <tr>
        <th rowspan="2">HR INI</th>
        <th colspan="3">LATEX</th>
        <th colspan="3">LUMP</th>
        <th rowspan="2">SHEET</th>
        <th rowspan="2">BR.CR</th>
        <th rowspan="2">JUMLAH</th>
      </tr>
      <tr>
        <th>KEBUN</th>
        <th>PABRIK</th>
        <th>%</th>
        <th>KEBUN</th>
        <th>PABRIK</th>
        <th>%</th>
      </tr>
    `;

    // Menambahkan data ke dalam tabel dan menghitung total
    groupedData[tipe].forEach(detail => {
      const row = document.createElement('tr');

      row.innerHTML = `
        <td>${detail.tahun_tanam}</td>
        <td>${detail.nik}</td>
        <td>${detail.mandor}</td>
        <td>${formatDecimal(detail.hko)}</td>
        <td>${formatDecimal(detail.jumlah_kebun_basah_latek)}</td>
        <td>${formatDecimal(detail.jumlah_pabrik_basah_latek)}</td>
        <td>${formatDecimal(detail.persentase_selisih_basah_latek)}</td>
        <td>${formatDecimal(detail.jumlah_kebun_basah_lump)}</td>
        <td>${formatDecimal(detail.jumlah_pabrik_basah_lump)}</td>
        <td>${formatDecimal(detail.persentase_selisih_basah_lump)}</td>
        <td>${formatDecimal(detail.k3_sheet)}</td>
        <td>${formatDecimal(detail.jumlah_sheet)}</td>
        <td>${formatDecimal(detail.jumlah_br_cr)}</td>
        <td>${formatDecimal(detail.jumlah_kering)}</td>
      `;

      // Menambahkan nilai ke total
      totalHko += parseFloat(detail.hko) || 0;
      totalJumlahPabrikBasahLatek += parseFloat(detail.jumlah_pabrik_basah_latek) || 0;
      totalJumlahKebunBasahLatek += parseFloat(detail.jumlah_kebun_basah_latek) || 0;
      totalSelisihBasahLatek += parseFloat(detail.selisih_basah_latek) || 0;
      totalJumlahSheet += parseFloat(detail.jumlah_sheet) || 0;
      totalK3Sheet += parseFloat(detail.k3_sheet) || 0;
      totalJumlahPabrikBasahLump += parseFloat(detail.jumlah_pabrik_basah_lump) || 0;
      totalJumlahKebunBasahLump += parseFloat(detail.jumlah_kebun_basah_lump) || 0;
      totalSelisihBasahLump += parseFloat(detail.selisih_basah_lump) || 0;
      totalPersentaseSelisihBasahLump += parseFloat(detail.persentase_selisih_basah_lump) || 0;
      totalJumlahBrCr += parseFloat(detail.jumlah_br_cr) || 0;

      tbody.appendChild(row);
    });

    // Membuat baris terakhir untuk jumlah
    const totalRow = document.createElement('tr');
    totalRow.innerHTML = `
      <td colspan="3" style="font-weight: bold;">TOTAL</td>
      <td>${formatDecimal(totalHko)}</td>
      <td>${formatDecimal(totalJumlahKebunBasahLatek)}</td>
      <td>${formatDecimal(totalJumlahPabrikBasahLatek)}</td>
      <td>${formatDecimal(totalSelisihBasahLatek)}</td>
      <td>${formatDecimal(totalJumlahKebunBasahLump)}</td>
      <td>${formatDecimal(totalJumlahPabrikBasahLump)}</td>
      <td>${formatDecimal(totalSelisihBasahLump)}</td>
      <td>${formatDecimal(totalK3Sheet)}</td>
      <td>${formatDecimal(totalJumlahSheet)}</td>
      <td>${formatDecimal(totalJumlahBrCr)}</td>
      <td>${formatDecimal(totalSelisihBasahLatek)}</td>
    `;
    tbody.appendChild(totalRow);
  });

  // Menambahkan baris "Jumlah Produksi" di bawah semua tipe
  const totalProductionRow = document.createElement('tr');
  totalProductionRow.innerHTML = `
    <td colspan="3" style="font-weight: bold;">JUMLAH PRODUKSI</td>
    <td>${formatDecimal(totalHko)}</td>
    <td>${formatDecimal(totalJumlahKebunBasahLatek)}</td>
    <td>${formatDecimal(totalJumlahPabrikBasahLatek)}</td>
    <td>${formatDecimal(totalSelisihBasahLatek)}</td>
    <td>${formatDecimal(totalJumlahKebunBasahLump)}</td>
    <td>${formatDecimal(totalJumlahPabrikBasahLump)}</td>
    <td>${formatDecimal(totalSelisihBasahLump)}</td>
    <td>${formatDecimal(totalK3Sheet)}</td>
    <td>${formatDecimal(totalJumlahSheet)}</td>
    <td>${formatDecimal(totalJumlahBrCr)}</td>
    <td>${formatDecimal(totalSelisihBasahLatek)}</td>
  `;
  container.appendChild(totalProductionRow);
}

// Memanggil fungsi fetchData saat halaman selesai dimuat
document.addEventListener('DOMContentLoaded', fetchData);
