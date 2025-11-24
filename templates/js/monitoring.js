// monitoring.js - Full JS for Penyadap/Mandor monitoring
// Default end date set to today (input[type="date"] expects YYYY-MM-DD)

'use strict';

let currentViewMode = 'penyadap';
let penyadapList = [];
let mandorList = [];

document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM Content Loaded');
    loadPenyadapList();
    loadMandorList();

    // set default end date to today for both forms
    setDefaultEndDates();

    const form = document.getElementById('monitoringForm');
    if (form) {
        form.addEventListener('submit', function(e) {
            e.preventDefault();
            console.log('Form submitted - calling searchData()');
            searchData();
        });
        console.log('Form submit handler attached successfully');
    } else {
        console.error('ERROR: Form element not found!');
    }

    addAutocompleteStyles();
});

// Helper: return today's date as YYYY-MM-DD (for input[type=date])
function getTodayISO() {
    const today = new Date();
    const yyyy = today.getFullYear();
    const mm = String(today.getMonth() + 1).padStart(2, '0');
    const dd = String(today.getDate()).padStart(2, '0');
    return `${yyyy}-${mm}-${dd}`;
}

function setDefaultEndDates() {
    try {
        const endDate = getTodayISO();
        const searchTanggalAkhir = document.getElementById('searchTanggalAkhir');
        const searchTanggalAkhirMandor = document.getElementById('searchTanggalAkhirMandor');

        if (searchTanggalAkhir && !searchTanggalAkhir.value) {
            searchTanggalAkhir.value = endDate;
            console.log('Default searchTanggalAkhir set to', endDate);
        }
        if (searchTanggalAkhirMandor && !searchTanggalAkhirMandor.value) {
            searchTanggalAkhirMandor.value = endDate;
            console.log('Default searchTanggalAkhirMandor set to', endDate);
        }
    } catch (err) {
        console.error('Error setting default end dates:', err);
    }
}

async function loadPenyadapList() {
    try {
        console.log('🔄 Fetching penyadap data from /api/penyadap...');
        const response = await fetch('/api/penyadap');

        if (!response.ok) {
            console.error('❌ Penyadap API error:', response.status);
            return;
        }

        const json = await response.json();

        if (json && Array.isArray(json)) {
            penyadapList = json;
        } else if (json && Array.isArray(json.data)) {
            penyadapList = json.data;
        } else {
            penyadapList = [];
        }

        console.log('✅ Penyadap list loaded:', penyadapList.length);
    } catch (error) {
        console.error('💥 Error loading penyadap list:', error);
    }
}

async function loadMandorList() {
    try {
        console.log('🔄 Fetching mandor data...');
        const response = await fetch('/api/mandor');
        console.log('📡 Status:', response.status, response.statusText);

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`HTTP ${response.status}: ${errorText}`);
        }

        const text = await response.text();
        console.log('📦 Raw API response:', text);

        let json;
        try {
            json = JSON.parse(text);
        } catch (err) {
            console.error('❌ JSON parse error:', err);
            return;
        }

        if (json && Array.isArray(json.data)) {
            mandorList = json.data;
            console.log('✅ Mandor list loaded:', mandorList.length);
            console.log('👤 Sample:', mandorList[0]);
        } else if (Array.isArray(json)) {
            mandorList = json;
            console.log('✅ Mandor list loaded (array root):', mandorList.length);
        } else {
            console.warn('⚠️ Unexpected response format:', json);
            mandorList = [];
        }

    } catch (error) {
        console.error('💥 Error loading mandor list:', error);
    }
}

// FIX: Perbaiki handlePenyadapInput - jangan hapus data-id saat user sudah select
function handlePenyadapInput(input) {
    const value = input.value.toLowerCase();
    const dropdown = document.getElementById('penyadapDropdown');

    // FIX: Hanya clear attributes jika user benar-benar mengubah text
    const storedDisplayValue = input.getAttribute('data-display-value');
    if (storedDisplayValue && input.value === storedDisplayValue) {
        // User tidak mengubah apa-apa, jangan hapus data
        dropdown.style.display = 'none';
        return;
    }

    // Jika user mengubah text, clear data attributes
    if (storedDisplayValue && input.value !== storedDisplayValue) {
        input.removeAttribute('data-id');
        input.removeAttribute('data-nama');
        input.removeAttribute('data-nik');
        input.removeAttribute('data-display-value');
    }

    if (value.length < 2) {
        dropdown.style.display = 'none';
        return;
    }

    if (penyadapList.length === 0) {
        dropdown.innerHTML = '<div class="autocomplete-item" style="color: #e74c3c; cursor: default;">⚠️ Data penyadap kosong</div>';
        dropdown.style.display = 'block';
        return;
    }

    const filtered = penyadapList.filter(p => 
        p.nama_penyadap.toLowerCase().includes(value) ||
        p.nik.toString().includes(value)
    );

    if (filtered.length === 0) {
        dropdown.innerHTML = '<div class="autocomplete-item" style="color: #999; cursor: default;">Tidak ada hasil ditemukan</div>';
        dropdown.style.display = 'block';
        return;
    }

    dropdown.innerHTML = filtered.map(p => {
        const escapedName = p.nama_penyadap.replace(/'/g, "\\'");
        const escapedNik = p.nik.toString().replace(/'/g, "\\'");

        return `
            <div class="autocomplete-item" onclick="selectPenyadap(${p.id}, '${escapedName}', '${escapedNik}')">
                <strong>${p.nama_penyadap}</strong><br>
                <small>NIK: ${p.nik}</small>
            </div>
        `;
    }).join('');

    dropdown.style.display = 'block';
}

// FIX: Perbaiki handleMandorInput - tampilkan setiap kombinasi mandor + tahun tanam tanpa duplikat
function handleMandorInput(input) {
    const value = input.value.toLowerCase();
    const dropdown = document.getElementById('mandorDropdown');

    console.log('🔍 handleMandorInput called with:', value);
    console.log('📊 mandorList length:', mandorList.length);

    // FIX: Hanya clear attributes jika user benar-benar mengubah text
    const storedDisplayValue = input.getAttribute('data-display-value');
    if (storedDisplayValue && input.value === storedDisplayValue) {
        // User tidak mengubah apa-apa, jangan hapus data
        dropdown.style.display = 'none';
        return;
    }

    // Jika user mengubah text, clear data attributes
    if (storedDisplayValue && input.value !== storedDisplayValue) {
        input.removeAttribute('data-id');
        input.removeAttribute('data-nama');
        input.removeAttribute('data-nik');
        input.removeAttribute('data-tahun-tanam');
        input.removeAttribute('data-display-value');
    }

    if (value.length < 2) {
        dropdown.style.display = 'none';
        return;
    }

    if (mandorList.length === 0) {
        dropdown.innerHTML = '<div class="autocomplete-item" style="color: #e74c3c; cursor: default;">⚠️ Data mandor kosong. Cek endpoint /api/mandor</div>';
        dropdown.style.display = 'block';
        return;
    }

    if (mandorList.length > 0) {
        console.log('🔑 First mandor object keys:', Object.keys(mandorList[0]));
        console.log('👤 First mandor full object:', mandorList[0]);
    }

    const firstItem = mandorList[0];

    const nameField = Object.keys(firstItem).find(key => {
        const lowerKey = key.toLowerCase();
        return lowerKey === 'mandor' || 
               lowerKey === 'nama_mandor' || 
               lowerKey === 'namamandor' ||
               lowerKey === 'nama';
    }) || 'mandor';

    const nikField = Object.keys(firstItem).find(key => {
        const lowerKey = key.toLowerCase();
        return lowerKey === 'nik';
    }) || 'nik';

    const idField = Object.keys(firstItem).find(key => {
        const lowerKey = key.toLowerCase();
        return lowerKey === 'id';
    }) || 'id';
    
    const tahunTanam = Object.keys(firstItem).find(key => {
        const lowerKey = key.toLowerCase();
        return lowerKey === 'tahun_tanam' || lowerKey === 'tahuntanam';
    }) || 'tahun_tanam';

    console.log('🎯 Detected fields:', {nameField, nikField, idField, tahunTanam});
    console.log('📝 Sample values:', {
        name: firstItem[nameField],
        nik: firstItem[nikField],
        id: firstItem[idField],
        tahunTanam: firstItem[tahunTanam]
    });

    const filtered = mandorList.filter(m => {
        const mandorName = m[nameField] ? String(m[nameField]) : '';
        const nikValue = m[nikField] ? String(m[nikField]) : '';

        const nameMatch = mandorName.toLowerCase().includes(value);
        const nikMatch = nikValue.includes(value);

        return nameMatch || nikMatch;
    });

    console.log('✅ Filtered results before deduplication:', filtered.length);

    // Remove exact duplicates based on NIK + Tahun Tanam combination
    const uniqueMap = new Map();
    filtered.forEach(m => {
        const mandorName = m[nameField] || 'N/A';
        const nikValue = m[nikField] || 'N/A';
        const idValue = m[idField] || 0;
        const tahunTanamValue = m[tahunTanam] || 'N/A';

        // Create unique key: NIK + Tahun Tanam
        const uniqueKey = `${nikValue}_${tahunTanamValue}`;
        
        // Only add if not already exists
        if (!uniqueMap.has(uniqueKey)) {
            uniqueMap.set(uniqueKey, {
                name: mandorName,
                nik: nikValue,
                id: idValue,
                tahunTanam: tahunTanamValue
            });
        }
    });

    const uniqueResults = Array.from(uniqueMap.values());
    
    // Sort by name, then by tahun tanam (descending)
    uniqueResults.sort((a, b) => {
        if (a.name !== b.name) {
            return a.name.localeCompare(b.name);
        }
        // Same name, sort by tahun tanam descending
        if (a.tahunTanam === 'N/A') return 1;
        if (b.tahunTanam === 'N/A') return -1;
        return b.tahunTanam - a.tahunTanam;
    });

    console.log('✅ Unique results after deduplication:', uniqueResults.length);

    if (uniqueResults.length === 0) {
        dropdown.innerHTML = `<div class="autocomplete-item" style="color: #999; cursor: default;">
            Tidak ada hasil ditemukan<br>
            <small>Total data: ${mandorList.length} | Mencari: "${value}"</small>
        </div>`;
        dropdown.style.display = 'block';
        return;
    }

    dropdown.innerHTML = uniqueResults.map(m => {
        const escapedName = String(m.name).replace(/'/g, "\\'");
        const escapedNik = String(m.nik).replace(/'/g, "\\'");

        return `
            <div class="autocomplete-item" onclick="selectMandor(${m.id}, '${escapedName}', '${escapedNik}', '${m.tahunTanam}')">
                <strong>${m.name}</strong><br>
                <small>NIK: ${m.nik}</small><br>
                <small style="color: #0093E9; font-weight: 600;">📅 Tahun Tanam: ${m.tahunTanam}</small>
            </div>
        `;
    }).join('');

    dropdown.style.display = 'block';
}

function selectPenyadap(id, nama, nik) {
    const input = document.getElementById('namaPenyadap');
    const displayValue = `${nama} (${nik})`;

    input.setAttribute('data-id', id);
    input.setAttribute('data-nama', nama);
    input.setAttribute('data-nik', nik);
    input.setAttribute('data-display-value', displayValue);
    input.value = displayValue;

    document.getElementById('penyadapDropdown').style.display = 'none';
    console.log('Penyadap selected:', {id, nama, nik});
}

function selectMandor(id, nama, nik, tahunTanam) {
    const input = document.getElementById('namaMandor');
    const displayValue = tahunTanam && tahunTanam !== 'N/A' 
        ? `${nama} (${nik}) - ${tahunTanam}`
        : `${nama} (${nik})`;

    input.setAttribute('data-id', id);
    input.setAttribute('data-nama', nama);
    input.setAttribute('data-nik', nik);
    input.setAttribute('data-tahun-tanam', tahunTanam || '');
    input.setAttribute('data-display-value', displayValue);
    input.value = displayValue;

    document.getElementById('mandorDropdown').style.display = 'none';
    console.log('Mandor selected:', {id, nama, nik, tahunTanam});
}

document.addEventListener('click', function(e) {
    if (!e.target.matches('#namaPenyadap')) {
        const dd = document.getElementById('penyadapDropdown');
        if (dd) dd.style.display = 'none';
    }
    if (!e.target.matches('#namaMandor')) {
        const dd2 = document.getElementById('mandorDropdown');
        if (dd2) dd2.style.display = 'none';
    }
});

function selectViewMode(mode) {
    currentViewMode = mode;
    const vmp = document.getElementById('viewModePopup');
    const mc = document.getElementById('monitoringContainer');
    if (vmp) vmp.classList.remove('active');
    if (mc) mc.classList.add('active');
    updateViewModeBadge();
    toggleFormFields();
    console.log('View mode selected:', mode);
}

function changeViewMode() {
    const mc = document.getElementById('monitoringContainer');
    const vmp = document.getElementById('viewModePopup');
    if (mc) mc.classList.remove('active');
    if (vmp) vmp.classList.add('active');
    clearAll();
}

function toggleFormFields() {
    const formPenyadap = document.getElementById('formPenyadap');
    const formMandor = document.getElementById('formMandor');
    const tablePenyadap = document.getElementById('tablePenyadap');
    const tableMandor = document.getElementById('tableMandor');
    const summaryPenyadap = document.getElementById('summaryPenyadap');
    const summaryMandor = document.getElementById('summaryMandor');

    if (!formPenyadap || !formMandor || !tablePenyadap || !tableMandor || !summaryPenyadap || !summaryMandor) return;

    if (currentViewMode === 'mandor') {
        formPenyadap.style.display = 'none';
        formMandor.style.display = 'grid';
        tablePenyadap.style.display = 'none';
        tableMandor.style.display = 'table';
        summaryPenyadap.style.display = 'none';
        summaryMandor.style.display = 'grid';
    } else {
        formPenyadap.style.display = 'grid';
        formMandor.style.display = 'none';
        tablePenyadap.style.display = 'table';
        tableMandor.style.display = 'none';
        summaryPenyadap.style.display = 'grid';
        summaryMandor.style.display = 'none';
    }
}

function updateViewModeBadge() {
    const badge = document.getElementById('viewModeBadge');
    if (!badge) return;
    if (currentViewMode === 'mandor') {
        badge.textContent = '👥 Mode: Mandor';
        badge.classList.add('mandor');
    } else {
        badge.textContent = '👤 Mode: Penyadap';
        badge.classList.remove('mandor');
    }
}

function formatDate(isoDate) {
    if (!isoDate) return '-';
    const date = new Date(isoDate);
    if (isNaN(date)) return isoDate;
    const day = String(date.getDate()).padStart(2, '0');
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const year = date.getFullYear();
    return `${day}/${month}/${year}`;
}

function calculateSummary(data, summary) {
    // Semua perhitungan sudah dilakukan di backend
    // Jika summary dari backend tersedia, gunakan langsung
    if (summary) {
        return summary;
    }

    // Fallback jika backend tidak mengirim summary (seharusnya tidak terjadi)
    return {
        total_records: data.length
    };
}

function renderTable(data, summary) {
    const summarySection = document.getElementById('summarySection');
    
    if (!data || data.length === 0) {
        if (currentViewMode === 'mandor') {
            const el = document.getElementById('mandorTableBody');
            if (el) el.innerHTML = '<tr><td colspan="16" style="text-align:center; color: #999;">Tidak ada data ditemukan</td></tr>';
        } else {
            const el = document.getElementById('bakuTableBody');
            if (el) el.innerHTML = '<tr><td colspan="12" style="text-align:center; color: #999;">Tidak ada data ditemukan</td></tr>';
        }
        if (summarySection) summarySection.style.display = 'none';
        return;
    }

    if (summarySection) summarySection.style.display = 'block';

    const summaryData = calculateSummary(data, summary);

    if (currentViewMode === 'mandor') {
        renderMandorTable(data, summaryData);
    } else {
        renderPenyadapTable(data, summaryData);
    }
}

function renderPenyadapTable(data, summary) {
    const tableBody = document.getElementById('bakuTableBody');
    if (!tableBody) return;
    
    const totalRecords = summary.total_records || summary.TotalRecords || 0;
    const totalBasahLatek = summary.total_basah_latek || summary.TotalLatek || 0;
    const totalSheet = summary.total_sheet || summary.TotalSheet || 0;
    const totalBasahLump = summary.total_basah_lump || summary.TotalLump || 0;
    const totalBrCr = summary.total_br_cr || summary.TotalBrCr || 0;
    const totalProduksi = summary.total_produksi || summary.TotalProduksi || 0;
    
    // Safely update summary elements with null checks
    const elTotalRecords = document.getElementById('summaryTotalRecords');
    if (elTotalRecords) elTotalRecords.textContent = totalRecords;
    
    const elBasahLatek = document.getElementById('summaryBasahLatek');
    if (elBasahLatek) elBasahLatek.textContent = totalBasahLatek;
    
    const elSheet = document.getElementById('summarySheet');
    if (elSheet) elSheet.textContent = totalSheet.toFixed(2);
    
    const elBasahLump = document.getElementById('summaryBasahLump');
    if (elBasahLump) elBasahLump.textContent = totalBasahLump;
    
    const elBrCr = document.getElementById('summaryBrCr');
    if (elBrCr) elBrCr.textContent = totalBrCr.toFixed(2);
    
    const elTotalProduksi = document.getElementById('summaryTotalProduksi');
    if (elTotalProduksi) elTotalProduksi.textContent = totalProduksi.toFixed(2);

    tableBody.innerHTML = '';
    data.forEach((item) => {
        const basahLatek = item.BasahLatek || item.basah_latek || 0;
        const sheet = item.Sheet || item.sheet || 0;
        const basahLump = item.BasahLump || item.basah_lump || 0;
        const brCr = item.BrCr || item.br_cr || 0;
        const totalProduksi = item.TotalProduksi || item.total_produksi || 0;
        
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${formatDate(item.Tanggal || item.tanggal)}</td>
            <td style="text-align: left;">${item.Mandor || item.mandor || '-'}</td>
            <td>${item.TipeProduksi || item.tipe_produksi || '-'}</td>
            <td>${item.TahunTanam || item.tahun_tanam || '-'}</td>
            <td>${item.Afdeling || item.afdeling || '-'}</td>
            <td>${item.NIK || item.nik || '-'}</td>
            <td style="text-align: left;">${item.NamaPenyadap || item.nama_penyadap || '-'}</td>
            <td style="font-weight: 600; color: #0093E9;">${basahLatek}</td>
            <td style="font-weight: 600; color: #0093E9;">${sheet.toFixed(2)}</td>
            <td style="font-weight: 600; color: #0093E9;">${basahLump}</td>
            <td style="font-weight: 600; color: #0093E9;">${brCr.toFixed(2)}</td>
            <td style="font-weight: 700; color: #27ae60; background-color: rgba(39, 174, 96, 0.1);">${totalProduksi.toFixed(2)}</td>
        `;
        tableBody.appendChild(row);
    });
}

function renderMandorTable(data, summary) {
    const tableBody = document.getElementById('mandorTableBody');
    if (!tableBody) return;
    
    const totalHKO = summary.total_hko || summary.TotalHKO || 0;
    const totalLatekKebun = summary.total_basah_latek_kebun || summary.TotalBasahLatekKebun || 0;
    const totalLatekPabrik = summary.total_basah_latek_pabrik || summary.TotalBasahLatekPabrik || 0;
    const totalSheetKering = summary.total_kering_sheet || summary.TotalKeringSheet || 0;
    const totalLumpKebun = summary.total_basah_lump_kebun || summary.TotalBasahLumpKebun || 0;
    const totalBrCr = summary.total_kering_br_cr || summary.TotalKeringBrCr || 0;
    const totalProduksi = summary.total_produksi || summary.TotalProduksi || 0;
    
    // Safely update summary elements with null checks
    const elTotalRecords = document.getElementById('summaryTotalRecordsMandor');
    if (elTotalRecords) elTotalRecords.textContent = summary.total_records || summary.TotalRecords || 0;
    
    const elHKO = document.getElementById('summaryHKO');
    if (elHKO) elHKO.textContent = totalHKO;
    
    const elLatekKebun = document.getElementById('summaryLatekKebun');
    if (elLatekKebun) elLatekKebun.textContent = totalLatekKebun;
    
    const elLatekPabrik = document.getElementById('summaryLatekPabrik');
    if (elLatekPabrik) elLatekPabrik.textContent = totalLatekPabrik;
    
    const elSheetKering = document.getElementById('summarySheetKering');
    if (elSheetKering) elSheetKering.textContent = totalSheetKering;
    
    const elRataPerTaper = document.getElementById('summaryRataPerTaper');
    if (elRataPerTaper) elRataPerTaper.textContent = (summary.rata_rata_produksi_per_taper || summary.RataRataProduksiPerTaper || 0).toFixed(2);
    
    const elTotalProduksi = document.getElementById('summaryTotalProduksiMandor');
    if (elTotalProduksi) elTotalProduksi.textContent = totalProduksi.toFixed(2);

    tableBody.innerHTML = '';
    data.forEach((item) => {
        const basahLatekKebun = item.hari_ini_basah_latek_kebun || 0;
        const keringSheet = item.hari_ini_kering_sheet || 0;
        const basahLumpKebun = item.hari_ini_basah_lump_kebun || 0;
        const keringBrCr = item.hari_ini_kering_br_cr || 0;
        const totalProduksi = item.total_produksi_hari_ini ||0;
        
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${formatDate(item.tanggal)}</td>
            <td style="text-align: left;">${item.mandor || '-'}</td>
            <td>${item.tipe_produksi || '-'}</td>
            <td>${item.tahun_tanam || '-'}</td>
            <td>${item.afdeling || '-'}</td>
            <td style="font-weight: 600; color: #667eea;">${item.hko_hari_ini || 0}</td>
            <td style="font-weight: 600; color: #0093E9;">${basahLatekKebun}</td>
            <td style="font-weight: 600; color: #0093E9;">${item.hari_ini_basah_latek_pabrik || 0}</td>
            <td style="font-weight: 600; color: ${item.hari_ini_basah_latek_persen > 5 ? '#e74c3c' : '#27ae60'};">${(item.hari_ini_basah_latek_persen || 0).toFixed(2)}%</td>
            <td style="font-weight: 600; color: #0093E9;">${basahLumpKebun}</td>
            <td style="font-weight: 600; color: #0093E9;">${item.hari_ini_basah_lump_pabrik || 0}</td>
            <td style="font-weight: 600; color: ${item.hari_ini_basah_lump_persen > 5 ? '#e74c3c' : '#27ae60'};">${(item.hari_ini_basah_lump_persen || 0).toFixed(2)}%</td>
            <td style="font-weight: 600; color: #0093E9;">${keringSheet}</td>
            <td style="font-weight: 600; color: #0093E9;">${keringBrCr}</td>
            <td style="font-weight: 600; color: #9b59b6;">${(item.produksi_per_taper_hari_ini || 0).toFixed(2)}</td>
            <td style="font-weight: 700; color: #27ae60; background-color: rgba(39, 174, 96, 0.1);">${totalProduksi.toFixed(2)}</td>
        `;
        tableBody.appendChild(row);
    });
}

async function fetchData(params) {
    const tableBody = currentViewMode === 'mandor' ? 
        document.getElementById('mandorTableBody') : 
        document.getElementById('bakuTableBody');
    const colspan = currentViewMode === 'mandor' ? '16' : '12';
    
    if (tableBody) tableBody.innerHTML = `<tr><td colspan="${colspan}" style="text-align:center;"><i>⏳ Memuat data...</i></td></tr>`;

    try {
        const queryParams = new URLSearchParams();
        
        if (!params.tanggalAwal || !params.tanggalAkhir) {
            if (tableBody) tableBody.innerHTML = `<tr><td colspan="${colspan}" style="text-align:center; color: #ff9800;">⚠️ Silakan pilih tanggal awal dan akhir terlebih dahulu</td></tr>`;
            return;
        }

        if (currentViewMode === 'mandor') {
            if (!params.idMandor) {
                if (tableBody) tableBody.innerHTML = `<tr><td colspan="${colspan}" style="text-align:center; color: #ff9800;">⚠️ Silakan pilih mandor terlebih dahulu</td></tr>`;
                return;
            }
            queryParams.append('idMandor', params.idMandor);
        } else {
            if (!params.idPenyadap) {
                if (tableBody) tableBody.innerHTML = `<tr><td colspan="${colspan}" style="text-align:center; color: #ff9800;">⚠️ Silakan pilih penyadap terlebih dahulu</td></tr>`;
                return;
            }
            queryParams.append('idPenyadap', params.idPenyadap);
        }

        queryParams.append('tanggalAwal', params.tanggalAwal);
        queryParams.append('tanggalAkhir', params.tanggalAkhir);

        if (params.tipeProduksi && params.tipeProduksi.trim() !== '') {
            queryParams.append('tipeProduksi', params.tipeProduksi);
        }
        if (params.afdeling && params.afdeling.trim() !== '') {
            queryParams.append('afdeling', params.afdeling);
        }

        const url = `/api/search?${queryParams.toString()}`;
        
        console.log('Fetching data from:', url);
        console.log('Params:', Object.fromEntries(queryParams));
        
        const response = await fetch(url);
        
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`HTTP ${response.status}: ${errorText}`);
        }
        
        const result = await response.json();
        console.log('API Response:', result);
        
        if (result.success === false) {
            if (tableBody) tableBody.innerHTML = `<tr><td colspan="${colspan}" style="text-align:center; color: #e74c3c;">❌ ${result.message || 'Gagal memuat data'}</td></tr>`;
            return;
        }
        
        const data = result.data || [];
        const summary = result.summary || null;
        
        renderTable(data, summary);
        
    } catch (error) {
        console.error('Error fetching data:', error);
        if (tableBody) tableBody.innerHTML = `<tr><td colspan="${colspan}" style="text-align:center; color: #e74c3c;">❌ Terjadi kesalahan: ${error.message}</td></tr>`;
    }
}

function searchData() {
    console.log('searchData() called');
    let params = {};
    
    if (currentViewMode === 'mandor') {
        const inputMandor = document.getElementById('namaMandor');
        const idMandor = inputMandor ? inputMandor.getAttribute('data-id') || '' : '';
        const afdeling = document.getElementById('afdelingMandor').value;
        const tanggalAwal = document.getElementById('searchTanggalAwalMandor').value;
        const tanggalAkhir = document.getElementById('searchTanggalAkhirMandor').value;
        const tipeProduksi = document.getElementById('filterJenisMandor').value;

        if (!idMandor) {
            alert('❌ Silakan pilih mandor dari daftar yang tersedia!');
            if (inputMandor) inputMandor.focus();
            return;
        }
        
        if (!/^\d+$/.test(idMandor)) {
            alert('❌ ID Mandor tidak valid. Silakan pilih dari daftar!');
            if (inputMandor) {
                inputMandor.value = '';
                inputMandor.removeAttribute('data-id');
                inputMandor.removeAttribute('data-display-value');
                inputMandor.focus();
            }
            return;
        }

        if (!tanggalAwal || !tanggalAkhir) {
            alert('❌ Silakan pilih tanggal awal dan akhir!');
            return;
        }

        params.idMandor = idMandor;
        if (afdeling && afdeling.trim() !== '') params.afdeling = afdeling;
        params.tanggalAwal = tanggalAwal;
        params.tanggalAkhir = tanggalAkhir;
        if (tipeProduksi && tipeProduksi.trim() !== '') params.tipeProduksi = tipeProduksi;
    } else {
        const inputPenyadap = document.getElementById('namaPenyadap');
        const idPenyadap = inputPenyadap ? inputPenyadap.getAttribute('data-id') || '' : '';
        const afdeling = document.getElementById('afdelingPenyadap').value;
        const tanggalAwal = document.getElementById('searchTanggalAwal').value;
        const tanggalAkhir = document.getElementById('searchTanggalAkhir').value;
        const tipeProduksi = document.getElementById('filterJenis').value;

        if (!idPenyadap) {
            alert('❌ Silakan pilih penyadap dari daftar yang tersedia!');
            if (inputPenyadap) inputPenyadap.focus();
            return;
        }
        
        if (!/^\d+$/.test(idPenyadap)) {
            alert('❌ NIK/ID Penyadap tidak valid. Silakan pilih dari daftar!');
            if (inputPenyadap) {
                inputPenyadap.value = '';
                inputPenyadap.removeAttribute('data-id');
                inputPenyadap.removeAttribute('data-display-value');
                inputPenyadap.focus();
            }
            return;
        }

        if (!tanggalAwal || !tanggalAkhir) {
            alert('❌ Silakan pilih tanggal awal dan akhir!');
            return;
        }

        params.idPenyadap = idPenyadap;
        if (afdeling && afdeling.trim() !== '') params.afdeling = afdeling;
        params.tanggalAwal = tanggalAwal;
        params.tanggalAkhir = tanggalAkhir;
        if (tipeProduksi && tipeProduksi.trim() !== '') params.tipeProduksi = tipeProduksi;
    }

    console.log('Search params:', params);
    fetchData(params);
}

function clearAll() {
    // Clear Penyadap form
    const inputPenyadap = document.getElementById('namaPenyadap');
    if (inputPenyadap) {
        inputPenyadap.value = '';
        inputPenyadap.removeAttribute('data-id');
        inputPenyadap.removeAttribute('data-nama');
        inputPenyadap.removeAttribute('data-nik');
        inputPenyadap.removeAttribute('data-display-value');
    }
    const afdelingPenyadap = document.getElementById('afdelingPenyadap');
    if (afdelingPenyadap) afdelingPenyadap.value = '';
    const searchTanggalAwal = document.getElementById('searchTanggalAwal');
    if (searchTanggalAwal) searchTanggalAwal.value = '';
    const searchTanggalAkhir = document.getElementById('searchTanggalAkhir');
    if (searchTanggalAkhir) searchTanggalAkhir.value = getTodayISO(); // set back to today
    const filterJenis = document.getElementById('filterJenis');
    if (filterJenis) filterJenis.value = '';
    
    // Clear Mandor form
    const inputMandor = document.getElementById('namaMandor');
    if (inputMandor) {
        inputMandor.value = '';
        inputMandor.removeAttribute('data-id');
        inputMandor.removeAttribute('data-nama');
        inputMandor.removeAttribute('data-nik');
        inputMandor.removeAttribute('data-tahun-tanam');
        inputMandor.removeAttribute('data-display-value');
    }
    const afdelingMandor = document.getElementById('afdelingMandor');
    if (afdelingMandor) afdelingMandor.value = '';
    const searchTanggalAwalMandor = document.getElementById('searchTanggalAwalMandor');
    if (searchTanggalAwalMandor) searchTanggalAwalMandor.value = '';
    const searchTanggalAkhirMandor = document.getElementById('searchTanggalAkhirMandor');
    if (searchTanggalAkhirMandor) searchTanggalAkhirMandor.value = getTodayISO(); // set back to today
    const filterJenisMandor = document.getElementById('filterJenisMandor');
    if (filterJenisMandor) filterJenisMandor.value = '';
    
    // Hide summary and reset tables
    const summarySection = document.getElementById('summarySection');
    if (summarySection) summarySection.style.display = 'none';
    
    if (currentViewMode === 'mandor') {
        const mandorTableBody = document.getElementById('mandorTableBody');
        if (mandorTableBody) mandorTableBody.innerHTML = '<tr><td colspan="16" style="text-align:center; color: #666; font-style: italic;">Silakan pilih filter dan klik tombol untuk menampilkan data</td></tr>';
    } else {
        const bakuTableBody = document.getElementById('bakuTableBody');
        if (bakuTableBody) bakuTableBody.innerHTML = '<tr><td colspan="12" style="text-align:center; color: #666; font-style: italic;">Silakan pilih filter dan klik tombol untuk menampilkan data</td></tr>';
    }
    
    console.log('All filters cleared');
}

function exportData() {
    alert('Export functionality for ' + currentViewMode + ' mode - Coming soon!');
}

function addAutocompleteStyles() {
    const style = document.createElement('style');
    style.textContent = `
    .autocomplete-dropdown {
        display: none;
        position: absolute;
        top: 100%;
        left: 0;
        right: 0;
        background: white;
        border: 2px solid #e8f4f8;
        border-top: none;
        border-radius: 0 0 10px 10px;
        max-height: 200px;
        overflow-y: auto;
        z-index: 1000;
        box-shadow: 0 4px 12px rgba(0,0,0,0.1);
    }

    .autocomplete-item {
        padding: 12px 14px;
        cursor: pointer;
        border-bottom: 1px solid #f0f0f0;
        transition: background-color 0.2s ease;
    }

    .autocomplete-item:hover {
        background-color: rgba(0, 147, 233, 0.1);
    }

    .autocomplete-item:last-child {
        border-bottom: none;
    }

    .autocomplete-item strong {
        color: #2c3e50;
        display: block;
        margin-bottom: 4px;
    }

    .autocomplete-item small {
        color: #7f8c8d;
        font-size: 0.85em;
    }
    `;
    document.head.appendChild(style);
}