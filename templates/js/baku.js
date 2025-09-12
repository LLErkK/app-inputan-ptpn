document.addEventListener("DOMContentLoaded", () => {
    let allBakuData = [];
    let editingId = null;

    // ========= Popup Mandor =========
    const tambahMandorBtn = document.getElementById("tambahMandor");
    const popupMandor = document.getElementById("popupMandor");
    const closePopupMandorBtn = document.getElementById("closePopupMandor");
    const mandorTableBody = document.getElementById("mandorTableBody");
    const formMandorBaru = document.getElementById("formMandorBaru");

    if (tambahMandorBtn) {
        tambahMandorBtn.addEventListener("click", () => {
            popupMandor.style.display = "flex";
            loadMandorList();
        });
    }
    if (closePopupMandorBtn) {
        closePopupMandorBtn.addEventListener("click", () => {
            popupMandor.style.display = "none";
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
            const res = await fetch("/api/mandor", {
                method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(payload)
            });
            const data = await res.json();
            if (data.success) {
                alert("Mandor ditambahkan!");
                formMandorBaru.reset();
                loadMandorList();
                loadMandorOptions();
            } else alert("Gagal: " + data.message);
        });
    }

    async function loadMandorList() {
        const res = await fetch("/api/mandor");
        const data = await res.json();
        mandorTableBody.innerHTML = "";
        if (data.success) {
            data.data.forEach(m => {
                const tr = document.createElement("tr");
                tr.innerHTML = `<td>${m.mandor}</td><td>${m.tahun_tanam}</td><td>${m.afdeling}</td>
          <td><button data-id="${m.id}" class="delete-mandor-btn">Hapus</button></td>`;
                mandorTableBody.appendChild(tr);
            });
            document.querySelectorAll(".delete-mandor-btn").forEach(btn => {
                btn.addEventListener("click", async () => {
                    const id = btn.getAttribute("data-id");
                    if (confirm("Hapus mandor?")) {
                        await fetch(`/api/mandor/${id}`, { method: "DELETE" });
                        loadMandorList(); loadMandorOptions();
                    }
                });
            });
        }
    }
    async function loadMandorOptions() {
        const res = await fetch("/api/mandor");
        const data = await res.json();
        const select = document.getElementById("mandor");
        select.innerHTML = '<option value="">-- Pilih Mandor --</option>';
        if (data.success) {
            data.data.forEach(m => {
                const opt = document.createElement("option");
                opt.value = m.id;
                opt.textContent = `${m.mandor} (${m.afdeling}) ${m.tahun_tanam}`;
                select.appendChild(opt);
            });
        }
    }

    // ========= Autocomplete Penyadap =========
    const inputNama = document.getElementById("namaPenyadap");
    const inputNik = document.getElementById("nik");
    const inputIdPenyadap = document.getElementById("idPenyadap");
    const dropdown = document.getElementById("namaDropdown");

    if (inputNama) {
        inputNama.addEventListener("input", async () => {
            const q = inputNama.value.trim();
            if (q.length < 2) { dropdown.style.display = "none"; return; }
            const res = await fetch(`/api/penyadap/search?nama=${encodeURIComponent(q)}`);
            const data = await res.json();
            dropdown.innerHTML = "";
            if (!data.success || !data.data.length) {
                dropdown.innerHTML = "<div style='padding:8px'>Tidak ditemukan</div>";
                dropdown.style.display = "block"; return;
            }
            data.data.forEach(item => {
                const opt = document.createElement("div");
                opt.textContent = `${item.nama_penyadap} (${item.nik})`;
                opt.className = "dropdown-item";
                opt.addEventListener("click", () => {
                    inputNama.value = item.nama_penyadap;
                    inputNik.value = item.nik;
                    inputIdPenyadap.value = item.id;
                    dropdown.style.display = "none";
                });
                dropdown.appendChild(opt);
            });
            dropdown.style.display = "block";
        });
        document.addEventListener("click", e => {
            if (!dropdown.contains(e.target) && e.target !== inputNama) dropdown.style.display = "none";
        });
    }

    // ========= Submit Form =========
    const form = document.getElementById("bakuForm");
    if (form) {
        form.addEventListener("submit", async e => {
            e.preventDefault();
            const payload = {
                idBakuMandor: parseInt(document.getElementById("mandor").value),
                idPenyadap: parseInt(inputIdPenyadap.value),
                tipe: document.getElementById("jenis").value,
                basahLatex: parseFloat(document.getElementById("latek").value),
                basahLump: parseFloat(document.getElementById("lump").value),
                sheet: parseFloat(document.getElementById("sheet").value),
                brCr: parseFloat(document.getElementById("brcr").value),
            };
            const url = editingId ? `/api/baku/${editingId}` : "/api/baku";
            const method = editingId ? "PUT" : "POST";
            const res = await fetch(url, { method, headers: { "Content-Type": "application/json" }, body: JSON.stringify(payload) });
            const data = await res.json();
            if (data.success) {
                alert(editingId ? "Data diperbarui" : "Data disimpan");
                editingId = null; form.reset(); inputNik.value = ""; inputIdPenyadap.value = "";
                loadSummaryTable(); loadBakuTable();
            } else alert("Gagal: " + data.message);
        });
    }

    // ========= Summary Table =========
    async function loadSummaryTable() {
        const today = new Date().toISOString().split("T")[0];
        const res = await fetch(`/api/reporting/mandor/${today}`);
        const data = await res.json();
        const body = document.querySelector("#summaryTable tbody");
        body.innerHTML = "";
        if (data.success && data.data.length) {
            data.data.forEach(it => {
                const tr = document.createElement("tr");
                tr.innerHTML = `<td>${it.mandor}</td><td>${it.jumlah_pabrik_basah_latek}</td>
          <td>${it.jumlah_kebun_basah_latek}</td><td>${it.jumlah_sheet}</td>
          <td>${it.k3_sheet}</td><td>${it.jumlah_pabrik_basah_lump}</td>
          <td>${it.jumlah_kebun_basah_lump}</td><td>${it.jumlah_br_cr}</td><td>${it.k3_br_cr}</td>`;
                body.appendChild(tr);
            });
        } else {
            body.innerHTML = `<tr><td colspan="9">Belum ada data hari ini (${today})</td></tr>`;
        }
    }

    // ========= Render Detail Per Mandor & Penyadap =========
    function renderMandorTables(dataArr) {
        const wrapper = document.getElementById("bakuTableWrapper");
        wrapper.innerHTML = "";
        const groups = {};
        dataArr.forEach(it => {
            const key = it.mandor ? it.mandor.mandor : "Unknown";
            if (!groups[key]) groups[key] = [];
            groups[key].push(it);
        });
        Object.keys(groups).forEach(mandor => {
            const table = document.createElement("table");
            table.className = "baku-table";
            table.innerHTML = `<caption>Mandor: ${mandor}</caption>
        <thead><tr><th>NIK</th><th>Penyadap</th><th>Latek</th><th>Lump</th><th>Sheet</th><th>Br.Cr</th><th>Action</th></tr></thead>
        <tbody></tbody>`;
            const tbody = table.querySelector("tbody");
            groups[mandor].forEach(it => {
                const tr = document.createElement("tr");
                tr.innerHTML = `<td>${it.penyadap?.nik || "-"}</td><td>${it.penyadap?.nama_penyadap || "-"}</td>
          <td>${it.basahLatex}</td><td>${it.basahLump}</td><td>${it.sheet}</td><td>${it.brCr}</td>
          <td><button class="edit-btn" data-id="${it.id}">Edit</button>
              <button class="delete-btn" data-id="${it.id}">Delete</button></td>`;
                tbody.appendChild(tr);
            });
            wrapper.appendChild(table);
        });
        addBakuActionListeners();
    }

    function addBakuActionListeners() {
        document.querySelectorAll(".edit-btn").forEach(btn => {
            btn.onclick = () => {
                const id = parseInt(btn.dataset.id);
                const item = allBakuData.find(d => d.id === id);
                if (item) {
                    document.getElementById("mandor").value = item.idBakuMandor;
                    document.getElementById("jenis").value = item.tipe;
                    if (item.penyadap) {
                        inputNama.value = item.penyadap.nama_penyadap;
                        inputNik.value = item.penyadap.nik;
                        inputIdPenyadap.value = item.penyadap.id;
                    }
                    document.getElementById("latek").value = item.basahLatex;
                    document.getElementById("lump").value = item.basahLump;
                    document.getElementById("sheet").value = item.sheet;
                    document.getElementById("brcr").value = item.brCr;
                    editingId = id;
                }
            };
        });
        document.querySelectorAll(".delete-btn").forEach(btn => {
            btn.onclick = async () => {
                const id = btn.dataset.id;
                if (confirm("Hapus data?")) {
                    await fetch(`/api/baku/${id}`, { method: "DELETE" });
                    loadBakuTable(); loadSummaryTable();
                }
            };
        });
    }

    // ========= Load Detail =========
    async function loadBakuTable() {
        const today = new Date().toISOString().split("T")[0];
        const res = await fetch(`/api/baku?tanggal=${today}`);
        const data = await res.json();
        if (data.success) {
            allBakuData = data.data;
            renderMandorTables(allBakuData);
        }
    }

    // ========= Toggle Button =========
    const showBtn = document.getElementById("showBakuTableBtn");
    if (showBtn) {
        showBtn.addEventListener("click", () => {
            const wrapper = document.getElementById("bakuTableWrapper");
            if (wrapper.style.display === "none" || wrapper.style.display === "") {
                wrapper.style.display = "block";
                showBtn.textContent = "Sembunyikan Data Produksi Baku";
                loadBakuTable();
            } else {
                wrapper.style.display = "none";
                showBtn.textContent = "Tampilkan Data Produksi Baku";
            }
        });
    }

    // ========= Init =========
    loadMandorOptions();
    loadSummaryTable();
    setInterval(loadSummaryTable, 30000);
});
