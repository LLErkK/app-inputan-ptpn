document.addEventListener("DOMContentLoaded", () => {
    // Tampilkan tabel data produksi baku saat tombol diklik
    setTimeout(function() {
        const showBakuTableBtn = document.getElementById("showBakuTableBtn");
        const bakuTableWrapper = document.getElementById("bakuTableWrapper");
        if (showBakuTableBtn && bakuTableWrapper) {
            showBakuTableBtn.addEventListener("click", function() {
                if (bakuTableWrapper.style.display === "none" || bakuTableWrapper.style.display === "") {
                    bakuTableWrapper.style.display = "block";
                    showBakuTableBtn.textContent = "Sembunyikan Data Produksi Baku";
                } else {
                    bakuTableWrapper.style.display = "none";
                    showBakuTableBtn.textContent = "Tampilkan Data Produksi Baku";
                }
            });
        }
    }, 0);
    // === POPUP MANDOR ===
    const tambahMandorBtn = document.getElementById("tambahMandor");
    const popupMandor = document.getElementById("popupMandor");
    const closePopupMandorBtn = document.getElementById("closePopupMandor");
    const mandorTableBody = document.getElementById("mandorTableBody");

    // Dummy data mandor, bisa diganti dengan data dari server
    let mandorList = [
        { nama: "Budi", tahun: "2015", afdeling: "A1" },
        { nama: "Siti", tahun: "2017", afdeling: "A2" },
        { nama: "Joko", tahun: "2016", afdeling: "A3" },
        { nama: "Siti", tahun: "2017", afdeling: "A4" },
        { nama: "Siti", tahun: "2017", afdeling: "A5" },
        { nama: "Siti", tahun: "2017", afdeling: "A6" }
    ];

    tambahMandorBtn.addEventListener("click", function() {
        popupMandor.style.display = "flex";
        document.body.classList.add("popup-mandor-active");
        renderMandorTable();
    });
    closePopupMandorBtn.addEventListener("click", function() {
        popupMandor.style.display = "none";
        document.body.classList.remove("popup-mandor-active");
    });

    function renderMandorTable() {
        mandorTableBody.innerHTML = "";
        mandorList.forEach((mandor, idx) => {
            const tr = document.createElement("tr");
            tr.innerHTML = `
                <td>${mandor.nama}</td>
                <td>${mandor.tahun}</td>
                <td>${mandor.afdeling}</td>
                <td>
                        <button class="action-btn update update-btn" data-idx="${idx}" title="Update">
                            <span class="action-icon"> 
                                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#0093E9" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 1 1 3 3L7 19l-4 1 1-4 12.5-12.5z"/></svg>
                            </span>
                        </button>
                        <button class="action-btn delete delete-btn" data-idx="${idx}" title="Delete">
                            <span class="action-icon">
                                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#ff3b3b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2v2"/><line x1="10" y1="11" x2="10" y2="17"/><line x1="14" y1="11" x2="14" y2="17"/></svg>
                            </span>
                        </button>
                    </td>
            `;
            mandorTableBody.appendChild(tr);
        });
    }
    const form = document.getElementById("bakuForm");
    const tableBody = document.querySelector("#bakuTableData tbody");
    // Fungsi untuk menampilkan tabel sesuai mandor yang dipilih
    function renderMandorTables(dataArr) {
        const wrapper = document.getElementById('bakuTableWrapper');
        wrapper.innerHTML = '';
        // Kelompokkan data berdasarkan mandor
        const mandorGroups = {};
        dataArr.forEach(item => {
            const mandorName = item.mandor ? item.mandor.mandor : '-';
            if (!mandorGroups[mandorName]) mandorGroups[mandorName] = [];
            mandorGroups[mandorName].push(item);
        });
        Object.keys(mandorGroups).forEach(mandorName => {
            const table = document.createElement('table');
            table.className = 'baku-table';
            table.innerHTML = `
                <caption style="font-weight:bold; text-align:left; margin-bottom:8px; color:#0093E9;">Mandor: ${mandorName}</caption>
                <thead>
                    <tr>
                        <th>NIK</th>
                        <th>Penyadap</th>
                        <th>Basah Latek</th>
                        <th>Basah Lump</th>
                        <th>Sheet</th>
                        <th>Br.Cr</th>
                        <th>Action</th>
                    </tr>
                </thead>
                <tbody></tbody>
            `;
            const tbody = table.querySelector('tbody');
            mandorGroups[mandorName].forEach((item, idx) => {
                const tr = document.createElement('tr');
                tr.innerHTML = `
                    <td>${item.penyadap ? item.penyadap.nik : '-'}</td>
                    <td>${item.penyadap ? item.penyadap.nama_penyadap : '-'}</td>
                    <td>${item.basahLatex}</td>
                    <td>${item.basahLump}</td>
                    <td>${item.sheet}</td>
                    <td>${item.brCr}</td>
                    <td>
                        <button class="action-btn update update-btn" data-idx="${idx}" title="Update">
                            <span class="action-icon"> 
                                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#0093E9" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 1 1 3 3L7 19l-4 1 1-4 12.5-12.5z"/></svg>
                            </span>
                        </button>
                        <button class="action-btn delete delete-btn" data-idx="${idx}" title="Delete">
                            <span class="action-icon">
                                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#ff3b3b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2v2"/><line x1="10" y1="11" x2="10" y2="17"/><line x1="14" y1="11" x2="14" y2="17"/></svg>
                            </span>
                        </button>
                    </td>
                `;
                tbody.appendChild(tr);
            });
            wrapper.appendChild(table);
        });
    }

    const inputNama = document.getElementById("namaPenyadap");
    const inputNik = document.getElementById("nik");
    const inputIdPenyadap = document.getElementById("idPenyadap");
    const dropdown = document.getElementById("namaDropdown");

    // ================== AUTOCOMPLETE PENYADAP ==================
    inputNama.addEventListener("input", async () => {
        const query = inputNama.value.trim();
        if (query.length < 2) {
            dropdown.style.display = "none";
            return;
        }

        try {
            // 🔥 update endpoint ke /api/penyadap/search
            const res = await fetch(`/api/penyadap/search?nama=${encodeURIComponent(query)}`);
            const data = await res.json();

            if (!data.success || !data.data || data.data.length === 0) {
                dropdown.innerHTML = "<div style='padding:5px;'>Tidak ditemukan</div>";
                dropdown.style.display = "block";
                return;
            }

            dropdown.innerHTML = "";
            data.data.forEach(item => {
                const option = document.createElement("div");
                option.textContent = `${item.nama_penyadap} (${item.nik})`;
                option.style.padding = "5px";
                option.style.cursor = "pointer";

                option.addEventListener("click", () => {
                    inputNama.value = item.nama_penyadap;
                    inputNik.value = item.nik;
                    inputIdPenyadap.value = item.id; // simpan ID penyadap
                    dropdown.style.display = "none";
                });

                dropdown.appendChild(option);
            });
            dropdown.style.display = "block";
        } catch (err) {
            console.error("Error fetching penyadap:", err);
        }
    });

    document.addEventListener("click", (e) => {
        if (!dropdown.contains(e.target) && e.target !== inputNama) {
            dropdown.style.display = "none";
        }
    });

    // ================== SUBMIT FORM ==================
    form.addEventListener("submit", async (e) => {
        e.preventDefault();

        const payload = {
            idBakuMandor: parseInt(document.getElementById("mandor").value) || 0,
            idPenyadap: parseInt(inputIdPenyadap.value) || 0,
            basahLatex: parseFloat(document.getElementById("latek").value) || 0,
            basahLump: parseFloat(document.getElementById("lump").value) || 0,
            sheet: parseFloat(document.getElementById("sheet").value) || 0,
            brCr: parseFloat(document.getElementById("brcr").value) || 0,
        };


        try {
            const res = await fetch("/api/baku", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(payload),
            });
            const data = await res.json();
            if (data.success) {
                alert("Data berhasil disimpan");
                form.reset();
                inputNik.value = "";
                inputIdPenyadap.value = "";
                loadTable();
            } else {
                alert("Gagal simpan: " + data.message);
            }
        } catch (err) {
            console.error("Error submit:", err);
        }
    });

    // ================== LOAD TABLE ==================
    async function loadTable() {
        try {
            const res = await fetch("/api/baku");
            const data = await res.json();

            if (!data.success) {
                console.error("Gagal load data:", data.message);
                return;
            }

    renderMandorTables(data.data);
    addActionListeners(data.data);
    // Tambahkan event handler untuk tombol update dan delete
    function addActionListeners(dataArr) {
        document.querySelectorAll('.update-btn').forEach(btn => {
            btn.onclick = function() {
                const idx = btn.getAttribute('data-idx');
                const item = dataArr[idx];
                // Isi form dengan data yang dipilih
                document.getElementById("mandor").value = item.mandor ? item.mandor.id : "";
                document.getElementById("namaPenyadap").value = item.penyadap ? item.penyadap.nama_penyadap : "";
                document.getElementById("nik").value = item.penyadap ? item.penyadap.nik : "";
                document.getElementById("idPenyadap").value = item.penyadap ? item.penyadap.id : "";
                document.getElementById("latek").value = item.basahLatex;
                document.getElementById("lump").value = item.basahLump;
                document.getElementById("sheet").value = item.sheet;
                document.getElementById("brcr").value = item.brCr;
                // Simpan index edit ke window (atau variabel global jika perlu)
                window.editBakuIndex = idx;
                form.querySelector('.save-btn').textContent = 'Update';
            };
        });
        document.querySelectorAll('.delete-btn').forEach(btn => {
            btn.onclick = async function() {
                const idx = btn.getAttribute('data-idx');
                if (confirm('Hapus data ini?')) {
                    // Hapus data via API
                    const item = dataArr[idx];
                    try {
                        const res = await fetch(`/api/baku/${item.id}`, { method: "DELETE" });
                        const data = await res.json();
                        if (data.success) {
                            alert("Data berhasil dihapus");
                            loadTable();
                        } else {
                            alert("Gagal hapus: " + data.message);
                        }
                    } catch (err) {
                        console.error("Error delete:", err);
                    }
                }
            };
        });
    }
        } catch (err) {
            console.error("Error loadTable:", err);
        }
    }

    // Initial load
    loadTable();
});
