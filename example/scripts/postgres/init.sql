CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    accommodation_id UUID NOT NULL,
    user_id UUID NOT NULL,
    start_at TIMESTAMP NOT NULL,
    end_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_dates CHECK (end_at > start_at)
);

CREATE INDEX idx_reservations_acommodation_id ON bookings(accommodation_id);
CREATE INDEX idx_reservations_user_id ON bookings(user_id);
CREATE INDEX idx_reservations_dates ON bookings(start_at, end_at);
