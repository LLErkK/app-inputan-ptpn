document.addEventListener("DOMContentLoaded", () => {
    const bakuTableBody = document.getElementById("bakuTableBody");
    const btnSearch = document.getElementById("searchBtn");
    const tglAwal = document.getElementById("searchTanggalAwal");
    const tglAkhir = document.getElementById("searchTanggalAkhir");
    const filterJenis = document.getElementById("filterJenis");
    const namaMandorInput = document.getElementById("namaMandor");
    const namaPenyadapInput = document.getElementById("namaPenyadap");

    // Debug: Cek apakah semua element ditemukan
    console.log("Elements found:", {
        bakuTableBody: !!bakuTableBody,
        btnSearch: !!btnSearch,
        tglAwal: !!tglAwal,
        tglAkhir: !!tglAkhir,
        filterJenis: !!filterJenis,
        namaMandorInput: !!namaMandorInput,
        namaPenyadapInput: !!namaPenyadapInput
    });

    // Fungsi untuk mengambil data dari backend
    async function fetchData(url) {
        console.log("Fetching data from:", url);
        try {
            const res = await fetch(url);
            console.log("Response status:", res.status);
            console.log("Response headers:", Object.fromEntries(res.headers.entries()));
            
            const json = await res.json();
            console.log("Response JSON:", json);
            
            if (json && json.success && Array.isArray(json.data)) {
                console.log("Data found:", json.data.length, "records");
                return json.data;
            } else {
                console.log("No valid data in response");
                return [];
            }
        } catch (e) {
            console.error("Fetch error:", e);
            return [];
        }
    }

    // Fungsi untuk render data ke dalam tabel
    function renderTable(dataArr) {
        console.log("Rendering table with data:", dataArr);
        
        if (!bakuTableBody) {
            console.error("bakuTableBody element not found!");
            return;
        }

        bakuTableBody.innerHTML = ""; // Clear previous data

        if (!dataArr || !dataArr.length) {
            bakuTableBody.innerHTML = `<tr><td colspan="11" style="text-align:center;">Data tidak ditemukan.</td></tr>`;
            return;
        }

        // Render data into the table
        dataArr.forEach((item, index) => {
            console.log(`Rendering row ${index}:`, item);
            const tr = document.createElement("tr");
            tr.innerHTML = `
                <td>${item.tanggal || "-"}</td>
                <td>${item.mandor || "-"}</td>
                <td>${item.tipe || "-"}</td>
                <td>${item.tahunTanam || "-"}</td>
                <td>${item.afdeling || "-"}</td>
                <td>${item.nik || "-"}</td>
                <td>${item.namaPenyadap || "-"}</td>
                <td>${item.basahLatex || 0}</td>
                <td>${item.sheet || 0}</td>
                <td>${item.basahLump || 0}</td>
                <td>${item.brCr || 0}</td>
            `;
            bakuTableBody.appendChild(tr);
        });
        
        console.log("Table rendered successfully");
    }

    // Fungsi untuk memicu pencarian
    async function searchData() {
        console.log("Starting search...");
        
        // Ambil input dari pengguna
        const namaMandor = namaMandorInput?.value.trim() || "";
        const namaPenyadap = namaPenyadapInput?.value.trim() || "";
        const awal = tglAwal?.value || "";
        const akhir = tglAkhir?.value || "";
        const jenis = filterJenis?.value || "";

        console.log("=== SEARCH PARAMETERS ===");
        console.log("Nama Mandor:", namaMandor);
        console.log("Nama Penyadap:", namaPenyadap);
        console.log("Tanggal Awal:", awal);
        console.log("Tanggal Akhir:", akhir);
        console.log("Jenis:", jenis);

        let url = "http://localhost:8080/api/monitoring/smart-search";
        let params = new URLSearchParams();
        
        // Tambahkan parameter sesuai dengan backend
        if (namaMandor) {
            params.append("namaMandor", namaMandor);
            console.log("✓ Added namaMandor parameter");
        }
        if (namaPenyadap) {
            params.append("namaPenyadap", namaPenyadap);
            console.log("✓ Added namaPenyadap parameter");
        }
        if (jenis) {
            params.append("tipe", jenis); // Backend menggunakan "tipe" bukan "filterJenis"
            console.log("✓ Added tipe parameter");
        }

        // PERBAIKAN: Filter berdasarkan rentang tanggal
        // Hanya tambahkan parameter tanggal jika ada input tanggal
        // PENTING: Backend menggunakan "tanggalAwal" dan "tanggalAkhir" (BUKAN "filterTanggalAwal")
        if (awal || akhir) {
            const tanggalAwal = awal || akhir; // Gunakan akhir jika awal kosong
            const tanggalAkhir = akhir || awal; // Gunakan awal jika akhir kosong
            
            params.append("tanggalAwal", tanggalAwal);
            params.append("tanggalAkhir", tanggalAkhir);
            
            console.log("✓ Added date range:", tanggalAwal, "to", tanggalAkhir);
        } else {
            console.log("⚠ No date filter applied - will return current month data");
        }

        // Gabungkan URL dengan query parameters
        if (params.toString()) {
            url += `?${params.toString()}`;
        }

        console.log("=== FINAL REQUEST ===");
        console.log("URL:", url);
        console.log("Parameters:", params.toString());
        console.log("=====================");

        // Tampilkan loading indicator
        if (bakuTableBody) {
            bakuTableBody.innerHTML = `<tr><td colspan="11" style="text-align:center;"><i>Memuat data...</i></td></tr>`;
        }

        // Ambil data dari API dan render
        const dataArr = await fetchData(url);
        renderTable(dataArr);
    }

    // Event search when clicking the search button
    if (btnSearch) {
        btnSearch.addEventListener("click", async (event) => {
            console.log("Search button clicked");
            event.preventDefault();
            await searchData();
        });
    }

    // Tambahkan event listener untuk form submit juga
    const monitoringForm = document.getElementById("monitoringForm");
    if (monitoringForm) {
        monitoringForm.addEventListener("submit", async (event) => {
            console.log("Form submitted");
            event.preventDefault();
            await searchData();
        });
    }

    // Event search when pressing Enter
    const inputElements = document.querySelectorAll("#namaMandor, #namaPenyadap, #searchTanggalAwal, #searchTanggalAkhir, #filterJenis");
    inputElements.forEach(input => {
        input.addEventListener("keydown", function(event) {
            if (event.key === 'Enter') {
                console.log("Enter pressed on input");
                event.preventDefault();
                searchData();
            }
        });
    });

    // Fungsi untuk memuat data berdasarkan hari ini
    function loadTodayData() {
        console.log("Loading today data...");
        const today = new Date().toISOString().split("T")[0];
        
        // Set nilai input tanggal ke hari ini
        if (tglAwal) tglAwal.value = today;
        if (tglAkhir) tglAkhir.value = today;
        
        // Trigger pencarian
        searchData();
    }

    // Fungsi untuk memuat data berdasarkan minggu ini
    function loadWeekData() {
        console.log("Loading week data...");
        const now = new Date();
        const dayOfWeek = now.getDay(); // 0 = Sunday, 1 = Monday, etc.
        
        // Hitung hari pertama minggu ini (Minggu)
        const firstDay = new Date(now);
        firstDay.setDate(now.getDate() - dayOfWeek);
        
        // Hitung hari terakhir minggu ini (Sabtu)
        const lastDay = new Date(now);
        lastDay.setDate(now.getDate() + (6 - dayOfWeek));

        const startDate = firstDay.toISOString().split("T")[0];
        const endDate = lastDay.toISOString().split("T")[0];

        console.log("Week range:", startDate, "to", endDate);

        // Set nilai input tanggal
        if (tglAwal) tglAwal.value = startDate;
        if (tglAkhir) tglAkhir.value = endDate;
        
        // Trigger pencarian
        searchData();
    }

    // Fungsi untuk memuat data berdasarkan bulan ini
    function loadMonthData() {
        console.log("Loading month data...");
        const now = new Date();
        const firstDayOfMonth = new Date(now.getFullYear(), now.getMonth(), 1);
        const lastDayOfMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0);

        const startDate = firstDayOfMonth.toISOString().split("T")[0];
        const endDate = lastDayOfMonth.toISOString().split("T")[0];

        console.log("Month range:", startDate, "to", endDate);

        // Set nilai input tanggal
        if (tglAwal) tglAwal.value = startDate;
        if (tglAkhir) tglAkhir.value = endDate;
        
        // Trigger pencarian
        searchData();
    }

    // Fungsi untuk reset form
    function clearAll() {
        console.log("Clearing all data...");
        if (tglAwal) tglAwal.value = "";
        if (tglAkhir) tglAkhir.value = "";
        if (filterJenis) filterJenis.value = "";
        if (namaMandorInput) namaMandorInput.value = "";
        if (namaPenyadapInput) namaPenyadapInput.value = "";
        
        // Reset ke pesan awal
        initializeTable();
    }

    // Expose fungsi ke global scope
    window.loadTodayData = loadTodayData;
    window.loadWeekData = loadWeekData;
    window.loadMonthData = loadMonthData;
    window.clearAll = clearAll;
    window.searchData = searchData;
    
    console.log("JavaScript loaded successfully");
    
    // Inisialisasi tabel dengan pesan kosong
    function initializeTable() {
        if (bakuTableBody) {
            bakuTableBody.innerHTML = `<tr><td colspan="11" style="text-align:center; color: #666; font-style: italic;">Silakan pilih filter dan klik tombol untuk menampilkan data</td></tr>`;
        }
    }
    
    // Inisialisasi tabel kosong saat halaman dimuat
    initializeTable();
});