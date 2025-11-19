// State untuk menyimpan data
let currentData = [];
let currentApiData = [];
let selectedVisualisationType = null; // 'rekap' atau 'produksi'
let penyadapList = [];
let mandorList = [];
let currentConfig = {
    threshold: 150,
    title: 'Grafik Produksi',
    field: 'hko',
    showTrend: true,
    showThreshold: true,
    satuan: 'day'
}

// ========== HELPER FUNCTIONS ==========

/**
 * Generate array of all dates between start and end
 */
function generateDateRange(startDate, endDate) {
    const dates = [];
    const current = new Date(startDate);
    const end = new Date(endDate);

    while (current <= end) {
        dates.push(new Date(current));
        current.setDate(current.getDate() + 1);
    }

    return dates;
}

/**
 * Format date to match API format (YYYY-MM-DD)
 */
function formatDateForAPI(date) {
    return date.toISOString().split('T')[0];
}

/**
 * Format date for display based on satuan
 */
function formatDateForDisplay(date, satuanWaktu) {
    if (satuanWaktu === 'day') {
        return date.toLocaleDateString('id-ID', {
            day: '2-digit',
            month: 'short',
            year: 'numeric'
        });
    } else if (satuanWaktu === 'week') {
        const weekNum = Math.ceil(date.getDate() / 7);
        return `Minggu ${weekNum} ${date.toLocaleDateString('id-ID', { month: 'short', year: 'numeric' })}`;
    } else if (satuanWaktu === 'month') {
        return date.toLocaleDateString('id-ID', {
            month: 'long',
            year: 'numeric'
        });
    }
    return formatDateForAPI(date);
}

// ========== HANDLER FUNCTIONS ==========

function handleTipeDataChange() {
    const tipeData = document.getElementById('tipeData').value;
    const afdelingGroup = document.getElementById('afdelingGroup');
    const mandorGroup = document.getElementById('mandorGroup');
    const penyadapGroup = document.getElementById('penyadapGroup');
    const tipeProduksiGroup = document.getElementById('tipeProduksiGroup');

    afdelingGroup.style.display = 'none';
    mandorGroup.style.display = 'none';
    penyadapGroup.style.display = 'none';
    tipeProduksiGroup.style.display = 'none';

    if (tipeData === 'afdeling') {
        afdelingGroup.style.display = 'block';
    } else if (tipeData === 'mandor') {
        mandorGroup.style.display = 'block';
    } else if (tipeData === 'penyadap') {
        penyadapGroup.style.display = 'block';
        tipeProduksiGroup.style.display = 'block';
    }
}

async function fetchVisualisasiData(params) {
    try {
        const queryString = new URLSearchParams(params).toString();
        const url = `/api/visualisasi?${queryString}`;

        console.log('Fetching from:', url);
        console.log('Params:', params);

        const response = await fetch(url);

        console.log('Response status:', response.status);
        console.log('Response headers:', response.headers.get('content-type'));

        if (!response.ok) {
            const errorText = await response.text();
            console.error('API Error Response:', errorText);
            throw new Error(`HTTP ${response.status}: ${errorText}`);
        }

        const responseText = await response.text();
        console.log('Raw response:', responseText.substring(0, 500));

        let result;
        try {
            result = JSON.parse(responseText);
        } catch (parseError) {
            console.error('JSON Parse Error:', parseError);
            console.error('Response text:', responseText);
            throw new Error(`Response bukan JSON valid. Response: ${responseText.substring(0, 100)}...`);
        }

        console.log('Parsed result:', result);
        console.log('Result type:', typeof result);
        console.log('Is array?', Array.isArray(result));
        if (result && typeof result === 'object') {
            console.log('Object keys:', Object.keys(result));
        }

        if (result && result.success && result.data) {
            console.log('✅ Format 1: {success, data}');
            return result.data;
        }
        else if (Array.isArray(result)) {
            console.log('✅ Format 2: Array langsung');
            return result;
        }
        else if (result && result.labels && Array.isArray(result.labels)) {
            console.log('✅ Format 3: {labels, data}');
            return result.labels.map((label, index) => ({
                tanggal: label,
                value: result.data && result.data[index] ? result.data[index].value : 0
            }));
        }
        else if (result && result.success === false) {
            throw new Error(result.message || 'Gagal mengambil data dari server');
        }
        else {
            console.error('❌ Format tidak dikenali. Result:', result);
            throw new Error('Format response tidak sesuai. Cek console untuk detail.');
        }
    } catch (error) {
        console.error('Error fetching data:', error);
        alert('Gagal mengambil data: ' + error.message);
        return [];
    }
}

function transformDataForChart(apiData, fieldName, satuanWaktu, startDate, endDate) {
    // Generate all dates in range
    const allDates = generateDateRange(new Date(startDate), new Date(endDate));
    console.log('📅 Total dates in range:', allDates.length);

    // Create map of API data by date
    const apiDataMap = {};

    if (apiData.length > 0 && apiData[0].hasOwnProperty('value') && apiData[0].hasOwnProperty('tanggal')) {
        console.log('✅ Data already aggregated from API');

        apiData.forEach(item => {
            const dateKey = formatDateForAPI(new Date(item.tanggal));
            if (!apiDataMap[dateKey]) {
                apiDataMap[dateKey] = [];
            }
            apiDataMap[dateKey].push(parseFloat(item.value) || 0);
        });
    } else {
        console.log('⚠️ Using fallback transformation');

        const fieldMappingAPI = {
            'hko': 'hko',
            'basah_latek_kebun': 'basahLatexKebun',
            'basah_latek_pabrik': 'basahLatexPabrik',
            'basah_latek_persen': 'basahLatexPersen',
            'basah_lump_kebun': 'basahLumpKebun',
            'basah_lump_pabrik': 'basahLumpPabrik',
            'basah_lump_persen': 'basahLumpPersen',
            'k3_sheet': 'k3Sheet',
            'kering_sheet': 'keringSheet',
            'kering_br_cr': 'keringBrCr',
            'kering_jumlah': 'keringJumlah',
            'produksi_per_taper': 'produksiPerTaper'
        };

        const apiFieldName = fieldMappingAPI[fieldName] || fieldName;

        apiData.forEach(item => {
            const dateKey = formatDateForAPI(new Date(item.tanggal));
            let value = item[apiFieldName] || item[fieldName] || 0;

            if (!apiDataMap[dateKey]) {
                apiDataMap[dateKey] = [];
            }
            apiDataMap[dateKey].push(parseFloat(value) || 0);
        });
    }

    // Build complete dataset with all dates
    const result = [];
    const grouped = {};

    allDates.forEach(date => {
        const dateKey = formatDateForAPI(date);
        const displayLabel = formatDateForDisplay(date, satuanWaktu);

        // Get value for this date or default to 0
        let value = 0;
        if (apiDataMap[dateKey] && apiDataMap[dateKey].length > 0) {
            // Sum all values for this date
            value = apiDataMap[dateKey].reduce((sum, val) => sum + val, 0);

            // For non-day views, we'll average later
            if (satuanWaktu !== 'day') {
                if (!grouped[displayLabel]) {
                    grouped[displayLabel] = { total: 0, count: 0 };
                }
                grouped[displayLabel].total += value;
                grouped[displayLabel].count += 1;
            }
        }

        if (satuanWaktu === 'day') {
            result.push({
                label: displayLabel,
                value: value
            });
        }
    });

    // For week/month aggregation
    if (satuanWaktu !== 'day') {
        const aggregated = Object.entries(grouped).map(([label, data]) => ({
            label: label,
            value: Math.round(data.total / data.count)
        }));
        console.log('📊 Aggregated data:', aggregated.length, 'points');
        return aggregated;
    }

    console.log('📊 Complete dataset:', result.length, 'days');
    console.log('📊 Days with data:', Object.keys(apiDataMap).length);
    console.log('📊 Days with zero:', result.filter(d => d.value === 0).length);

    return result;
}

function drawBarChart() {
    const canvas = document.getElementById('barChart');
    const ctx = canvas.getContext('2d');

    if (currentData.length === 0) {
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        ctx.fillStyle = '#666';
        ctx.font = '16px Arial';
        ctx.textAlign = 'center';
        ctx.fillText('Tidak ada data untuk ditampilkan', canvas.width / 2, canvas.height / 2);
        return;
    }

    const threshold = currentConfig.threshold;
    const showTrend = currentConfig.showTrend;
    const showThreshold = currentConfig.showThreshold;

    canvas.width = Math.max(600, currentData.length * 80);
    canvas.height = 300;

    ctx.clearRect(0, 0, canvas.width, canvas.height);

    const barWidth = 50;
    const barSpacing = 80;
    const chartHeight = 250;
    const chartTop = 30;
    const values = currentData.map(d => d.value);
    const maxValue = Math.max(...values, threshold) * 1.1;

    currentData.forEach((item, index) => {
        const value = item.value;
        const barHeight = (value / maxValue) * chartHeight;
        const x = index * barSpacing + 30;
        const y = chartTop + chartHeight - barHeight;
        const prevValue = index > 0 ? currentData[index - 1].value : null;

        let color = '#28a745';

        // Color logic
        if (value === 0) {
            color = '#cccccc'; // Gray for zero values
        } else if (value < threshold) {
            color = '#d43636';
        } else if (prevValue !== null && value < prevValue) {
            color = '#FFD700';
        }

        ctx.fillStyle = color;
        ctx.fillRect(x, y, barWidth, barHeight);

        ctx.strokeStyle = 'rgba(0,0,0,0.1)';
        ctx.lineWidth = 1;
        ctx.strokeRect(x, y, barWidth, barHeight);

        // Only show value if > 0
        if (value > 0) {
            ctx.fillStyle = '#333';
            ctx.font = 'bold 12px Arial';
            ctx.textAlign = 'center';
            ctx.fillText(value, x + barWidth/2, y - 5);
        }

        ctx.font = '10px Arial';
        const maxLabelWidth = barSpacing - 10;
        const label = item.label;

        const words = label.split(' ');
        let line = '';
        let lineY = chartTop + chartHeight + 15;

        words.forEach((word, i) => {
            const testLine = line + word + ' ';
            const metrics = ctx.measureText(testLine);

            if (metrics.width > maxLabelWidth && i > 0) {
                ctx.fillText(line, x + barWidth/2, lineY);
                line = word + ' ';
                lineY += 12;
            } else {
                line = testLine;
            }
        });
        ctx.fillText(line, x + barWidth/2, lineY);

        if (showTrend && prevValue !== null && value > 0) {
            const trendY = y - 20;
            ctx.font = '14px Arial';
            if (value > prevValue) {
                ctx.fillStyle = '#28a745';
                ctx.fillText('↗', x + barWidth/2, trendY);
            } else if (value < prevValue) {
                ctx.fillStyle = '#ff6b35';
                ctx.fillText('↘', x + barWidth/2, trendY);
            } else {
                ctx.fillStyle = '#6c757d';
                ctx.fillText('→', x + barWidth/2, trendY);
            }
        }
    });

    if (showThreshold) {
        const thresholdY = chartTop + chartHeight - (threshold / maxValue) * chartHeight;
        ctx.strokeStyle = '#ff6b35';
        ctx.lineWidth = 3;
        ctx.setLineDash([8, 4]);
        ctx.beginPath();
        ctx.moveTo(20, thresholdY);
        ctx.lineTo(canvas.width - 20, thresholdY);
        ctx.stroke();
        ctx.setLineDash([]);

        ctx.fillStyle = 'rgba(255, 107, 53, 0.9)';
        ctx.fillRect(25, thresholdY - 18, 100, 16);
        ctx.fillStyle = 'white';
        ctx.font = 'bold 11px Arial';
        ctx.textAlign = 'left';
        ctx.fillText(`Target: ${threshold}`, 28, thresholdY - 7);
    }
}

async function handleUpdateGrafik() {
    const tipeData = document.getElementById('tipeData').value;
    const kodeAfdeling = document.getElementById('kodeAfdeling').value;
    const idMandor = document.getElementById('idMandor').value;
    const idPenyadap = document.getElementById('idPenyadap').value;
    const tipeProduksi = document.getElementById('tipeProduksi').value;
    const field = document.getElementById('fieldSelect').value;
    const tanggalAwal = document.getElementById('tanggalAwal').value;
    const tanggalAkhir = document.getElementById('tanggalAkhir').value;
    const satuanWaktu = currentConfig.satuan;

    if (!tanggalAwal || !tanggalAkhir) {
        alert('Tanggal awal dan akhir harus diisi');
        return;
    }

    if (tipeData === 'afdeling' && !kodeAfdeling) {
        alert('Silakan pilih Kode Afdeling terlebih dahulu');
        return;
    }

    if (tipeData === 'mandor' && !idMandor) {
        alert('Silakan pilih Mandor terlebih dahulu');
        return;
    }

    if (tipeData === 'penyadap' && !idPenyadap) {
        alert('Silakan pilih Penyadap terlebih dahulu');
        return;
    }

    currentConfig.field = field;
    currentConfig.threshold = parseInt(document.getElementById('barThreshold').value) || 150;

    const canvas = document.getElementById('barChart');
    const ctx = canvas.getContext('2d');
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    ctx.fillStyle = '#666';
    ctx.font = '16px Arial';
    ctx.textAlign = 'center';
    ctx.fillText('Memuat data...', canvas.width / 2, canvas.height / 2);

    let satuan;

    if (tipeData === 'penyadap') {
        satuan = field;
        console.log('📊 PRODUKSI mode - satuan = field:', satuan);
    } else {
        const rekapSatuanMap = {
            'hko': 'hko',
            'basah_latek_kebun': 'basah_latek_kebun',
            'basah_latek_pabrik': 'basah_latek_pabrik',
            'basah_latek_persen': 'basah_latek_persen',
            'basah_lump_kebun': 'basah_lump_kebun',
            'basah_lump_pabrik': 'basah_lump_pabrik',
            'basah_lump_persen': 'basah_lump_persen',
            'k3_sheet': 'k3_sheet',
            'kering_sheet': 'kering_sheet',
            'kering_br_cr': 'kering_br_cr',
            'kering_jumlah': 'kering_jumlah',
            'produksi_per_taper': 'produksi_per_taper'
        };

        satuan = rekapSatuanMap[field];

        if (!satuan) {
            console.error('⚠️ Satuan tidak ditemukan untuk field:', field);
            satuan = field;
        }

        console.log('📊 REKAP mode - satuan mapped:', satuan);
    }

    let params = {};

    if (tipeData === 'penyadap') {
        params = {
            tipeData: 'penyadap',
            idPenyadap: idPenyadap,
            tanggalAwal: tanggalAwal,
            tanggalAkhir: tanggalAkhir,
            satuan: satuan
        };

        if (tipeProduksi) {
            params.tipeProduksi = tipeProduksi;
        }
    } else {
        params = {
            tipeData: tipeData,
            tanggalAwal: tanggalAwal,
            tanggalAkhir: tanggalAkhir,
            satuan: satuan
        };

        if (tipeData === 'mandor') {
            params.idMandor = idMandor;
        }

        if (tipeData === 'afdeling') {
            params.afdeling = kodeAfdeling;
        }
    }

    console.log('=== VISUALISASI REQUEST ===');
    console.log('Tipe Visualisasi:', selectedVisualisationType);
    console.log('Tipe Data:', tipeData);
    console.log('Field Selected:', field);
    console.log('Satuan Extracted:', satuan);
    console.log('Date Range:', tanggalAwal, 'to', tanggalAkhir);
    console.log('Parameters:', params);
    console.log('URL:', `/api/visualisasi?${new URLSearchParams(params).toString()}`);
    console.log('===========================');

    const apiData = await fetchVisualisasiData(params);

    // Always transform with complete date range
    currentApiData = apiData;
    currentData = transformDataForChart(apiData, field, satuanWaktu, tanggalAwal, tanggalAkhir);

    console.log('Transformed data:', currentData);

    const fieldNames = {
        'hko': 'HKO',
        'basah_latek_kebun': 'Basah Latek Kebun',
        'basah_latek_pabrik': 'Basah Latek Pabrik',
        'basah_latek_persen': 'Basah Latek Persen',
        'basah_lump_kebun': 'Basah Lump Kebun',
        'basah_lump_pabrik': 'Basah Lump Pabrik',
        'basah_lump_persen': 'Basah Lump Persen',
        'k3_sheet': 'K3 Sheet',
        'kering_sheet': 'Kering Sheet',
        'kering_br_cr': 'Kering BR/CR',
        'kering_jumlah': 'Kering Jumlah',
        'produksi_per_taper': 'Produksi Per Taper',
        'basah_latek': 'Basah Latek',
        'sheet': 'Sheet',
        'basah_lump': 'Basah Lump',
        'br_cr': 'BR/CR'
    };

    let titleSuffix = '';
    if (tipeData === 'afdeling') {
        const selectedOption = document.getElementById('kodeAfdeling').selectedOptions[0];
        titleSuffix = ` - Afdeling: ${selectedOption.textContent}`;
    } else if (tipeData === 'mandor') {
        const mandorName = document.getElementById('mandorSearch').value;
        titleSuffix = ` - Mandor: ${mandorName}`;
    } else if (tipeData === 'penyadap') {
        const penyadapName = document.getElementById('penyadapSearch').value;
        titleSuffix = ` - Penyadap: ${penyadapName}`;
    } else {
        titleSuffix = ' - Total';
    }

    currentConfig.title = `${fieldNames[field]}${titleSuffix}`;
    document.getElementById('dynamicChartTitle').textContent = currentConfig.title;

    drawBarChart();
}

function handleFieldChange() {
    if (currentApiData.length === 0) {
        alert('Silakan update grafik terlebih dahulu dengan klik tombol "Update Grafik"');
        return;
    }

    const field = document.getElementById('fieldSelect').value;
    const tanggalAwal = document.getElementById('tanggalAwal').value;
    const tanggalAkhir = document.getElementById('tanggalAkhir').value;

    currentConfig.field = field;
    currentData = transformDataForChart(currentApiData, field, currentConfig.satuan, tanggalAwal, tanggalAkhir);

    const fieldNames = {
        'hko': 'HKO',
        'basah_latek_kebun': 'Basah Latek Kebun',
        'basah_latek_pabrik': 'Basah Latek Pabrik',
        'basah_latek_persen': 'Basah Latek Persen',
        'basah_lump_kebun': 'Basah Lump Kebun',
        'basah_lump_pabrik': 'Basah Lump Pabrik',
        'basah_lump_persen': 'Basah Lump Persen',
        'k3_sheet': 'K3 Sheet',
        'kering_sheet': 'Kering Sheet',
        'kering_br_cr': 'Kering BR/CR',
        'kering_jumlah': 'Kering Jumlah',
        'produksi_per_taper': 'Produksi Per Taper',
        'basah_latek': 'Basah Latek',
        'sheet': 'Sheet',
        'basah_lump': 'Basah Lump',
        'br_cr': 'BR/CR'
    };

    const tipeData = document.getElementById('tipeData').value;
    let titleSuffix = '';

    if (tipeData === 'afdeling') {
        const selectedOption = document.getElementById('kodeAfdeling').selectedOptions[0];
        titleSuffix = ` - Afdeling: ${selectedOption.textContent}`;
    } else if (tipeData === 'mandor') {
        const mandorName = document.getElementById('mandorSearch').value;
        titleSuffix = ` - Mandor: ${mandorName}`;
    } else if (tipeData === 'penyadap') {
        const penyadapName = document.getElementById('penyadapSearch').value;
        titleSuffix = ` - Penyadap: ${penyadapName}`;
    } else {
        titleSuffix = ' - Total';
    }

    currentConfig.title = `${fieldNames[field]}${titleSuffix}`;
    document.getElementById('dynamicChartTitle').textContent = currentConfig.title;

    drawBarChart();
}

function setVisualRange(range) {
    const container = document.querySelector('.time-toggle');
    if (!container) return;

    const buttons = Array.from(container.querySelectorAll('.time-btn'));
    buttons.forEach(btn => {
        if (btn.getAttribute('data-range') === range) {
            btn.classList.add('active');
        } else {
            btn.classList.remove('active');
        }
    });

    currentConfig.satuan = range;

    if (currentApiData.length > 0) {
        const tanggalAwal = document.getElementById('tanggalAwal').value;
        const tanggalAkhir = document.getElementById('tanggalAkhir').value;
        currentData = transformDataForChart(currentApiData, currentConfig.field, range, tanggalAwal, tanggalAkhir);
        drawBarChart();
    }
}

// ========== EVENT LISTENERS ==========

document.getElementById('tipeData').addEventListener('change', handleTipeDataChange);
document.getElementById('fieldSelect').addEventListener('change', handleFieldChange);

document.getElementById('barThreshold').addEventListener('input', function() {
    currentConfig.threshold = parseInt(this.value) || 150;
    drawBarChart();
});

document.getElementById('showTrend').addEventListener('change', function() {
    currentConfig.showTrend = this.checked;
    drawBarChart();
});

document.getElementById('showThreshold').addEventListener('change', function() {
    currentConfig.showThreshold = this.checked;
    drawBarChart();
});

document.addEventListener('DOMContentLoaded', function() {
    const container = document.querySelector('.time-toggle');
    if (!container) return;

    container.addEventListener('click', function(e){
        const btn = e.target.closest('.time-btn');
        if (!btn) return;
        const range = btn.getAttribute('data-range');
        if (range) setVisualRange(range);
    });

    const titleEl = document.getElementById('dynamicChartTitle');
    if (titleEl) {
        titleEl.style.userSelect = 'none';
        titleEl.style.webkitUserSelect = 'none';
        titleEl.addEventListener('click', function(e){
            e.stopPropagation();
            e.preventDefault();
        });
        titleEl.addEventListener('dblclick', function(e){
            e.stopPropagation();
            e.preventDefault();
        });
    }
});

// ========== INITIALIZATION ==========

window.onload = async function() {
    const now = new Date();
    const today = now.toISOString().split('T')[0];

    document.getElementById('tanggalAwal').value = today;
    document.getElementById('tanggalAkhir').value = today;

    await populateMandorDropdown();
    await populatePenyadapDropdown();

    showModal();
};


// ========== POPUP MODAL FUNCTIONS ==========

function showModal() {
    const modal = document.getElementById('typeModal');
    modal.classList.add('active');

    document.querySelectorAll('.option-card').forEach(card => {
        card.classList.remove('selected');
    });
}

function closeModal() {
    const modal = document.getElementById('typeModal');
    modal.classList.remove('active');
}

function selectAndConfirm(type) {
    selectedVisualisationType = type;

    document.querySelectorAll('.option-card').forEach(card => {
        card.classList.remove('selected');
    });
    document.querySelector(`[data-type="${type}"]`).classList.add('selected');

    setTimeout(() => {
        closeModal();
        initializeVisualization();
    }, 300);
}

// ========== FIELD OPTIONS MANAGEMENT ==========

function updateFieldOptions() {
    const fieldSelect = document.getElementById('fieldSelect');

    if (selectedVisualisationType === 'produksi') {
        fieldSelect.innerHTML = `
            <option value="basah_latek">Basah Latek</option>
            <option value="sheet">Sheet</option>
            <option value="basah_lump">Basah Lump</option>
            <option value="br_cr">BR/CR</option>
        `;
        console.log('✅ Field options set to PRODUKSI mode');
    } else {
        fieldSelect.innerHTML = `
            <option value="hko">HKO</option>
            <option value="basah_latek_kebun">Basah Latek Kebun</option>
            <option value="basah_latek_pabrik">Basah Latek Pabrik</option>
            <option value="basah_latek_persen">Basah Latek Persen</option>
            <option value="basah_lump_kebun">Basah Lump Kebun</option>
            <option value="basah_lump_pabrik">Basah Lump Pabrik</option>
            <option value="basah_lump_persen">Basah Lump Persen</option>
            <option value="k3_sheet">K3 Sheet</option>
            <option value="kering_sheet">Kering Sheet</option>
            <option value="kering_br_cr">Kering BR/CR</option>
            <option value="kering_jumlah">Kering Jumlah</option>
            <option value="produksi_per_taper">Produksi Per Taper</option>
        `;
        console.log('✅ Field options set to REKAP mode');
    }
}

function initializeVisualization() {
    const mainContainer = document.getElementById('mainContainer');
    const typeLabel = document.getElementById('typeLabel');
    const tipeDataSelect = document.getElementById('tipeData');
    const tipeDataGroup = document.getElementById('tipeDataGroup');

    mainContainer.style.display = 'block';

    if (selectedVisualisationType === 'rekap') {
        typeLabel.textContent = '- REKAP';
        typeLabel.style.color = '#0093E9';

        tipeDataSelect.innerHTML = `
            <option value="total">Total</option>
            <option value="afdeling">Afdeling</option>
            <option value="mandor">Mandor</option>
        `;
        tipeDataGroup.style.display = 'block';

        setupMandorAutocomplete();
        updateFieldOptions();
        handleTipeDataChange();

    } else if (selectedVisualisationType === 'produksi') {
        typeLabel.textContent = '- PRODUKSI';
        typeLabel.style.color = '#FF6B6B';

        tipeDataSelect.innerHTML = `
            <option value="penyadap">Penyadap</option>
        `;
        tipeDataGroup.style.display = 'none';

        setupPenyadapAutocomplete();
        updateFieldOptions();
        handleTipeDataChange();
    }
}

// ========== API FUNCTIONS ==========

async function populateMandorDropdown() {
    try {
        const response = await fetch('/api/mandor');
        const result = await response.json();

        if (result.success && result.data) {
            mandorList = result.data;
            console.log('Mandor list loaded:', mandorList.length, 'items');
        }
    } catch (error) {
        console.error('Error fetching mandor list:', error);
        mandorList = [];
    }
}

function setupMandorAutocomplete() {
    const searchInput = document.getElementById('mandorSearch');
    const hiddenInput = document.getElementById('idMandor');
    const dropdown = document.getElementById('mandorDropdown');

    if (!searchInput || !dropdown) return;

    searchInput.addEventListener('input', function() {
        const value = this.value.toLowerCase();

        const storedDisplayValue = searchInput.getAttribute('data-display-value');
        if (storedDisplayValue && searchInput.value === storedDisplayValue) {
            dropdown.style.display = 'none';
            return;
        }

        if (storedDisplayValue && searchInput.value !== storedDisplayValue) {
            searchInput.removeAttribute('data-nik');
            searchInput.removeAttribute('data-nama');
            searchInput.removeAttribute('data-display-value');
            hiddenInput.value = '';
        }

        if (value.length < 2) {
            dropdown.style.display = 'none';
            return;
        }

        if (mandorList.length === 0) {
            dropdown.innerHTML = '<div class="autocomplete-item" style="color: #e74c3c; cursor: default;">⚠️ Data mandor kosong</div>';
            dropdown.style.display = 'block';
            return;
        }

        const filtered = mandorList.filter(m =>
            m.nama.toLowerCase().includes(value) ||
            (m.nik && m.nik.toLowerCase().includes(value))
        );

        if (filtered.length === 0) {
            dropdown.innerHTML = '<div class="autocomplete-item" style="color: #999; cursor: default;">Tidak ada hasil ditemukan</div>';
            dropdown.style.display = 'block';
            return;
        }

        dropdown.innerHTML = '';
        filtered.forEach(mandor => {
            const item = document.createElement('div');
            item.className = 'autocomplete-item';
            item.innerHTML = `
                <strong>${mandor.nama}</strong><br>
                <small>NIK: ${mandor.nik}</small>
            `;

            item.onclick = function() {
                const displayValue = `${mandor.nama} (${mandor.nik})`;
                searchInput.setAttribute('data-nik', mandor.nik);
                searchInput.setAttribute('data-nama', mandor.nama);
                searchInput.setAttribute('data-display-value', displayValue);
                searchInput.value = displayValue;
                hiddenInput.value = mandor.id;
                dropdown.style.display = 'none';
                console.log('Mandor selected:', {nik: mandor.nik, nama: mandor.nama});
            };

            dropdown.appendChild(item);
        });

        dropdown.style.display = 'block';
    });

    document.addEventListener('click', function(e) {
        if (!searchInput.contains(e.target) && !dropdown.contains(e.target)) {
            dropdown.style.display = 'none';
        }
    });

    searchInput.addEventListener('keydown', function(e) {
        if (e.key === 'Enter') {
            e.preventDefault();
            const firstItem = dropdown.querySelector('.autocomplete-item');
            if (firstItem && firstItem.onclick) {
                firstItem.click();
            }
        }
    });
}

async function populatePenyadapDropdown() {
    try {
        const response = await fetch('/api/penyadap');
        const result = await response.json();

        if (result.success && result.data) {
            penyadapList = result.data;
            console.log('Penyadap list loaded:', penyadapList.length, 'items');
        } else if (Array.isArray(result)) {
            penyadapList = result;
            console.log('Penyadap list loaded:', penyadapList.length, 'items');
        }
    } catch (error) {
        console.error('Error fetching penyadap list:', error);
        penyadapList = [];
    }
}

function setupPenyadapAutocomplete() {
    const searchInput = document.getElementById('penyadapSearch');
    const hiddenInput = document.getElementById('idPenyadap');
    const dropdown = document.getElementById('penyadapDropdown');

    if (!searchInput || !dropdown) {
        console.error('❌ Penyadap autocomplete elements not found!');
        return;
    }

    console.log('✅ Penyadap autocomplete initialized');

    searchInput.addEventListener('input', function() {
        const value = this.value.toLowerCase();

        console.log('🔍 Penyadap search input:', value);
        console.log('📊 Penyadap list length:', penyadapList.length);

        const storedDisplayValue = searchInput.getAttribute('data-display-value');
        if (storedDisplayValue && searchInput.value === storedDisplayValue) {
            dropdown.style.display = 'none';
            return;
        }

        if (storedDisplayValue && searchInput.value !== storedDisplayValue) {
            searchInput.removeAttribute('data-id');
            searchInput.removeAttribute('data-nama');
            searchInput.removeAttribute('data-nik');
            searchInput.removeAttribute('data-display-value');
            hiddenInput.value = '';
        }

        if (value.length < 2) {
            dropdown.style.display = 'none';
            return;
        }

        if (penyadapList.length === 0) {
            dropdown.innerHTML = '<div class="autocomplete-item" style="color: #e74c3c; cursor: default;">⚠️ Data penyadap kosong. Cek endpoint /api/penyadap</div>';
            dropdown.style.display = 'block';
            console.warn('⚠️ Penyadap list is empty!');
            return;
        }

        if (penyadapList.length > 0) {
            console.log('🔑 First penyadap object keys:', Object.keys(penyadapList[0]));
            console.log('👤 First penyadap full object:', penyadapList[0]);
        }

        const firstItem = penyadapList[0];
        const nameField = Object.keys(firstItem).find(key => {
            const lowerKey = key.toLowerCase();
            return lowerKey === 'nama' ||
                lowerKey === 'nama_penyadap' ||
                lowerKey === 'namapenyadap';
        }) || 'nama';

        const nikField = Object.keys(firstItem).find(key => {
            const lowerKey = key.toLowerCase();
            return lowerKey === 'nik';
        }) || 'nik';

        const idField = Object.keys(firstItem).find(key => {
            const lowerKey = key.toLowerCase();
            return lowerKey === 'id';
        }) || 'id';

        console.log('🎯 Detected fields:', {nameField, nikField, idField});

        const filtered = penyadapList.filter(p => {
            const penyadapName = p[nameField] ? String(p[nameField]) : '';
            const nikValue = p[nikField] ? String(p[nikField]) : '';

            const nameMatch = penyadapName.toLowerCase().includes(value);
            const nikMatch = nikValue.includes(value);

            return nameMatch || nikMatch;
        });

        console.log('✅ Filtered penyadap results:', filtered.length);

        if (filtered.length === 0) {
            dropdown.innerHTML = `<div class="autocomplete-item" style="color: #999; cursor: default;">
                Tidak ada hasil ditemukan<br>
                <small>Total data: ${penyadapList.length} | Mencari: "${value}"</small>
            </div>`;
            dropdown.style.display = 'block';
            return;
        }

        dropdown.innerHTML = '';
        filtered.forEach(penyadap => {
            const penyadapName = penyadap[nameField] || 'N/A';
            const nikValue = penyadap[nikField] || '';
            const idValue = penyadap[idField] || 0;

            const item = document.createElement('div');
            item.className = 'autocomplete-item';
            item.innerHTML = `
                <strong>${penyadapName}</strong><br>
                <small>${nikValue ? 'NIK: ' + nikValue : 'ID: ' + idValue}</small>
            `;

            item.onclick = function() {
                const displayValue = nikValue ? `${penyadapName} (${nikValue})` : penyadapName;
                searchInput.setAttribute('data-id', idValue);
                searchInput.setAttribute('data-nama', penyadapName);
                if (nikValue) searchInput.setAttribute('data-nik', nikValue);
                searchInput.setAttribute('data-display-value', displayValue);
                searchInput.value = displayValue;
                hiddenInput.value = idValue;
                dropdown.style.display = 'none';
                console.log('✅ Penyadap selected:', {id: idValue, nama: penyadapName, nik: nikValue});
            };

            dropdown.appendChild(item);
        });

        dropdown.style.display = 'block';
    });

    document.addEventListener('click', function(e) {
        if (!searchInput.contains(e.target) && !dropdown.contains(e.target)) {
            dropdown.style.display = 'none';
        }
    });

    searchInput.addEventListener('keydown', function(e) {
        if (e.key === 'Enter') {
            e.preventDefault();
            const firstItem = dropdown.querySelector('.autocomplete-item');
            if (firstItem && firstItem.onclick) {
                firstItem.click();
            }
        }
    });
}