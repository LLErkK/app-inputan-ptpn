// /public/js/baku.js
document.addEventListener("DOMContentLoaded", () => {
    // ================= Helpers =================
    function todayLocalYYYYMMDD() {
        const d = new Date();
        const y = d.getFullYear();
        const m = String(d.getMonth() + 1).padStart(2, "0");
        const day = String(d.getDate()).padStart(2, "0");
        return `${y}-${m}-${day}`;
    }
    function safeText(v, fallback = "-") {
        return (v === undefined || v === null || v === "") ? fallback : v;
    }
    function n(v, def = 0) {
        const x = parseFloat(v);
        return Number.isFinite(x) ? x : def;
    }

    let allBakuData = [];   // DETAIL per penyadap untuk Edit/Delete → dari /api/baku?tanggal=YYYY-MM-DD
    let editingId = null;   // simpan sebagai STRING

    // ========= Popup Mandor =========
    const tambahMandorBtn = document.getElementById("tambahMandor");
    const popupMandor = document.getElementById("popupMandor");
    const closePopupMandorBtn = document.getElementById("closePopupMandor");
    const mandorTableBody = document.getElementById("mandorTableBody");
    const formMandorBaru = document.getElementById("formMandorBaru");

    if (tambahMandorBtn) {
        tambahMandorBtn.addEventListener("click", () => {
            if (popupMandor) popupMandor.style.display = "flex";
            loadMandorList();
        });
    }
    if (closePopupMandorBtn) {
        closePopupMandorBtn.addEventListener("click", () => {
            if (popupMandor) popupMandor.style.display = "none";
        });
    }
    if (formMandorBaru) {
        formMandorBaru.addEventListener("submit", async e => {
            e.preventDefault();
            const payload = {
                mandor: document.getElementById("inputNamaMandor").value.trim(),
                tahun_tanam: parseInt(document.getElementById("inputTahunTanam").value) || 0,
                afdeling: document.getElementById("inputAfdeling").value.trim()
            };
            try {
                const res = await fetch("/api/mandor", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(payload)
                });
                const data = await res.json();
                if (data.success) {
                    alert("Mandor ditambahkan!");
                    formMandorBaru.reset();
                    loadMandorList();
                    loadMandorOptions();
                } else {
                    alert("Gagal: " + data.message);
                }
            } catch (err) {
                alert("Error: " + err.message);
            }
        });
    }

    async function loadMandorList() {
        if (!mandorTableBody) return;
        mandorTableBody.innerHTML = "";
        try {
            const res = await fetch("/api/mandor");
            const data = await res.json();
            if (data.success && Array.isArray(data.data)) {
                data.data.forEach(m => {
                    const tr = document.createElement("tr");
                    tr.innerHTML = `
            <td>${safeText(m.mandor)}</td>
            <td>${safeText(m.tahun_tanam, "-")}</td>
            <td>${safeText(m.afdeling, "-")}</td>
            <td><button data-id="${String(m.id)}" class="delete-mandor-btn">Hapus</button></td>`;
                    mandorTableBody.appendChild(tr);
                });
                mandorTableBody.querySelectorAll(".delete-mandor-btn").forEach(btn => {
                    btn.addEventListener("click", async () => {
                        const id = String(btn.getAttribute("data-id"));
                        if (confirm("Hapus mandor?")) {
                            await fetch(`/api/mandor/${encodeURIComponent(id)}`, { method: "DELETE" });
                            loadMandorList();
                            loadMandorOptions();
                        }
                    });
                });
            } else {
                mandorTableBody.innerHTML = `<tr><td colspan="4">Tidak ada data mandor.</td></tr>`;
            }
        } catch (e) {
            mandorTableBody.innerHTML = `<tr><td colspan="4">Error: ${e.message}</td></tr>`;
        }
    }

    async function loadMandorOptions() {
        const select = document.getElementById("mandor");
        if (!select) return;
        select.innerHTML = '<option value="">-- Pilih Mandor --</option>';
        try {
            const res = await fetch("/api/mandor");
            const data = await res.json();
            if (data.success && Array.isArray(data.data)) {
                data.data.forEach(m => {
                    const opt = document.createElement("option");
                    opt.value = String(m.id);
                    opt.textContent = `${safeText(m.mandor)} (${safeText(m.afdeling, "-")}) ${safeText(m.tahun_tanam, "-")}`;
                    select.appendChild(opt);
                });
            }
        } catch (e) {
            // biarkan kosong
        }
    }

    // ========= Autocomplete Penyadap =========
    const inputNama = document.getElementById("namaPenyadap");
    const inputNik = document.getElementById("nik");
    const inputIdPenyadap = document.getElementById("idPenyadap");
    const dropdown = document.getElementById("namaDropdown");

    if (inputNama && dropdown) {
        inputNama.addEventListener("input", async () => {
            const q = inputNama.value.trim();
            if (q.length < 2) { dropdown.style.display = "none"; return; }
            try {
                const res = await fetch(`/api/penyadap/search?nama=${encodeURIComponent(q)}`);
                const data = await res.json();
                dropdown.innerHTML = "";
                if (!data.success || !data.data || !data.data.length) {
                    dropdown.innerHTML = "<div style='padding:8px'>Tidak ditemukan</div>";
                    dropdown.style.display = "block"; return;
                }
                data.data.forEach(item => {
                    const opt = document.createElement("div");
                    opt.textContent = `${safeText(item.nama_penyadap)} (${safeText(item.nik)})`;
                    opt.className = "dropdown-item";
                    opt.addEventListener("click", () => {
                        inputNama.value = safeText(item.nama_penyadap, "");
                        inputNik.value = safeText(item.nik, "");
                        inputIdPenyadap.value = String(item.id || "");
                        dropdown.style.display = "none";
                    });
                    dropdown.appendChild(opt);
                });
                dropdown.style.display = "block";
            } catch (e) {
                dropdown.innerHTML = "<div style='padding:8px'>Gagal memuat</div>";
                dropdown.style.display = "block";
            }
        });
        document.addEventListener("click", e => {
            if (!dropdown.contains(e.target) && e.target !== inputNama) dropdown.style.display = "none";
        });
    }

    // ========= Submit Form (Create/Update) =========
    const form = document.getElementById("bakuForm");
    if (form) {
        form.addEventListener("submit", async e => {
            e.preventDefault();
            const idMandorStr = document.getElementById("mandor").value;
            const idPenyadapStr = inputIdPenyadap ? inputIdPenyadap.value : "";

            const payload = {
                idBakuMandor: idMandorStr ? parseInt(idMandorStr) : null,
                idPenyadap: idPenyadapStr ? parseInt(idPenyadapStr) : null,
                tipe: document.getElementById("jenis").value,
                basahLatex: n(document.getElementById("latek").value),
                basahLump: n(document.getElementById("lump").value),
                sheet: n(document.getElementById("sheet").value),
                brCr: n(document.getElementById("brcr").value),
            };

            const url = editingId ? `/api/baku/${encodeURIComponent(editingId)}` : "/api/baku";
            const method = editingId ? "PUT" : "POST";

            try {
                const res = await fetch(url, {
                    method,
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(payload)
                });
                const data = await res.json();
                if (data.success) {
                    alert(editingId ? "Data diperbarui" : "Data disimpan");
                    editingId = null;
                    form.reset();
                    if (inputNik) inputNik.value = "";
                    if (inputIdPenyadap) inputIdPenyadap.value = "";
                    const submitBtn = document.getElementById("submitBtn");
                    if (submitBtn) submitBtn.textContent = "Simpan";
                    // refresh kedua tampilan
                    renderRekapMandor();   // ambil dari /api/baku/detail
                    renderDetailBaku();    // ambil dari /api/baku?tanggal=YYYY-MM-DD
                } else {
                    alert("Gagal: " + data.message);
                }
            } catch (err) {
                alert("Error: " + err.message);
            }
        });
    }

    // ========= Rekap Mandor (PAKAI /api/baku/detail) =========
    // ========= Rekap Mandor (PAKAI /api/baku/detail/YYYY-MM-DD) =========
    async function renderRekapMandor() {
        const body = document.querySelector("#summaryTable tbody");
        if (!body) return;

        body.innerHTML = "<tr><td colspan='11'>Loading...</td></tr>";

        const today = todayLocalYYYYMMDD(); // ambil tanggal hari ini

        try {
            const res = await fetch(`/api/baku/detail/${encodeURIComponent(today)}`);
            const json = await res.json();

            if (json && json.success && Array.isArray(json.data) && json.data.length) {
                body.innerHTML = ""; // reset

                json.data.forEach(row => {
                    const tr = document.createElement("tr");
                    tr.innerHTML = `
                    <td>${safeText(row.mandor)}</td>
                    <td>${safeText(row.afdeling)}</td>
                    <td>${safeText(row.tipe)}</td>
                    <td>${n(row.jumlah_pabrik_basah_latek)}</td>
                    <td>${n(row.jumlah_kebun_basah_latek)}</td>
                    <td>${n(row.jumlah_sheet)}</td>
                    <td>${n(row.k3_sheet)}</td>
                    <td>${n(row.jumlah_pabrik_basah_lump)}</td>
                    <td>${n(row.jumlah_kebun_basah_lump)}</td>
                    <td>${n(row.jumlah_br_cr)}</td>
                    <td>${n(row.k3_br_cr)}</td>
                  
                `;
                    body.appendChild(tr);
                });
            } else {
                body.innerHTML = `<tr><td colspan='15'>Belum ada data rekap untuk hari ini (${today}).</td></tr>`;
            }
        } catch (e) {
            body.innerHTML = `<tr><td colspan='15'>Gagal memuat rekap: ${e.message}</td></tr>`;
            console.error(e);
        }
    }

    // ========= DETAIL per Penyadap (PAKAI /api/baku?tanggal=YYYY-MM-DD) =========
    async function renderDetailBaku() {
        const wrapper = document.getElementById("bakuTableWrapper");
        if (!wrapper) return;
        wrapper.innerHTML = "";

        const today = todayLocalYYYYMMDD();
        let dataArr = [];
        try {
            const res = await fetch(`/api/baku/rekap/today`);
            const json = await res.json();
            if (json && json.success && Array.isArray(json.data)) {
                dataArr = json.data;
            }
        } catch (e) {
            // biarkan kosong; akan tampil empty-state
        }

        allBakuData = Array.isArray(dataArr) ? dataArr : [];

        // group by mandor
        const groups = {};
        allBakuData.forEach(it => {
            const key = (it && it.mandor && it.mandor.mandor) ? it.mandor.mandor : "Unknown";
            if (!groups[key]) groups[key] = [];
            groups[key].push(it);
        });

        const mandorNames = Object.keys(groups);
        if (!mandorNames.length) {
            wrapper.innerHTML = `<div class="empty-state">Belum ada data detail untuk hari ini (${today}).</div>`;
            return;
        }

        mandorNames.forEach(mandor => {
            const table = document.createElement("table");
            table.className = "baku-table";
            table.innerHTML = `
        <caption>Mandor: ${safeText(mandor)}</caption>
        <thead>
          <tr><th>NIK</th><th>Penyadap</th><th>Latek</th><th>Lump</th><th>Sheet</th><th>Br.Cr</th><th>Action</th></tr>
        </thead>
        <tbody></tbody>`;
            const tbody = table.querySelector("tbody");

            groups[mandor].forEach(it => {
                const nik  = (it && it.penyadap && it.penyadap.nik) ? it.penyadap.nik : "-";
                const nama = (it && it.penyadap && it.penyadap.nama_penyadap) ? it.penyadap.nama_penyadap : "-";
                const tr = document.createElement("tr");
                tr.innerHTML = `
          <td>${nik}</td>
          <td>${nama}</td>
          <td>${safeText(it.basahLatex, 0)}</td>
          <td>${safeText(it.basahLump, 0)}</td>
          <td>${safeText(it.sheet, 0)}</td>
          <td>${safeText(it.brCr, 0)}</td>
          <td>
            <button class="edit-btn" data-id="${String(it.id)}">Edit</button>
            <button class="delete-btn" data-id="${String(it.id)}">Delete</button>
          </td>`;
                tbody.appendChild(tr);
            });
            wrapper.appendChild(table);
        });
    }

    // ========= Event Delegation untuk Edit/Delete =========
    const bakuTableWrapper = document.getElementById("bakuTableWrapper");
    if (bakuTableWrapper) {
        bakuTableWrapper.addEventListener("click", async (e) => {
            const editBtn = e.target.closest && e.target.closest(".edit-btn");
            const delBtn  = e.target.closest && e.target.closest(".delete-btn");

            if (editBtn) {
                const id = String(editBtn.dataset.id);
                const item = allBakuData.find(d => String(d.id) === id);
                if (!item) { alert("Data tidak ditemukan."); return; }

                const mandorSelect = document.getElementById("mandor");
                if (mandorSelect) mandorSelect.value = String(item.idBakuMandor);

                const jenisSelect = document.getElementById("jenis");
                if (jenisSelect) {
                    jenisSelect.value = item.tipe || "";
                    if (!jenisSelect.value && item.tipe) {
                        const opt = document.createElement("option");
                        opt.value = item.tipe;
                        opt.textContent = String(item.tipe).replace(/_/g, " ").toUpperCase();
                        jenisSelect.appendChild(opt);
                        jenisSelect.value = item.tipe;
                    }
                }

                if (item.penyadap) {
                    if (inputNama) inputNama.value = safeText(item.penyadap.nama_penyadap, "");
                    if (inputNik) inputNik.value = safeText(item.penyadap.nik, "");
                    if (inputIdPenyadap) inputIdPenyadap.value = String(item.penyadap.id || "");
                }
                document.getElementById("latek").value = item.basahLatex ?? 0;
                document.getElementById("lump").value  = item.basahLump  ?? 0;
                document.getElementById("sheet").value = item.sheet      ?? 0;
                document.getElementById("brcr").value  = item.brCr       ?? 0;

                editingId = id;
                const submitBtn = document.getElementById("submitBtn");
                if (submitBtn) submitBtn.textContent = "Perbarui";
            }

            if (delBtn) {
                const id = String(delBtn.dataset.id);
                if (confirm("Hapus data?")) {
                    await fetch(`/api/baku/${encodeURIComponent(id)}`, { method: "DELETE" });
                    renderRekapMandor();
                    renderDetailBaku();
                }
            }
        });
    }

    // ========= Toggle Button untuk menampilkan detail =========
    const showBtn = document.getElementById("showBakuTableBtn");
    if (showBtn) {
        showBtn.addEventListener("click", () => {
            const wrapper = document.getElementById("bakuTableWrapper");
            if (!wrapper) return;
            if (wrapper.style.display === "none" || wrapper.style.display === "") {
                wrapper.style.display = "block";
                showBtn.textContent = "Sembunyikan Data Produksi Baku";
                renderDetailBaku();
            } else {
                wrapper.style.display = "none";
                showBtn.textContent = "Tampilkan Data Produksi Baku";
            }
        });
    }

    // ========= Init =========
    loadMandorOptions();
    renderRekapMandor();
    renderDetailBaku();
    setInterval(renderRekapMandor, 30000); // auto refresh setiap 30 detik
});
