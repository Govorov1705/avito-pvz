BEGIN;

CREATE TYPE pvz_city AS ENUM ('Москва', 'Санкт-Петербург', 'Казань');

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE pvzs(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    registration_date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    city pvz_city NOT NULL
);

CREATE TYPE reception_status AS ENUM ('in_progress', 'close');

CREATE TABLE receptions(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    pvz_id UUID NOT NULL REFERENCES pvzs(id),
    status reception_status NOT NULL DEFAULT 'in_progress'
);

CREATE UNIQUE INDEX one_open_reception_per_pvz
ON receptions(pvz_id)
WHERE status = 'in_progress';

CREATE TYPE product_type AS ENUM ('электроника', 'одежда', 'обувь');

CREATE TABLE products(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    type product_type NOT NULL,
    reception_id UUID NOT NULL REFERENCES receptions(id)
);

CREATE TYPE user_role AS ENUM ('employee', 'moderator');

CREATE TABLE users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(60) NOT NULL,
    role user_role NOT NULL
);

COMMIT;

