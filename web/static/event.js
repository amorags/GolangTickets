const API_URL = 'http://localhost:8080';
const eventDetails = document.getElementById('eventDetails');

// Get event ID from URL
const urlParams = new URLSearchParams(window.location.search);
const eventId = urlParams.get('id');

// Update nav based on auth status
const token = localStorage.getItem('token');
if (token) {
    document.getElementById('authLink').textContent = 'Logout';
    document.getElementById('authLink').href = '#';
    document.getElementById('authLink').addEventListener('click', (e) => {
        e.preventDefault();
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.href = '/';
    });
} else {
    document.getElementById('profileLink').style.display = 'none';
}

if (!eventId) {
    window.location.href = '/home';
}

// Load event details
async function loadEvent() {
    try {
        const response = await fetch(`${API_URL}/events/${eventId}`);

        if (response.ok) {
            const data = await response.json();
            displayEvent(data.data);
        } else {
            eventDetails.innerHTML = '<div class="empty-state">Event not found</div>';
        }
    } catch (error) {
        console.error('Error loading event:', error);
        eventDetails.innerHTML = '<div class="empty-state">Error loading event</div>';
    }
}

function displayEvent(event) {
    const eventDate = new Date(event.date).toLocaleDateString('en-US', {
        weekday: 'long',
        month: 'long',
        day: 'numeric',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });

    const bookingSection = token ? `
        <div class="booking-section">
            <h3>Book Tickets</h3>
            <div id="bookingMessage" class="message"></div>
            <form id="bookingForm" class="booking-form">
                <div class="form-group">
                    <label for="quantity">Number of Tickets</label>
                    <input type="number" id="quantity" min="1" max="${event.available_tickets}" value="1" required>
                </div>
                <button type="submit" class="btn-primary">Book Now - $${event.price.toFixed(2)} each</button>
            </form>
        </div>
    ` : `
        <div class="booking-section">
            <h3>Book Tickets</h3>
            <p>Please <a href="/">login</a> to book tickets</p>
        </div>
    `;

    eventDetails.innerHTML = `
        <div class="event-header">
            <div class="event-detail-image" style="${event.image_url ? `background-image: url(${event.image_url}); background-size: cover;` : ''}"></div>
            <h1 class="event-title">${event.name}</h1>
            <div class="event-type">${event.event_type || 'event'}</div>
        </div>

        <div class="event-meta">
            <div>📅 ${eventDate}</div>
            <div>📍 ${event.venue_name}</div>
            <div>🏙️ ${event.city}, ${event.address}</div>
            <div>👥 ${event.available_tickets}/${event.capacity} tickets available</div>
            <div>💰 $${event.price.toFixed(2)}</div>
        </div>

        <div class="event-description">
            <h3>About this event</h3>
            <p>${event.description || 'No description available'}</p>
        </div>

        ${bookingSection}
    `;

    // Add booking form handler if user is logged in
    if (token) {
        document.getElementById('bookingForm').addEventListener('submit', handleBooking);
    }
}

async function handleBooking(e) {
    e.preventDefault();

    const quantity = parseInt(document.getElementById('quantity').value);
    const messageDiv = document.getElementById('bookingMessage');

    try {
        const response = await fetch(`${API_URL}/bookings`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                event_id: parseInt(eventId),
                quantity: quantity
            })
        });

        const data = await response.json();

        if (response.ok) {
            messageDiv.textContent = 'Booking successful! Redirecting to your profile...';
            messageDiv.className = 'message success show';

            setTimeout(() => {
                window.location.href = '/profile-page';
            }, 2000);
        } else {
            messageDiv.textContent = data.error || 'Booking failed';
            messageDiv.className = 'message error show';
        }
    } catch (error) {
        console.error('Booking error:', error);
        messageDiv.textContent = 'Network error. Please try again.';
        messageDiv.className = 'message error show';
    }
}

// Load event on page load
loadEvent();