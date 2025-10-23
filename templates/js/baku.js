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
  
  // Fungsi untuk format angka: hilangkan .00, tampilkan max 2 desimal
  function formatNumber(value) {
    const num = parseFloat(value);
    if (!Number.isFinite(num)) return "0";
    
    // Jika angka bulat (tidak ada desimal), tampilkan tanpa desimal
    if (num === Math.floor(num)) {
      return num.toString();
    }
    
    // Jika ada desimal, tampilkan maksimal 2 digit dan hapus trailing zeros
    return num.toFixed(2).replace(/\.?0+$/, '');
  }

  let allBakuData = [];
  let editingId = null;
  let mandorDataCache = [];
  let penyadapDataCache = [];

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
        mandorTableBody.innerHTML = `<tr><td colspan="6">Tidak ada data mandor.</td></tr>`;
      }
    } catch (e) {
      mandorTableBody.innerHTML = `<tr><td colspan="6">Error: ${e.message}</td></tr>`;
    }
  }

  async function loadMandorOptions() {
    const input = document.getElementById("mandor");
    const datalist = document.getElementById("mandor-list");
    if (!input || !datalist) return;
    
    datalist.innerHTML = '';
    
    try {
      const res = await fetch("/api/mandor");
      const data = await res.json();
      if (data.success && Array.isArray(data.data)) {
        mandorDataCache = data.data;

        data.data.forEach(m => {
          const option = document.createElement("option");
          option.value = `${safeText(m.mandor)} (${safeText(m.afdeling, "-")}) ${safeText(m.tahun_tanam, "-")} - Tipe: ${safeText(m.tipe)}`;
          option.setAttribute('data-id', m.id);
          datalist.appendChild(option);
        });
      }
    } catch (e) {
      console.error("Error loading mandor:", e);
    }
  }

  const inputMandor = document.getElementById("mandor");
  const datalistMandor = document.getElementById("mandor-list");

  if (inputMandor && datalistMandor) {
    inputMandor.addEventListener("input", function() {
      const value = this.value;
      const options = datalistMandor.querySelectorAll('option');
      
      options.forEach(option => {
        if (option.value === value) {
          const id = option.getAttribute('data-id');
          this.setAttribute('data-mandor-id', id);
        }
      });
    });
    
    inputMandor.addEventListener("change", function() {
      const value = this.value;
      const options = datalistMandor.querySelectorAll('option');
      
      options.forEach(option => {
        if (option.value === value) {
          const id = option.getAttribute('data-id');
          this.setAttribute('data-mandor-id', id);
        }
      });
    });
  }

  function getMandorById(id) {
    return mandorDataCache.find(m => String(m.id) === String(id));
  }

  // ========= Popup Penyadap =========
  const tambahPenyadapBtn = document.getElementById("tambahPenyadap");
  const popupPenyadap = document.getElementById("popupPenyadap");
  const closePopupPenyadapBtn = document.getElementById("closePopupPenyadap");
  const formPenyadapBaru = document.getElementById("formPenyadapBaru");
  const penyadapTableBody = document.getElementById("penyadapTableBody");

  if (tambahPenyadapBtn) {
    tambahPenyadapBtn.addEventListener("click", () => {
      if (popupPenyadap) popupPenyadap.style.display = "flex";
      loadPenyadapList();
    });
  }

  if (closePopupPenyadapBtn) {
    closePopupPenyadapBtn.addEventListener("click", () => {
      if (popupPenyadap) popupPenyadap.style.display = "none";
    });
  }

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
          loadPenyadapOptions();
        } else {
          alert("Gagal: " + (data.message || "Terjadi kesalahan"));
        }
      } catch (err) {
        console.error("Error:", err);
        alert("Error: " + err.message);
      }
    });
  }

  async function loadPenyadapList() {
    if (!penyadapTableBody) return;
    
    penyadapTableBody.innerHTML = '<tr><td colspan="3">Memuat data...</td></tr>';
    
    try {
      const res = await fetch("/api/penyadap");
      
      if (!res.ok) {
        throw new Error(`HTTP ${res.status}`);
      }
      
      const data = await res.json();
      penyadapTableBody.innerHTML = "";
      
      if (data.success && Array.isArray(data.data) && data.data.length > 0) {
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
        
        penyadapTableBody.querySelectorAll(".delete-penyadap-btn").forEach(btn => {
          btn.addEventListener("click", async function() {
            const id = String(this.getAttribute("data-id"));
            
            if (confirm("Hapus penyadap ini?")) {
              try {
                const res = await fetch(`/api/penyadap/${encodeURIComponent(id)}`, { 
                  method: "DELETE" 
                });
                
                const data = await res.json();
                
                if (res.ok && data.success) {
                  alert("Penyadap berhasil dihapus!");
                  loadPenyadapList();
                  loadPenyadapOptions();
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
        penyadapTableBody.innerHTML = `<tr><td colspan="3">Belum ada data penyadap.</td></tr>`;
      }
    } catch (e) {
      console.error("Error loading penyadap:", e);
      penyadapTableBody.innerHTML = `<tr><td colspan="3" style="color: red;">Error: ${safeText(e.message)}</td></tr>`;
    }
  }

  async function loadPenyadapOptions() {
    const input = document.getElementById("namaPenyadap");
    const datalist = document.getElementById("penyadap-list");
    if (!input || !datalist) return;
    
    datalist.innerHTML = '';
    
    try {
      const res = await fetch("/api/penyadap");
      const data = await res.json();
      if (data.success && Array.isArray(data.data)) {
        penyadapDataCache = data.data;

        data.data.forEach(p => {
          const option = document.createElement("option");
          option.value = `${safeText(p.nama_penyadap)} - NIK: ${safeText(p.nik)}`;
          option.setAttribute('data-id', p.id);
          option.setAttribute('data-nama', p.nama_penyadap);
          option.setAttribute('data-nik', p.nik);
          datalist.appendChild(option);
        });
      }
    } catch (e) {
      console.error("Error loading penyadap options:", e);
    }
  }

  const inputNamaPenyadap = document.getElementById("namaPenyadap");
  const inputNik = document.getElementById("nik");
  const inputIdPenyadap = document.getElementById("idPenyadap");
  const datalistPenyadap = document.getElementById("penyadap-list");

  if (inputNamaPenyadap && datalistPenyadap) {
    inputNamaPenyadap.addEventListener("input", function() {
      const value = this.value;
      const options = datalistPenyadap.querySelectorAll('option');
      let found = false;
      
      options.forEach(option => {
        if (option.value === value) {
          const id = option.getAttribute('data-id');
          const nama = option.getAttribute('data-nama');
          const nik = option.getAttribute('data-nik');
          
          if (inputIdPenyadap) inputIdPenyadap.value = id || "";
          if (inputNik) inputNik.value = nik || "";
          
          found = true;
        }
      });
      
      if (!found) {
        if (inputIdPenyadap) inputIdPenyadap.value = "";
        if (inputNik) inputNik.value = "";
      }
    });
    
    inputNamaPenyadap.addEventListener("change", function() {
      const value = this.value;
      const options = datalistPenyadap.querySelectorAll('option');
      
      options.forEach(option => {
        if (option.value === value) {
          const id = option.getAttribute('data-id');
          const nama = option.getAttribute('data-nama');
          const nik = option.getAttribute('data-nik');
          
          if (inputIdPenyadap) inputIdPenyadap.value = id || "";
          if (inputNik) inputNik.value = nik || "";
          this.value = nama;
        }
      });
    });
  }

  function getPenyadapById(id) {
    return penyadapDataCache.find(p => String(p.id) === String(id));
  }

  // ========= Submit Form (Create/Update) =========
  const form = document.getElementById("bakuForm");
  if (form) {
    form.addEventListener("submit", async e => {
      e.preventDefault();
      
      const mandorInput = document.getElementById("mandor");
      const idMandorStr = mandorInput ? mandorInput.getAttribute('data-mandor-id') : "";
      const idPenyadapStr = inputIdPenyadap ? inputIdPenyadap.value : "";

      const selectedMandor = getMandorById(idMandorStr);
      const tahunTanam = selectedMandor ? selectedMandor.tahun_tanam : null;

      const payload = {
        idBakuMandor: idMandorStr ? parseInt(idMandorStr) : null,
        idPenyadap: idPenyadapStr ? parseInt(idPenyadapStr) : null,
        tahunTanam: tahunTanam,
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
          if (mandorInput) mandorInput.removeAttribute('data-mandor-id');
          const submitBtn = document.querySelector('button[type="submit"]');
          if (submitBtn) submitBtn.textContent = "Save";
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

  // ========= Rekap Mandor (DENGAN FORMAT NUMBER) =========
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
            <td>${formatNumber(row.jumlah_pabrik_basah_latek)}</td>
            <td>${formatNumber(row.jumlah_kebun_basah_latek)}</td>
            <td>${formatNumber(row.jumlah_sheet)}</td>
            <td>${formatNumber(row.k3_sheet)}</td>
            <td>${formatNumber(row.jumlah_pabrik_basah_lump)}</td>
            <td>${formatNumber(row.jumlah_kebun_basah_lump)}</td>
            <td>${formatNumber(row.jumlah_br_cr)}</td>
            <td>${formatNumber(row.k3_br_cr)}</td>`;
          body.appendChild(tr);
        });
      } else {
        body.innerHTML = `<tr><td colspan="11">Belum ada data rekap.</td></tr>`;
      }
    } catch (e) {
      body.innerHTML = `<tr><td colspan="11">Gagal memuat rekap.</td></tr>`;
    }
  }

  // ========= Detail Baku (DENGAN FORMAT NUMBER) =========
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
      // biarkan kosong
    }

    allBakuData = Array.isArray(dataArr) ? dataArr : [];

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
          <td>${formatNumber(it.basahLatex)}</td>
          <td>${formatNumber(it.basahLump)}</td>
          <td>${formatNumber(it.sheet)}</td>
          <td>${formatNumber(it.brCr)}</td>
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

        window.scrollTo({
          top: 0,
          behavior: 'smooth'
        });

        const mandorInput = document.getElementById("mandor");
        const selectedMandor = getMandorById(String(item.idBakuMandor));
        if (mandorInput && selectedMandor) {
          const mandorText = `${safeText(selectedMandor.mandor)} (${safeText(selectedMandor.afdeling, "-")}) ${safeText(selectedMandor.tahun_tanam, "-")} - Tipe: ${safeText(selectedMandor.tipe)}`;
          mandorInput.value = mandorText;
          mandorInput.setAttribute('data-mandor-id', String(item.idBakuMandor));
        }

        if (item.penyadap) {
          if (inputNamaPenyadap) inputNamaPenyadap.value = safeText(item.penyadap.nama_penyadap, "");
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

        editingId = id;
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
  loadPenyadapOptions();
  renderRekapMandor();
  renderDetailBaku();
  setInterval(renderRekapMandor, 300000);
});