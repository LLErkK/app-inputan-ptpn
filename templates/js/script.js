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
        // window.location.href = "dashboard.html";
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
    const eyeOpen = document.getElementById("icon-eye");
    const eyeOff = document.getElementById("icon-eye-off");

    if (passwordField.type === "password") {
        passwordField.type = "text";
        eyeOpen.style.display = "none";
        eyeOff.style.display = "block";
    } else {
        passwordField.type = "password";
        eyeOpen.style.display = "block";
        eyeOff.style.display = "none";
    }
});
