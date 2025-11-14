function showChoiceModal() {
    document.getElementById('choiceModal').classList.add('active');
}

function closeModal(modalId) {
    document.getElementById(modalId).classList.remove('active');
}

function showUsernameForm() {
    closeModal('choiceModal');
    setTimeout(() => {
        document.getElementById('usernameModal').classList.add('active');
    }, 300);
}

function showPasswordForm() {
    closeModal('choiceModal');
    setTimeout(() => {
        document.getElementById('passwordModal').classList.add('active');
    }, 300);
}

function backToChoice(currentModal) {
    closeModal(currentModal);
    setTimeout(() => {
        showChoiceModal();
    }, 300);
}

async function handleUsernameSubmit(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    const data = {
        oldUsername: formData.get('oldUsername'),
        newUsername: formData.get('newUsername'),
        password: formData.get('passwordUsername')
    };
    
    console.log('Mengirim data ke:', '/api/manajemen/change-username');
    console.log('Data:', data);
    
    try {
        const response = await fetch('/api/manajemen/change-username', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify(data)
        });

        console.log('Response status:', response.status);
        
        const result = await response.json();
        console.log('Response data:', result);

        if (result.success) {
            alert('✓ ' + result.message + '\n\nUsername Lama: ' + data.oldUsername + '\nUsername Baru: ' + data.newUsername);
            event.target.reset();
            closeModal('usernameModal');
            
            // Optional: Redirect atau reload setelah berhasil
            setTimeout(() => {
                window.location.reload();
                // atau window.location.href = '/dashboard';
            }, 1000);
        } else {
            alert('✗ ' + result.message);
        }
    } catch (error) {
        console.error('Error:', error);
        alert('✗ Terjadi kesalahan saat mengubah username. Silakan coba lagi.');
    }
}

async function handlePasswordSubmit(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    const data = {
        username: formData.get('usernamePassword'),
        oldPassword: formData.get('oldPassword'),
        newPassword: formData.get('newPassword')
    };
    
    console.log('Mengirim data ke:', '/api/manajemen/change-password');
    console.log('Data:', data);
    
    try {
        const response = await fetch('/api/manajemen/change-password', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify(data)
        });

        console.log('Response status:', response.status);
        
        const result = await response.json();
        console.log('Response data:', result);

        if (result.success) {
            alert('✓ ' + result.message + '\n\nUsername: ' + data.username);
            event.target.reset();
            closeModal('passwordModal');
            
            // Optional: Redirect ke login setelah ganti password
            setTimeout(() => {
                alert('Silakan login kembali dengan password baru Anda.');
                window.location.href = '/login';
            }, 1000);
        } else {
            alert('✗ ' + result.message);
        }
    } catch (error) {
        console.error('Error:', error);
        alert('✗ Terjadi kesalahan saat mengubah password. Silakan coba lagi.');
    }
}

// Close modal saat klik di luar modal content
document.querySelectorAll('.modal-overlay').forEach(overlay => {
    overlay.addEventListener('click', function(e) {
        if (e.target === this) {
            this.classList.remove('active');
        }
    });
});

// Close modal dengan ESC key
document.addEventListener('keydown', function(e) {
    if (e.key === 'Escape') {
        document.querySelectorAll('.modal-overlay.active').forEach(modal => {
            modal.classList.remove('active');
        });
    }
});