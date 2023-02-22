CREATE TABLE IF NOT EXISTS blacklisted_users (
     id BIGSERIAL PRIMARY KEY,
     phone TEXT NOT NULL,
     username TEXT NOT NULL,
     cause TEXT NOT NULL,
     timestamp TIMESTAMP NOT NULL,
     caused_by TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
     id BIGSERIAL PRIMARY KEY,
     username TEXT NOT NULL,
     password TEXT NOT NULL
);

-- Add test user with credentials: admin/admin (encoded)
INSERT INTO users (username, password)
VALUES ('admin', '646173646a6173646a61736a6a323133d033e22ae348aeb5660fc2140aec35850c4da997')
