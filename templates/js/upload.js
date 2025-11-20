// DOM Elements
const fileInput = document.getElementById('fileInput');
const fileName = document.getElementById('fileName');
const fileMeta = document.getElementById('fileMeta');
const uploadForm = document.getElementById('uploadForm');
const result = document.getElementById('result');
const loading = document.getElementById('loading');
const submitBtn = document.getElementById('submitBtn');

// Modal Elements
const openDeleteMasterModalBtn = document.getElementById('openDeleteMasterModalBtn');
const deleteMasterModal = document.getElementById('deleteMasterModal');
const deleteMasterModalClose = document.getElementById('deleteMasterModalClose');
const masterListBody = document.getElementById('masterListBody');
const masterModalLoading = document.getElementById('masterModalLoading');

// File Input Change Handler
fileInput.addEventListener('change', function() {
    if (this.files && this.files[0]) {
        const file = this.files[0];
        const sizeKB = (file.size / 1024).toFixed(2);
        const sizeMB = (file.size / 1024 / 1024).toFixed(2);
        const name = file.name;
        fileName.querySelector('.file-badge').textContent = name;
        fileMeta.textContent = `${sizeMB} MB (${sizeKB} KB)`;
    } else {
        fileName.querySelector('.file-badge').textContent = 'Belum ada file dipilih';
        fileMeta.textContent = '';
    }
});

// Form Submit Handler
uploadForm.addEventListener('submit', async function(e) {
    e.preventDefault();

    const tanggal = document.getElementById('tanggal').value;
    const afdeling = document.getElementById('afdeling').value;
    const file = fileInput.files[0];

    if (!tanggal || !file || !afdeling) {
        showResult(false, 'Tanggal, afdeling, dan file wajib diisi');
        return;
    }

    if (file.size > 10 * 1024 * 1024) {
        showResult(false, 'Ukuran file terlalu besar (maksimal 10MB)');
        return;
    }

    const formData = new FormData();
    formData.append('tanggal', tanggal);
    formData.append('afdeling', afdeling);
    formData.append('file', file);

    loading.classList.add('show');
    submitBtn.disabled = true;
    result.style.display = 'none';

    try {
        const response = await fetch('/api/upload', {
            method: 'POST',
            body: formData
        });

        const data = await response.json();

        if (data.success) {
            showResult(true, data.message, {
                tanggal: data.tanggal,
                afdeling: data.afdeling,
                fileName: data.fileName,
                fileSize: data.fileSize
            });
            uploadForm.reset();
            fileName.querySelector('.file-badge').textContent = 'Belum ada file dipilih';
            fileMeta.textContent = '';
        } else {
            showResult(false, data.message);
        }
    } catch (error) {
        showResult(false, 'Terjadi kesalahan saat upload: ' + error.message);
    } finally {
        loading.classList.remove('show');
        submitBtn.disabled = false;
    }
});

// Show Result Function
function showResult(success, message, data = null) {
    result.classList.remove('error');
    if (!success) {
        result.classList.add('error');
        document.getElementById('resultTitle').textContent = '❌ ' + message;
        document.getElementById('resultDate').textContent = '-';
        document.getElementById('resultAfdeling').textContent = '-';
        document.getElementById('resultFile').textContent = '-';
        document.getElementById('resultSize').textContent = '-';
    } else {
        document.getElementById('resultTitle').textContent = '✅ ' + message;
        if (data) {
            document.getElementById('resultDate').textContent = data.tanggal;
            document.getElementById('resultAfdeling').textContent = data.afdeling;
            document.getElementById('resultFile').textContent = data.fileName;
            document.getElementById('resultSize').textContent = (data.fileSize / 1024).toFixed(2) + ' KB';
        }
    }

    result.style.display = 'block';
}

// Modal Functions
function openModal() {
    deleteMasterModal.style.display = 'block';
    loadMasterList();
}

function closeModal() {
    deleteMasterModal.style.display = 'none';
    masterListBody.innerHTML = '';
}

// Modal Event Listeners
openDeleteMasterModalBtn.addEventListener('click', openModal);
deleteMasterModalClose.addEventListener('click', closeModal);
deleteMasterModal.addEventListener('click', function(e) {
    if (e.target === deleteMasterModal) closeModal();
});

// Load Master List
async function loadMasterList() {
    masterListBody.innerHTML = '<tr><td colspan="5" style="text-align:center; padding:12px; color:#666;">Memuat...</td></tr>';
    masterModalLoading.style.display = 'block';
    
    try {
        const res = await fetch('/api/master', { credentials: 'same-origin' });
        if (!res.ok) throw new Error('HTTP ' + res.status);
        
        const masters = await res.json();
        
        if (!Array.isArray(masters) || masters.length === 0) {
            masterListBody.innerHTML = '<tr><td colspan="5" style="text-align:center; padding:12px; color:#666;">Tidak ada master</td></tr>';
            return;
        }
        
        masterListBody.innerHTML = '';
        masters.forEach(m => {
            const id = m.ID || m.id || m.Id || m.id_master || m.IDMaster || '';
            const tanggal = m.Tanggal || m.tanggal || m.tanggal_str || '';
            const afdeling = m.Afdeling || m.afdeling || '';
            const namaFile = m.NamaFile || m.nama_file || m.FileName || m.fileName || '';

            const tr = document.createElement('tr');
            tr.innerHTML = `
                <td style="padding:8px; border-bottom:1px solid #eee;">${id}</td>
                <td style="padding:8px; border-bottom:1px solid #eee;">${formatDate(tanggal)}</td>
                <td style="padding:8px; border-bottom:1px solid #eee;">${afdeling}</td>
                <td style="padding:8px; border-bottom:1px solid #eee;">${namaFile}</td>
                <td style="padding:8px; border-bottom:1px solid #eee; text-align:right;">
                    <button class="delete-row-btn" data-id="${id}" style="background:#ff4d4f; color:#fff; border:0; padding:6px 10px; border-radius:6px; cursor:pointer;">Hapus</button>
                </td>
            `;
            masterListBody.appendChild(tr);
        });

        // Attach delete handlers
        Array.from(masterListBody.querySelectorAll('.delete-row-btn')).forEach(btn => {
            btn.addEventListener('click', async function() {
                const id = this.dataset.id;
                if (!id) return alert('ID tidak ditemukan untuk item ini');
                if (!confirm(`Hapus master dengan ID ${id}?`)) return;
                
                const btnRef = this;
                btnRef.disabled = true;
                const oldText = btnRef.textContent;
                btnRef.textContent = 'Menghapus...';
                
                try {
                    const resp = await fetch('/api/master/' + encodeURIComponent(id), {
                        method: 'DELETE',
                        credentials: 'same-origin'
                    });
                    
                    if (!resp.ok) {
                        const txt = await resp.text().catch(() => 'HTTP ' + resp.status);
                        alert('Gagal menghapus: ' + txt);
                        btnRef.disabled = false;
                        btnRef.textContent = oldText;
                        return;
                    }
                    
                    // Remove row
                    const row = btnRef.closest('tr');
                    if (row) row.remove();
                    
                } catch (err) {
                    alert('Kesalahan jaringan: ' + err.message);
                    btnRef.disabled = false;
                    btnRef.textContent = oldText;
                }
            });
        });
        
    } catch (err) {
        masterListBody.innerHTML = '<tr><td colspan="5" style="text-align:center; padding:12px; color:#c00;">Gagal memuat daftar master</td></tr>';
        console.error(err);
    } finally {
        masterModalLoading.style.display = 'none';
    }
}

// Format Date Helper - FIXED VERSION
function formatDate(d) {
    if (!d) return '-';
    
    try {
        // If it's already in YYYY-MM-DD format, return as is
        if (typeof d === 'string' && /^\d{4}-\d{2}-\d{2}$/.test(d)) {
            return d;
        }
        
        // Parse the date
        const dt = new Date(d);
        if (isNaN(dt)) return d;
        
        // Use local timezone instead of UTC to avoid -1 day issue
        const year = dt.getFullYear();
        const month = String(dt.getMonth() + 1).padStart(2, '0');
        const day = String(dt.getDate()).padStart(2, '0');
        
        return `${year}-${month}-${day}`;
    } catch (e) {
        console.error('Error formatting date:', e);
        return d;
    }
}