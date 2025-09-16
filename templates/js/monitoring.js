document.addEventListener("DOMContentLoaded", () => {
    const wrapper = document.getElementById("bakuTableWrapper");
    const btnSearch = document.getElementById("searchBtn");
    const tglAwal = document.getElementById("searchTanggalAwal");
    const tglAkhir = document.getElementById("searchTanggalAkhir");

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
        wrapper.innerHTML = "";

        if (!dataArr.length) {
            wrapper.innerHTML = `<div class="empty-state">Data tidak ditemukan.</div>`;
            return;
        }

        // Buat tabel
        const table = document.createElement("table");
        table.className = "baku-table";
        table.innerHTML = `
          <thead>
            <tr>
              <th>Tanggal</th>
              <th>Mandor</th>
              <th>Afdeling</th>
              <th>NIK</th>
              <th>Nama Penyadap</th>
              <th>Periode</th>
              <th>Basah Latek</th>
              <th>Sheet</th>
              <th>Basah Lump</th>
              <th>Br.Cr</th>
            </tr>
          </thead>
          <tbody></tbody>
        `;
        const tbody = table.querySelector("tbody");

        dataArr.forEach(it => {
            const tr = document.createElement("tr");
            tr.innerHTML = `
              <td>${it.tanggal || "-"}</td>
              <td>${it.mandor?.mandor || "-"}</td>
              <td>${it.mandor?.afdeling || "-"}</td>
              <td>${it.penyadap?.nik || "-"}</td>
              <td>${it.penyadap?.nama_penyadap || "-"}</td>
              <td>${it.periode || "-"}</td>
              <td>${it.basahLatex || 0}</td>
              <td>${it.sheet || 0}</td>
              <td>${it.basahLump || 0}</td>
              <td>${it.brCr || 0}</td>
            `;
            tbody.appendChild(tr);
        });

        wrapper.appendChild(table);
    }

    // Event search
    btnSearch.addEventListener("click", async () => {
        let dataArr = [];

        const awal = tglAwal.value;
        const akhir = tglAkhir.value;

        if (awal && akhir) {
            // Cari berdasarkan range tanggal
            dataArr = await fetchData(`/api/baku/detail/range?tanggal_mulai=${awal}&tanggal_selesai=${akhir}`);
        } else if (awal && !akhir) {
            // Cari 1 tanggal
            dataArr = await fetchData(`/api/baku/detail/${awal}`);
        } else {
            // Default: hari ini
            dataArr = await fetchData(`/api/baku/rekap/today`);
        }

        renderTable(dataArr);
    });

    // Load awal (hari ini)
    (async () => {
        const dataArr = await fetchData(`/api/baku/rekap/today`);
        renderTable(dataArr);
    })();
});
