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
        const namaMandor = namaMandorInput?.value || "";
        const namaPenyadap = namaPenyadapInput?.value || "";
        const awal = tglAwal?.value || "";
        const akhir = tglAkhir?.value || "";
        const jenis = filterJenis?.value || "";

        console.log("Search parameters:", {
            namaMandor, namaPenyadap, awal, akhir, jenis
        });

        let url = "http://localhost:8080/api/monitoring/smart-search";
        let params = new URLSearchParams();
        
        // Tambahkan parameter sesuai dengan backend
        if (namaMandor) params.append("namaMandor", namaMandor);
        if (namaPenyadap) params.append("namaPenyadap", namaPenyadap);
        if (jenis) params.append("filterJenis", jenis);

        // Filter berdasarkan rentang tanggal
        if (awal && akhir) {
            params.append("filterTanggalAwal", awal);
            params.append("filterTanggalAkhir", akhir);
        } else if (awal && !akhir) {
            params.append("filterTanggalAwal", awal);
            params.append("filterTanggalAkhir", awal);
        } else {
            // Jika tidak ada tanggal, gunakan tanggal hari ini
            const today = new Date().toISOString().split("T")[0];
            params.append("filterTanggalAwal", today);
            params.append("filterTanggalAkhir", today);
        }

        // Gabungkan URL dengan query parameters
        if (params.toString()) {
            url += `?${params.toString()}`;
        }

        console.log("Final URL:", url);

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
        
        const url = `http://localhost:8080/api/monitoring/today/summary?tipe=${filterJenis.value}&namaMandor=${namaMandorInput.value}&penyadap=${namaPenyadapInput.value}`;
        console.log("Today data URL:", url);
        fetchData(url).then(data => renderTable(data));
    }

    // Fungsi untuk memuat data berdasarkan minggu ini
    function loadWeekData() {
        console.log("Loading week data...");
        const currentDate = new Date();
        const firstDayOfWeek = new Date(currentDate.setDate(currentDate.getDate() - currentDate.getDay()));
        const lastDayOfWeek = new Date(currentDate.setDate(currentDate.getDate() - currentDate.getDay() + 6));

        const startDate = firstDayOfWeek.toISOString().split("T")[0];
        const endDate = lastDayOfWeek.toISOString().split("T")[0];

        const url = `http://localhost:8080/api/monitoring/week/summary?tipe=${filterJenis.value}&namaMandor=${namaMandorInput.value}&penyadap=${namaPenyadapInput.value}`;
        console.log("Week data URL:", url);
        fetchData(url).then(data => renderTable(data));
    }

    // Fungsi untuk memuat data berdasarkan bulan ini
    function loadMonthData() {
        console.log("Loading month data...");
        const currentDate = new Date();
        const firstDayOfMonth = new Date(currentDate.getFullYear(), currentDate.getMonth(), 1);
        const lastDayOfMonth = new Date(currentDate.getFullYear(), currentDate.getMonth() + 1, 0);

        const startDate = firstDayOfMonth.toISOString().split("T")[0];
        const endDate = lastDayOfMonth.toISOString().split("T")[0];

        const url = `http://localhost:8080/api/monitoring/month/summary?tipe=${filterJenis.value}&namaMandor=${namaMandorInput.value}&penyadap=${namaPenyadapInput.value}`;
        console.log("Month data URL:", url);
        fetchData(url).then(data => renderTable(data));
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
    window.searchData = searchData; // Tambahkan ini juga
    
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