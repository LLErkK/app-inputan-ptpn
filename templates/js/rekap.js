// Fungsi untuk memformat angka dengan satu digit di belakang koma
function formatDecimal(value) {
  if (isNaN(value)) return value;
  const formattedValue = parseFloat(value).toFixed(1);
  if (formattedValue.endsWith(".0")) {
    return formattedValue.slice(0, -2);
  }
  return formattedValue;
}

// Fungsi untuk mengambil data dari API
async function fetchData() {
  try {
    console.log("=== START LOADING DATA ===");
    const response = await fetch('/rekap/today');
    
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.Message || 'Gagal memuat data');
    }

    const data = await response.json();
    console.log("Parsed JSON data:", data);
    
    if (data.success && Array.isArray(data.data)) {
      console.log("Data type: " + typeof data.data);
      console.log("Is array? " + Array.isArray(data.data));
      console.log("Data length: " + data.data.length);
      populateTable(data.data);
    } else {
      document.getElementById('rekap-tbody').innerHTML = `<tr><td colspan="14" style="text-align: center;">${data.message}</td></tr>`;
    }
  } catch (error) {
    console.error('Error:', error);
    document.getElementById('rekap-tbody').innerHTML = `<tr><td colspan="14" style="text-align: center;">${error.message}</td></tr>`;
  }
}

// Fungsi untuk mengelompokkan data berdasarkan tipe
function groupByTipe(details) {
  return details.reduce((result, detail) => {
    if (!result[detail.tipe]) {
      result[detail.tipe] = [];
    }
    result[detail.tipe].push(detail);
    return result;
  }, {});
}

// Fungsi untuk membuat tabel dan menampilkan data berdasarkan tipe
function populateTable(details) {
  const container = document.getElementById('rekap-tables-container');
  container.innerHTML = '';

  const groupedData = groupByTipe(details);

  // Inisialisasi GRAND TOTAL untuk semua tipe
  let grandTotalHko = 0;
  let grandTotalJumlahKebunBasahLatek = 0;
  let grandTotalJumlahPabrikBasahLatek = 0;
  let grandTotalJumlahKebunBasahLump = 0;
  let grandTotalJumlahPabrikBasahLump = 0;
  let grandTotalK3Sheet = 0;
  let grandTotalJumlahSheet = 0;
  let grandTotalJumlahBrCr = 0;
  let grandTotalJumlahKering = 0;

  // Loop untuk setiap tipe dan membuat tabel
  Object.keys(groupedData).forEach(tipe => {
    const tipeHeader = document.createElement('h2');
    tipeHeader.textContent = `Tipe: ${tipe}`;
    
    const tableContainer = document.createElement('div');
    const table = document.createElement('table');
    const tbody = document.createElement('tbody');
    const thead = document.createElement('thead');
    
    tableContainer.classList.add('rekap-table-container');
    table.classList.add('rekap-table');
    table.appendChild(thead);
    table.appendChild(tbody);
    tableContainer.appendChild(tipeHeader);
    tableContainer.appendChild(table);
    container.appendChild(tableContainer);

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

    // Inisialisasi total PER TIPE (reset untuk setiap tipe)
    let tipeHko = 0;
    let tipeJumlahKebunBasahLatek = 0;
    let tipeJumlahPabrikBasahLatek = 0;
    let tipeJumlahKebunBasahLump = 0;
    let tipeJumlahPabrikBasahLump = 0;
    let tipeK3Sheet = 0;
    let tipeJumlahSheet = 0;
    let tipeJumlahBrCr = 0;
    let tipeJumlahKering = 0;

    // Menambahkan data ke dalam tabel
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

      // Akumulasi untuk total per tipe
      tipeHko += parseFloat(detail.hko) || 0;
      tipeJumlahKebunBasahLatek += parseFloat(detail.jumlah_kebun_basah_latek) || 0;
      tipeJumlahPabrikBasahLatek += parseFloat(detail.jumlah_pabrik_basah_latek) || 0;
      tipeJumlahKebunBasahLump += parseFloat(detail.jumlah_kebun_basah_lump) || 0;
      tipeJumlahPabrikBasahLump += parseFloat(detail.jumlah_pabrik_basah_lump) || 0;
      tipeK3Sheet += parseFloat(detail.k3_sheet) || 0;
      tipeJumlahSheet += parseFloat(detail.jumlah_sheet) || 0;
      tipeJumlahBrCr += parseFloat(detail.jumlah_br_cr) || 0;
      tipeJumlahKering += parseFloat(detail.jumlah_kering) || 0;

      tbody.appendChild(row);
    });

    // Hitung persentase untuk total per tipe
    const tipePersenLatek = tipeJumlahPabrikBasahLatek > 0 
      ? ((tipeJumlahKebunBasahLatek - tipeJumlahPabrikBasahLatek) / tipeJumlahPabrikBasahLatek * 100)
      : 0;
    
    const tipePersenLump = tipeJumlahPabrikBasahLump > 0
      ? ((tipeJumlahKebunBasahLump - tipeJumlahPabrikBasahLump) / tipeJumlahPabrikBasahLump * 100)
      : 0;

    // Baris TOTAL per tipe
    const totalRow = document.createElement('tr');
    totalRow.innerHTML = `
      <td colspan="3" style="font-weight: bold;">TOTAL</td>
      <td>${formatDecimal(tipeHko)}</td>
      <td>${formatDecimal(tipeJumlahKebunBasahLatek)}</td>
      <td>${formatDecimal(tipeJumlahPabrikBasahLatek)}</td>
      <td>${formatDecimal(tipePersenLatek)}</td>
      <td>${formatDecimal(tipeJumlahKebunBasahLump)}</td>
      <td>${formatDecimal(tipeJumlahPabrikBasahLump)}</td>
      <td>${formatDecimal(tipePersenLump)}</td>
      <td>${formatDecimal(tipeK3Sheet)}</td>
      <td>${formatDecimal(tipeJumlahSheet)}</td>
      <td>${formatDecimal(tipeJumlahBrCr)}</td>
      <td>${formatDecimal(tipeJumlahKering)}</td>
    `;
    tbody.appendChild(totalRow);

    // Akumulasi ke grand total
    grandTotalHko += tipeHko;
    grandTotalJumlahKebunBasahLatek += tipeJumlahKebunBasahLatek;
    grandTotalJumlahPabrikBasahLatek += tipeJumlahPabrikBasahLatek;
    grandTotalJumlahKebunBasahLump += tipeJumlahKebunBasahLump;
    grandTotalJumlahPabrikBasahLump += tipeJumlahPabrikBasahLump;
    grandTotalK3Sheet += tipeK3Sheet;
    grandTotalJumlahSheet += tipeJumlahSheet;
    grandTotalJumlahBrCr += tipeJumlahBrCr;
    grandTotalJumlahKering += tipeJumlahKering;
  });

  // Hitung persentase untuk grand total
  const grandPersenLatek = grandTotalJumlahPabrikBasahLatek > 0
    ? ((grandTotalJumlahKebunBasahLatek - grandTotalJumlahPabrikBasahLatek) / grandTotalJumlahPabrikBasahLatek * 100)
    : 0;
  
  const grandPersenLump = grandTotalJumlahPabrikBasahLump > 0
    ? ((grandTotalJumlahKebunBasahLump - grandTotalJumlahPabrikBasahLump) / grandTotalJumlahPabrikBasahLump * 100)
    : 0;

  // Membuat tabel terpisah untuk JUMLAH PRODUKSI
  const grandTotalContainer = document.createElement('div');
  grandTotalContainer.classList.add('rekap-table-container');
  
  const grandTotalHeader = document.createElement('h2');
  grandTotalHeader.textContent = 'JUMLAH PRODUKSI';
  grandTotalHeader.style.marginTop = '20px';
  
  const grandTotalTable = document.createElement('table');
  grandTotalTable.classList.add('rekap-table');
  
  // Membuat thead yang sama dengan tabel di atas
  const grandTotalThead = document.createElement('thead');
  grandTotalThead.innerHTML = `
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
  
  const grandTotalTbody = document.createElement('tbody');
  const grandTotalRow = document.createElement('tr');
  grandTotalRow.innerHTML = `
    <td colspan="3" style="font-weight: bold;">TOTAL SEMUA TIPE</td>
    <td>${formatDecimal(grandTotalHko)}</td>
    <td>${formatDecimal(grandTotalJumlahKebunBasahLatek)}</td>
    <td>${formatDecimal(grandTotalJumlahPabrikBasahLatek)}</td>
    <td>${formatDecimal(grandPersenLatek)}</td>
    <td>${formatDecimal(grandTotalJumlahKebunBasahLump)}</td>
    <td>${formatDecimal(grandTotalJumlahPabrikBasahLump)}</td>
    <td>${formatDecimal(grandPersenLump)}</td>
    <td>${formatDecimal(grandTotalK3Sheet)}</td>
    <td>${formatDecimal(grandTotalJumlahSheet)}</td>
    <td>${formatDecimal(grandTotalJumlahBrCr)}</td>
    <td>${formatDecimal(grandTotalJumlahKering)}</td>
  `;
  
  grandTotalTbody.appendChild(grandTotalRow);
  grandTotalTable.appendChild(grandTotalThead);
  grandTotalTable.appendChild(grandTotalTbody);
  grandTotalContainer.appendChild(grandTotalHeader);
  grandTotalContainer.appendChild(grandTotalTable);
  container.appendChild(grandTotalContainer);
}

// Memanggil fungsi fetchData saat halaman selesai dimuat
document.addEventListener('DOMContentLoaded', fetchData);