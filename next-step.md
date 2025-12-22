# Next Steps: Feature & Refactoring Roadmap

This document outlines the immediate next steps for the project, formatted as requirement specifications suitable for Jira tickets. These tasks focus on high-impact UX improvements and critical backend feature additions.

## Epic 1: Frontend User Experience Polish
**Goal:** Transform the raw functional UI into a professional, responsive, and user-friendly experience.

### Ticket FE-001: Implement Skeleton Loading States
**Priority:** High
**Description:** Currently, the application shows blank screens or abrupt content shifts while fetching data from the API. We need to implement skeleton loaders to indicate progress and improve perceived performance.
**Scope:**
- `web/pages/index.vue` (Event list)
- `web/pages/events/[id].vue` (Event details)
- `web/pages/profile.vue` (User profile & bookings)
**Acceptance Criteria:**
1.  **Event List:** Display a grid of 3-6 skeleton cards (image placeholder + title bar + text bar) while the `GET /events` request is pending.
2.  **Event Details:** Display a skeleton layout matching the event detail page structure while data loads.
3.  **Transitions:** When data arrives, the skeleton should fade out and actual content should fade in (or switch instantly without layout shift).
4.  **No "Flash of Unstyled Content":** Ensure skeletons appear immediately on mount.

### Ticket FE-002: Enhanced Error Handling & Toast Notifications
**Priority:** High
**Description:** Error handling is currently minimal. Users need clear visual feedback when actions succeed or fail.
**Scope:** Global application (Plugin/Component)
**Acceptance Criteria:**
1.  **Toast Component:** Create a global toast/notification component that can be triggered from anywhere.
    - Types: Success (Green), Error (Red), Info (Blue).
    - Behavior: Auto-dismiss after 5 seconds.
2.  **API Integration:** Update `useApi` composable to automatically trigger error toasts on non-200 responses (unless handled locally).
3.  **Specific Scenarios:**
    - **Login Failed:** Show "Invalid email or password".
    - **Booking Success:** Show "Booking confirmed! Check your profile."
    - **Network Error:** Show "Unable to connect to server."

### Ticket FE-003: Search and Filter UI
**Priority:** Medium
**Description:** Users currently have to scroll through all events to find what they want. Add a search bar and basic filters.
**Scope:** `web/pages/index.vue`
**Acceptance Criteria:**
1.  **Search Bar:** Add a text input to filter events by title or description.
2.  **Date Filter:** Add "From" and "To" date pickers to filter events by date range.
3.  **Reactivity:** The event list should update automatically (or upon clicking "Search") based on inputs.
4.  **Empty State:** Display a friendly "No events found matching your criteria" message if results are empty.

---

## Epic 2: Backend Core Features
**Goal:** Add essential functionality that depends on backend processing and robust data handling.

### Ticket BE-001: Email Notification Service
**Priority:** High
**Description:** The system currently lacks any communication channel with users. We need to send transactional emails for account actions.
**Technical Recommendation:** Use `gomail` for sending. For scalability, consider integrating a Redis-based job queue (e.g., `asynq`) later, but start with a simple async goroutine for the MVP.
**Acceptance Criteria:**
1.  **Setup:** Configure SMTP settings (use Mailtrap or similar for dev/test) in `config.go`.
2.  **Welcome Email:** Trigger an email to the user upon successful Signup (`POST /signup`).
3.  **Booking Confirmation:** Trigger an email containing Event Name, Date, and Quantity upon successful Booking (`POST /bookings`).
4.  **Template:** Use basic HTML templates for the emails (not just plain text).

### Ticket BE-002: Backend Search & Filter API
**Priority:** Medium
**Description:** Support the frontend search requirements by updating the API to handle query parameters.
**Scope:** `internal/handlers/event_handler.go` -> `GetEvents`
**Acceptance Criteria:**
1.  **Query Parameters:** Update `GetEvents` to accept optional query params:
    - `q` (string): Search term for partial match on title/description.
    - `start_date` (string, RFC3339): Filter events starting after this date.
    - `end_date` (string, RFC3339): Filter events starting before this date.
2.  **Logic:** Build the SQL query dynamically based on provided parameters.
3.  **Validation:** Return 400 Bad Request if date formats are invalid.

### Ticket BE-003: Redis Caching for Event List
**Priority:** Medium
**Description:** The `/events` endpoint is read-heavy and data doesn't change every second. Cache the response to improve performance.
**Technical Recommendation:** Use `go-redis/redis`.
**Acceptance Criteria:**
1.  **Cache Strategy:** Cache the output of `GetEvents` for 1-5 minutes.
2.  **Invalidation:** Invalidate (delete) the cache key whenever an event is Created, Updated, or Deleted.
3.  **Fallback:** If Redis is down, the system should gracefully fall back to querying the database directly.

---

## Epic 3: Technical Debt & Refactoring
**Goal:** Improve code quality and performance to ensure scalability.

### Ticket TD-001: Fix N+1 Query in `GetEvents`
**Priority:** Critical
**Description:** The current `GetEvents` implementation likely iterates through events and queries available tickets for *each* event individually (N+1 problem).
**Scope:** `internal/repository/event_repository.go`
**Acceptance Criteria:**
1.  **Query Optimization:** Rewrite the database query to fetch events and their booking counts/availability in a single SQL query (using `JOIN` or `GROUP BY`).
2.  **Verification:** Verify that loading the event list triggers only 1 database query instead of 1 + N queries.

### Ticket TD-002: Security - Fix Delete Authorization
**Priority:** Critical
**Description:** Currently, it appears any authenticated user might be able to delete events. Only the event creator or an admin should be able to do this.
**Scope:** `internal/handlers/event_handler.go` -> `DeleteEvent`
**Acceptance Criteria:**
1.  **Check:** Before deletion, query the event to check the `creator_id`.
2.  **Verify:** Compare `creator_id` with the `user_id` from the JWT context.
3.  **Enforce:** Return `403 Forbidden` if they do not match.
