// peta.js
// Leaflet + omnivore KML dengan sistem debug dan handshake
// Mengirim event 'afdelingClicked' saat polygon diklik
// Menjawab 'rekapPing' dengan 'petaPong'

// -----------------------------
// Debug helper
// -----------------------------
function sendDebug(updates) {
    try {
        window.dispatchEvent(new CustomEvent('debugUpdate', { detail: updates }));
    } catch(e) {
        console.warn('Gagal mengirim debug update:', e);
    }
}

// -----------------------------
// Fungsi gaya berdasarkan nama fitur
// -----------------------------
function styleAfdeling(feature) {
    let nama = (feature && feature.properties && (feature.properties.name || feature.properties.NAME || feature.properties.Nama)) || "";
    console.log("🎨 Memeriksa nama untuk style:", nama);

    // Gaya berdasarkan kategori nama
    if (/SET|ST|SETRO/i.test(nama)) {
        return { color: "#ff0000", fillColor: "#ffaaaa", weight: 2, fillOpacity: 0.5 };
    }
    if (/KLP|KLEPU|YP-IE10-17-KLSG[0-9]+/i.test(nama)) {
        return { color: "#00b300", fillColor: "#aaffaa", weight: 2, fillOpacity: 0.5 };
    }
    if (/JR|JTR|JATIRUNG|JATIROENGGO/i.test(nama)) {
        return { color: "#0000ff", fillColor: "#aaaaff", weight: 2, fillOpacity: 0.5 };
    }

    // Afdeling Gebugan (beberapa format)
    if (
        /GEB|GB|GEBUG/i.test(nama) ||
        /FM-IE10-\d{2}-AR\d+/i.test(nama) ||
        /FI-IE10-\d{2}-AR\d+/i.test(nama) ||
        /YP-IE10-19-GESR[0-9]+/i.test(nama) ||
        /FM-IE10-\d{2}-RO\d+/i.test(nama)
    ) {
        return { color: "#ffa500", fillColor: "#ffe5b4", weight: 2, fillOpacity: 0.5 };
    }

    // default style
    return { color: "#555", fillColor: "#ddd", weight: 1, fillOpacity: 0.3 };
}

// -----------------------------
// Parse description KML -> object key:value
// -----------------------------
function parseDescription(desc) {
    if (!desc) return {};
    let teks = desc.replace(/<[^>]+>/g, " ").replace(/\s+/g, " ").trim();
    let parts = teks.split(/\s{2,}|\n/);
    let data = {};
    parts.forEach(p => {
        const pair = p.split(":");
        if (pair.length === 2) {
            data[pair[0].trim()] = pair[1].trim();
        }
    });
    return data;
}

// -----------------------------
// Deteksi label afdeling dari nama feature
// -----------------------------
function detectAfdelingFromName(namaText) {
    if (!namaText) return { key: null, label: "Tidak diketahui" };
    const n = namaText.toString().toLowerCase();

    // Mapping regex ke key dan label
    const mappings = [
        {
            regex: /set|st|setro/i,
            key: 'setro',
            label: 'Setro'
        },
        {
            regex: /klp|klepu|yp-ie10-17-klsg[0-9]+/i,
            key: 'klepu',
            label: 'Klepu'
        },
        {
            regex: /jr|jtr|jatirung|jatiroenggo/i,
            key: 'jatiroenggo',
            label: 'Jatiroenggo'
        },
        {
            regex: /geb|gb|gebug|fm-ie10-\d{2}-ar\d+|fi-ie10-\d{2}-ar\d+|yp-ie10-19-gesr[0-9]+|fm-ie10-\d{2}-ro\d+/i,
            key: 'gebugan',
            label: 'Gebugan'
        }
    ];

    // Cari yang cocok
    const match = mappings.find(m => m.regex.test(n));
    return match || { key: null, label: "Tidak diketahui" };
}

// -----------------------------
// onEachFeature: bind popup & pasang klik handler
// -----------------------------
function onEachFeature(feature, layer) {
    const nama = feature && feature.properties && (feature.properties.name || feature.properties.NAME || feature.properties.Nama) || "Tanpa Nama";
    console.log("📋 Memproses nama fitur (inisialisasi):", nama);

    const descData = parseDescription(feature && feature.properties && feature.properties.description);
    const luas = descData["Luas"] || "Tidak diketahui";
    const tahun = descData["Tahun Tanam"] || "Tidak diketahui";
    const lokasi = descData["Lokasi"] || "Tidak diketahui";

    // Tentukan afdeling label dan key untuk popup / event
    const afdelingInfo = detectAfdelingFromName(nama);
    const afdeling = afdelingInfo.label;
    const afdKey = afdelingInfo.key;
    console.log("🏷️ Afdeling yang terdeteksi:", afdelingInfo);

    const popupHtml = `<b>${nama}</b><br>
        <b>${afdelingInfo.label}</b><br>
        Lokasi: ${lokasi}<br>
        Luas: ${luas}<br>
        Tahun Tanam: ${tahun}`;

    try {
        layer.bindPopup(popupHtml);
    } catch (err) {
        console.warn('⚠️ Gagal bindPopup untuk feature:', err);
    }

    // simpan style awal supaya bisa di-reset
    const originalStyle = styleAfdeling(feature);

    function highlightLayerTemporary() {
        try {
            if (layer.setStyle) {
                layer.setStyle({ weight: 4, color: '#FFD700', fillOpacity: 0.7 });
            }
            if (layer.bringToFront) layer.bringToFront();
        } catch (err) {
            console.warn('⚠️ Gagal highlight layer:', err);
        }
    }

    function resetLayerStyle() {
        try {
            if (layer.setStyle && originalStyle) {
                layer.setStyle(originalStyle);
            }
        } catch (err) {
            console.warn('⚠️ Gagal reset layer style:', err);
        }
    }

    // KIRIM EVENT SAAT USER KLIK POLYGON
    layer.on('click', function (e) {
        console.log('🗺️ Polygon diklik:', nama, '→', afdeling, '(key:', afdKey, ')');
        sendDebug({ lastClick: `${nama} (${afdeling})` });

        // buka popup di posisi klik
        try {
            if (layer.openPopup) {
                layer.openPopup(e.latlng);
            }
        } catch (err) {
            console.warn('⚠️ Gagal membuka popup pada klik:', err);
        }

    // simpan terakhir ke global (fallback)
    window.lastDetectedAfdeling = afdeling;
    window.lastAfdelingKey = afdKey;

        // kirim CustomEvent agar rekap.js bisa tangkap
        try {
            // gunakan key yang sudah dideteksi sebelumnya
            const payload = { 
                afdeling: afdeling,
                afdelingKey: afdKey, // key untuk rekap.js (setro, klepu, dll)
                name: nama,
                featureId: feature && (feature.id || (feature.properties && (feature.properties.id || feature.properties.ID))) || null,
                descData: descData // data mentah dari KML untuk fallback
            };
            const evt = new CustomEvent('afdelingClicked', { detail: payload });
            window.dispatchEvent(evt);
            console.log('✅ Dispatched event afdelingClicked:', payload, 'dengan afdelingKey:', afdKey);
        } catch (err) {
            console.error('❌ Gagal dispatch afdelingClicked:', err);
            sendDebug({ mapStatus: 'Event Error ❌' });
        }

        // highlight sementara
        highlightLayerTemporary();
        setTimeout(() => {
            resetLayerStyle();
        }, 2000);
    });

    // Optional UX: mouseover / mouseout
    layer.on('mouseover', function (e) {
        try { 
            if (layer.openPopup) layer.openPopup(); 
            highlightLayerTemporary(); 
        } catch (err) {}
    });
    
    layer.on('mouseout', function () {
        try { 
            resetLayerStyle(); 
        } catch (err) {}
    });
}

// -----------------------------
// Inisialisasi peta Leaflet
// -----------------------------
console.log('🗺️ Inisialisasi peta Leaflet...');
sendDebug({ mapStatus: 'Initializing...' });

var map = L.map('map').setView([-7.5, 110.3], 12);

// Tile layer OSM
L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
}).addTo(map);

console.log('✅ Base map loaded');

// -----------------------------
// Memuat KML via omnivore
// -----------------------------
var kmlUrl = '/kml/KEBUN_NGOBO.kml';
console.log('📡 Meminta KML dari:', kmlUrl);
sendDebug({ mapStatus: 'Loading KML...' });

// Debug: Test fetch KML terlebih dahulu
fetch(kmlUrl)
    .then(response => {
        console.log('📊 KML Fetch status:', response.status, response.statusText);
        console.log('📊 KML Content-Type:', response.headers.get('Content-Type'));
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
        return response.text();
    })
    .then(text => {
        console.log('✅ KML content preview (200 chars):', text.substring(0, 200));
        sendDebug({ mapStatus: 'KML Fetched ✅' });
    })
    .catch(err => {
        console.error('❌ KML Fetch error:', err);
        sendDebug({ mapStatus: `KML Error: ${err.message} ❌` });
        alert(`Gagal fetch KML: ${err.message}\n\nCek:\n1. File ada di ./templates/kml/KEBUN_NGOBO.kml\n2. Server Go running\n3. Path benar\n4. Network tab untuk detail`);
    });

// Create layer via omnivore
var kmlLayer = omnivore.kml(kmlUrl, null, L.geoJSON(null, {
    style: styleAfdeling,
    onEachFeature: onEachFeature
}));

kmlLayer.on('ready', function() {
    try {
        var bounds = kmlLayer.getBounds();
        if (bounds && bounds.isValid && bounds.isValid()) {
            map.fitBounds(bounds);
            console.log('✅ Map fitBounds ke KML');
            sendDebug({ mapStatus: 'KML Loaded & Ready ✅' });
        } else {
            map.fitWorld();
            console.warn('⚠️ Bounds KML tidak valid, menggunakan fitWorld');
            sendDebug({ mapStatus: 'KML Invalid Bounds ⚠️' });
        }
    } catch (err) {
        console.warn('⚠️ Gagal menyesuaikan bounds KML:', err);
        sendDebug({ mapStatus: 'KML Bounds Error ⚠️' });
    }
});

kmlLayer.on('error', function(e) {
    console.error('❌ Gagal memuat KML:', e);
    sendDebug({ mapStatus: 'KML Load Failed ❌' });
    alert('Gagal memuat file KML. Periksa console Network untuk detail (404/CORS).');
});

// tambahkan layer ke peta
kmlLayer.addTo(map);

// -----------------------------
// Expose globals untuk fallback/debug
// -----------------------------
window._leafletMap = map;
window._kmlLayer = kmlLayer;
window.lastDetectedAfdeling = window.lastDetectedAfdeling || null;
window.getLastDetectedAfdeling = function() { return window.lastDetectedAfdeling || null; };

// -----------------------------
// RESPOND to rekapPing -> petaPong (handshake)
// -----------------------------
window.addEventListener('rekapPing', function (e) {
    try {
        const info = {
            ready: true,
            hasMap: !!window._leafletMap,
            lastDetectedAfdeling: window.lastDetectedAfdeling || null,
            ts: Date.now(),
            receivedDetail: e && e.detail ? e.detail : null
        };
        window.dispatchEvent(new CustomEvent('petaPong', { detail: info }));
        console.log('🤝 petaPong dikirim sebagai jawaban rekapPing:', info);
        sendDebug({ handshake: 'Connected ✅' });
    } catch (err) {
        console.error('❌ Gagal merespon rekapPing:', err);
        sendDebug({ handshake: 'Failed ❌' });
    }
});

console.log('✅ peta.js loaded successfully');