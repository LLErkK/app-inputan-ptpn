// JS for dashboard with API integration + Debug Logging
// Mapping quadrant keys to afdeling names
const afdelingMap = {
    setro: 'Setro',
    jatiroenggo: 'Jatiroenggo',
    klepu: 'Klepu',
    gebugan: 'Gebugan'
};

// Cache untuk menyimpan data yang sudah diambil
const dataCache = {};

// Loading state
let isLoading = false;

// Format number dengan pemisah ribuan
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

// Fungsi untuk fetch data dari API
async function fetchDashboardData(afdeling) {
    const url = `/api/dashboard?afdeling=${encodeURIComponent(afdeling)}`;
    console.log('🔄 Fetching data from:', url);
    
    try {
        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'same-origin'
        });
        
        console.log('📡 Response status:', response.status, response.statusText);
        
        if (!response.ok) {
            const errorText = await response.text();
            console.error('❌ Error response:', errorText);
            throw new Error(`HTTP error! status: ${response.status} - ${errorText}`);
        }
        
        const data = await response.json();
        console.log('✅ Data received:', data);
        
        // Validate data structure
        if (!data || typeof data !== 'object') {
            throw new Error('Invalid data format received from API');
        }
        
        return data;
    } catch (error) {
        console.error('❌ Error fetching dashboard data:', error);
        throw error;
    }
}

// Transform API response ke format yang dibutuhkan UI
function transformApiData(apiData) {
    console.log('🔄 Transforming API data:', apiData);
    
    // Handle both camelCase and lowercase property names from API
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

// Set loading state pada info boxes
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

// Set error state
function setErrorState(errorMessage = 'Error') {
    console.error('❌ Setting error state:', errorMessage);
    const infoValues = document.querySelectorAll('.info-value');
    infoValues.forEach(el => {
        el.innerText = '—';
        el.classList.add('error');
    });
    
    // Show error notification
    showNotification(`Error: ${errorMessage}`, 'error');
}

// Show notification to user
function showNotification(message, type = 'info') {
    // Check if notification element exists, if not create it
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
    
    // Set color based on type
    const colors = {
        error: '#ef4444',
        success: '#10b981',
        info: '#3b82f6'
    };
    
    notification.style.borderLeft = `4px solid ${colors[type] || colors.info}`;
    notification.textContent = message;
    notification.style.opacity = '1';
    
    // Auto hide after 5 seconds
    setTimeout(() => {
        notification.style.opacity = '0';
    }, 5000);
}

async function setActiveQuadrant(key) {
    console.log('🎯 Active quadrant changed to:', key);
    
    if (isLoading) {
        console.warn('⚠️ Already loading, skipping request');
        return;
    }
    
    // Toggle active class & aria-pressed
    document.querySelectorAll('.quad-btn').forEach(btn => {
        const isActive = btn.dataset.key === key;
        btn.classList.toggle('active', isActive);
        btn.setAttribute('aria-pressed', isActive ? 'true' : 'false');
    });

    const afdeling = afdelingMap[key];
    if (!afdeling) {
        console.error('❌ Afdeling not found for key:', key);
        setErrorState('Afdeling tidak ditemukan');
        return;
    }
    
    console.log('📍 Afdeling:', afdeling);

    // Check cache first
    if (dataCache[key]) {
        console.log('💾 Using cached data for:', key);
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
        setErrorState(error.message);
        
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
    
    // Remove the class after animation ends
    setTimeout(() => el.classList.remove('updated'), 550);
}

// Initialize dashboard when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    console.log('🚀 Dashboard initialized');
    
    // Add click handlers to quadrant buttons
    document.querySelectorAll('.quad-btn').forEach(btn => {
        btn.addEventListener('click', () => setActiveQuadrant(btn.dataset.key));
        btn.addEventListener('keydown', (e) => {
            if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault();
                setActiveQuadrant(btn.dataset.key);
            }
        });
    });

    // Default selection - load Setro data on page load
    console.log('📍 Loading default quadrant: setro');
    setActiveQuadrant('setro');
    
    // Add refresh button functionality if exists
    const refreshBtn = document.getElementById('refresh-dashboard');
    if (refreshBtn) {
        refreshBtn.addEventListener('click', () => {
            console.log('🔄 Manual refresh triggered');
            // Clear cache
            Object.keys(dataCache).forEach(key => delete dataCache[key]);
            // Reload current active quadrant
            const activeBtn = document.querySelector('.quad-btn.active');
            if (activeBtn) {
                setActiveQuadrant(activeBtn.dataset.key);
            }
        });
    }
});