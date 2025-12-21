const API_URL = 'http://localhost:8080';

// Check authentication
const token = localStorage.getItem('token');
if (!token) {
    window.location.href = '/';
}

// DOM elements
const profileContent = document.getElementById('profileContent');
const bookingsContent = document.getElementById('bookingsContent');
const logoutBtn = document.getElementById('logoutBtn');

// Fetch profile on load
async function loadProfile() {
    try {
        const response = await fetch(`${API_URL}/profile`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        });

        if (response.ok) {
            const data = await response.json();
            displayProfile(data.data);
        } else {
            // Token might be invalid
            logout();
        }
    } catch (error) {
        console.error('Error loading profile:', error);
        profileContent.innerHTML = '<div class="error">Failed to load profile</div>';
    }
}

// Fetch bookings
async function loadBookings() {
    try {
        const response = await fetch(`${API_URL}/bookings`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        });

        if (response.ok) {
            const data = await response.json();
            displayBookings(data.data || []);
        } else {
            bookingsContent.innerHTML = '<div class="empty-state">Failed to load bookings</div>';
        }
    } catch (error) {
        console.error('Error loading bookings:', error);
        bookingsContent.innerHTML = '<div class="empty-state">Error loading bookings</div>';
    }
}

function displayProfile(user) {
    const createdDate = new Date(user.created_at).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });

    profileContent.innerHTML = `
        <div class="profile-info">
            <div class="profile-field">
                <span class="profile-label">Username:</span>
                <span class="profile-value">${user.username}</span>
            </div>
            <div class="profile-field">
                <span class="profile-label">Email:</span>
                <span class="profile-value">${user.email}</span>
            </div>
            <div class="profile-field">
                <span class="profile-label">User ID:</span>
                <span class="profile-value">${user.id}</span>
            </div>
            <div class="profile-field">
                <span class="profile-label">Member Since:</span>
                <span class="profile-value">${createdDate}</span>
            </div>
        </div>
    `;
}

function displayBookings(bookings) {
    if (!bookings || bookings.length === 0) {
        bookingsContent.innerHTML = '<div class="empty-state">No bookings yet. <a href="/home">Browse events</a></div>';
        return;
    }

    bookingsContent.innerHTML = `
        <div class="bookings-grid">
            ${bookings.map(booking => {
                const bookingDate = new Date(booking.CreatedAt).toLocaleDateString('en-US', {
                    month: 'short',
                    day: 'numeric',
                    year: 'numeric'
                });

                return `
                    <div class="booking-card">
                        <div class="booking-info">
                            <h4>${booking.event_name || 'Event'}</h4>
                            <div class="booking-details">
                                <div>Quantity: ${booking.quantity} tickets</div>
                                <div>Total: $${booking.total_price.toFixed(2)}</div>
                                <div>Booked: ${bookingDate}</div>
                                <div>Status: ${booking.status}</div>
                            </div>
                        </div>
                        <div class="booking-actions">
                            ${booking.status === 'confirmed' ?
                                `<button class="btn-small btn-danger" onclick="cancelBooking(${booking.ID})">Cancel</button>`
                                : ''
                            }
                        </div>
                    </div>
                `;
            }).join('')}
        </div>
    `;
}

async function cancelBooking(bookingId) {
    if (!confirm('Are you sure you want to cancel this booking?')) {
        return;
    }

    try {
        const response = await fetch(`${API_URL}/bookings/${bookingId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        });

        if (response.ok) {
            alert('Booking cancelled successfully');
            loadBookings(); // Reload bookings
        } else {
            const data = await response.json();
            alert(data.error || 'Failed to cancel booking');
        }
    } catch (error) {
        console.error('Error cancelling booking:', error);
        alert('Network error. Please try again.');
    }
}

function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = '/';
}

// Logout button handler
logoutBtn.addEventListener('click', logout);

// Load profile and bookings on page load
loadProfile();
loadBookings();