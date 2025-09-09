document.getElementById("loginForm").addEventListener("submit", function(e) {
    e.preventDefault();

    const username = document.getElementById("username").value.trim();
    const password = document.getElementById("password").value.trim();

    // akun demo
    const validUser = "admin";
    const validPass = "12345";

    if (username === "" && password === "") {
        alert("Username dan Password wajib diisi!");
    } else if (username === "" && password !== "") {
        alert("Username wajib diisi!");
    } else if (username !== "" && password === "") {
        alert("Password wajib diisi!");
    } else if (username === validUser && password === validPass) {
        alert("Login berhasil!");
        window.location.href = "/dashboard"; // Fixed: gunakan route yang benar
    } else if (username === validUser && password !== validPass) {
        alert("Password salah!");
    } else if (username !== validUser && password === validPass) {
        alert("Username salah!");
    } else {
        alert("Username dan Password salah!");
    }
});

// 👁️ Toggle Password
document.getElementById("togglePassword").addEventListener("click", function() {
    const passwordField = document.getElementById("password");

    if (passwordField.type === "password") {
        passwordField.type = "text";
        this.textContent = "🙈"; // Ubah icon saat password terlihat
    } else {
        passwordField.type = "password";
        this.textContent = "👁️"; // Kembalikan icon mata
    }
});