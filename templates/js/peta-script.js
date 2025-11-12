const API_BASE = '/api'; // Sesuaikan dengan base URL API Anda
let currentData = null;
let allPetaData = []; // Menyimpan semua data untuk filtering

function showAlert(message, type) {
    const alertSuccess = document.getElementById('alertSuccess');
    const alertError = document.getElementById('alertError');

    if (type === 'success') {
        alertSuccess.textContent = message;
        alertSuccess.classList.add('show');
        alertError.classList.remove('show');
    } else {
        alertError.textContent = message;
        alertError.classList.add('show');
        alertSuccess.classList.remove('show');
    }

    setTimeout(() => {
        alertSuccess.classList.remove('show');
        alertError.classList.remove('show');
    }, 5000);
}

function showLoading(show) {
    document.getElementById('loading').classList.toggle('show', show);
}

async function searchPeta(code) {
    showLoading(true);
    try {
        const response = await fetch(`${API_BASE}/peta?code=${encodeURIComponent(code)}`);
        if (!response.ok) {
            throw new Error('Data tidak ditemukan');
        }

        currentData = await response.json();
        displayData(currentData);
        document.getElementById('tableContainer').style.display = 'none';
        showAlert('Data berhasil dimuat', 'success');
    } catch (error) {
        showAlert('Error: ' + error.message, 'error');
        document.getElementById('dataDisplay').classList.remove('show');
    } finally {
        showLoading(false);
    }
}

function displayData(data) {
    const dataGrid = document.getElementById('dataGrid');
    dataGrid.innerHTML = `
        <div class="data-item">
            <label>ID</label>
            <div class="value">${data.ID || '-'}</div>
        </div>
        <div class="data-item">
            <label>Kode</label>
            <div class="value">${data.Code || '-'}</div>
        </div>
        <div class="data-item">
            <label>Blok</label>
            <div class="value">${data.Blok || '-'}</div>
        </div>
        <div class="data-item">
            <label>Afdeling</label>
            <div class="value">${data.Afdeling || '-'}</div>
        </div>
        <div class="data-item">
            <label>Luas (ha)</label>
            <div class="value">${data.Luas || 0}</div>
        </div>
        <div class="data-item">
            <label>Jumlah Pohon</label>
            <div class="value">${data.JumlahPohon || 0}</div>
        </div>
        <div class="data-item">
            <label>Jenis Kebun</label>
            <div class="value">${data.JenisKebun || '-'}</div>
        </div>
        <div class="data-item">
            <label>Tahun Tanam</label>
            <div class="value">${data.TahunTanam || '-'}</div>
        </div>
        <div class="data-item">
            <label>Kloon</label>
            <div class="value">${data.Kloon || '-'}</div>
        </div>
    `;
    document.getElementById('dataDisplay').classList.add('show');
}

function showEditForm() {
    if (!currentData) return;

    document.getElementById('editId').value = currentData.ID;
    document.getElementById('editCode').value = currentData.Code || '';
    document.getElementById('editBlok').value = currentData.Blok || '';
    document.getElementById('editAfdeling').value = currentData.Afdeling || '';
    document.getElementById('editLuas').value = currentData.Luas || '';
    document.getElementById('editJumlahPohon').value = currentData.JumlahPohon || '';
    document.getElementById('editJenisKebun').value = currentData.JenisKebun || '';
    document.getElementById('editTahunTanam').value = currentData.TahunTanam || '';
    document.getElementById('editKloon').value = currentData.Kloon || '';

    document.getElementById('editForm').classList.add('show');
}

function cancelEdit() {
    document.getElementById('editForm').classList.remove('show');
}

async function saveEdit() {
    const id = document.getElementById('editId').value;
    const code = document.getElementById('editCode').value.trim();
    const afdeling = document.getElementById('editAfdeling').value.trim();

    if (!code || !afdeling) {
        showAlert('Kode dan Afdeling wajib diisi', 'error');
        return;
    }

    const updatedData = {
        Code: code,
        Blok: document.getElementById('editBlok').value.trim(),
        Afdeling: afdeling,
        Luas: parseFloat(document.getElementById('editLuas').value) || 0,
        JumlahPohon: parseInt(document.getElementById('editJumlahPohon').value) || 0,
        JenisKebun: document.getElementById('editJenisKebun').value.trim(),
        TahunTanam: parseInt(document.getElementById('editTahunTanam').value) || 0,
        Kloon: document.getElementById('editKloon').value.trim()
    };

    showLoading(true);
    try {
        const response = await fetch(`${API_BASE}/peta/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(updatedData)
        });

        if (!response.ok) {
            throw new Error('Gagal menyimpan perubahan');
        }

        const result = await response.json();
        currentData = result;
        displayData(currentData);
        cancelEdit();
        showAlert('Data berhasil diperbarui', 'success');
    } catch (error) {
        showAlert('Error: ' + error.message, 'error');
    } finally {
        showLoading(false);
    }
}

async function loadAllData() {
    showLoading(true);
    try {
        const response = await fetch(`${API_BASE}/all/peta`);
        if (!response.ok) {
            throw new Error('Gagal memuat data');
        }

        const data = await response.json();
        allPetaData = data; // Simpan data untuk filtering
        displayTable(data);
        document.getElementById('dataDisplay').classList.remove('show');
        document.getElementById('editForm').classList.remove('show');
        showAlert(`Berhasil memuat ${data.length} data`, 'success');
        updateSearchInfo(data.length, data.length);
    } catch (error) {
        showAlert('Error: ' + error.message, 'error');
    } finally {
        showLoading(false);
    }
}

function displayTable(data) {
    const tbody = document.getElementById('dataTableBody');
    tbody.innerHTML = '';

    data.forEach(item => {
        const row = tbody.insertRow();
        row.innerHTML = `
            <td>${item.Code || '-'}</td>
            <td>${item.Blok || '-'}</td>
            <td>${item.Afdeling || '-'}</td>
            <td>${item.Luas || 0}</td>
            <td>${item.JumlahPohon || 0}</td>
            <td>${item.JenisKebun || '-'}</td>
            <td>${item.TahunTanam || '-'}</td>
            <td>${item.Kloon || '-'}</td>
            <td>
                <div class="action-buttons">
                    <button class="btn btn-primary btn-small" onclick="viewDetail('${item.Code}')">👁️ Lihat</button>
                </div>
            </td>
        `;
    });

    document.getElementById('tableContainer').style.display = 'block';
}

async function viewDetail(code) {
    await searchPeta(code);
    window.scrollTo({ top: 0, behavior: 'smooth' });
}

// Fungsi untuk filter data tabel
function filterTable() {
    const searchValue = document.getElementById('searchTable').value.toLowerCase().trim();

    if (!searchValue) {
        displayTable(allPetaData);
        updateSearchInfo(allPetaData.length, allPetaData.length);
        return;
    }

    const filteredData = allPetaData.filter(item => {
        const code = (item.Code || '').toLowerCase();
        const blok = (item.Blok || '').toLowerCase();
        const afdeling = (item.Afdeling || '').toLowerCase();
        const jenisKebun = (item.JenisKebun || '').toLowerCase();
        const kloon = (item.Kloon || '').toLowerCase();
        const tahunTanam = (item.TahunTanam || '').toString();

        return code.includes(searchValue) ||
            blok.includes(searchValue) ||
            afdeling.includes(searchValue) ||
            jenisKebun.includes(searchValue) ||
            kloon.includes(searchValue) ||
            tahunTanam.includes(searchValue);
    });

    displayTable(filteredData);
    updateSearchInfo(filteredData.length, allPetaData.length);
}

// Update info hasil pencarian
function updateSearchInfo(shown, total) {
    const searchInfo = document.getElementById('searchInfo');
    if (shown === total) {
        searchInfo.textContent = `Menampilkan ${total} data`;
    } else {
        searchInfo.textContent = `Menampilkan ${shown} dari ${total} data`;
        searchInfo.style.color = '#667eea';
        searchInfo.style.fontWeight = '600';
    }
}

// Event listener untuk search input
document.addEventListener('DOMContentLoaded', function() {
    const searchInput = document.getElementById('searchTable');
    if (searchInput) {
        searchInput.addEventListener('input', filterTable);

        // Clear search saat tekan Escape
        searchInput.addEventListener('keydown', function(e) {
            if (e.key === 'Escape') {
                this.value = '';
                filterTable();
            }
        });
    }
});