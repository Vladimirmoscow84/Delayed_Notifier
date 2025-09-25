BEGIN;

CREATE TABLE IS NOT EXISTS notifications(
    id SERIAL PRIMARY KEY,
    notice_uid TEXT,
    body TEXT,
    date_created TIMESTAMP,
    send_date TIMESTAMP,
    send_atempts INT,
    send_status TEXT

);

COMMIT;