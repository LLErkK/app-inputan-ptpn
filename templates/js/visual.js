// State untuk menyimpan data
let currentData = [];
let currentApiData = [];
let currentConfig = {
    threshold: 150,
    title: 'Grafik Produksi',
    jenisData: '',
    field: 'jumlah_pabrik_basah_latek',
    showTrend: true,
    showThreshold: true
};

// Fungsi untuk fetch data dari API berdasarkan jenis
async function fetchDataByJenis(jenis, bulan, tahun) {
    try {
        let url = '/api/visualisasi/default';
        const params = [];
        
        if (jenis) {
            params.push(`tipe=${jenis}`);
        }
        if (bulan) {
            params.push(`bulan=${bulan}`);
        }
        if (tahun) {
            params.push(`tahun=${tahun}`);
        }
        
        if (params.length > 0) {
            url += '?' + params.join('&');
        }
        
        const response = await fetch(url);
        const result = await response.json();
        
        if (result.success && result.data) {
            return result.data;
        } else {
            throw new Error(result.message || 'Gagal mengambil data');
        }
    } catch (error) {
        console.error('Error fetching data:', error);
        alert('Gagal mengambil data: ' + error.message);
        return [];
    }
}

// Fungsi untuk mengubah data API menjadi format chart
function transformDataForChart(apiData, fieldName) {
    // Group by tanggal dan sum nilai
    const grouped = {};
    
    apiData.forEach(item => {
        const date = new Date(item.tanggal).toLocaleDateString('id-ID', { 
            day: '2-digit', 
            month: 'short' 
        });
        
        if (!grouped[date]) {
            grouped[date] = 0;
        }
        
        grouped[date] += item[fieldName] || 0;
    });
    
    return Object.values(grouped).map(val => Math.round(val));
}

// Fungsi menggambar grafik batang dengan logika warna yang disempurnakan
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
    const title = currentConfig.title;
    const data = currentData;
    const showTrend = currentConfig.showTrend;
    const showThreshold = currentConfig.showThreshold;
    
    // Perbarui judul grafik
    document.getElementById('dynamicChartTitle').textContent = title;
    
    // Set dimensi kanvas
    canvas.width = Math.max(600, data.length * 60);
    canvas.height = 300;
    
    // Clear canvas
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    
    // Chart settings
    const barWidth = 40;
    const barSpacing = 60;
    const chartHeight = 250;
    const chartTop = 30;
    const maxValue = Math.max(...data, threshold) * 1.1;
    
    // Gambar batang dengan logika warna yang disempurnakan
    data.forEach((value, index) => {
        const barHeight = (value / maxValue) * chartHeight;
        const x = index * barSpacing + 30;
        const y = chartTop + chartHeight - barHeight;
        const prevValue = index > 0 ? data[index - 1] : null;
        
        // Logika warna 
        let color = '#28a745'; // default hijau (normal)
        
        if (value < threshold) {
            // Jika lebih kecil dari threshold = merah (bahaya)
            color = '#d43636';
        } else if (prevValue !== null && value < prevValue) {
            // Jika turun dari nilai sebelumnya tapi masih >= threshold = kuning (peringatan)
            color = '#FFD700';
        }
        
        // Gambar batang dengan warna yang ditentukan
        ctx.fillStyle = color;
        ctx.fillRect(x, y, barWidth, barHeight);
        
        // Tambahkan border tipis pada batang
        ctx.strokeStyle = 'rgba(0,0,0,0.1)';
        ctx.lineWidth = 1;
        ctx.strokeRect(x, y, barWidth, barHeight);
        
        // Value labels
        ctx.fillStyle = '#333';
        ctx.font = 'bold 12px Arial';
        ctx.textAlign = 'center';
        ctx.fillText(value, x + barWidth/2, y - 5);
        
        // Index labels
        ctx.font = '11px Arial';
        ctx.fillText(`P${index + 1}`, x + barWidth/2, chartTop + chartHeight + 20);
        
        // Tambahkan indikator tren untuk peningkatan visual
        if (showTrend && prevValue !== null) {
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
    
    // Gambar garis batas minimum
    if (showThreshold) {
        const thresholdY = chartTop + chartHeight - (threshold / maxValue) * chartHeight;
        ctx.strokeStyle = '#ff6b35';
        ctx.lineWidth = 3;
        ctx.setLineDash([8, 4]);
        ctx.beginPath();
        ctx.moveTo(20, thresholdY);
        ctx.lineTo(canvas.width - 20, thresholdY);
        ctx.stroke();
        ctx.setLineDash([]); // Reset dash
        
        // Label batas dengan latar belakang
        ctx.fillStyle = 'rgba(255, 107, 53, 0.9)';
        ctx.fillRect(25, thresholdY - 18, 100, 16);
        ctx.fillStyle = 'white';
        ctx.font = 'bold 11px Arial';
        ctx.textAlign = 'left';
        ctx.fillText(`Ambang Batas: ${threshold}`, 28, thresholdY - 7);
    }
}

// Handler untuk perubahan jenis data
async function handleJenisChange() {
    const jenis = document.getElementById('jenis').value;
    const field = document.getElementById('fieldSelect').value;
    const bulan = document.getElementById('bulan').value;
    const tahun = document.getElementById('tahun').value;
    
    if (!jenis) {
        alert('Silakan pilih jenis data terlebih dahulu');
        return;
    }
    
    // Update config
    currentConfig.field = field;
    currentConfig.jenisData = jenis;
    
    // Tampilkan loading
    const canvas = document.getElementById('barChart');
    const ctx = canvas.getContext('2d');
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    ctx.fillStyle = '#666';
    ctx.font = '16px Arial';
    ctx.textAlign = 'center';
    ctx.fillText('Memuat data...', canvas.width / 2, canvas.height / 2);
    
    // Fetch data dari API
    const apiData = await fetchDataByJenis(jenis, bulan, tahun);
    
    if (apiData.length > 0) {
        currentApiData = apiData;
        currentData = transformDataForChart(apiData, field);
        
        // Update title dengan nama field yang lebih friendly
        const fieldNames = {
            'jumlah_pabrik_basah_latek': 'Pabrik Basah Latek',
            'jumlah_kebun_basah_latek': 'Kebun Basah Latek',
            'jumlah_sheet': 'Jumlah Sheet',
            'k3_sheet': 'K3 Sheet',
            'jumlah_pabrik_basah_lump': 'Pabrik Basah Lump',
            'jumlah_kebun_basah_lump': 'Kebun Basah Lump',
            'jumlah_br_cr': 'Jumlah BR/CR',
            'k3_br_cr': 'K3 BR/CR'
        };
        
        currentConfig.title = `${fieldNames[field]} - ${jenis.replace(/_/g, ' ')}`;
        
        drawBarChart();
    } else {
        currentData = [];
        currentApiData = [];
        drawBarChart();
    }
}

// Handler untuk perubahan field
function handleFieldChange() {
    if (currentApiData.length === 0) {
        return;
    }
    
    const field = document.getElementById('fieldSelect').value;
    currentConfig.field = field;
    currentData = transformDataForChart(currentApiData, field);
    
    // Update title
    const fieldNames = {
        'jumlah_pabrik_basah_latek': 'Pabrik Basah Latek',
        'jumlah_kebun_basah_latek': 'Kebun Basah Latek',
        'jumlah_sheet': 'Jumlah Sheet',
        'k3_sheet': 'K3 Sheet',
        'jumlah_pabrik_basah_lump': 'Pabrik Basah Lump',
        'jumlah_kebun_basah_lump': 'Kebun Basah Lump',
        'jumlah_br_cr': 'Jumlah BR/CR',
        'k3_br_cr': 'K3 BR/CR'
    };
    
    currentConfig.title = `${fieldNames[field]} - ${currentConfig.jenisData.replace(/_/g, ' ')}`;
    
    drawBarChart();
}

// Event listeners
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

// Inisialisasi saat halaman dimuat
window.onload = function() {
    // Set bulan dan tahun saat ini sebagai default
    const now = new Date();
    const tahunInput = document.getElementById('tahun');
    if (tahunInput && !tahunInput.value) {
        tahunInput.value = now.getFullYear();
    }
    
    // Gambar chart kosong pertama kali
    drawBarChart();
};