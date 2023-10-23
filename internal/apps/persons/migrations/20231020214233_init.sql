-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

SET TIMEZONE='Europe/Moscow';

CREATE TABLE IF NOT EXISTS persons(
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    patronymic VARCHAR(255) NOT NULL,
    age INT CHECK (age>0),
    gender VARCHAR(100) CHECK(gender IN ('male', 'female')),
    nation VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS persons;
-- +goose StatementEnd
