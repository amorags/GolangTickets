const API_URL = 'http://localhost:8080';
const eventsContainer = document.getElementById('eventsContainer');

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

// Load events
async function loadEvents() {
    try {
        const response = await fetch(`${API_URL}/events`);

        if (response.ok) {
            const data = await response.json();
            displayEvents(data.data || []);
        } else {
            eventsContainer.innerHTML = '<div class="empty-state">Failed to load events</div>';
        }
    } catch (error) {
        console.error('Error loading events:', error);
        eventsContainer.innerHTML = '<div class="empty-state">Error loading events</div>';
    }
}

function displayEvents(events) {
    if (!events || events.length === 0) {
        eventsContainer.innerHTML = '<div class="empty-state">No events available</div>';
        return;
    }

    eventsContainer.innerHTML = events.map(event => {
        const eventDate = new Date(event.date).toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric',
            year: 'numeric'
        });

        return `
            <div class="event-card" onclick="window.location.href='/event?id=${event.ID}'">
                <div class="event-image" style="${event.image_url ? `background-image: url(${event.image_url}); background-size: cover;` : ''}"></div>
                <div class="event-content">
                    <div class="event-type">${event.event_type || 'event'}</div>
                    <div class="event-name">${event.name}</div>
                    <div class="event-info">
                        <div class="event-info-item">📅 ${eventDate}</div>
                        <div class="event-info-item">📍 ${event.venue_name}, ${event.city}</div>
                        <div class="event-info-item">👥 ${event.available_tickets}/${event.capacity} available</div>
                    </div>
                    <div class="event-price">$${event.price.toFixed(2)}</div>
                </div>
            </div>
        `;
    }).join('');
}

// Load events on page load
loadEvents();