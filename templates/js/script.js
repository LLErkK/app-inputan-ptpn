document.getElementById("loginForm").addEventListener("submit", async function(e) {
    e.preventDefault();

    const username = document.getElementById("username").value.trim();
    const password = document.getElementById("password").value.trim();
    const errorDiv = document.getElementById("errorMessage");

    // Reset error message
    errorDiv.style.display = "none";

    // Basic validation
    if (username === "" || password === "") {
        showError("Username dan Password wajib diisi!");
        return;
    }

    // Disable button while processing
    const submitBtn = document.querySelector(".btn-login");
    const originalText = submitBtn.textContent;
    submitBtn.textContent = "LOADING...";
    submitBtn.disabled = true;

    try {
        const response = await fetch('/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                username: username,
                password: password
            })
        });

        const data = await response.json();

        if (data.success) {
            // Login berhasil, redirect ke dashboard
            window.location.href = "/dashboard";
        } else {
            // Login gagal, tampilkan pesan error
            showError(data.message || "Login gagal!");
        }

    } catch (error) {
        console.error('Error:', error);
        showError("Terjadi kesalahan koneksi!");
    } finally {
        // Re-enable button
        submitBtn.textContent = originalText;
        submitBtn.disabled = false;
    }
});

// Function to show error message
function showError(message) {
    const errorDiv = document.getElementById("errorMessage");
    errorDiv.textContent = message;
    errorDiv.style.display = "block";
}

// Toggle Password Visibility
document.getElementById("togglePassword").addEventListener("click", function() {
    const passwordField = document.getElementById("password");
    const toggleBtn = this;

    if (passwordField.type === "password") {
        passwordField.type = "text";
        toggleBtn.textContent = "🙈";
    } else {
        passwordField.type = "password";
        toggleBtn.textContent = "👁️";
    }
});