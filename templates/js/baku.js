document.addEventListener("DOMContentLoaded", () => {
    const form = document.getElementById("bakuForm");
    const tableBody = document.querySelector("#bakuTable tbody");

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

            tableBody.innerHTML = "";
            data.data.forEach(item => {
                const tr = document.createElement("tr");
                tr.innerHTML = `
          <td>${item.mandor ? item.mandor.mandor : "-"}</td>
          <td>${item.penyadap ? item.penyadap.nik : "-"}</td>
          <td>${item.penyadap ? item.penyadap.nama_penyadap : "-"}</td>
          <td>${item.basahLatex}</td>
          <td>${item.basahLump}</td>
          <td>${item.sheet}</td>
          <td>${item.brCr}</td>
        `;
                tableBody.appendChild(tr);
            });
        } catch (err) {
            console.error("Error loadTable:", err);
        }
    }

    // Initial load
    loadTable();
});
