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

CREATE TABLE IF NOT EXISTS subject_queue (
    subject_id INT NOT NULL REFERENCES subjects(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY(subject_id, user_id)
);
