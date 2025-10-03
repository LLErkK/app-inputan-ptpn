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
  let mandorDataCache = []; // Cache untuk data mandor

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
        nik: document.getElementById("inputNIKMandor").value.trim(),
        tahun_tanam: parseInt(document.getElementById("inputTahunTanam").value) || 0,
        afdeling: document.getElementById("inputAfdeling").value.trim(),
        tipe: document.getElementById("jenis").value
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
            <td>${safeText(m.nik, "-")}</td>
            <td>${safeText(m.tahun_tanam, "-")}</td>
            <td>${safeText(m.afdeling, "-")}</td>
            <td>${safeText(m.tipe,"-")}</td>
            <td><button data-id="${String(m.id)}" class="delete-mandor-btn">
            <span class="action-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#ff3b3b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2v2"/><line x1="10" y1="11" x2="10" y2="17"/><line x1="14" y1="11" x2="14" y2="17"/></svg>
            </span>
            </button></td>`;
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
        mandorTableBody.innerHTML = `<tr><td colspan="5">Tidak ada data mandor.</td></tr>`;
      }
    } catch (e) {
      mandorTableBody.innerHTML = `<tr><td colspan="5">Error: ${e.message}</td></tr>`;
    }
  }

async function loadMandorOptions() {
    const input = document.getElementById("mandor");
    const datalist = document.getElementById("mandor-list");
    if (!input || !datalist) return;
    
    // Kosongkan datalist sebelum mengisi ulang
    datalist.innerHTML = '';
    
    try {
        const res = await fetch("/api/mandor");
        const data = await res.json();
        if (data.success && Array.isArray(data.data)) {
            // Cache data mandor untuk digunakan nanti
            mandorDataCache = data.data;

            data.data.forEach(m => {
                const option = document.createElement("option");
                option.value = `${safeText(m.mandor)} (${safeText(m.afdeling, "-")}) ${safeText(m.tahun_tanam, "-")} - Tipe: ${safeText(m.tipe)}`;
                option.setAttribute('data-id', m.id); // Menyimpan ID mandor
                datalist.appendChild(option);
            });
        }
    } catch (e) {
        // biarkan kosong
    }
}



  // ========= Function untuk mendapatkan data mandor berdasarkan ID =========
  function getMandorById(id) {
    return mandorDataCache.find(m => String(m.id) === String(id));
  }

  // ========= Popup Penyadap =========
// ========= Popup Penyadap =========
const tambahPenyadapBtn = document.getElementById("tambahPenyadap");
const popupPenyadap = document.getElementById("popupPenyadap");
const closePopupPenyadapBtn = document.getElementById("closePopupPenyadap");
const formPenyadapBaru = document.getElementById("formPenyadapBaru");
const penyadapTableBody = document.getElementById("penyadapTableBody");

// Buka popup
if (tambahPenyadapBtn) {
  tambahPenyadapBtn.addEventListener("click", () => {
    if (popupPenyadap) popupPenyadap.style.display = "flex";
    loadPenyadapList();
  });
}

// Tutup popup
if (closePopupPenyadapBtn) {
  closePopupPenyadapBtn.addEventListener("click", () => {
    if (popupPenyadap) popupPenyadap.style.display = "none";
  });
}

// Submit form penyadap baru
if (formPenyadapBaru) {
  formPenyadapBaru.addEventListener("submit", async e => {
    e.preventDefault();
    
    const nama = document.getElementById("inputNamaPenyadap").value.trim();
    const nik = document.getElementById("inputNIK").value.trim();
    
    if (!nama || !nik) {
      alert("Nama Penyadap dan NIK wajib diisi!");
      return;
    }

    const payload = {
      nama_penyadap: nama,
      nik: nik
    };

    try {
      const res = await fetch("/api/penyadap", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      
      const data = await res.json();
      
      if (res.ok && data.success) {
        alert("Penyadap berhasil ditambahkan!");
        formPenyadapBaru.reset();
        loadPenyadapList();
      } else {
        alert("Gagal: " + (data.message || "Terjadi kesalahan"));
      }
    } catch (err) {
      console.error("Error:", err);
      alert("Error: " + err.message);
    }
  });
}

// Load daftar penyadap
async function loadPenyadapList() {
  if (!penyadapTableBody) {
    console.error("penyadapTableBody tidak ditemukan!");
    return;
  }
  
  console.log("Loading penyadap list...");
  penyadapTableBody.innerHTML = '<tr><td colspan="3">Memuat data...</td></tr>';
  
  try {
    const res = await fetch("/api/penyadap");
    console.log("Response status:", res.status);
    
    if (!res.ok) {
      throw new Error(`HTTP ${res.status}`);
    }
    
    const data = await res.json();
    console.log("Data received:", data);
    
    penyadapTableBody.innerHTML = "";
    
    if (data.success && Array.isArray(data.data) && data.data.length > 0) {
      console.log("Rendering", data.data.length, "penyadaps");
      
      data.data.forEach(p => {
        const tr = document.createElement("tr");
        tr.innerHTML = `
          <td>${safeText(p.nama_penyadap)}</td>
          <td>${safeText(p.nik)}</td>
          <td>
            <button data-id="${String(p.id)}" class="delete-penyadap-btn">
              <span class="action-icon">
                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#ff3b3b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="3 6 5 6 21 6"/>
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2v2"/>
                  <line x1="10" y1="11" x2="10" y2="17"/>
                  <line x1="14" y1="11" x2="14" y2="17"/>
                </svg>
              </span>
            </button>
          </td>
        `;
        penyadapTableBody.appendChild(tr);
      });
      
      // Pasang event listener setelah semua row ditambahkan
      const deleteButtons = penyadapTableBody.querySelectorAll(".delete-penyadap-btn");
      console.log("Found delete buttons:", deleteButtons.length);
      
      deleteButtons.forEach(btn => {
        btn.addEventListener("click", async function() {
          const id = String(this.getAttribute("data-id"));
          console.log("Delete clicked for ID:", id);
          
          if (confirm("Hapus penyadap ini?")) {
            try {
              console.log("Sending DELETE request to:", `/api/penyadap/${id}`);
              
              const res = await fetch(`/api/penyadap/${encodeURIComponent(id)}`, { 
                method: "DELETE" 
              });
              
              console.log("DELETE response status:", res.status);
              const data = await res.json();
              console.log("DELETE response data:", data);
              
              if (res.ok && data.success) {
                alert("Penyadap berhasil dihapus!");
                loadPenyadapList();
              } else {
                alert("Gagal menghapus: " + (data.message || "Unknown error"));
              }
            } catch (err) {
              console.error("Delete error:", err);
              alert("Error: " + err.message);
            }
          }
        });
      });
      
    } else {
      console.log("No data found");
      penyadapTableBody.innerHTML = `<tr><td colspan="3">Belum ada data penyadap.</td></tr>`;
    }
  } catch (e) {
    console.error("Error loading penyadap:", e);
    penyadapTableBody.innerHTML = `<tr><td colspan="3" style="color: red;">Error: ${safeText(e.message)}</td></tr>`;
  }
}  // ========= Autocomplete Penyadap =========
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

      // Dapatkan data mandor untuk mengambil tahun tanam
      const selectedMandor = getMandorById(idMandorStr);
      const tahunTanam = selectedMandor ? selectedMandor.tahun_tanam : null;

      const payload = {
        idBakuMandor: idMandorStr ? parseInt(idMandorStr) : null,
        idPenyadap: idPenyadapStr ? parseInt(idPenyadapStr) : null,
        tahunTanam: tahunTanam, // Tambahkan tahun tanam ke payload
        basahLatex: n(document.getElementById("latek").value),
        basahLump: n(document.getElementById("lump").value),
        sheet: n(document.getElementById("sheet").value),
        brCr: n(document.getElementById("brcr").value),
      };

      console.log("Payload yang dikirim:", payload); // Debug log

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
          const submitBtn = document.querySelector('button[type="submit"]');
          if (submitBtn) submitBtn.textContent = "Save";
          // refresh kedua tampilan
          renderRekapMandor();
          renderDetailBaku();
        } else {
          alert("Gagal: " + data.message);
        }
      } catch (err) {
        alert("Error: " + err.message);
      }
    });
  }

  // ========= Rekap Mandor (PAKAI /api/baku/detail) =========
  async function renderRekapMandor() {
    const body = document.querySelector("#summaryTable tbody");
    if (!body) return;

    body.innerHTML = "";
    try {
      const res = await fetch(`/api/baku/detail/${encodeURIComponent(todayLocalYYYYMMDD())}`);
      const json = await res.json();

      if (json && json.success && Array.isArray(json.data) && json.data.length) {
        json.data.forEach(row => {
          const tr = document.createElement("tr");
          tr.innerHTML = `
            <td>${safeText(row.mandor)}</td>
            <td>${safeText(row.afdeling)}</td>
            <td>${safeText(row.tipe)}</td>
            <td>${safeText(row.jumlah_pabrik_basah_latek, 0)}</td>
            <td>${safeText(row.jumlah_kebun_basah_latek, 0)}</td>
            <td>${safeText(row.jumlah_sheet, 0)}</td>
            <td>${safeText(row.k3_sheet, 0)}</td>
            <td>${safeText(row.jumlah_pabrik_basah_lump, 0)}</td>
            <td>${safeText(row.jumlah_kebun_basah_lump, 0)}</td>
            <td>${safeText(row.jumlah_br_cr, 0)}</td>
            <td>${safeText(row.k3_br_cr, 0)}</td>`;
          body.appendChild(tr);
        });
      } else {
        body.innerHTML = `<tr><td colspan="11">Belum ada data rekap.</td></tr>`;
      }
    } catch (e) {
      body.innerHTML = `<tr><td colspan="11">Gagal memuat rekap.</td></tr>`;
    }
          function safeText(value, defaultValue = "") {
      if (value === null || value === undefined || value === "") {
        return defaultValue;
      }
      
      // Jika nilai berupa angka, batasi angka desimal hingga 2
      if (!isNaN(value)) {
        return parseFloat(value).toFixed(2); // Ubah '2' sesuai jumlah angka desimal yang diinginkan
      }

      return value;
    }

  }

async function renderDetailBaku() {
  const wrapper = document.getElementById("bakuTableWrapper");
  if (!wrapper) return;
  wrapper.innerHTML = "";

  const today = todayLocalYYYYMMDD();
  let dataArr = [];
  try {
    const res = await fetch(`/api/baku?tanggal=${encodeURIComponent(today)}`);
    const json = await res.json();
    if (json && json.success && Array.isArray(json.data)) {
      dataArr = json.data;
    }
  } catch (e) {
    // biarkan kosong; akan tampil empty-state
  }

  allBakuData = Array.isArray(dataArr) ? dataArr : [];

  // group by tahunTanam
  const groupsByTahunTanam = {};
  allBakuData.forEach(it => {
    const key = it.tahunTanam || it.tahun_tanam || "Unknown";
    if (!groupsByTahunTanam[key]) groupsByTahunTanam[key] = [];
    groupsByTahunTanam[key].push(it);
  });

  const tahunTanamNames = Object.keys(groupsByTahunTanam);
  if (!tahunTanamNames.length) {
    wrapper.innerHTML = `<div class="empty-state">Belum ada data detail untuk hari ini (${today}).</div>`;
    return;
  }

  tahunTanamNames.forEach(tahunTanam => {
    const table = document.createElement("table");
    table.className = "baku-table";
    table.innerHTML = `
      <caption>Tahun Tanam: ${safeText(tahunTanam)}</caption>
      <thead>
        <tr><th>Mandor</th><th>NIK</th><th>Penyadap</th><th>Latek</th><th>Lump</th><th>Sheet</th><th>Br.Cr</th><th>Action</th></tr>
      </thead>
      <tbody></tbody>`;
    const tbody = table.querySelector("tbody");

    groupsByTahunTanam[tahunTanam].forEach(it => {
      const nik  = (it && it.penyadap && it.penyadap.nik) ? it.penyadap.nik : "-";
      const nama = (it && it.penyadap && it.penyadap.nama_penyadap) ? it.penyadap.nama_penyadap : "-";
      const mandor = (it && it.mandor && it.mandor.mandor) ? it.mandor.mandor : "Unknown";

      const tr = document.createElement("tr");
      tr.innerHTML = `
        <td>${safeText(mandor, "-")}</td>
        <td>${nik}</td>
        <td>${nama}</td>
        <td>${safeText(it.basahLatex, 0)}</td>
        <td>${safeText(it.basahLump, 0)}</td>
        <td>${safeText(it.sheet, 0)}</td>
        <td>${safeText(it.brCr, 0)}</td>
        <td>
          <button class="edit-btn" data-id="${String(it.id)}">
            <span class="action-icon"> 
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#0093E9" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 1 1 3 3L7 19l-4 1 1-4 12.5-12.5z"/></svg>
            </span>
          </button>
          <button class="delete-btn" data-id="${String(it.id)}">
            <span class="action-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#ff3b3b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2v2"/><line x1="10" y1="11" x2="10" y2="17"/><line x1="14" y1="11" x2="14" y2="17"/></svg>
            </span>
          </button>
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
        if (!item) { 
          alert("Data tidak ditemukan.");
          return; 
        }

        // Men-scroll halaman ke atas dengan animasi smooth
        window.scrollTo({
          top: 0,
          behavior: 'smooth' // Scroll dengan animasi smooth
        });

        // set mandor (select value = STRING)
        const mandorSelect = document.getElementById("mandor");
        if (mandorSelect) mandorSelect.value = String(item.idBakuMandor);

        // set penyadap & angka
        const inputNama = document.getElementById("namaPenyadap");
        const inputNik = document.getElementById("nik");
        const inputIdPenyadap = document.getElementById("idPenyadap");
        if (item.penyadap) {
          if (inputNama) inputNama.value = safeText(item.penyadap.nama_penyadap, "");
          if (inputNik) inputNik.value = safeText(item.penyadap.nik, "");
          if (inputIdPenyadap) inputIdPenyadap.value = String(item.penyadap.id || "");
        }
        const latek = document.getElementById("latek");
        const lump  = document.getElementById("lump");
        const sheet = document.getElementById("sheet");
        const brcr  = document.getElementById("brcr");
        if (latek) latek.value = item.basahLatex ?? 0;
        if (lump)  lump.value  = item.basahLump  ?? 0;
        if (sheet) sheet.value = item.sheet      ?? 0;
        if (brcr)  brcr.value  = item.brCr       ?? 0;

        editingId = id; // STRING
        const submitBtn = document.querySelector('button[type="submit"]');
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
  setInterval(renderRekapMandor, 300000);
});