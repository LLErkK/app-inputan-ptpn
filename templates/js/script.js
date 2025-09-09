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
