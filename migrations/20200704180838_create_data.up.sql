CREATE TABLE securebin (
    id bigserial not null primary key,
    img bytea not null,
    encrypted_password varchar not null,
    lifetime timestamp without time zone DEFAULT (CURRENT_TIMESTAMP + '1 day'::interval) NOT NULL
);