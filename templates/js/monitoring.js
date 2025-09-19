document.addEventListener("DOMContentLoaded", () => {
    const wrapper = document.getElementById("bakuTableWrapper");
    const btnSearch = document.getElementById("searchBtn");
    const tglAwal = document.getElementById("searchTanggalAwal");
    const tglAkhir = document.getElementById("searchTanggalAkhir");
    const filterJenis = document.getElementById("filterJenis");
    const bakuTableBody = document.getElementById("bakuTableBody");

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
                <td>${item.mandor?.mandor || "-"}</td>
                <td>${item.mandor?.afdeling || "-"}</td>
                <td>${item.penyadap?.nik || "-"}</td>
                <td>${item.penyadap?.nama_penyadap || "-"}</td>
                <td>${item.periode || "-"}</td>
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
        
        const awal = tglAwal.value;
        const akhir = tglAkhir.value;
        const jenis = filterJenis.value;

        let url = "/api/baku/detail";  // Default API endpoint

        // Filter berdasarkan rentang tanggal
        if (awal && akhir) {
            url = `/api/baku/detail/range?tanggal_mulai=${awal}&tanggal_selesai=${akhir}&tipe=${jenis}`;
        } else if (awal && !akhir) {
            // Cari berdasarkan satu tanggal
            url = `/api/baku/detail/${awal}?tipe=${jenis}`;
        } else if (!awal && !akhir) {
            // Jika tidak ada tanggal, gunakan tanggal hari ini
            url = `/api/baku/detail/today?tipe=${jenis}`;
        }

        // Ambil data dan render
        dataArr = await fetchData(url);
        renderTable(dataArr);
    });

    // Load awal (hari ini)
    (async () => {
        const dataArr = await fetchData(`/api/baku/detail/today`);
        renderTable(dataArr);
    })();
});
