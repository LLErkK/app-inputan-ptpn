let currentViewMode = 'penyadap';
let penyadapList = [];
let mandorList = [];
let comparisonCount = 0;
let comparisonData = {};
let chartInstance = null; 

// Initialize when page loads
document.addEventListener('DOMContentLoaded', function() {
    console.log('Perbandingan page loaded');
    loadPenyadapList();
    loadMandorList();
});

// Load Penyadap List from API
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

// Load Mandor List from API
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
        } else {
            console.warn('⚠️ Unexpected response format:', json);
            mandorList = [];
        }

    } catch (error) {
        console.error('💥 Error loading mandor list:', error);
    }
}

// Select View Mode (Penyadap or Mandor)
function selectViewMode(mode) {
    currentViewMode = mode;
    document.getElementById('viewModePopup').classList.remove('active');
    updateViewModeBadge();
    resetAll();
}

// Change View Mode
function changeViewMode() {
    document.getElementById('viewModePopup').classList.add('active');
}

// Update View Mode Badge
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

// Update Layout Columns based on card count
function updateLayoutColumns() {
    const wrapper = document.getElementById('comparisonWrapper');
    const cards = wrapper.querySelectorAll('.comparison-card');
    const count = cards.length;
    
    wrapper.className = 'comparison-wrapper';
    if (count === 1) wrapper.classList.add('cols-1');
    else if (count === 2) wrapper.classList.add('cols-2');
    else if (count === 3) wrapper.classList.add('cols-3');
    else wrapper.classList.add('cols-4');
}

// Update Add Button State
function updateAddButtonState() {
    const currentCards = document.querySelectorAll('.comparison-card').length;
    const addBtn = document.querySelector('.add-btn');
    
    if (currentCards >= 4) {
        addBtn.disabled = true;
        addBtn.style.opacity = '0.5';
        addBtn.style.cursor = 'not-allowed';
        addBtn.innerHTML = '🚫 Maksimal 4 Perbandingan';
    } else {
        addBtn.disabled = false;
        addBtn.style.opacity = '1';
        addBtn.style.cursor = 'pointer';
        addBtn.innerHTML = '➕ Tambah Perbandingan';
    }
}

// Update Result Columns based on result count
function updateResultColumns() {
    const wrapper = document.getElementById('resultWrapper');
    const cards = wrapper.querySelectorAll('.result-card');
    const count = cards.length;
    
    wrapper.className = 'result-wrapper';
    if (count === 1) wrapper.classList.add('cols-1');
    else if (count === 2) wrapper.classList.add('cols-2');
    else if (count === 3) wrapper.classList.add('cols-3');
    else wrapper.classList.add('cols-4');
}

// Add Comparison Card
function addComparisonCard() {
    const currentCards = document.querySelectorAll('.comparison-card').length;
    
    // Check if maximum cards reached
    if (currentCards >= 4) {
        alert('⚠️ Maksimal 4 perbandingan sudah tercapai!');
        return;
    }

    comparisonCount++;
    const cardId = `card-${comparisonCount}`;
    const container = document.getElementById('comparisonWrapper');

    const cardHTML = currentViewMode === 'penyadap' 
        ? createPenyadapCard(cardId, comparisonCount)
        : createMandorCard(cardId, comparisonCount);

    container.insertAdjacentHTML('beforeend', cardHTML);
    setupAutocomplete(cardId);
    updateLayoutColumns();
    updateAddButtonState();
}

// Create Penyadap Card HTML
function createPenyadapCard(cardId, number) {
    const today = new Date().toISOString().split('T')[0];
    
    return `
        <div class="comparison-card" id="${cardId}">
            <div class="card-header">
                <div class="card-number">${number}</div>
                <button class="remove-card-btn" onclick="removeCard('${cardId}')" ${number === 1 ? 'style="display:none"' : ''}>✕</button>
                <div class="card-title">Penyadap ${number}</div>
            </div>
            <div class="card-form">
                <div class="form-group">
                    <label>Nama Penyadap <span style="color: red;">*</span></label>
                    <div class="autocomplete-wrapper">
                        <input type="text" 
                               id="${cardId}-nama" 
                               placeholder="Ketik nama penyadap..." 
                               autocomplete="off"
                               oninput="handleAutocomplete('${cardId}', 'penyadap')">
                        <div id="${cardId}-dropdown" class="autocomplete-dropdown"></div>
                    </div>
                </div>
                <div class="form-group">
                    <label>Tanggal Awal <span style="color: red;">*</span></label>
                    <input type="date" id="${cardId}-tanggalAwal">
                </div>
                <div class="form-group">
                    <label>Tanggal Akhir <span style="color: red;">*</span></label>
                    <input type="date" id="${cardId}-tanggalAkhir" value="${today}">
                </div>
                <div class="form-group">
                    <label>Type Produksi</label>
                    <select id="${cardId}-tipe">
                        <option value="">-- Semua Type --</option>
                        <option value="BAKU">Baku</option>
                        <option value="BAKU_BORONG">Baku Borong</option>
                        <option value="BORONG_EXTERNAL">Borong Eksternal</option>
                        <option value="BORONG_INTERNAL">Borong Internal</option>
                        <option value="BORONG_MINGGU">Borong Minggu</option>
                        <option value="TETES_LANJUT">Tetes
                        <option value="TETES_LANJUT">Tetes Lanjut</option>
                    </select>
                </div>
            </div>
        </div>
    `;
}

// Create Mandor Card HTML
function createMandorCard(cardId, number) {
    const today = new Date().toISOString().split('T')[0];
    
    return `
        <div class="comparison-card" id="${cardId}">
            <div class="card-header">
                <div class="card-number">${number}</div>
                <button class="remove-card-btn" onclick="removeCard('${cardId}')" ${number === 1 ? 'style="display:none"' : ''}>✕</button>
                <div class="card-title">Mandor ${number}</div>
            </div>
            <div class="card-form">
                <div class="form-group">
                    <label>Nama Mandor <span style="color: red;">*</span></label>
                    <div class="autocomplete-wrapper">
                        <input type="text" 
                               id="${cardId}-nama" 
                               placeholder="Ketik nama mandor..." 
                               autocomplete="off"
                               oninput="handleAutocomplete('${cardId}', 'mandor')">
                        <div id="${cardId}-dropdown" class="autocomplete-dropdown"></div>
                    </div>
                </div>
                <div class="form-group">
                    <label>Tanggal Awal <span style="color: red;">*</span></label>
                    <input type="date" id="${cardId}-tanggalAwal">
                </div>
                <div class="form-group">
                    <label>Tanggal Akhir <span style="color: red;">*</span></label>
                    <input type="date" id="${cardId}-tanggalAkhir" value="${today}">
                </div>
                <div class="form-group">
                    <label>Type Produksi</label>
                    <select id="${cardId}-tipe">
                        <option value="">-- Semua Type --</option>
                        <option value="PRODUKSI BAKU">Baku</option>
                        <option value="PRODUKSI BAKU BORONG">Baku Borong</option>
                        <option value="PRODUKSI BORONG EXTERNAL">Borong Eksternal</option>
                        <option value="PRODUKSI BORONG INTERNAL">Borong Internal</option>
                        <option value="PRODUKSI BORONG MINGGU">Borong Minggu</option>
                        <option value="PRODUKSI TETES LANJUT">Tetes Lanjut</option>
                    </select>
                </div>
            </div>
        </div>
    `;
}

// Setup Autocomplete Event Listeners
function setupAutocomplete(cardId) {
    const input = document.getElementById(`${cardId}-nama`);
    const dropdown = document.getElementById(`${cardId}-dropdown`);

    // Close dropdown when clicking outside
    document.addEventListener('click', function(e) {
        if (!e.target.closest(`#${cardId}`)) {
            dropdown.style.display = 'none';
        }
    });
}

// Handle Autocomplete Input
function handleAutocomplete(cardId, type) {
    const input = document.getElementById(`${cardId}-nama`);
    const dropdown = document.getElementById(`${cardId}-dropdown`);
    const value = input.value.toLowerCase();

    // Check if value matches stored display value
    const storedValue = input.getAttribute('data-display-value');
    if (storedValue && input.value === storedValue) {
        dropdown.style.display = 'none';
        return;
    }

    // Reset data attributes if value changed
    if (storedValue && input.value !== storedValue) {
        input.removeAttribute('data-id');
        input.removeAttribute('data-display-value');
    }

    if (value.length < 2) {
        dropdown.style.display = 'none';
        return;
    }

    const list = type === 'penyadap' ? penyadapList : mandorList;
    
    if (list.length === 0) {
        dropdown.innerHTML = '<div class="autocomplete-item" style="color: #e74c3c;">⚠️ Data tidak tersedia</div>';
        dropdown.style.display = 'block';
        return;
    }

    // Determine field names dynamically
    let nameField, nikField, idField;
    
    if (type === 'penyadap') {
        nameField = 'nama_penyadap';
        nikField = 'nik';
        idField = 'id';
    } else {
        const firstItem = mandorList[0];
        nameField = Object.keys(firstItem).find(key => {
            const lowerKey = key.toLowerCase();
            return lowerKey === 'mandor' || 
                   lowerKey === 'nama_mandor' || 
                   lowerKey === 'namamandor' ||
                   lowerKey === 'nama';
        }) || 'mandor';
        
        nikField = Object.keys(firstItem).find(key => {
            const lowerKey = key.toLowerCase();
            return lowerKey === 'nik';
        }) || 'nik';
        
        idField = Object.keys(firstItem).find(key => {
            const lowerKey = key.toLowerCase();
            return lowerKey === 'id';
        }) || 'id';
    }

    // Filter list based on search value
    const filtered = list.filter(item => {
        const name = item[nameField] ? String(item[nameField]).toLowerCase() : '';
        const nik = item[nikField] ? String(item[nikField]) : '';
        return name.includes(value) || nik.includes(value);
    });

    if (filtered.length === 0) {
        dropdown.innerHTML = '<div class="autocomplete-item" style="color: #999;">Tidak ada hasil</div>';
        dropdown.style.display = 'block';
        return;
    }

    // Build dropdown HTML
    dropdown.innerHTML = filtered.map(item => {
        const name = item[nameField] || 'N/A';
        const nik = item[nikField] || 'N/A';
        const id = item[idField] || 0;
        
        const escapedName = String(name).replace(/'/g, "\\'");
        const escapedNik = String(nik).replace(/'/g, "\\'");
        
        return `
            <div class="autocomplete-item" onclick="selectItem('${cardId}', ${id}, '${escapedName}', '${escapedNik}')">
                <strong>${name}</strong><br>
                <small>NIK: ${nik}</small>
            </div>
        `;
    }).join('');

    dropdown.style.display = 'block';
}

// Select Item from Autocomplete
function selectItem(cardId, id, nama, nik) {
    const input = document.getElementById(`${cardId}-nama`);
    const dropdown = document.getElementById(`${cardId}-dropdown`);
    const displayValue = `${nama} (${nik})`;
    
    input.value = displayValue;
    input.setAttribute('data-id', id);
    input.setAttribute('data-display-value', displayValue);
    dropdown.style.display = 'none';
    
    console.log(`Selected for ${cardId}:`, {id, nama, nik});
}

// Remove Card
function removeCard(cardId) {
    const card = document.getElementById(cardId);
    if (card) {
        card.style.transition = 'all 0.3s ease';
        card.style.opacity = '0';
        card.style.transform = 'scale(0.8)';
        
        setTimeout(() => {
            card.remove();
            delete comparisonData[cardId];
            updateLayoutColumns();
            updateAddButtonState();
        }, 300);
    }
}

// Reset All Cards
function resetAll() {
    document.getElementById('comparisonWrapper').innerHTML = '';
    document.getElementById('comparisonResult').style.display = 'none';
    comparisonCount = 0;
    comparisonData = {};
    addComparisonCard();
    updateAddButtonState();
}

// Compare Data
async function compareData() {
    const cards = document.querySelectorAll('.comparison-card');
    
    if (cards.length < 2) {
        alert('❌ Tambahkan minimal 2 data untuk dibandingkan!');
        return;
    }

    comparisonData = {};
    let isValid = true;

    console.log('Starting comparison for', cards.length, 'cards');

    for (const card of cards) {
        const cardId = card.id;
        const input = document.getElementById(`${cardId}-nama`);
        const id = input.getAttribute('data-id');
        const tanggalAwal = document.getElementById(`${cardId}-tanggalAwal`).value;
        const tanggalAkhir = document.getElementById(`${cardId}-tanggalAkhir`).value;

        if (!id || !tanggalAwal || !tanggalAkhir) {
            alert(`❌ Lengkapi semua field yang wajib pada kartu ${cardId.replace('card-', '')}!`);
            isValid = false;
            break;
        }

        // Build API parameters
        const params = {
            [currentViewMode === 'penyadap' ? 'idPenyadap' : 'idMandor']: id,
            tanggalAwal: tanggalAwal,
            tanggalAkhir: tanggalAkhir
        };

        const tipeProduksi = document.getElementById(`${cardId}-tipe`).value;

        if (tipeProduksi && tipeProduksi.trim() !== '') {
            params.tipeProduksi = tipeProduksi;
        }

        try {
            const queryParams = new URLSearchParams(params);
            const url = `/api/search?${queryParams}`;
            console.log(`Fetching data for ${cardId}:`, url);
            
            const response = await fetch(url);
            
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}`);
            }
            
            const result = await response.json();
            console.log(`Data received for ${cardId}:`, result);

            comparisonData[cardId] = {
                nama: input.value,
                data: result.data || [],
                summary: result.summary || {}
            };
        } catch (error) {
            console.error(`Error fetching data for ${cardId}:`, error);
            alert(`❌ Gagal mengambil data untuk kartu ${cardId.replace('card-', '')}: ${error.message}`);
            isValid = false;
            break;
        }
    }

    if (isValid) {
        console.log('All data fetched successfully:', comparisonData);
        displayComparisonResult();
    }
}

// Create Pie Chart
function createComparisonChart() {
    // Check if Chart.js is loaded
    if (typeof Chart === 'undefined') {
        console.error('Chart.js not loaded');
        return;
    }

    const canvas = document.getElementById('comparisonChart');
    if (!canvas) {
        console.error('Canvas element not found');
        return;
    }

    const ctx = canvas.getContext('2d');

    // Destroy existing chart if any
    if (chartInstance) {
        chartInstance.destroy();
    }

    const labels = [];
    const data = [];
    const backgroundColors = [
        'rgba(102, 126, 234, 0.8)',
        'rgba(245, 87, 108, 0.8)',
        'rgba(79, 172, 254, 0.8)',
        'rgba(240, 147, 251, 0.8)'
    ];
    const borderColors = [
        'rgba(102, 126, 234, 1)',
        'rgba(245, 87, 108, 1)',
        'rgba(79, 172, 254, 1)',
        'rgba(240, 147, 251, 1)'
    ];

    // Collect data based on mode
    Object.values(comparisonData).forEach((item, index) => {
        labels.push(item.nama);
        
        if (currentViewMode === 'penyadap') {
            const value = item.summary.total_basah_latek || item.summary.TotalLatek || 0;
            data.push(value);
        } else {
            const value = item.summary.total_hko || item.summary.TotalHKO || 0;
            data.push(value);
        }
    });

    try {
        chartInstance = new Chart(ctx, {
            type: 'pie',
            data: {
                labels: labels,
                datasets: [{
                    label: currentViewMode === 'penyadap' ? 'Basah Latek' : 'Total HKO',
                    data: data,
                    backgroundColor: backgroundColors,
                    borderColor: borderColors,
                    borderWidth: 2
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: true,
                plugins: {
                    legend: {
                        position: 'bottom',
                        labels: {
                            padding: 15,
                            font: {
                                size: 12,
                                family: "'Segoe UI', Tahoma, Geneva, Verdana, sans-serif"
                            }
                        }
                    },
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                let label = context.label || '';
                                if (label) {
                                    label += ': ';
                                }
                                label += context.parsed.toLocaleString();
                                
                                // Calculate percentage
                                const total = context.dataset.data.reduce((a, b) => a + b, 0);
                                const percentage = ((context.parsed / total) * 100).toFixed(2);
                                label += ` (${percentage}%)`;
                                
                                return label;
                            }
                        }
                    }
                }
            }
        });
        
        console.log('Chart created successfully');
    } catch (error) {
        console.error('Error creating chart:', error);
    }
}


// Display Comparison Result
function displayComparisonResult() {
    const resultDiv = document.getElementById('comparisonResult');
    const resultContainer = document.getElementById('resultWrapper');

    resultContainer.innerHTML = '';

    if (currentViewMode === 'penyadap') {
        const totals = { latek: [], sheet: [] };

        Object.values(comparisonData).forEach(item => {
            const s = item.summary;
            totals.latek.push(s.total_basah_latek || s.TotalLatek || 0);
            totals.sheet.push(parseFloat(s.total_sheet || s.TotalSheet || 0));
        });

        const maxLatek = Math.max(...totals.latek);
        const maxSheet = Math.max(...totals.sheet);

        let rank = 0;
        Object.values(comparisonData).forEach(item => {
            rank++;
            const s = item.summary;
            const currentLatek = s.total_basah_latek || s.TotalLatek || 0;
            const currentSheet = parseFloat(s.total_sheet || s.TotalSheet || 0);
            
            const isTopLatek = currentLatek === maxLatek && currentLatek > 0;
            const isTopSheet = currentSheet === maxSheet && currentSheet > 0;
            
            const cardHTML = `
                <div class="result-card">
                    <div class="result-card-header">
                        <div class="result-rank">${rank}</div>
                        <div class="result-card-title">${item.nama}</div>
                    </div>
                    <div class="result-card-body">
                        <div class="result-item">
                            <div class="result-item-label">Total Records</div>
                            <div class="result-item-value">${s.total_records || s.TotalRecords || 0}</div>
                        </div>
                        <div class="result-item ${isTopLatek ? 'highlight' : ''}">
                            <div class="result-item-label">Basah Latek</div>
                            <div class="result-item-value">
                                ${currentLatek}
                                ${isTopLatek ? '<span class="winner-badge">🏆 Terbaik</span>' : ''}
                            </div>
                        </div>
                        <div class="result-item ${isTopSheet ? 'highlight' : ''}">
                            <div class="result-item-label">Sheet</div>
                            <div class="result-item-value">
                                ${currentSheet.toFixed(2)}
                                ${isTopSheet ? '<span class="winner-badge">🏆 Terbaik</span>' : ''}
                            </div>
                        </div>
                        <div class="result-item">
                            <div class="result-item-label">Basah Lump</div>
                            <div class="result-item-value">${s.total_basah_lump || s.TotalLump || 0}</div>
                        </div>
                        <div class="result-item">
                            <div class="result-item-label">Br.Cr</div>
                            <div class="result-item-value">${parseFloat(s.total_br_cr || s.TotalBrCr || 0).toFixed(2)}</div>
                        </div>
                    </div>
                </div>
            `;
            
            resultContainer.insertAdjacentHTML('beforeend', cardHTML);
        });
    } else {
        const totals = { hko: [], rata: [] };

        Object.values(comparisonData).forEach(item => {
            const s = item.summary;
            totals.hko.push(s.total_hko || s.TotalHKO || 0);
            totals.rata.push(parseFloat(s.rata_rata_produksi_per_taper || s.RataRataProduksiPerTaper || 0));
        });

        const maxHKO = Math.max(...totals.hko);
        const maxRata = Math.max(...totals.rata);

        let rank = 0;
        Object.values(comparisonData).forEach(item => {
            rank++;
            const s = item.summary;
            const currentHKO = s.total_hko || s.TotalHKO || 0;
            const currentRata = parseFloat(s.rata_rata_produksi_per_taper || s.RataRataProduksiPerTaper || 0);
            
            const isTopHKO = currentHKO === maxHKO && currentHKO > 0;
            const isTopRata = currentRata === maxRata && currentRata > 0;
            
            const cardHTML = `
                <div class="result-card">
                    <div class="result-card-header">
                        <div class="result-rank">${rank}</div>
                        <div class="result-card-title">${item.nama}</div>
                    </div>
                    <div class="result-card-body">
                        <div class="result-item">
                            <div class="result-item-label">Total Records</div>
                            <div class="result-item-value">${s.total_records || s.TotalRecords || 0}</div>
                        </div>
                        <div class="result-item ${isTopHKO ? 'highlight' : ''}">
                            <div class="result-item-label">Total HKO</div>
                            <div class="result-item-value">
                                ${currentHKO}
                                ${isTopHKO ? '<span class="winner-badge">🏆 Terbaik</span>' : ''}
                            </div>
                        </div>
                        <div class="result-item">
                            <div class="result-item-label">Latek Kebun</div>
                            <div class="result-item-value">${s.total_basah_latek_kebun || s.TotalBasahLatekKebun || 0}</div>
                        </div>
                        <div class="result-item">
                            <div class="result-item-label">Latek Pabrik</div>
                            <div class="result-item-value">${s.total_basah_latek_pabrik || s.TotalBasahLatekPabrik || 0}</div>
                        </div>
                        <div class="result-item">
                            <div class="result-item-label">Sheet Kering</div>
                            <div class="result-item-value">${s.total_kering_sheet || s.TotalKeringSheet || 0}</div>
                        </div>
                        <div class="result-item ${isTopRata ? 'highlight' : ''}">
                            <div class="result-item-label">Rata² Per Taper</div>
                            <div class="result-item-value">
                                ${currentRata.toFixed(2)}
                                ${isTopRata ? '<span class="winner-badge">🏆 Terbaik</span>' : ''}
                            </div>
                        </div>
                    </div>
                </div>
            `;
            
            resultContainer.insertAdjacentHTML('beforeend', cardHTML);
        });
    }

    updateResultColumns();
    resultDiv.style.display = 'block';
    
    // Create pie chart with slight delay to ensure DOM is ready
    setTimeout(() => {
        createComparisonChart();
    }, 100);
    
    resultDiv.scrollIntoView({ behavior: 'smooth', block: 'start' });
    
    console.log('Comparison result displayed');
}
