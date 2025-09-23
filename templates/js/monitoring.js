document.addEventListener("DOMContentLoaded", () => {
    const bakuTableBody = document.getElementById("bakuTableBody");
    const btnSearch = document.getElementById("searchBtn");
    const tglAwal = document.getElementById("searchTanggalAwal");
    const tglAkhir = document.getElementById("searchTanggalAkhir");
    const filterJenis = document.getElementById("filterJenis");
    const namaMandorInput = document.getElementById("namaMandor");
    const namaPenyadapInput = document.getElementById("namaPenyadap");

    // Fungsi untuk mengambil data dari backend
    async function fetchData(url) {
        try {
            const res = await fetch(url);
            const json = await res.json();
            if (json && json.success && Array.isArray(json.data)) {
                return json.data;
            }
        } catch (e) {
            console.error("Fetch error:", e);
        }
        return [];
    }

    // Fungsi untuk render data ke dalam tabel
    function renderTable(dataArr) {
        bakuTableBody.innerHTML = ""; // Clear previous data

        if (!dataArr.length) {
            bakuTableBody.innerHTML = `<tr><td colspan="10" style="text-align:center;">Data tidak ditemukan.</td></tr>`;
            return;
        }

        // Render data into the table
        dataArr.forEach(item => {
            const tr = document.createElement("tr");
            tr.innerHTML = `
                <td>${item.tanggal || "-"}</td>
                <td>${item.mandor || "-"}</td>
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
    }

    // Event search
    btnSearch.addEventListener("click", async (event) => {
        event.preventDefault(); // Prevent form submission
        let dataArr = [];

        // Ambil input dari pengguna
        const namaMandor = namaMandorInput.value;
        const namaPenyadap = namaPenyadapInput.value;
        const awal = tglAwal.value;
        const akhir = tglAkhir.value;
        const jenis = filterJenis.value;

        let url = "http://localhost:8080/api/monitoring/smart-search";  // API endpoint baru

        // Membuat parameter query string
        let params = new URLSearchParams();
        params.append("namaMandor", namaMandor);
        params.append("penyadap", namaPenyadap);
        params.append("tipe", jenis);

        // Filter berdasarkan rentang tanggal
        if (awal && akhir) {
            params.append("tanggalAwal", awal);
            params.append("tanggalAkhir", akhir);
        } else if (awal && !akhir) {
            // Jika hanya ada tanggal awal
            params.append("tanggalAwal", awal);
            params.append("tanggalAkhir", awal); // Bisa disesuaikan dengan logika lain
        } else {
            // Jika tidak ada tanggal, gunakan tanggal hari ini
            const today = new Date().toISOString().split("T")[0]; // format YYYY-MM-DD
            params.append("tanggalAwal", today);
            params.append("tanggalAkhir", today);
        }

        // Gabungkan URL dengan query parameters
        url += `?${params.toString()}`;

        // Ambil data dari API dan render
        dataArr = await fetchData(url);
        renderTable(dataArr);
    });

    // Fungsi untuk memuat data berdasarkan hari ini
    function loadTodayData() {
        const today = new Date().toISOString().split("T")[0];
        tglAwal.value = today;
        tglAkhir.value = today;
        const url = `http://localhost:8080/api/monitoring/smart-search?namaMandor=${namaMandorInput.value}&penyadap=${namaPenyadapInput.value}&tanggalAwal=${today}&tanggalAkhir=${today}&tipe=${filterJenis.value}`;
        fetchData(url).then(data => renderTable(data));
    }

    // Fungsi untuk memuat data berdasarkan minggu ini
    function loadWeekData() {
        const currentDate = new Date();
        const firstDayOfWeek = new Date(currentDate.setDate(currentDate.getDate() - currentDate.getDay()));
        const lastDayOfWeek = new Date(currentDate.setDate(currentDate.getDate() - currentDate.getDay() + 6));

        tglAwal.value = firstDayOfWeek.toISOString().split("T")[0];
        tglAkhir.value = lastDayOfWeek.toISOString().split("T")[0];

        const url = `http://localhost:8080/api/monitoring/smart-search?namaMandor=${namaMandorInput.value}&penyadap=${namaPenyadapInput.value}&tanggalAwal=${tglAwal.value}&tanggalAkhir=${tglAkhir.value}&tipe=${filterJenis.value}`;
        fetchData(url).then(data => renderTable(data));
    }

    // Fungsi untuk memuat data berdasarkan bulan ini
    function loadMonthData() {
        const currentDate = new Date();
        const firstDayOfMonth = new Date(currentDate.getFullYear(), currentDate.getMonth(), 1);
        const lastDayOfMonth = new Date(currentDate.getFullYear(), currentDate.getMonth() + 1, 0);

        tglAwal.value = firstDayOfMonth.toISOString().split("T")[0];
        tglAkhir.value = lastDayOfMonth.toISOString().split("T")[0];

        const url = `http://localhost:8080/api/monitoring/smart-search?namaMandor=${namaMandorInput.value}&penyadap=${namaPenyadapInput.value}&tanggalAwal=${tglAwal.value}&tanggalAkhir=${tglAkhir.value}&tipe=${filterJenis.value}`;
        fetchData(url).then(data => renderTable(data));
    }

    // Fungsi untuk reset form
    function clearAll() {
        tglAwal.value = "";
        tglAkhir.value = "";
        filterJenis.value = "";
        namaMandorInput.value = "";
        namaPenyadapInput.value = "";
        renderTable([]); // Clear table
    }
});
