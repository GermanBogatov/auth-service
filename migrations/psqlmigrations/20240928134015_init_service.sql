-- +goose Up
-- +goose StatementBegin

CREATE TYPE roleType AS ENUM (
    'super-admin',
    'admin',
    'user'
    );

CREATE TABLE IF NOT EXISTS users (
    id                  UUID NOT NULL PRIMARY KEY,
    name                VARCHAR(255) NOT NULL,
    surname             VARCHAR(255) NOT NULL,
    email               VARCHAR(255) UNIQUE NOT NULL,
    password            VARCHAR(255) NOT NULL,
    role                roleType NOT NULL,
    created_date        TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    updated_date        TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_email_password
    ON users(email,password);

CREATE INDEX IF NOT EXISTS idx_users_email
    ON users(email);

INSERT INTO users (id,name,surname,email,password,role,created_date) VALUES ('c1cfe4b9-f7c2-423c-abfa-6ed1c05a15c5','Герман','Богатов','bogatovgrmn@gmail.com','7361643334326d736c6664323334313273646673646631323334686766f6ee94ecb014f74f887b9dcc52daecf73ab3e3333320cadd98bcb59d895c52f5','super-admin',current_timestamp);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_users_email;
DROP TABLE users;
DROP TYPE roleType;
-- +goose StatementEnd
