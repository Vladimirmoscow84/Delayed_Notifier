BEGIN;

CREATE TABLE IF NOT EXISTS notifications(
    id SERIAL PRIMARY KEY,
    body TEXT,
    date_created TIMESTAMP,
    send_date TIMESTAMP,
    send_attempts INT,
    send_status TEXT

);

COMMIT;