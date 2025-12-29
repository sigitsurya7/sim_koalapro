CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS m_user (
    id BIGSERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid(),
    username VARCHAR(100) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role VARCHAR(10) NOT NULL CHECK (role IN ('admin', 'viewer')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(100),
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR(100),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(100),
    last_seen TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_m_user_uid ON m_user (uid);
CREATE INDEX IF NOT EXISTS idx_m_user_username ON m_user (username);

CREATE TABLE IF NOT EXISTS t_user_token (
    id BIGSERIAL PRIMARY KEY,
    user_uid UUID NOT NULL REFERENCES m_user(uid),
    token TEXT NOT NULL,
    last_seen TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_t_user_token_user_uid ON t_user_token (user_uid);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'e_jenis_pengguna') THEN
        CREATE TYPE e_jenis_pengguna AS ENUM ('stockity', 'binomo', 'olymptrade');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS m_pengguna (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    id_pengguna BIGINT NOT NULL,
    telegram VARCHAR(100),
    jenis e_jenis_pengguna NOT NULL DEFAULT 'stockity',
    active BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

ALTER TABLE m_pengguna
    ADD COLUMN IF NOT EXISTS jenis e_jenis_pengguna NOT NULL DEFAULT 'stockity';

CREATE UNIQUE INDEX IF NOT EXISTS idx_m_pengguna_uuid ON m_pengguna (uuid);
CREATE UNIQUE INDEX IF NOT EXISTS idx_m_pengguna_id_pengguna ON m_pengguna (id_pengguna);

CREATE TABLE IF NOT EXISTS t_pengguna_detail (
    id BIGINT PRIMARY KEY,
    pengguna_id BIGINT NOT NULL REFERENCES m_pengguna(id),
    avatar TEXT,
    first_name TEXT,
    last_name TEXT,
    nickname TEXT,
    balance NUMERIC NOT NULL DEFAULT 0,
    balance_version BIGINT NOT NULL DEFAULT 0,
    bonus NUMERIC NOT NULL DEFAULT 0,
    gender TEXT,
    email TEXT,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    phone TEXT,
    phone_verified BOOLEAN NOT NULL DEFAULT FALSE,
    phone_prefix TEXT,
    receive_news BOOLEAN NOT NULL DEFAULT TRUE,
    receive_sms BOOLEAN NOT NULL DEFAULT TRUE,
    receive_notification BOOLEAN NOT NULL DEFAULT TRUE,
    country TEXT,
    country_name TEXT,
    currency TEXT,
    birthday TEXT,
    activate BOOLEAN NOT NULL DEFAULT TRUE,
    password_is_set BOOLEAN NOT NULL DEFAULT FALSE,
    tutorial BOOLEAN NOT NULL DEFAULT FALSE,
    coupons JSONB,
    free_deals JSONB,
    blocked BOOLEAN NOT NULL DEFAULT FALSE,
    agree_risk BOOLEAN NOT NULL DEFAULT FALSE,
    agreed BOOLEAN NOT NULL DEFAULT TRUE,
    status_group TEXT,
    docs_verified BOOLEAN NOT NULL DEFAULT FALSE,
    registered_at TIMESTAMPTZ,
    status_by_deposit TEXT,
    status_id INTEGER,
    deposits_sum NUMERIC NOT NULL DEFAULT 0,
    push_notification_categories JSONB,
    preserve_name BOOLEAN NOT NULL DEFAULT FALSE,
    registration_country_iso TEXT
);

CREATE INDEX IF NOT EXISTS idx_t_pengguna_detail_pengguna_id ON t_pengguna_detail (pengguna_id);
