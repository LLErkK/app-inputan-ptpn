let currentViewMode = 'penyadap';
let penyadapList = [];
let mandorList = [];

document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM Content Loaded');
    loadPenyadapList();
    loadMandorList();
    
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

        // ✅ Ambil data yang sesuai dengan format API kamu
        if (json && Array.isArray(json.data)) {
            mandorList = json.data;
            console.log('✅ Mandor list loaded:', mandorList.length);
            console.log('👤 Sample:', mandorList[0]);
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

// FIX: Perbaiki handleMandorInput - jangan hapus data-id saat user sudah select
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
    
    console.log('🎯 Detected fields:', {nameField, nikField, idField});
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
    
    console.log('✅ Filtered results:', filtered.length);
    
    if (filtered.length === 0) {
        dropdown.innerHTML = `<div class="autocomplete-item" style="color: #999; cursor: default;">
            Tidak ada hasil ditemukan<br>
            <small>Total data: ${mandorList.length} | Mencari: "${value}"</small>
        </div>`;
        dropdown.style.display = 'block';
        return;
    }
    
    dropdown.innerHTML = filtered.map(m => {
        const mandorName = m[nameField] || 'N/A';
        const nikValue = m[nikField] || 'N/A';
        const idValue = m[idField] || 0;
        const tahunTanamValue = m[tahunTanam] || 'N/A';
        
        const escapedName = String(mandorName).replace(/'/g, "\\'");
        const escapedNik = String(nikValue).replace(/'/g, "\\'");
        
        return `
            <div class="autocomplete-item" onclick="selectMandor(${idValue}, '${escapedName}', '${escapedNik}')">
                <strong>${mandorName}</strong><br>
                <small>NIK: ${nikValue}</small>
                <strong><br>Tahun Tanam: ${tahunTanamValue}</strong>
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

function selectMandor(id, nama, nik) {
    const input = document.getElementById('namaMandor');
    const displayValue = `${nama} (${nik})`;
    
    input.setAttribute('data-id', id);
    input.setAttribute('data-nama', nama);
    input.setAttribute('data-nik', nik);
    input.setAttribute('data-display-value', displayValue);
    input.value = displayValue;
    
    document.getElementById('mandorDropdown').style.display = 'none';
    console.log('Mandor selected:', {id, nama, nik});
}

document.addEventListener('click', function(e) {
    if (!e.target.matches('#namaPenyadap')) {
        document.getElementById('penyadapDropdown').style.display = 'none';
    }
    if (!e.target.matches('#namaMandor')) {
        document.getElementById('mandorDropdown').style.display = 'none';
    }
});

function selectViewMode(mode) {
    currentViewMode = mode;
    document.getElementById('viewModePopup').classList.remove('active');
    document.getElementById('monitoringContainer').classList.add('active');
    updateViewModeBadge();
    toggleFormFields();
    console.log('View mode selected:', mode);
}

function changeViewMode() {
    document.getElementById('monitoringContainer').classList.remove('active');
    document.getElementById('viewModePopup').classList.add('active');
    clearAll();
}

function toggleFormFields() {
    const formPenyadap = document.getElementById('formPenyadap');
    const formMandor = document.getElementById('formMandor');
    const tablePenyadap = document.getElementById('tablePenyadap');
    const tableMandor = document.getElementById('tableMandor');
    const summaryPenyadap = document.getElementById('summaryPenyadap');
    const summaryMandor = document.getElementById('summaryMandor');
    
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
    const day = String(date.getDate()).padStart(2, '0');
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const year = date.getFullYear();
    return `${day}/${month}/${year}`;
}

function calculateSummary(data, summary) {
    if (summary) {
        return summary;
    }
    
    const calculated = {
        total_records: data.length
    };

    if (currentViewMode === 'mandor') {
        calculated.total_hko = data.reduce((sum, item) => sum + (item.hko_hari_ini || 0), 0);
        calculated.total_basah_latek_kebun = data.reduce((sum, item) => sum + (item.hari_ini_basah_latek_kebun || 0), 0);
        calculated.total_basah_latek_pabrik = data.reduce((sum, item) => sum + (item.hari_ini_basah_latek_pabrik || 0), 0);
        calculated.total_kering_sheet = data.reduce((sum, item) => sum + (item.hari_ini_kering_sheet || 0), 0);
        calculated.rata_rata_produksi_per_taper = calculated.total_hko > 0 
            ? (calculated.total_basah_latek_kebun + calculated.total_kering_sheet) / calculated.total_hko 
            : 0;
    } else {
        calculated.total_basah_latek = data.reduce((sum, item) => sum + (item.BasahLatek || item.basah_latek || 0), 0);
        calculated.total_sheet = data.reduce((sum, item) => sum + (item.Sheet || item.sheet || 0), 0);
        calculated.total_basah_lump = data.reduce((sum, item) => sum + (item.BasahLump || item.basah_lump || 0), 0);
        calculated.total_br_cr = data.reduce((sum, item) => sum + (item.BrCr || item.br_cr || 0), 0);
    }

    return calculated;
}

function renderTable(data, summary) {
    const summarySection = document.getElementById('summarySection');
    
    if (!data || data.length === 0) {
        if (currentViewMode === 'mandor') {
            document.getElementById('mandorTableBody').innerHTML = '<tr><td colspan="15" style="text-align:center; color: #999;">Tidak ada data ditemukan</td></tr>';
        } else {
            document.getElementById('bakuTableBody').innerHTML = '<tr><td colspan="11" style="text-align:center; color: #999;">Tidak ada data ditemukan</td></tr>';
        }
        summarySection.style.display = 'none';
        return;
    }

    summarySection.style.display = 'block';

    const summaryData = calculateSummary(data, summary);

    if (currentViewMode === 'mandor') {
        renderMandorTable(data, summaryData);
    } else {
        renderPenyadapTable(data, summaryData);
    }
}

function renderPenyadapTable(data, summary) {
    const tableBody = document.getElementById('bakuTableBody');
    
    document.getElementById('summaryTotalRecords').textContent = summary.total_records || summary.TotalRecords || 0;
    document.getElementById('summaryBasahLatek').textContent = summary.total_basah_latek || summary.TotalLatek || 0;
    document.getElementById('summarySheet').textContent = (summary.total_sheet || summary.TotalSheet || 0).toFixed(2);
    document.getElementById('summaryBasahLump').textContent = summary.total_basah_lump || summary.TotalLump || 0;
    document.getElementById('summaryBrCr').textContent = (summary.total_br_cr || summary.TotalBrCr || 0).toFixed(2);

    tableBody.innerHTML = '';
    data.forEach((item) => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${formatDate(item.Tanggal || item.tanggal)}</td>
            <td style="text-align: left;">${item.Mandor || item.mandor || '-'}</td>
            <td>${item.TipeProduksi || item.tipe_produksi || '-'}</td>
            <td>${item.TahunTanam || item.tahun_tanam || '-'}</td>
            <td>${item.Afdeling || item.afdeling || '-'}</td>
            <td>${item.NIK || item.nik || '-'}</td>
            <td style="text-align: left;">${item.NamaPenyadap || item.nama_penyadap || '-'}</td>
            <td style="font-weight: 600; color: #0093E9;">${item.BasahLatek || item.basah_latek || 0}</td>
            <td style="font-weight: 600; color: #0093E9;">${(item.Sheet || item.sheet || 0).toFixed(2)}</td>
            <td style="font-weight: 600; color: #0093E9;">${item.BasahLump || item.basah_lump || 0}</td>
            <td style="font-weight: 600; color: #0093E9;">${(item.BrCr || item.br_cr || 0).toFixed(2)}</td>
        `;
        tableBody.appendChild(row);
    });
}

function renderMandorTable(data, summary) {
    const tableBody = document.getElementById('mandorTableBody');
    
    document.getElementById('summaryTotalRecordsMandor').textContent = summary.total_records || summary.TotalRecords || 0;
    document.getElementById('summaryHKO').textContent = summary.total_hko || summary.TotalHKO || 0;
    document.getElementById('summaryLatekKebun').textContent = summary.total_basah_latek_kebun || summary.TotalBasahLatekKebun || 0;
    document.getElementById('summaryLatekPabrik').textContent = summary.total_basah_latek_pabrik || summary.TotalBasahLatekPabrik || 0;
    document.getElementById('summarySheetKering').textContent = summary.total_kering_sheet || summary.TotalKeringSheet || 0;
    document.getElementById('summaryRataPerTaper').textContent = (summary.rata_rata_produksi_per_taper || summary.RataRataProduksiPerTaper || 0).toFixed(2);

    tableBody.innerHTML = '';
    data.forEach((item) => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${formatDate(item.tanggal)}</td>
            <td style="text-align: left;">${item.mandor || '-'}</td>
            <td>${item.tipe_produksi || '-'}</td>
            <td>${item.tahun_tanam || '-'}</td>
            <td>${item.afdeling || '-'}</td>
            <td style="font-weight: 600; color: #667eea;">${item.hko_hari_ini || 0}</td>
            <td style="font-weight: 600; color: #0093E9;">${item.hari_ini_basah_latek_kebun || 0}</td>
            <td style="font-weight: 600; color: #0093E9;">${item.hari_ini_basah_latek_pabrik || 0}</td>
            <td style="font-weight: 600; color: ${item.hari_ini_basah_latek_persen > 5 ? '#e74c3c' : '#27ae60'};">${(item.hari_ini_basah_latek_persen || 0).toFixed(2)}%</td>
            <td style="font-weight: 600; color: #0093E9;">${item.hari_ini_basah_lump_kebun || 0}</td>
            <td style="font-weight: 600; color: #0093E9;">${item.hari_ini_basah_lump_pabrik || 0}</td>
            <td style="font-weight: 600; color: ${item.hari_ini_basah_lump_persen > 5 ? '#e74c3c' : '#27ae60'};">${(item.hari_ini_basah_lump_persen || 0).toFixed(2)}%</td>
            <td style="font-weight: 600; color: #0093E9;">${item.hari_ini_kering_sheet || 0}</td>
            <td style="font-weight: 600; color: #0093E9;">${item.hari_ini_kering_br_cr || 0}</td>
            <td style="font-weight: 600; color: #9b59b6;">${(item.produksi_per_taper_hari_ini || 0).toFixed(2)}</td>
        `;
        tableBody.appendChild(row);
    });
}

async function fetchData(params) {
    const tableBody = currentViewMode === 'mandor' ? 
        document.getElementById('mandorTableBody') : 
        document.getElementById('bakuTableBody');
    const colspan = currentViewMode === 'mandor' ? '15' : '11';
    
    tableBody.innerHTML = `<tr><td colspan="${colspan}" style="text-align:center;"><i>⏳ Memuat data...</i></td></tr>`;

    try {
        const queryParams = new URLSearchParams();
        
        if (!params.tanggalAwal || !params.tanggalAkhir) {
            tableBody.innerHTML = `<tr><td colspan="${colspan}" style="text-align:center; color: #ff9800;">⚠️ Silakan pilih tanggal awal dan akhir terlebih dahulu</td></tr>`;
            return;
        }
        
        if (currentViewMode === 'mandor') {
            if (!params.idMandor) {
                tableBody.innerHTML = `<tr><td colspan="${colspan}" style="text-align:center; color: #ff9800;">⚠️ Silakan pilih mandor terlebih dahulu</td></tr>`;
                return;
            }
            queryParams.append('idMandor', params.idMandor);
        } else {
            if (!params.idPenyadap) {
                tableBody.innerHTML = `<tr><td colspan="${colspan}" style="text-align:center; color: #ff9800;">⚠️ Silakan pilih penyadap terlebih dahulu</td></tr>`;
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
            tableBody.innerHTML = `<tr><td colspan="${colspan}" style="text-align:center; color: #e74c3c;">❌ ${result.message || 'Gagal memuat data'}</td></tr>`;
            return;
        }
        
        const data = result.data || [];
        const summary = result.summary || null;
        
        renderTable(data, summary);
        
    } catch (error) {
        console.error('Error fetching data:', error);
        tableBody.innerHTML = `<tr><td colspan="${colspan}" style="text-align:center; color: #e74c3c;">❌ Terjadi kesalahan: ${error.message}</td></tr>`;
    }
}

function searchData() {
    console.log('searchData() called');
    let params = {};
    
    if (currentViewMode === 'mandor') {
        const inputMandor = document.getElementById('namaMandor');
        const idMandor = inputMandor.getAttribute('data-id') || '';
        const afdeling = document.getElementById('afdelingMandor').value;
        const tanggalAwal = document.getElementById('searchTanggalAwalMandor').value;
        const tanggalAkhir = document.getElementById('searchTanggalAkhirMandor').value;
        const tipeProduksi = document.getElementById('filterJenisMandor').value;

        if (!idMandor) {
            alert('❌ Silakan pilih mandor dari daftar yang tersedia!');
            inputMandor.focus();
            return;
        }
        
        if (!/^\d+$/.test(idMandor)) {
            alert('❌ ID Mandor tidak valid. Silakan pilih dari daftar!');
            inputMandor.value = '';
            inputMandor.removeAttribute('data-id');
            inputMandor.removeAttribute('data-display-value');
            inputMandor.focus();
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
        const idPenyadap = inputPenyadap.getAttribute('data-id') || '';
        const afdeling = document.getElementById('afdelingPenyadap').value;
        const tanggalAwal = document.getElementById('searchTanggalAwal').value;
        const tanggalAkhir = document.getElementById('searchTanggalAkhir').value;
        const tipeProduksi = document.getElementById('filterJenis').value;

        if (!idPenyadap) {
            alert('❌ Silakan pilih penyadap dari daftar yang tersedia!');
            inputPenyadap.focus();
            return;
        }
        
        if (!/^\d+$/.test(idPenyadap)) {
            alert('❌ NIK/ID Penyadap tidak valid. Silakan pilih dari daftar!');
            inputPenyadap.value = '';
            inputPenyadap.removeAttribute('data-id');
            inputPenyadap.removeAttribute('data-display-value');
            inputPenyadap.focus();
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
    inputPenyadap.value = '';
    inputPenyadap.removeAttribute('data-id');
    inputPenyadap.removeAttribute('data-nama');
    inputPenyadap.removeAttribute('data-nik');
    inputPenyadap.removeAttribute('data-display-value');
    document.getElementById('afdelingPenyadap').value = '';
    document.getElementById('searchTanggalAwal').value = '';
    document.getElementById('searchTanggalAkhir').value = '';
    document.getElementById('filterJenis').value = '';
    
    // Clear Mandor form
    const inputMandor = document.getElementById('namaMandor');
    inputMandor.value = '';
    inputMandor.removeAttribute('data-id');
    inputMandor.removeAttribute('data-nama');
    inputMandor.removeAttribute('data-nik');
    inputMandor.removeAttribute('data-display-value');
    document.getElementById('afdelingMandor').value = '';
    document.getElementById('searchTanggalAwalMandor').value = '';
    document.getElementById('searchTanggalAkhirMandor').value = '';
    document.getElementById('filterJenisMandor').value = '';
    
    // Hide summary and reset tables
    document.getElementById('summarySection').style.display = 'none';
    
    if (currentViewMode === 'mandor') {
        document.getElementById('mandorTableBody').innerHTML = '<tr><td colspan="15" style="text-align:center; color: #666; font-style: italic;">Silakan pilih filter dan klik tombol untuk menampilkan data</td></tr>';
    } else {
        document.getElementById('bakuTableBody').innerHTML = '<tr><td colspan="11" style="text-align:center; color: #666; font-style: italic;">Silakan pilih filter dan klik tombol untuk menampilkan data</td></tr>';
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