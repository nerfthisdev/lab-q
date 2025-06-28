CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    username TEXT,
    is_admin BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS subjects (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS subject_schedule (
    id SERIAL PRIMARY KEY,
    subject_id INT NOT NULL REFERENCES subjects(id),
    day_of_week INT NOT NULL,
    time_of_day TIME NOT NULL,
    start_date DATE NOT NULL,
    interval_weeks INT NOT NULL
);

CREATE TABLE IF NOT EXISTS subject_queue (
    subject_id INT NOT NULL REFERENCES subjects(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY(subject_id, user_id)
);
