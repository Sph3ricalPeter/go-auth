CREATE TABLE users
(
    id       SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE refresh_tokens
(
    id    SERIAL PRIMARY KEY,
    token VARCHAR(8192) NOT NULL
);