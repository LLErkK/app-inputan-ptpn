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
            window.location.href = "/dashboard";
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

// Modern Password Toggle dengan Icon SVG
document.getElementById("togglePassword").addEventListener("click", function() {
    const passwordField = document.getElementById("password");
    const isPassword = passwordField.type === "password";
    
    // Toggle tipe input
    passwordField.type = isPassword ? "text" : "password";
    
    // Update icon dengan SVG yang lebih modern
    if (isPassword) {
        // Icon eye-off (password terlihat)
        this.innerHTML = `
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path>
                <line x1="1" y1="1" x2="23" y2="23"></line>
            </svg>
        `;
    } else {
        // Icon eye (password tersembunyi)
        this.innerHTML = `
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
                <circle cx="12" cy="12" r="3"></circle>
            </svg>
        `;
    }
    
    // Tambahkan class untuk animasi
    this.style.transform = "scale(0.9)";
    setTimeout(() => {
        this.style.transform = "scale(1)";
    }, 100);
});

// Menambahkan event listener untuk form input
document.getElementById('bakuForm')?.addEventListener('keydown', function(e) {
  if (e.key === 'Enter') {
    e.preventDefault();
    let currentElement = document.activeElement;
    let nextElement = null;

    if (currentElement && (currentElement.tagName === 'INPUT' || currentElement.tagName === 'SELECT')) {
      nextElement = getNextInput(currentElement);
    }

    if (nextElement) {
      nextElement.focus();
    }
  }
});

function getNextInput(currentElement) {
  let allInputs = document.querySelectorAll('#bakuForm input, #bakuForm select');
  let currentIndex = Array.prototype.indexOf.call(allInputs, currentElement);
  return allInputs[currentIndex + 1] || null;
}