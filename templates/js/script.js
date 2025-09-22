document.getElementById("loginForm").addEventListener("submit", async function(e) {
    e.preventDefault();

    const username = document.getElementById("username").value.trim();
    const password = document.getElementById("password").value.trim();

    if (username === "" && password === "") {
        alert("Username dan Password wajib diisi!");
        return;
    } else if (username === "" && password !== "") {
        alert("Username wajib diisi!");
        return;
    } else if (username !== "" && password === "") {
        alert("Password wajib diisi!");
        return;
    }

    try {
        const res = await fetch("/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ username, password }),
        });

        const data = await res.json();

        if (data.success) {
            alert("Login berhasil!");
            window.location.href = "/dashboard"; // redirect ke dashboard
        } else {
            alert("Login gagal: " + data.message);
        }
    } catch (err) {
        console.error("Error login:", err);
        alert("Terjadi kesalahan pada server.");
    }
});

document.addEventListener('DOMContentLoaded', function() {
    var btnLogout = document.getElementById('btn-logout');
    if (btnLogout) {
        btnLogout.addEventListener('click', function() {
            window.top.location.href = '/logout';
        });
    }
});

// 👁️ Toggle Password
document.getElementById("togglePassword").addEventListener("click", function() {
    const passwordField = document.getElementById("password");

    if (passwordField.type === "password") {
        passwordField.type = "text";
        this.textContent = "🙈";
    } else {
        passwordField.type = "password";
        this.textContent = "👁️";
    }
});

// Menambahkan event listener untuk form input
document.getElementById('bakuForm').addEventListener('keydown', function(e) {
  if (e.key === 'Enter') {
    // Mencegah form disubmit saat menekan Enter
    e.preventDefault();

    // Mendapatkan elemen input yang sedang aktif
    let currentElement = document.activeElement;
    let nextElement = null;

    // Mengecek apakah elemen yang aktif adalah input atau select
    if (currentElement && (currentElement.tagName === 'INPUT' || currentElement.tagName === 'SELECT')) {
      // Mendapatkan elemen input/select berikutnya
      nextElement = getNextInput(currentElement);
    }

    // Fokus ke elemen berikutnya jika ada
    if (nextElement) {
      nextElement.focus();
    }
  }
});

// Fungsi untuk mendapatkan elemen input berikutnya
function getNextInput(currentElement) {
  let allInputs = document.querySelectorAll('#bakuForm input, #bakuForm select');
  let currentIndex = Array.prototype.indexOf.call(allInputs, currentElement);
  return allInputs[currentIndex + 1] || null;  // Mengembalikan input berikutnya
}
