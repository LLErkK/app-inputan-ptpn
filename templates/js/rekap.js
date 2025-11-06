// rekap.js
// Dashboard JS with API integration, debug logging, and handshake

// -----------------------------
// Mapping and global cache/state
// -----------------------------
const afdelingMap = {
    setro: 'Setro',
    jatiroenggo: 'Jatiroenggo',
    jatirunggo: 'Jatiroenggo',  // alias
    klepu: 'Klepu',
    gebugan: 'Gebugan'
};

const dataCache = {};
let isLoading = false;

// -----------------------------
// Debug helper
// -----------------------------
function updateDebugUI(updates) {
    try {
        window.dispatchEvent(new CustomEvent('debugUpdate', { detail: updates }));
        console.log('🐛 Debug update:', updates);
    } catch(e) {
        console.warn('⚠️ Gagal update debug UI:', e);
    }
}

// -----------------------------
// Utility: format number
// -----------------------------
function formatNumber(num) {
    if (num === null || num === undefined || num === '') {
        console.log('formatNumber: value is null/undefined/empty', num);
        return '0';
    }
    const parsed = parseFloat(num);
    if (isNaN(parsed)) {
        console.log('formatNumber: value is NaN', num);
        return '0';
    }
    const formatted = parsed.toLocaleString('id-ID', { 
        maximumFractionDigits: 2, 
        minimumFractionDigits: 0 
    });
    console.log(`formatNumber: ${num} -> ${formatted}`);
    return formatted;
}

// -----------------------------
// Fetch dari API
// -----------------------------
async function fetchDashboardData(afdeling) {
    const url = `/api/dashboard?afdeling=${encodeURIComponent(afdeling)}`;
    console.log('🔄 Fetching data from:', url);
    updateDebugUI({ apiCall: `GET ${url}` });
    
    try {
        const response = await fetch(url, { 
            method: 'GET', 
            headers: { 'Content-Type': 'application/json' }, 
            credentials: 'same-origin' 
        });
        
        console.log('📡 Response status:', response.status, response.statusText);
        
        if (!response.ok) {
            const errorText = await response.text();
            console.error('❌ Error response:', errorText);
            updateDebugUI({ apiCall: `ERROR ${response.status}` });
            throw new Error(`HTTP error! status: ${response.status} - ${errorText}`);
        }
        
        const data = await response.json();
        console.log('✅ Data received:', data);
        updateDebugUI({ apiCall: 'Success ✅' });
        
        if (!data || typeof data !== 'object') {
            throw new Error('Invalid data format received from API');
        }
        
        return data;
    } catch (error) {
        console.error('❌ Error fetching dashboard data:', error);
        updateDebugUI({ apiCall: `Error: ${error.message}` });
        throw error;
    }
}

// -----------------------------
// Transform API data
// -----------------------------
function transformApiData(apiData) {
    console.log('🔄 Transforming API data:', apiData);
    
    const getData = (key) => {
        return apiData[key] !== undefined ? apiData[key] : 
               apiData[key.toLowerCase()] !== undefined ? apiData[key.toLowerCase()] : 
               null;
    };

    const transformed = {
        basahLatekKebun: formatNumber(getData('totalHariIniBasahLatekKebun') || getData('totalhariinibasahlatekKebun') || 0),
        basahLatekPabrik: formatNumber(getData('totalHariIniBasahLatekPabrik') || getData('totalhariinibasahlatekpabrik') || 0),
        basahLumpKebun: formatNumber(getData('totalHariIniBasahLumpKebun') || getData('totalhariiinibasahlumpkebun') || 0),
        basahLumpPabrik: formatNumber(getData('totalHariIniBasahLumpPabrik') || getData('totalhariinibasahlumppabrik') || 0),
        k3Sheet: formatNumber(getData('totalHariIniK3Sheet') || getData('totalhariinik3sheet') || 0),
        jumlahKering: formatNumber(getData('totalHariIniKeringJumlah') || getData('totalhariinikeringjumlah') || 0)
    };

    console.log('✅ Transformed data:', transformed);
    return transformed;
}

// -----------------------------
// UI helpers
// -----------------------------
function setLoadingState(isLoadingState) {
    console.log('⏳ Setting loading state:', isLoadingState);
    isLoading = isLoadingState;
    const infoValues = document.querySelectorAll('.info-value');
    infoValues.forEach(el => {
        if (isLoadingState) {
            el.classList.add('loading');
            el.innerText = '⏳';
        } else {
            el.classList.remove('loading');
        }
    });
}

function setErrorState(errorMessage = 'Error') {
    console.error('❌ Setting error state:', errorMessage);
    const infoValues = document.querySelectorAll('.info-value');
    infoValues.forEach(el => {
        el.innerText = '—';
        el.classList.add('error');
    });
    showNotification(`Error: ${errorMessage}`, 'error');
    updateDebugUI({ apiCall: `ERROR: ${errorMessage}` });
}

function showNotification(message, type = 'info') {
    let notification = document.getElementById('dashboard-notification');
    if (!notification) {
        notification = document.createElement('div');
        notification.id = 'dashboard-notification';
        notification.style.cssText = `
            position: fixed;
            top: 20px;
            right: 20px;
            padding: 12px 20px;
            border-radius: 8px;
            background: white;
            box-shadow: 0 4px 12px rgba(0,0,0,0.15);
            z-index: 10000;
            font-size: 14px;
            font-weight: 600;
            max-width: 300px;
            opacity: 0;
            transition: opacity 0.3s ease;
        `;
        document.body.appendChild(notification);
    }
    
    const colors = { 
        error: '#ef4444', 
        success: '#10b981', 
        info: '#3b82f6' 
    };
    
    notification.style.borderLeft = `4px solid ${colors[type] || colors.info}`;
    notification.textContent = message;
    notification.style.opacity = '1';
    
    setTimeout(() => { 
        notification.style.opacity = '0'; 
    }, 5000);
}

// -----------------------------
// setActiveQuadrant - MAIN FUNCTION
// -----------------------------
async function setActiveQuadrant(key, forceFetch = false) {
    console.log('🎯 Active quadrant changed to:', key);
    updateDebugUI({ activeAfdeling: key });
    
    if (isLoading) { 
        console.warn('⚠️ Already loading, skipping request'); 
        return; 
    }

    // Toggle active class pada tombol (jika ada)
    document.querySelectorAll('.quad-btn').forEach(btn => {
        const isActive = btn.dataset.key === key;
        btn.classList.toggle('active', isActive);
        btn.setAttribute('aria-pressed', isActive ? 'true' : 'false');
    });

    const afdeling = afdelingMap[key];
    if (!afdeling) { 
        console.error('❌ Afdeling not found for key:', key); 
        updateDebugUI({ apiCall: 'ERROR: Afdeling not found' });
        setErrorState('Afdeling tidak ditemukan'); 
        return; 
    }
    
    console.log('📍 Afdeling:', afdeling);
    updateDebugUI({ apiCall: `Preparing fetch for ${afdeling}...` });

    // Check cache first (unless forced)
    if (!forceFetch && dataCache[key]) {
        console.log('💾 Using cached data for:', key);
        updateDebugUI({ apiCall: 'Using cache 💾' });
        updateUIWithData(dataCache[key]);
        return;
    }

    // Fetch data from API
    setLoadingState(true);
    
    try {
        const apiData = await fetchDashboardData(afdeling);
        
        // Check if data is empty or all zeros
        const hasData = Object.values(apiData).some(val => 
            val !== null && val !== undefined && val !== 0 && val !== '0'
        );
        console.log('📊 Has non-zero data:', hasData);
        
        const transformedData = transformApiData(apiData);
        
        // Store in cache
        dataCache[key] = transformedData;
        
        // Update UI
        updateUIWithData(transformedData);
        
        // Show success notification
        showNotification(`Data ${afdeling} berhasil dimuat`, 'success');
    } catch (error) {
        console.error('💥 Failed to load data:', error);
        setErrorState(error.message || 'Gagal memuat data');
        
        // Clear cache for this key on error
        delete dataCache[key];
    } finally {
        setLoadingState(false);
    }
}

function updateUIWithData(data) {
    console.log('🖼️ Updating UI with data:', data);
    updateValue('basah-latek-kebun', data.basahLatekKebun);
    updateValue('basah-latek-pabrik', data.basahLatekPabrik);
    updateValue('basah-lump-kebun', data.basahLumpKebun);
    updateValue('basah-lump-pabrik', data.basahLumpPabrik);
    updateValue('k3-sheet', data.k3Sheet);
    updateValue('jumlah-kering', data.jumlahKering);
}

// Set semua info box ke 0 (default) — gunakan saat load awal
function setAllInfoToZero() {
    updateValue('basah-latek-kebun', '0');
    updateValue('basah-latek-pabrik', '0');
    updateValue('basah-lump-kebun', '0');
    updateValue('basah-lump-pabrik', '0');
    updateValue('k3-sheet', '0');
    updateValue('jumlah-kering', '0');
}

function updateValue(id, text) {
    const el = document.getElementById(id);
    if (!el) { 
        console.error('❌ Element not found:', id); 
        return; 
    }
    if (el.innerText === text) { 
        console.log('⏭️ No change for:', id); 
        return; 
    }
    
    console.log(`✏️ Updating ${id}: "${el.innerText}" -> "${text}"`);
    
    el.classList.remove('error');
    el.innerText = text;
    el.classList.add('updated');
    
    setTimeout(() => el.classList.remove('updated'), 550);
}

// -----------------------------
// Handshake: ping peta.js and wait for response
// -----------------------------
function pingPeta(timeoutMs = 2000) {
    return new Promise((resolve) => {
        let handled = false;

        function onPong(e) {
            handled = true;
            window.removeEventListener('petaPong', onPong);
            console.log('📬 petaPong diterima:', e.detail);
            updateDebugUI({ handshake: 'Connected ✅' });
            resolve({ ok: true, detail: e.detail });
        }

        window.addEventListener('petaPong', onPong);

        // kirim ping
        try {
            window.dispatchEvent(new CustomEvent('rekapPing', { detail: { ts: Date.now() } }));
            console.log('📣 rekapPing dikirim ke peta.js');
            updateDebugUI({ handshake: 'Pinging...' });
        } catch (err) {
            console.warn('Gagal dispatch rekapPing:', err);
        }

        // timeout fallback
        setTimeout(() => {
            if (handled) return;
            window.removeEventListener('petaPong', onPong);
            
            // fallback: cek global objects peta mungkin sudah expose
            const fallbackInfo = {
                hasMap: !!window._leafletMap,
                lastDetectedAfdeling: window.lastDetectedAfdeling || null
            };
            
            if (fallbackInfo.hasMap || fallbackInfo.lastDetectedAfdeling) {
                console.log('🛠️ Fallback: peta terdeteksi via global variables:', fallbackInfo);
                updateDebugUI({ handshake: 'Fallback ⚠️' });
                resolve({ ok: true, detail: { fallback: true, ...fallbackInfo } });
            } else {
                console.warn('❌ Tidak ada respons dari peta.js (timeout)');
                updateDebugUI({ handshake: 'Failed ❌' });
                resolve({ ok: false });
            }
        }, timeoutMs);
    });
}

// -----------------------------
// Event listener: menerima dari peta.js ketika polygon diklik
// -----------------------------
window.addEventListener('afdelingClicked', (event) => {
    // 1. Log event untuk debugging
    console.log('📍 Event dari peta.js:', event?.detail);
    
    // 2. Validasi payload
    if (!event?.detail) {
        console.warn('❌ Event tanpa detail');
        showNotification('Data tidak lengkap', 'error');
        return;
    }
    
    // 3. Cek dan gunakan key yang dikirim dari peta
    const receivedKey = event.detail.afdelingKey;
    if (receivedKey && afdelingMap[receivedKey]) {
        console.log('✅ Menggunakan key dari peta:', receivedKey);
        setActiveQuadrant(receivedKey, true); // force fresh fetch when user clicks on map
        return;
    }
    
    // 4. Fallback ke nama afdeling jika tidak ada key
    const receivedName = event.detail.afdeling;
    if (!receivedName) {
        console.warn('❌ Tidak ada key atau nama afdeling');
        showNotification('Data afdeling tidak valid', 'error');
        return;
    }
    
    // 5. Coba temukan key dari nama
    const normalizedName = receivedName.toLowerCase();
    let mappedKey = null;
    
    // 6. Coba exact match dulu
    for (const [key, value] of Object.entries(afdelingMap)) {
        if (value.toLowerCase() === normalizedName) {
            mappedKey = key;
            break;
        }
    }
    
    // 7. Jika tidak ada exact match, coba dengan includes
    if (!mappedKey) {
        if (normalizedName.includes('setro')) mappedKey = 'setro';
        else if (normalizedName.includes('klepu')) mappedKey = 'klepu';
        else if (normalizedName.includes('gebugan')) mappedKey = 'gebugan';
        else if (normalizedName.includes('jatirung')) mappedKey = 'jatiroenggo';
    }
    
    // 8. Gunakan key yang ditemukan atau tampilkan error
    if (mappedKey) {
        console.log('✅ Key ditemukan dari nama:', mappedKey);
        setActiveQuadrant(mappedKey, true); // force fresh fetch when user clicks on map
    } else {
        console.warn('❌ Tidak dapat menentukan key untuk:', receivedName);
        showNotification(`Afdeling "${receivedName}" tidak dikenali`, 'error');
    }
});

// -----------------------------
// DOM ready: inisialisasi tombol & handshake
// -----------------------------
document.addEventListener('DOMContentLoaded', async () => {
    console.log('🚀 Dashboard initialized');
    updateDebugUI({ 
        activeAfdeling: '-', 
        lastClick: '-', 
        apiCall: 'Ready', 
        handshake: 'Checking...' 
    });

    // Add click / keyboard handlers to quadrant buttons (jika ada)
    document.querySelectorAll('.quad-btn').forEach(btn => {
        btn.addEventListener('click', () => setActiveQuadrant(btn.dataset.key));
        btn.addEventListener('keydown', (e) => {
            if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault();
                setActiveQuadrant(btn.dataset.key);
            }
        });
    });

    // lakukan ping ke peta.js untuk memastikan integrasi
    const pingResult = await pingPeta(2000);
    
    if (pingResult.ok) {
        console.log('🤝 Handshake dengan peta.js sukses:', pingResult.detail);
        showNotification('Terhubung ke peta', 'success');

        // Jika peta memberikan lastDetectedAfdeling, gunakan untuk load awal
        const last = pingResult.detail.lastDetectedAfdeling || window.lastDetectedAfdeling;
        if (last) {
            const afdKey = Object.keys(afdelingMap).find(
                k => afdelingMap[k].toLowerCase() === (last.toLowerCase && last.toLowerCase())
            );
            if (afdKey) {
                console.log('📍 Menggunakan lastDetectedAfdeling dari peta:', afdKey);
                setActiveQuadrant(afdKey);
                return;
            }
        }
    } else {
        console.warn('⚠️ Tidak terhubung ke peta (handshake gagal). Akan menggunakan default.');
        showNotification('Tidak terhubung ke peta - menggunakan default', 'info');
    }
    
    // Default: jangan auto-fetch; tampilkan 0 sampai user klik afdeling
    console.log('📍 Default: menampilkan 0 pada info-box (menunggu klik user)');
    setAllInfoToZero();

    // Refresh button jika ada
    const refreshBtn = document.getElementById('refresh-dashboard');
    if (refreshBtn) {
        refreshBtn.addEventListener('click', () => {
            console.log('🔄 Manual refresh triggered');
            Object.keys(dataCache).forEach(key => delete dataCache[key]);
            const activeBtn = document.querySelector('.quad-btn.active');
            if (activeBtn) {
                setActiveQuadrant(activeBtn.dataset.key);
            } else {
                setActiveQuadrant('setro');
            }
        });
    }
});

console.log('✅ rekap.js loaded successfully');